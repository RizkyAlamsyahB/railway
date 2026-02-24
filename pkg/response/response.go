package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Envelope is the standard JSON response wrapper used by all API endpoints.
type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  interface{} `json:"errors"`
	Meta    interface{} `json:"meta"`
}

// Success sends a successful JSON envelope response.
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Envelope{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
		Meta:    nil,
	})
}

// SuccessWithMeta sends a successful JSON envelope response with metadata.
func SuccessWithMeta(c *gin.Context, statusCode int, message string, data interface{}, meta interface{}) {
	c.JSON(statusCode, Envelope{
		Success: true,
		Message: message,
		Data:    data,
		Errors:  nil,
		Meta:    meta,
	})
}

// Error sends an error JSON envelope response.
func Error(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.JSON(statusCode, Envelope{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  errors,
		Meta:    nil,
	})
}

// Abort sends an error JSON envelope response and aborts the request chain.
// This is intended for use in middleware where further handler execution
// should be prevented.
func Abort(c *gin.Context, statusCode int, message string, errors interface{}) {
	c.AbortWithStatusJSON(statusCode, Envelope{
		Success: false,
		Message: message,
		Data:    nil,
		Errors:  errors,
		Meta:    nil,
	})
}

// OK is a convenience wrapper for 200 OK responses.
func OK(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusOK, message, data)
}

// Created is a convenience wrapper for 201 Created responses.
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, http.StatusCreated, message, data)
}

// BadRequest is a convenience wrapper for 400 Bad Request responses.
func BadRequest(c *gin.Context, message string, errors interface{}) {
	Error(c, http.StatusBadRequest, message, errors)
}

// Unauthorized is a convenience wrapper for 401 Unauthorized responses.
func Unauthorized(c *gin.Context, message string, errors interface{}) {
	Error(c, http.StatusUnauthorized, message, errors)
}

// Forbidden is a convenience wrapper for 403 Forbidden responses.
func Forbidden(c *gin.Context, message string, errors interface{}) {
	Error(c, http.StatusForbidden, message, errors)
}

// NotFound is a convenience wrapper for 404 Not Found responses.
func NotFound(c *gin.Context, message string, errors interface{}) {
	Error(c, http.StatusNotFound, message, errors)
}

// InternalServerError is a convenience wrapper for 500 Internal Server Error responses.
func InternalServerError(c *gin.Context, message string, errors interface{}) {
	Error(c, http.StatusInternalServerError, message, errors)
}
