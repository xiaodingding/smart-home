// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// ResponseType Тип ответов
// swagger:model ResponseType
type ResponseType string

const (
	// ResponseTypeSuccess captures enum value "success"
	ResponseTypeSuccess ResponseType = "success"
	// ResponseTypeBusinessConflict captures enum value "business_conflict"
	ResponseTypeBusinessConflict ResponseType = "business_conflict"
	// ResponseTypeUnprocessableEntity captures enum value "unprocessable_entity"
	ResponseTypeUnprocessableEntity ResponseType = "unprocessable_entity"
	// ResponseTypeBadParameters captures enum value "bad_parameters"
	ResponseTypeBadParameters ResponseType = "bad_parameters"
	// ResponseTypeInternalError captures enum value "internal_error"
	ResponseTypeInternalError ResponseType = "internal_error"
	// ResponseTypeNotFound captures enum value "not_found"
	ResponseTypeNotFound ResponseType = "not_found"
	// ResponseTypeSecurityError captures enum value "security_error"
	ResponseTypeSecurityError ResponseType = "security_error"
	// ResponseTypePermissionError captures enum value "permission_error"
	ResponseTypePermissionError ResponseType = "permission_error"
)

// for schema
var responseTypeEnum []interface{}

func init() {
	var res []ResponseType
	if err := json.Unmarshal([]byte(`["success","business_conflict","unprocessable_entity","bad_parameters","internal_error","not_found","security_error","permission_error"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		responseTypeEnum = append(responseTypeEnum, v)
	}
}

func (m ResponseType) validateResponseTypeEnum(path, location string, value ResponseType) error {
	if err := validate.Enum(path, location, value, responseTypeEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this response type
func (m ResponseType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateResponseTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}