// Copyright 2023 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tektoncompiler

import (
	"encoding/json"
	"fmt"

	"github.com/kubeflow/pipelines/api/v2alpha1/go/pipelinespec"
	"github.com/kubeflow/pipelines/backend/src/v2/compiler"
	"github.com/kubeflow/pipelines/backend/src/v2/tekton-kfptask/apis/kfptask"
	ktv1alpha1 "github.com/kubeflow/pipelines/backend/src/v2/tekton-kfptask/apis/kfptask/v1alpha1"
	pipelineapi "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"google.golang.org/protobuf/types/known/structpb"
	k8score "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	volumeNameKFPLauncher = "kfp-launcher"
	kfpLauncherPath       = "/tekton/home/launch"
)

// add KubernetesSpec for the container of the component
func (c *pipelinerunCompiler) AddKubernetesSpec(name string, kubernetesSpec *structpb.Struct) error {
	err := c.saveKubernetesSpec(name, kubernetesSpec)
	if err != nil {
		return err
	}
	return nil
}

func (c *pipelinerunCompiler) Container(taskName, compRef string,
	task *pipelinespec.PipelineTaskSpec,
	component *pipelinespec.ComponentSpec,
	container *pipelinespec.PipelineDeploymentConfig_PipelineContainerSpec,
) error {

	err := c.saveComponentSpec(compRef, component)
	if err != nil {
		return err
	}
	err = c.saveComponentImpl(compRef, container)
	if err != nil {
		return err
	}

	componentSpec, err := c.useComponentSpec(compRef)
	if err != nil {
		return fmt.Errorf("component spec for %q not found", compRef)
	}
	taskSpecJson, err := stablyMarshalJSON(task)
	if err != nil {
		return err
	}
	containerImpl, err := c.useComponentImpl(compRef)
	if err != nil {
		return err
	}

	exitHandler := false
	if task.GetTriggerPolicy().GetStrategy().String() == "ALL_UPSTREAM_TASKS_COMPLETED" {
		exitHandler = true
	}
	kubernetesConfigPlaceholder, _ := c.useKubernetesImpl(taskName)
	return c.containerDriverTask(taskName, &containerDriverInputs{
		component:        componentSpec,
		task:             taskSpecJson,
		container:        containerImpl,
		parentDag:        c.CurrentDag(),
		taskDef:          task,
		containerDef:     container,
		exitHandler:      exitHandler,
		kubernetesConfig: kubernetesConfigPlaceholder,
	})
}

type containerDriverOutputs struct {
	// break down podSpecPath to the following
	executionId    string
	executiorInput string
	cached         string
	condition      string
	podSpecPatch   string
}

type containerDriverInputs struct {
	component        string
	task             string
	taskDef          *pipelinespec.PipelineTaskSpec
	container        string
	containerDef     *pipelinespec.PipelineDeploymentConfig_PipelineContainerSpec
	parentDag        string
	iterationIndex   string // optional, when this is an iteration task
	exitHandler      bool
	kubernetesConfig string
}

func (i *containerDriverInputs) getParentDagID(isExitHandler bool) string {
	if i.parentDag == "" {
		return "0"
	}
	if isExitHandler && i.parentDag == compiler.RootComponentName {
		return fmt.Sprintf("$(params.%s)", paramParentDagID)
	} else {
		return taskOutputParameter(getDAGDriverTaskName(i.parentDag), paramExecutionID)
	}
}

func (i *containerDriverInputs) getParentDagCondition(isExitHandler bool) string {
	if i.parentDag == "" {
		return "0"
	}
	if isExitHandler && i.parentDag == compiler.RootComponentName {
		return fmt.Sprintf("$(params.%s)", paramCondition)
	} else {
		return taskOutputParameter(getDAGDriverTaskName(i.parentDag), paramCondition)
	}
}

func (c *pipelinerunCompiler) containerDriverTask(name string, inputs *containerDriverInputs) error {

	containerDriverName := getContainerDriverTaskName(name)
	driverTask := &pipelineapi.PipelineTask{
		Name: containerDriverName,
		TaskRef: &pipelineapi.TaskRef{
			APIVersion: "kfp-driver.tekton.dev/v1alpha1",
			Kind:       "KFPDriver",
		},
		Params: []pipelineapi.Param{
			// "--type", "CONTAINER",
			{
				Name:  paramNameType,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: "CONTAINER"},
			},
			// "--pipeline_name", c.spec.GetPipelineInfo().GetName(),
			{
				Name:  paramNamePipelineName,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: c.spec.GetPipelineInfo().GetName()},
			},
			// "--run_id", runID(),
			{
				Name:  paramNameRunId,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: runID()},
			},
			// "--dag_execution_id"
			{
				Name:  paramNameDagExecutionId,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: inputs.getParentDagID(c.ExitHandlerScope())},
			},
			// "--component"
			{
				Name:  paramComponent,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: inputs.component},
			},
			// "--task"
			{
				Name:  paramTask,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: inputs.task},
			},
			// "--container"
			{
				Name:  paramContainer,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: inputs.container},
			},
			// "--iteration_index", inputValue(paramIterationIndex),
			{
				Name:  paramNameIterationIndex,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: inputs.iterationIndex},
			},
			// "--kubernetes_config"
			{
				Name:  paramKubernetesConfig,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: inputs.kubernetesConfig},
			},
			// produce the following outputs:
			// - execution-id
			// - executor-input
			// - cached-decision
			// - condition
		},
	}

	if len(inputs.taskDef.GetDependentTasks()) > 0 {
		driverTask.RunAfter = inputs.taskDef.GetDependentTasks()
	}

	// adding WhenExpress for condition only if the task belongs to a DAG had a condition TriggerPolicy
	if c.ConditionScope() {
		driverTask.WhenExpressions = pipelineapi.WhenExpressions{
			pipelineapi.WhenExpression{
				Input:    inputs.getParentDagCondition(c.ExitHandlerScope()),
				Operator: selection.NotIn,
				Values:   []string{"false"},
			},
		}
	}

	c.addPipelineTask(driverTask)

	// need container driver's output for executor
	containerDriverOutputs := containerDriverOutputs{
		executionId:    taskOutputParameter(containerDriverName, paramExecutionID),
		condition:      taskOutputParameter(containerDriverName, paramCondition),
		executiorInput: taskOutputParameter(containerDriverName, paramExecutorInput),
		cached:         taskOutputParameter(containerDriverName, paramCachedDecision),
		podSpecPatch:   taskOutputParameter(containerDriverName, paramPodSpecPatch),
	}

	t, err := c.containerExecutorTemplate(name, inputs.containerDef, c.spec.PipelineInfo.GetName())

	if err != nil {
		return err
	}

	executorTask := &pipelineapi.PipelineTask{
		Name:     name,
		TaskSpec: t,
		WhenExpressions: pipelineapi.WhenExpressions{
			{
				Input:    containerDriverOutputs.cached,
				Operator: "in",
				Values:   []string{"false"},
			},
		},
		Params: []pipelineapi.Param{
			{
				Name:  paramExecutorInput,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: containerDriverOutputs.executiorInput},
			},
			{
				Name:  paramExecutionID,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: containerDriverOutputs.executionId},
			},
			{
				Name:  paramRunId,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: runID()},
			},
			{
				Name:  paramComponentSpec,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: inputs.component},
			},
			{
				Name:  paramPodSpecPatch,
				Value: pipelineapi.ParamValue{Type: "string", StringVal: containerDriverOutputs.podSpecPatch},
			},
		},
	}

	c.addPipelineTask(executorTask)

	return nil
}

func (c *pipelinerunCompiler) containerExecutorTemplate(
	name string, container *pipelinespec.PipelineDeploymentConfig_PipelineContainerSpec,
	pipelineName string,
) (*pipelineapi.EmbeddedTask, error) {
	userCmdArgs := make([]string, 0, len(container.Command)+len(container.Args))
	userCmdArgs = append(userCmdArgs, container.Command...)
	userCmdArgs = append(userCmdArgs, container.Args...)
	// userCmdArgs = append(userCmdArgs, "--executor_input", "{{$}}", "--function_to_execute", inputValue(paramFunctionToExecute))
	launcherCmd := []string{
		kfpLauncherPath,
		"--pipeline_name", pipelineName,
		"--run_id", inputValue(paramRunId),
		"--execution_id", inputValue(paramExecutionID),
		"--executor_input", inputValue(paramExecutorInput),
		"--component_spec", inputValue(paramComponentSpec),
		"--pod_name",
		"$(KFP_POD_NAME)",
		"--pod_uid",
		"$(KFP_POD_UID)",
		"--mlmd_server_address", // METADATA_GRPC_SERVICE_* come from metadata-grpc-configmap
		"$(METADATA_GRPC_SERVICE_HOST)",
		"--mlmd_server_port",
		"$(METADATA_GRPC_SERVICE_PORT)",
		"--", // separater before user command and args
	}
	mlmdConfigOptional := true
	kfpTaskSpec := ktv1alpha1.KfpTaskSpec{
		TaskSpec: &pipelineapi.TaskSpec{
			Params: []pipelineapi.ParamSpec{
				{Name: paramExecutorInput, Type: "string"}, // --executor_input
				{Name: paramExecutionID, Type: "string"},   // --execution_id
				{Name: paramRunId, Type: "string"},         // --run_id
				{Name: paramComponentSpec, Type: "string"}, // --component_spec
			},
			Steps: []pipelineapi.Step{
				// step 1: copy launcher
				{
					Name:            "kfp-launcher",
					Image:           c.launcherImage,
					Command:         []string{"launcher-v2", "--copy", kfpLauncherPath},
					ImagePullPolicy: "Always",
				},
				// wrap user program with executor
				{
					Name:    "user-main",
					Image:   container.Image,
					Command: launcherCmd,
					Args:    userCmdArgs,
					EnvFrom: []k8score.EnvFromSource{{
						ConfigMapRef: &k8score.ConfigMapEnvSource{
							LocalObjectReference: k8score.LocalObjectReference{
								Name: "metadata-grpc-configmap",
							},
							Optional: &mlmdConfigOptional,
						},
					}},
					Env: []k8score.EnvVar{{
						Name: "KFP_POD_NAME",
						ValueFrom: &k8score.EnvVarSource{
							FieldRef: &k8score.ObjectFieldSelector{
								FieldPath: "metadata.name",
							},
						},
					}, {
						Name: "KFP_POD_UID",
						ValueFrom: &k8score.EnvVarSource{
							FieldRef: &k8score.ObjectFieldSelector{
								FieldPath: "metadata.uid",
							},
						},
					}, {
						Name:  "METADATA_GRPC_SERVICE_HOST",
						Value: "metadata-grpc-service.kubeflow.svc.cluster.local",
					}, {
						Name:  "METADATA_GRPC_SERVICE_PORT",
						Value: "8080",
					}},
				},
			},
		},
	}

	raw, err := json.Marshal(kfpTaskSpec)
	if err != nil {
		return nil, fmt.Errorf("unable to Marshal KfpTaskSpec:%v", err)
	}

	return &pipelineapi.EmbeddedTask{
		Metadata: pipelineapi.PipelineTaskMetadata{
			Annotations: map[string]string{
				"pipelines.kubeflow.org/v2_pipeline": "true",
			},
			Labels: map[string]string{
				"pipelines.kubeflow.org/v2_component": "true",
			},
		},
		TypeMeta: runtime.TypeMeta{
			Kind:       kfptask.Kind,
			APIVersion: ktv1alpha1.SchemeGroupVersion.String(),
		},
		Spec: runtime.RawExtension{
			Raw: raw,
		},
	}, nil
}