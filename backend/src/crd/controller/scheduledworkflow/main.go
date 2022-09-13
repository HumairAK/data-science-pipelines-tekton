// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"strings"
	"time"

	commonutil "github.com/kubeflow/pipelines/backend/src/common/util"
	"github.com/kubeflow/pipelines/backend/src/crd/controller/scheduledworkflow/util"
	swfclientset "github.com/kubeflow/pipelines/backend/src/crd/pkg/client/clientset/versioned"
	swfinformers "github.com/kubeflow/pipelines/backend/src/crd/pkg/client/informers/externalversions"
	"github.com/kubeflow/pipelines/backend/src/crd/pkg/signals"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	masterURL     string
	kubeconfig    string
	namespace     string
	location      *time.Location
	clientQPS     float64
	clientBurst   int
	executionType string
)

func main() {
	flag.Parse()

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	// Use the commonutil to store the ExecutionType
	commonutil.SetExecutionType(commonutil.ExecutionType(executionType))

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	cfg.QPS = float32(clientQPS)
	cfg.Burst = clientBurst

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	scheduleClient, err := swfclientset.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("Error building schedule clientset: %s", err.Error())
	}

	clientParam := commonutil.ClientParameters{QPS: float64(cfg.QPS), Burst: cfg.Burst}
	execClient := commonutil.NewExecutionClientOrFatal(commonutil.CurrentExecutionType(), time.Second*30, clientParam)

	var scheduleInformerFactory swfinformers.SharedInformerFactory
	execInformer := commonutil.NewExecutionInformerOrFatal(commonutil.CurrentExecutionType(), namespace, time.Second*30, clientParam)
	if namespace == "" {
		scheduleInformerFactory = swfinformers.NewSharedInformerFactory(scheduleClient, time.Second*30)
	} else {
		scheduleInformerFactory = swfinformers.NewFilteredSharedInformerFactory(scheduleClient, time.Second*30, namespace, nil)
	}

	controller := NewController(
		kubeClient,
		scheduleClient,
		execClient,
		scheduleInformerFactory,
		execInformer,
		commonutil.NewRealTime(),
		location)

	go scheduleInformerFactory.Start(stopCh)
	go execInformer.InformerFactoryStart(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		log.Fatalf("Error running controller: %s", err.Error())
	}
}

func initEnv() {
	// Import environment variable, support nested vars e.g. OBJECTSTORECONFIG_ACCESSKEY
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)
}

func init() {
	initEnv()

	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&namespace, "namespace", "", "The namespace name used for Kubernetes informers to obtain the listers.")
	// Use default value of client QPS (5) & burst (10) defined in
	// k8s.io/client-go/rest/config.go#RESTClientFor
	flag.Float64Var(&clientQPS, "clientQPS", 5, "The maximum QPS to the master from this client.")
	flag.IntVar(&clientBurst, "clientBurst", 10, "Maximum burst for throttle from this client.")
	flag.StringVar(&executionType, "executionType", "Workflow", "Custom Resource's name of the backend Orchestration Engine")
	var err error
	location, err = util.GetLocation()
	if err != nil {
		log.Fatalf("Error running controller: %s", err.Error())
	}
	log.Infof("Location: %s", location.String())
}
