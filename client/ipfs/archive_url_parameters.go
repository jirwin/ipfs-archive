// Code generated by go-swagger; DO NOT EDIT.

package ipfs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"
	"time"

	"golang.org/x/net/context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/jirwin/ipfs-archive/models"
)

// NewArchiveURLParams creates a new ArchiveURLParams object
// with the default values initialized.
func NewArchiveURLParams() *ArchiveURLParams {
	var ()
	return &ArchiveURLParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewArchiveURLParamsWithTimeout creates a new ArchiveURLParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewArchiveURLParamsWithTimeout(timeout time.Duration) *ArchiveURLParams {
	var ()
	return &ArchiveURLParams{

		timeout: timeout,
	}
}

// NewArchiveURLParamsWithContext creates a new ArchiveURLParams object
// with the default values initialized, and the ability to set a context for a request
func NewArchiveURLParamsWithContext(ctx context.Context) *ArchiveURLParams {
	var ()
	return &ArchiveURLParams{

		Context: ctx,
	}
}

// NewArchiveURLParamsWithHTTPClient creates a new ArchiveURLParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewArchiveURLParamsWithHTTPClient(client *http.Client) *ArchiveURLParams {
	var ()
	return &ArchiveURLParams{
		HTTPClient: client,
	}
}

/*ArchiveURLParams contains all the parameters to send to the API endpoint
for the archive Url operation typically these are written to a http.Request
*/
type ArchiveURLParams struct {

	/*Body
	  The URL to archive

	*/
	Body *models.ArchiveRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the archive Url params
func (o *ArchiveURLParams) WithTimeout(timeout time.Duration) *ArchiveURLParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the archive Url params
func (o *ArchiveURLParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the archive Url params
func (o *ArchiveURLParams) WithContext(ctx context.Context) *ArchiveURLParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the archive Url params
func (o *ArchiveURLParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the archive Url params
func (o *ArchiveURLParams) WithHTTPClient(client *http.Client) *ArchiveURLParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the archive Url params
func (o *ArchiveURLParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the archive Url params
func (o *ArchiveURLParams) WithBody(body *models.ArchiveRequest) *ArchiveURLParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the archive Url params
func (o *ArchiveURLParams) SetBody(body *models.ArchiveRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *ArchiveURLParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
