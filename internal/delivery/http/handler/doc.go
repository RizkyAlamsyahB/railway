// Package handler contains HTTP request handlers for the API.
//
// Each handler struct receives its dependencies (usecases) via constructor
// injection and exposes methods that serve as Gin handler functions.
// Handlers are responsible for parsing request parameters, calling the
// appropriate usecase, and returning a JSON envelope response.
package handler
