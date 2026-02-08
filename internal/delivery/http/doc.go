// Package http contains the HTTP delivery layer of the application.
//
// It is organized into three sub-packages:
//   - handler: HTTP request handlers that invoke usecases and return responses.
//   - middleware: HTTP middleware for cross-cutting concerns (CORS, recovery, request ID).
//   - router: Route registration and Gin engine setup.
package http
