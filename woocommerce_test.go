package woocommerce

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestClient creates a test HTTP server and a Client pointing at it.
// The server is closed automatically when the test ends.
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	client, err := New(srv.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	client.Authenticate("ck_test", "cs_test")
	return client
}

// writeJSON encodes v as JSON and writes it to w with a 200 status.
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}

// writeAPIError writes a WooCommerce-style error response.
func writeAPIError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, `{"code":"woocommerce_rest_error","message":%q,"data":{"status":%d}}`, msg, statusCode)
}

// assertMethod fails the test if the request method doesn't match.
func assertMethod(t *testing.T, r *http.Request, method string) {
	t.Helper()
	if r.Method != method {
		t.Errorf("method: got %s, want %s", r.Method, method)
	}
}

// assertPathSuffix fails the test if the request path doesn't end with suffix.
func assertPathSuffix(t *testing.T, r *http.Request, suffix string) {
	t.Helper()
	if !strings.HasSuffix(r.URL.Path, suffix) {
		t.Errorf("path %q does not end with %q", r.URL.Path, suffix)
	}
}

func TestNew(t *testing.T) {
	t.Run("valid URL", func(t *testing.T) {
		c, err := New("https://example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c == nil {
			t.Fatal("expected non-nil client")
		}
	})

	t.Run("empty URL", func(t *testing.T) {
		_, err := New("")
		if err == nil {
			t.Fatal("expected error for empty URL")
		}
	})
}

func TestNewWithConfig(t *testing.T) {
	t.Run("custom http client", func(t *testing.T) {
		httpClient := &http.Client{}
		c, err := NewWithConfig(ClientConfig{
			RestEndpointURL: "https://example.com",
			HttpClient:      httpClient,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.client != httpClient {
			t.Error("expected custom http client to be used")
		}
	})

	t.Run("custom endpoint version", func(t *testing.T) {
		c, err := NewWithConfig(ClientConfig{
			RestEndpointURL:     "https://example.com",
			RestEndpointVersion: "v2",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if c.config.RestEndpointVersion != "v2" {
			t.Errorf("expected v2, got %s", c.config.RestEndpointVersion)
		}
	})

	t.Run("empty URL", func(t *testing.T) {
		_, err := NewWithConfig(ClientConfig{})
		if err == nil {
			t.Fatal("expected error for empty URL")
		}
	})
}

func TestAuthenticate(t *testing.T) {
	var gotAuth string
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		writeJSON(w, struct{}{})
	})

	client.Authenticate("ck_key", "cs_secret")
	req, err := client.NewRequest(context.Background(), http.MethodGet, "/orders", nil, nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}
	client.Do(req, nil) //nolint

	if !strings.HasPrefix(gotAuth, "Basic ") {
		t.Errorf("expected Basic auth header, got %q", gotAuth)
	}
}

func TestNewRequest_QueryParams(t *testing.T) {
	client, err := New("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	opts := &ListOrdersParams{Page: 2, PerPage: 10}
	req, err := client.NewRequest(context.Background(), http.MethodGet, "/orders", opts, nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	q := req.URL.Query()
	if q.Get("page") != "2" {
		t.Errorf("expected page=2, got %q", q.Get("page"))
	}
	if q.Get("per_page") != "10" {
		t.Errorf("expected per_page=10, got %q", q.Get("per_page"))
	}
}

func TestNewRequest_Headers(t *testing.T) {
	client, err := New("https://example.com")
	if err != nil {
		t.Fatal(err)
	}
	client.Authenticate("ck_key", "cs_secret")

	req, err := client.NewRequest(context.Background(), http.MethodGet, "/orders", nil, nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	if req.Header.Get("Accept") != "application/json" {
		t.Errorf("missing Accept header")
	}
	if req.Header.Get("Content-type") != "application/json" {
		t.Errorf("missing Content-type header")
	}
	if !strings.HasPrefix(req.Header.Get("User-Agent"), "go-woocommerce-api") {
		t.Errorf("unexpected User-Agent: %s", req.Header.Get("User-Agent"))
	}
}

func TestDo_ErrorResponse(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusNotFound, "Resource not found")
	})

	req, err := client.NewRequest(context.Background(), http.MethodGet, "/orders/999", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*ErrorResponse)
	if !ok {
		t.Fatalf("expected *ErrorResponse, got %T", err)
	}
	if apiErr.Message != "Resource not found" {
		t.Errorf("unexpected message: %s", apiErr.Message)
	}
	if apiErr.Data.Status != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", apiErr.Data.Status)
	}
}

func TestDo_ContextCancellation(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, struct{}{})
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	req, err := client.NewRequest(ctx, http.MethodGet, "/orders", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(req, nil)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestCheckResponse_2xx(t *testing.T) {
	for _, code := range []int{200, 201, 204} {
		resp := &http.Response{StatusCode: code}
		if err := checkResponse(resp); err != nil {
			t.Errorf("code %d: unexpected error: %v", code, err)
		}
	}
}
