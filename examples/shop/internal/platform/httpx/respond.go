// Package httpx contains small HTTP helpers shared by every module's transport
// layer: JSON encoding/decoding, a uniform error envelope, and validation error
// mapping.
package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// ErrorResponse is the uniform error envelope returned by the API.
type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// JSON writes v as a JSON response with the given status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// Error writes a JSON error envelope with the given status and message.
func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, ErrorResponse{Error: msg})
}

// Decode reads and strictly decodes a JSON request body into v.
func Decode(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// QueryInt reads an integer query parameter, returning def when missing or
// malformed.
func QueryInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// ValidationError writes a 422 response, expanding validator.ValidationErrors
// into a per-field map when possible.
func ValidationError(w http.ResponseWriter, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make(map[string]string, len(ve))
		for _, fe := range ve {
			fields[lowerFirst(fe.Field())] = messageFor(fe)
		}
		JSON(w, http.StatusUnprocessableEntity, ErrorResponse{Error: "validation failed", Fields: fields})
		return
	}
	Error(w, http.StatusBadRequest, "invalid request body")
}

func messageFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "must be at least " + fe.Param() + " characters"
	case "max":
		return "must be at most " + fe.Param() + " characters"
	default:
		return "is invalid"
	}
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	b := []rune(s)
	if b[0] >= 'A' && b[0] <= 'Z' {
		b[0] += 'a' - 'A'
	}
	return string(b)
}
