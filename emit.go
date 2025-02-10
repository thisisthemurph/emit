package emit

import (
	"encoding/json"
	"net/http"
)

// A ResponseBuilder is used to chain functions for crafting HTTP responses.
type ResponseBuilder struct {
	w          http.ResponseWriter
	statusCode int
	headers    map[string]string
}

// New creates an initial ResponseBuilder.
func New(w http.ResponseWriter) *ResponseBuilder {
	return &ResponseBuilder{
		w:          w,
		statusCode: http.StatusOK,
		headers:    make(map[string]string),
	}
}

// Status sets the HTTP status code of the response.
func (rb *ResponseBuilder) Status(code int) *ResponseBuilder {
	rb.statusCode = code
	return rb
}

// Header sets a response header.
func (rb *ResponseBuilder) Header(key, value string) *ResponseBuilder {
	rb.headers[key] = value
	return rb
}

// Cookie sets a cookie on the response.
//
// If a nil cookie is passed in, it is ignored and no error is surfaced.
func (rb *ResponseBuilder) Cookie(cookie *http.Cookie) *ResponseBuilder {
	if cookie != nil {
		http.SetCookie(rb.w, cookie)
	}
	return rb
}

// Text sends the given string as a text/plain response.
//
// Uses a default status code of 200 if one was not set using the Status method.
func (rb *ResponseBuilder) Text(text string) {
	rb.headers["Content-Type"] = "text/plain"
	rb.applyHeaders()
	rb.w.WriteHeader(rb.statusCode)
	_, _ = rb.w.Write([]byte(text))
}

// JSON sends a JSON response.
//
// Uses a default status code of 200 if one was not set using the Status method.
//
// If nil data is passed, Content-Type will be set to application/json, but no JSON body will be set.
func (rb *ResponseBuilder) JSON(data interface{}) {
	rb.headers["Content-Type"] = "application/json"
	rb.applyHeaders()
	rb.w.WriteHeader(rb.statusCode)

	if data == nil {
		return
	}

	if err := json.NewEncoder(rb.w).Encode(data); err != nil {
		http.Error(rb.w, `{"error": "failed to encode response"}`, http.StatusInternalServerError)
	}
}

// ErrorJSON sends a JSON structure response with an error message.
//
// The JSON is sent with a 500 status code if no status code has set or if it was manually set to 200.
func (rb *ResponseBuilder) ErrorJSON(message string) {
	rb.headers["X-Content-Type-Options"] = "nosniff"

	if rb.statusCode == http.StatusOK {
		rb.statusCode = http.StatusInternalServerError
	}

	rb.JSON(map[string]string{"error": message})
}

// NoContent sends a 204 no content response with no data in the body.
func (rb *ResponseBuilder) NoContent() {
	rb.applyHeaders()
	rb.w.WriteHeader(http.StatusNoContent)
}

// applyHeaders applies the user set headers using the Header method.
func (rb *ResponseBuilder) applyHeaders() {
	for k, v := range rb.headers {
		rb.w.Header().Set(k, v)
	}
}
