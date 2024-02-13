// Code generated by go-swagger; DO NOT EDIT.

package pipeline_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewPipelineServiceListPipelineVersionsParams creates a new PipelineServiceListPipelineVersionsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewPipelineServiceListPipelineVersionsParams() *PipelineServiceListPipelineVersionsParams {
	return &PipelineServiceListPipelineVersionsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewPipelineServiceListPipelineVersionsParamsWithTimeout creates a new PipelineServiceListPipelineVersionsParams object
// with the ability to set a timeout on a request.
func NewPipelineServiceListPipelineVersionsParamsWithTimeout(timeout time.Duration) *PipelineServiceListPipelineVersionsParams {
	return &PipelineServiceListPipelineVersionsParams{
		timeout: timeout,
	}
}

// NewPipelineServiceListPipelineVersionsParamsWithContext creates a new PipelineServiceListPipelineVersionsParams object
// with the ability to set a context for a request.
func NewPipelineServiceListPipelineVersionsParamsWithContext(ctx context.Context) *PipelineServiceListPipelineVersionsParams {
	return &PipelineServiceListPipelineVersionsParams{
		Context: ctx,
	}
}

// NewPipelineServiceListPipelineVersionsParamsWithHTTPClient creates a new PipelineServiceListPipelineVersionsParams object
// with the ability to set a custom HTTPClient for a request.
func NewPipelineServiceListPipelineVersionsParamsWithHTTPClient(client *http.Client) *PipelineServiceListPipelineVersionsParams {
	return &PipelineServiceListPipelineVersionsParams{
		HTTPClient: client,
	}
}

/*
PipelineServiceListPipelineVersionsParams contains all the parameters to send to the API endpoint

	for the pipeline service list pipeline versions operation.

	Typically these are written to a http.Request.
*/
type PipelineServiceListPipelineVersionsParams struct {

	/* Filter.

	     A base-64 encoded, JSON-serialized Filter protocol buffer (see
	filter.proto).
	*/
	Filter *string

	/* PageSize.

	     The number of pipeline versions to be listed per page. If there are more
	pipeline versions than this number, the response message will contain a
	nextPageToken field you can use to fetch the next page.

	     Format: int32
	*/
	PageSize *int32

	/* PageToken.

	     A page token to request the next page of results. The token is acquried
	from the nextPageToken field of the response from the previous
	ListPipelineVersions call or can be omitted when fetching the first page.
	*/
	PageToken *string

	/* ResourceKeyID.

	   The ID of the resource that referred to.
	*/
	ResourceKeyID *string

	/* ResourceKeyType.

	   The type of the resource that referred to.

	   Default: "UNKNOWN_RESOURCE_TYPE"
	*/
	ResourceKeyType *string

	/* SortBy.

	     Can be format of "field_name", "field_name asc" or "field_name desc"
	Ascending by default.
	*/
	SortBy *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the pipeline service list pipeline versions params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PipelineServiceListPipelineVersionsParams) WithDefaults() *PipelineServiceListPipelineVersionsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the pipeline service list pipeline versions params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *PipelineServiceListPipelineVersionsParams) SetDefaults() {
	var (
		resourceKeyTypeDefault = string("UNKNOWN_RESOURCE_TYPE")
	)

	val := PipelineServiceListPipelineVersionsParams{
		ResourceKeyType: &resourceKeyTypeDefault,
	}

	val.timeout = o.timeout
	val.Context = o.Context
	val.HTTPClient = o.HTTPClient
	*o = val
}

// WithTimeout adds the timeout to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) WithTimeout(timeout time.Duration) *PipelineServiceListPipelineVersionsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) WithContext(ctx context.Context) *PipelineServiceListPipelineVersionsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) WithHTTPClient(client *http.Client) *PipelineServiceListPipelineVersionsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFilter adds the filter to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) WithFilter(filter *string) *PipelineServiceListPipelineVersionsParams {
	o.SetFilter(filter)
	return o
}

// SetFilter adds the filter to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) SetFilter(filter *string) {
	o.Filter = filter
}

// WithPageSize adds the pageSize to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) WithPageSize(pageSize *int32) *PipelineServiceListPipelineVersionsParams {
	o.SetPageSize(pageSize)
	return o
}

// SetPageSize adds the pageSize to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) SetPageSize(pageSize *int32) {
	o.PageSize = pageSize
}

// WithPageToken adds the pageToken to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) WithPageToken(pageToken *string) *PipelineServiceListPipelineVersionsParams {
	o.SetPageToken(pageToken)
	return o
}

// SetPageToken adds the pageToken to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) SetPageToken(pageToken *string) {
	o.PageToken = pageToken
}

// WithResourceKeyID adds the resourceKeyID to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) WithResourceKeyID(resourceKeyID *string) *PipelineServiceListPipelineVersionsParams {
	o.SetResourceKeyID(resourceKeyID)
	return o
}

// SetResourceKeyID adds the resourceKeyId to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) SetResourceKeyID(resourceKeyID *string) {
	o.ResourceKeyID = resourceKeyID
}

// WithResourceKeyType adds the resourceKeyType to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) WithResourceKeyType(resourceKeyType *string) *PipelineServiceListPipelineVersionsParams {
	o.SetResourceKeyType(resourceKeyType)
	return o
}

// SetResourceKeyType adds the resourceKeyType to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) SetResourceKeyType(resourceKeyType *string) {
	o.ResourceKeyType = resourceKeyType
}

// WithSortBy adds the sortBy to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) WithSortBy(sortBy *string) *PipelineServiceListPipelineVersionsParams {
	o.SetSortBy(sortBy)
	return o
}

// SetSortBy adds the sortBy to the pipeline service list pipeline versions params
func (o *PipelineServiceListPipelineVersionsParams) SetSortBy(sortBy *string) {
	o.SortBy = sortBy
}

// WriteToRequest writes these params to a swagger request
func (o *PipelineServiceListPipelineVersionsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Filter != nil {

		// query param filter
		var qrFilter string

		if o.Filter != nil {
			qrFilter = *o.Filter
		}
		qFilter := qrFilter
		if qFilter != "" {

			if err := r.SetQueryParam("filter", qFilter); err != nil {
				return err
			}
		}
	}

	if o.PageSize != nil {

		// query param page_size
		var qrPageSize int32

		if o.PageSize != nil {
			qrPageSize = *o.PageSize
		}
		qPageSize := swag.FormatInt32(qrPageSize)
		if qPageSize != "" {

			if err := r.SetQueryParam("page_size", qPageSize); err != nil {
				return err
			}
		}
	}

	if o.PageToken != nil {

		// query param page_token
		var qrPageToken string

		if o.PageToken != nil {
			qrPageToken = *o.PageToken
		}
		qPageToken := qrPageToken
		if qPageToken != "" {

			if err := r.SetQueryParam("page_token", qPageToken); err != nil {
				return err
			}
		}
	}

	if o.ResourceKeyID != nil {

		// query param resource_key.id
		var qrResourceKeyID string

		if o.ResourceKeyID != nil {
			qrResourceKeyID = *o.ResourceKeyID
		}
		qResourceKeyID := qrResourceKeyID
		if qResourceKeyID != "" {

			if err := r.SetQueryParam("resource_key.id", qResourceKeyID); err != nil {
				return err
			}
		}
	}

	if o.ResourceKeyType != nil {

		// query param resource_key.type
		var qrResourceKeyType string

		if o.ResourceKeyType != nil {
			qrResourceKeyType = *o.ResourceKeyType
		}
		qResourceKeyType := qrResourceKeyType
		if qResourceKeyType != "" {

			if err := r.SetQueryParam("resource_key.type", qResourceKeyType); err != nil {
				return err
			}
		}
	}

	if o.SortBy != nil {

		// query param sort_by
		var qrSortBy string

		if o.SortBy != nil {
			qrSortBy = *o.SortBy
		}
		qSortBy := qrSortBy
		if qSortBy != "" {

			if err := r.SetQueryParam("sort_by", qSortBy); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}