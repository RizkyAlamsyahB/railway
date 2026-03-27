package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/media-inovasi-strategis/haji-umroh-store-be/pkg/response"
)

type requestErrors struct {
	Fields  map[string][]string `json:"fields"`
	Request []string            `json:"request"`
}

func newRequestErrors() requestErrors {
	return requestErrors{
		Fields:  map[string][]string{},
		Request: []string{},
	}
}

func respondValidationError(c *gin.Context, err error, payload interface{}) {
	details := buildRequestErrors(err, payload)
	response.BadRequest(c, "validation failed", details)
}

func respondErrorWithRequestDetails(c *gin.Context, status int, message string, details requestErrors) {
	response.Error(c, status, message, details)
}

func buildRequestErrors(err error, payload interface{}) requestErrors {
	details := newRequestErrors()

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		fieldMap := jsonFieldMap(payload)
		for _, validationErr := range validationErrs {
			fieldName := fieldMap[validationErr.StructField()]
			if fieldName == "" {
				fieldName = toSnakeCase(validationErr.Field())
			}
			details.Fields[fieldName] = append(details.Fields[fieldName], validationMessage(validationErr))
		}
		return details
	}

	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeErr) {
		fieldName := unmarshalTypeErr.Field
		if fieldName == "" {
			details.Request = append(details.Request, "request body contains an invalid value")
			return details
		}

		jsonName := jsonFieldMap(payload)[fieldName]
		if jsonName == "" {
			jsonName = toSnakeCase(fieldName)
		}
		details.Fields[jsonName] = append(details.Fields[jsonName], fmt.Sprintf("must be a valid %s", unmarshalTypeErr.Type.String()))
		return details
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		details.Request = append(details.Request, "request body must be valid JSON")
		return details
	}

	if errors.Is(err, io.ErrUnexpectedEOF) {
		details.Request = append(details.Request, "request body must be valid JSON")
		return details
	}

	details.Request = append(details.Request, "request body is invalid")
	return details
}

func validationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", err.Param())
	case "max":
		return fmt.Sprintf("must not exceed %s characters", err.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", strings.ReplaceAll(err.Param(), " ", ", "))
	default:
		return "is invalid"
	}
}

func jsonFieldMap(payload interface{}) map[string]string {
	fieldMap := map[string]string{}
	t := reflect.TypeOf(payload)
	if t == nil {
		return fieldMap
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fieldMap
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		jsonName := strings.Split(jsonTag, ",")[0]
		if jsonName == "" {
			jsonName = toSnakeCase(field.Name)
		}
		fieldMap[field.Name] = jsonName
	}

	return fieldMap
}

func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}
