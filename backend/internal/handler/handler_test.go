package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		method  string
		body    string
		status  int
		code    string
		message string
		result  *float64
		allow   string
	}{
		{
			name:   "add",
			method: http.MethodPost,
			body:   `{"operation":"add","a":19,"b":9}`,
			status: http.StatusOK,
			result: ptr(28),
		},
		{
			name:   "percentage",
			method: http.MethodPost,
			body:   `{"operation":"percentage","a":200,"b":15}`,
			status: http.StatusOK,
			result: ptr(30),
		},
		{
			name:   "sqrt",
			method: http.MethodPost,
			body:   `{"operation":"sqrt","a":9}`,
			status: http.StatusOK,
			result: ptr(3),
		},
		{
			name:   "divide decimals",
			method: http.MethodPost,
			body:   `{"operation":"divide","a":-9,"b":4}`,
			status: http.StatusOK,
			result: ptr(-2.25),
		},
		{
			name:    "malformed json",
			method:  http.MethodPost,
			body:    `{`,
			status:  http.StatusBadRequest,
			code:    "invalid_json",
			message: "request body must be one JSON object",
		},
		{
			name:    "empty body",
			method:  http.MethodPost,
			body:    ``,
			status:  http.StatusBadRequest,
			code:    "invalid_json",
			message: "request body must be one JSON object",
		},
		{
			name:    "json array",
			method:  http.MethodPost,
			body:    `[1,2]`,
			status:  http.StatusBadRequest,
			code:    "invalid_json",
			message: "request body must be one JSON object",
		},
		{
			name:    "wrong operand type",
			method:  http.MethodPost,
			body:    `{"operation":"add","a":"no","b":1}`,
			status:  http.StatusBadRequest,
			code:    "invalid_json",
			message: "request body must be one JSON object",
		},
		{
			name:    "unknown field",
			method:  http.MethodPost,
			body:    `{"operation":"add","a":1,"b":2,"c":3}`,
			status:  http.StatusBadRequest,
			code:    "invalid_json",
			message: "request body contains an unknown field",
		},
		{
			name:    "trailing json",
			method:  http.MethodPost,
			body:    `{"operation":"add","a":1,"b":2}{"operation":"add"}`,
			status:  http.StatusBadRequest,
			code:    "invalid_json",
			message: "request body must not contain trailing data",
		},
		{
			name:    "null operand",
			method:  http.MethodPost,
			body:    `{"operation":"add","a":null,"b":2}`,
			status:  http.StatusBadRequest,
			code:    "missing_operand",
			message: "operand a is required",
		},
		{
			name:    "missing operand",
			method:  http.MethodPost,
			body:    `{"operation":"subtract","a":10}`,
			status:  http.StatusBadRequest,
			code:    "missing_operand",
			message: "operand b is required",
		},
		{
			name:    "missing operation",
			method:  http.MethodPost,
			body:    `{"a":1,"b":2}`,
			status:  http.StatusBadRequest,
			code:    "missing_operation",
			message: "operation is required",
		},
		{
			name:    "unsupported operation",
			method:  http.MethodPost,
			body:    `{"operation":"modulo","a":1,"b":2}`,
			status:  http.StatusBadRequest,
			code:    "unsupported_operation",
			message: `operation is not supported: "modulo"`,
		},
		{
			name:    "division by zero",
			method:  http.MethodPost,
			body:    `{"operation":"divide","a":1,"b":0}`,
			status:  http.StatusBadRequest,
			code:    "division_by_zero",
			message: "division by zero",
		},
		{
			name:    "negative square root",
			method:  http.MethodPost,
			body:    `{"operation":"sqrt","a":-4}`,
			status:  http.StatusBadRequest,
			code:    "negative_sqrt",
			message: "square root of a negative number is not defined",
		},
		{
			name:    "zero to the zero",
			method:  http.MethodPost,
			body:    `{"operation":"power","a":0,"b":0}`,
			status:  http.StatusBadRequest,
			code:    "undefined_power",
			message: "zero raised to the power of zero is undefined",
		},
		{
			name:    "overflow",
			method:  http.MethodPost,
			body:    `{"operation":"multiply","a":1e308,"b":10}`,
			status:  http.StatusBadRequest,
			code:    "non_finite_result",
			message: "calculation result is not a finite number",
		},
		{
			name:    "get method",
			method:  http.MethodGet,
			status:  http.StatusMethodNotAllowed,
			code:    "method_not_allowed",
			message: "method not allowed",
			allow:   http.MethodPost,
		},
		{
			name:    "put method",
			method:  http.MethodPut,
			body:    `{"operation":"add","a":1,"b":2}`,
			status:  http.StatusMethodNotAllowed,
			code:    "method_not_allowed",
			message: "method not allowed",
			allow:   http.MethodPost,
		},
		{
			name:    "delete method",
			method:  http.MethodDelete,
			status:  http.StatusMethodNotAllowed,
			code:    "method_not_allowed",
			message: "method not allowed",
			allow:   http.MethodPost,
		},
	}

	mux := NewMux()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.method, "/api/calculate", strings.NewReader(tt.body))
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != tt.status {
				t.Fatalf("status = %d, want %d; body = %s", rr.Code, tt.status, rr.Body.String())
			}
			if got := rr.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", got)
			}
			if tt.allow != "" && rr.Header().Get("Allow") != tt.allow {
				t.Fatalf("Allow = %q, want %q", rr.Header().Get("Allow"), tt.allow)
			}

			if tt.result != nil {
				got := decodeResult(t, rr.Body.Bytes())
				if got != *tt.result {
					t.Fatalf("result = %v, want %v", got, *tt.result)
				}
				return
			}

			code, message := decodeError(t, rr.Body.Bytes())
			if code != tt.code {
				t.Fatalf("code = %q, want %q", code, tt.code)
			}
			if message != tt.message {
				t.Fatalf("message = %q, want %q", message, tt.message)
			}
		})
	}
}

func TestHealth(t *testing.T) {
	t.Parallel()

	mux := NewMux()

	t.Run("ok", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body = %s", rr.Code, rr.Body.String())
		}
		if got := rr.Header().Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want application/json", got)
		}

		dec := json.NewDecoder(strings.NewReader(rr.Body.String()))
		dec.DisallowUnknownFields()
		var body struct {
			Status string `json:"status"`
		}
		if err := dec.Decode(&body); err != nil {
			t.Fatalf("decode health: %v", err)
		}
		if body.Status != "ok" {
			t.Fatalf("status field = %q, want ok", body.Status)
		}
		if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
			t.Fatalf("trailing health body: %v", err)
		}
	})

	t.Run("post not allowed", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodPost, "/health", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Fatalf("status = %d, want 405", rr.Code)
		}
		if rr.Header().Get("Allow") != http.MethodGet {
			t.Fatalf("Allow = %q, want GET", rr.Header().Get("Allow"))
		}
		code, message := decodeError(t, rr.Body.Bytes())
		if code != "method_not_allowed" || message != "method not allowed" {
			t.Fatalf("error = %s %s", code, message)
		}
	})
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	mux := NewMux()
	tests := []struct {
		name   string
		method string
		path   string
	}{
		{name: "root", method: http.MethodGet, path: "/"},
		{name: "unknown path", method: http.MethodGet, path: "/no-such-route"},
		{name: "unknown api path", method: http.MethodPost, path: "/api/other"},
		{name: "trailing slash", method: http.MethodGet, path: "/health/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want 404; body = %s", rr.Code, rr.Body.String())
			}
			if got := rr.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", got)
			}
			code, message := decodeError(t, rr.Body.Bytes())
			if code != "not_found" || message != "not found" {
				t.Fatalf("error = %s %s", code, message)
			}
		})
	}
}

func TestCalculateRejectsOversizedBody(t *testing.T) {
	t.Parallel()

	// A long string stays valid JSON until the size limit stops the read.
	body := `{"operation":"` + strings.Repeat("a", maxBodyBytes) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(body))
	rr := httptest.NewRecorder()
	NewMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rr.Code, rr.Body.String())
	}
	code, message := decodeError(t, rr.Body.Bytes())
	if code != "invalid_json" || message != "request body exceeds 1 MiB" {
		t.Fatalf("error = %s %s", code, message)
	}
}

func ptr(v float64) *float64 {
	return &v
}

func decodeResult(t *testing.T, body []byte) float64 {
	t.Helper()
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.DisallowUnknownFields()
	var resp resultResponse
	if err := dec.Decode(&resp); err != nil {
		t.Fatalf("decode result: %v; body = %s", err, body)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		t.Fatalf("trailing result body: %v", err)
	}
	return resp.Result
}

func decodeError(t *testing.T, body []byte) (string, string) {
	t.Helper()
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.DisallowUnknownFields()
	var resp errorResponse
	if err := dec.Decode(&resp); err != nil {
		t.Fatalf("decode error: %v; body = %s", err, body)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		t.Fatalf("trailing error body: %v", err)
	}
	if resp.Code == "" || resp.Message == "" {
		t.Fatalf("error response missing code or message: %s", body)
	}
	return resp.Code, resp.Message
}
