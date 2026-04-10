package woocommerce

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestServerFn is like newTestServer but lets you provide a full handler
// for cases that need to inspect query params, decode the request body, etc.
func newTestServerFn(t *testing.T, handler http.HandlerFunc) *WooCommerce {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	woo, err := New(Config{
		ShopURL:        srv.URL,
		ConsumerKey:    "ck_test",
		ConsumerSecret: "cs_test",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	return woo
}

func assertMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if r.Method != want {
		t.Errorf("method: got %s, want %s", r.Method, want)
	}
}

func assertPathSuffix(t *testing.T, r *http.Request, suffix string) {
	t.Helper()
	if !strings.HasSuffix(r.URL.Path, suffix) {
		t.Errorf("path: got %s, want suffix %s", r.URL.Path, suffix)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func writeAPIError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Code:    "error",
		Message: message,
		Data:    ErrorData{Status: status},
	})
}
