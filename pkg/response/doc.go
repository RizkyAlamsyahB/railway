// Package response provides a standardized JSON envelope for all API responses.
//
// Every API response is wrapped in a consistent structure containing:
//   - success: boolean indicating if the request was successful
//   - message: human-readable status message
//   - data: the response payload (null on error)
//   - errors: error details (null on success)
//   - meta: optional metadata such as pagination info
//
// This ensures a uniform response format across all endpoints.
package response
