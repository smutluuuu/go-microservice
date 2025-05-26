package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// readJSON: Reads JSON data from the HTTP request body and stores it in the provided 'data' variable.
// This function:
// - Limits the body size to 1MB (for security)
// - Decodes the JSON data
// - Checks that the body contains only a single JSON object
// - Returns appropriate error messages if anything goes wrong
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1 MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}
	return nil
}

// writeJSON: Converts the provided 'data' to JSON format and sends it as an HTTP response.
// This function:
// - Converts the data to JSON format
// - Adds optional HTTP headers if provided
// - Sets the Content-Type header to "application/json"
// - Sets the HTTP status code
// - Writes the JSON data to the response body
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {

	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

// errorJSON: Creates and sends a standardized JSON error response for error situations.
// This function:
// - Uses HTTP 400 Bad Request status code by default
// - Allows specifying a different status code if needed
// - Converts the error message to a JSON structure
// - Uses the writeJSON function to send the error response
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {

	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse

	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
