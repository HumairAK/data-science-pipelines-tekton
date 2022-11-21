// Code generated by go-swagger; DO NOT EDIT.

package run_model

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// APIPipelineRuntime api pipeline runtime
// swagger:model apiPipelineRuntime
type APIPipelineRuntime struct {

	// Output. The runtime JSON manifest of the pipeline, including the status
	// of pipeline steps and fields need for UI visualization etc.
	PipelineManifest string `json:"pipeline_manifest,omitempty"`

	// Output. The runtime JSON manifest of the argo workflow.
	// This is deprecated after pipeline_runtime_manifest is in use.
	WorkflowManifest string `json:"workflow_manifest,omitempty"`
}

// Validate validates this api pipeline runtime
func (m *APIPipelineRuntime) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *APIPipelineRuntime) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *APIPipelineRuntime) UnmarshalBinary(b []byte) error {
	var res APIPipelineRuntime
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}