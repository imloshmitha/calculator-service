// Package handler exposes the calculator over HTTP.
package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"calculator-service/internal/calculator"
)

// maxBodyBytes limits a calculate request. One mebibyte is enough for two
// numbers and an operation name, and it rejects accidental large bodies.
const maxBodyBytes = 1 << 20

type calculateRequest struct {
	Operation string   `json:"operation"`
	A         *float64 `json:"a"`
	B         *float64 `json:"b"`
}

type resultResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type statusError struct {
	status  int
	code    string
	message string
}

func (e *statusError) Error() string {
	return e.message
}

var (
	errInvalidJSON = &statusError{
		status:  http.StatusBadRequest,
		code:    "invalid_json",
		message: "request body must be one JSON object",
	}
	errUnknownField = &statusError{
		status:  http.StatusBadRequest,
		code:    "invalid_json",
		message: "request body contains an unknown field",
	}
	errTrailingJSON = &statusError{
		status:  http.StatusBadRequest,
		code:    "invalid_json",
		message: "request body must not contain trailing data",
	}
	errRequestTooLarge = &statusError{
		status:  http.StatusBadRequest,
		code:    "invalid_json",
		message: "request body exceeds 1 MiB",
	}
)

// NewMux returns the HTTP routes served by the calculator API.
func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/calculate", Calculate)
	mux.HandleFunc("/health", Health)
	return mux
}

// Health reports that the process is serving requests.
func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeMethodNotAllowed(w, http.MethodGet)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Calculate handles POST /api/calculate.
func Calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeMethodNotAllowed(w, http.MethodPost)
		return
	}

	req, err := decodeCalculateRequest(w, r)
	if err != nil {
		writeStatusError(w, err)
		return
	}

	result, err := calculator.Calculate(req.Operation, req.A, req.B)
	if err != nil {
		writeCalculationError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resultResponse{Result: result})
}

func decodeCalculateRequest(w http.ResponseWriter, r *http.Request) (calculateRequest, error) {
	if r.Body == nil {
		return calculateRequest{}, errInvalidJSON
	}
	defer r.Body.Close()

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var req calculateRequest
	if err := dec.Decode(&req); err != nil {
		return calculateRequest{}, classifyDecodeError(err)
	}

	// A second value means the body contained more than one JSON document.
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if isMaxBytes(err) {
			return calculateRequest{}, errRequestTooLarge
		}
		return calculateRequest{}, errTrailingJSON
	}
	return req, nil
}

func classifyDecodeError(err error) error {
	if isMaxBytes(err) {
		return errRequestTooLarge
	}

	var syntax *json.SyntaxError
	var typeErr *json.UnmarshalTypeError
	switch {
	case errors.Is(err, io.EOF), errors.Is(err, io.ErrUnexpectedEOF):
		return errInvalidJSON
	case errors.As(err, &syntax), errors.As(err, &typeErr):
		return errInvalidJSON
	case strings.Contains(err.Error(), "unknown field"):
		return errUnknownField
	default:
		return errInvalidJSON
	}
}

func isMaxBytes(err error) bool {
	var maxBytes *http.MaxBytesError
	return errors.As(err, &maxBytes)
}

func writeCalculationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calculator.ErrMissingOperation):
		writeError(w, http.StatusBadRequest, "missing_operation", err.Error())
	case errors.Is(err, calculator.ErrUnsupportedOperation):
		writeError(w, http.StatusBadRequest, "unsupported_operation", err.Error())
	case errors.Is(err, calculator.ErrMissingOperandA), errors.Is(err, calculator.ErrMissingOperandB):
		writeError(w, http.StatusBadRequest, "missing_operand", err.Error())
	case errors.Is(err, calculator.ErrDivisionByZero):
		writeError(w, http.StatusBadRequest, "division_by_zero", err.Error())
	case errors.Is(err, calculator.ErrNegativeSqrt):
		writeError(w, http.StatusBadRequest, "negative_sqrt", err.Error())
	case errors.Is(err, calculator.ErrUndefinedPower):
		writeError(w, http.StatusBadRequest, "undefined_power", err.Error())
	case errors.Is(err, calculator.ErrNonFiniteOperand), errors.Is(err, calculator.ErrNonFiniteResult):
		writeError(w, http.StatusBadRequest, "non_finite_result", err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal_error", "internal error")
	}
}

func writeStatusError(w http.ResponseWriter, err error) {
	var statusErr *statusError
	if errors.As(err, &statusErr) {
		writeError(w, statusErr.status, statusErr.code, statusErr.message)
		return
	}
	writeError(w, http.StatusBadRequest, "invalid_json", errInvalidJSON.message)
}

func writeMethodNotAllowed(w http.ResponseWriter, allow string) {
	w.Header().Set("Allow", allow)
	writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Code: code, Message: message})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(body)
}
