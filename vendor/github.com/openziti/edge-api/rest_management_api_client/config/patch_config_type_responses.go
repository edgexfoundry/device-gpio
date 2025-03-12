// Code generated by go-swagger; DO NOT EDIT.

//
// Copyright NetFoundry Inc.
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
//
// __          __              _
// \ \        / /             (_)
//  \ \  /\  / /_ _ _ __ _ __  _ _ __   __ _
//   \ \/  \/ / _` | '__| '_ \| | '_ \ / _` |
//    \  /\  / (_| | |  | | | | | | | | (_| | : This file is generated, do not edit it.
//     \/  \/ \__,_|_|  |_| |_|_|_| |_|\__, |
//                                      __/ |
//                                     |___/

package config

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/openziti/edge-api/rest_model"
)

// PatchConfigTypeReader is a Reader for the PatchConfigType structure.
type PatchConfigTypeReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PatchConfigTypeReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPatchConfigTypeOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewPatchConfigTypeBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 401:
		result := NewPatchConfigTypeUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewPatchConfigTypeNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 429:
		result := NewPatchConfigTypeTooManyRequests()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 503:
		result := NewPatchConfigTypeServiceUnavailable()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewPatchConfigTypeOK creates a PatchConfigTypeOK with default headers values
func NewPatchConfigTypeOK() *PatchConfigTypeOK {
	return &PatchConfigTypeOK{}
}

/* PatchConfigTypeOK describes a response with status code 200, with default header values.

The patch request was successful and the resource has been altered
*/
type PatchConfigTypeOK struct {
	Payload *rest_model.Empty
}

func (o *PatchConfigTypeOK) Error() string {
	return fmt.Sprintf("[PATCH /config-types/{id}][%d] patchConfigTypeOK  %+v", 200, o.Payload)
}
func (o *PatchConfigTypeOK) GetPayload() *rest_model.Empty {
	return o.Payload
}

func (o *PatchConfigTypeOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.Empty)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchConfigTypeBadRequest creates a PatchConfigTypeBadRequest with default headers values
func NewPatchConfigTypeBadRequest() *PatchConfigTypeBadRequest {
	return &PatchConfigTypeBadRequest{}
}

/* PatchConfigTypeBadRequest describes a response with status code 400, with default header values.

The supplied request contains invalid fields or could not be parsed (json and non-json bodies). The error's code, message, and cause fields can be inspected for further information
*/
type PatchConfigTypeBadRequest struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *PatchConfigTypeBadRequest) Error() string {
	return fmt.Sprintf("[PATCH /config-types/{id}][%d] patchConfigTypeBadRequest  %+v", 400, o.Payload)
}
func (o *PatchConfigTypeBadRequest) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *PatchConfigTypeBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchConfigTypeUnauthorized creates a PatchConfigTypeUnauthorized with default headers values
func NewPatchConfigTypeUnauthorized() *PatchConfigTypeUnauthorized {
	return &PatchConfigTypeUnauthorized{}
}

/* PatchConfigTypeUnauthorized describes a response with status code 401, with default header values.

The supplied session does not have the correct access rights to request this resource
*/
type PatchConfigTypeUnauthorized struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *PatchConfigTypeUnauthorized) Error() string {
	return fmt.Sprintf("[PATCH /config-types/{id}][%d] patchConfigTypeUnauthorized  %+v", 401, o.Payload)
}
func (o *PatchConfigTypeUnauthorized) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *PatchConfigTypeUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchConfigTypeNotFound creates a PatchConfigTypeNotFound with default headers values
func NewPatchConfigTypeNotFound() *PatchConfigTypeNotFound {
	return &PatchConfigTypeNotFound{}
}

/* PatchConfigTypeNotFound describes a response with status code 404, with default header values.

The requested resource does not exist
*/
type PatchConfigTypeNotFound struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *PatchConfigTypeNotFound) Error() string {
	return fmt.Sprintf("[PATCH /config-types/{id}][%d] patchConfigTypeNotFound  %+v", 404, o.Payload)
}
func (o *PatchConfigTypeNotFound) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *PatchConfigTypeNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchConfigTypeTooManyRequests creates a PatchConfigTypeTooManyRequests with default headers values
func NewPatchConfigTypeTooManyRequests() *PatchConfigTypeTooManyRequests {
	return &PatchConfigTypeTooManyRequests{}
}

/* PatchConfigTypeTooManyRequests describes a response with status code 429, with default header values.

The resource requested is rate limited and the rate limit has been exceeded
*/
type PatchConfigTypeTooManyRequests struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *PatchConfigTypeTooManyRequests) Error() string {
	return fmt.Sprintf("[PATCH /config-types/{id}][%d] patchConfigTypeTooManyRequests  %+v", 429, o.Payload)
}
func (o *PatchConfigTypeTooManyRequests) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *PatchConfigTypeTooManyRequests) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPatchConfigTypeServiceUnavailable creates a PatchConfigTypeServiceUnavailable with default headers values
func NewPatchConfigTypeServiceUnavailable() *PatchConfigTypeServiceUnavailable {
	return &PatchConfigTypeServiceUnavailable{}
}

/* PatchConfigTypeServiceUnavailable describes a response with status code 503, with default header values.

The request could not be completed due to the server being busy or in a temporarily bad state
*/
type PatchConfigTypeServiceUnavailable struct {
	Payload *rest_model.APIErrorEnvelope
}

func (o *PatchConfigTypeServiceUnavailable) Error() string {
	return fmt.Sprintf("[PATCH /config-types/{id}][%d] patchConfigTypeServiceUnavailable  %+v", 503, o.Payload)
}
func (o *PatchConfigTypeServiceUnavailable) GetPayload() *rest_model.APIErrorEnvelope {
	return o.Payload
}

func (o *PatchConfigTypeServiceUnavailable) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(rest_model.APIErrorEnvelope)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
