package woocommerce

import (
	"context"
	_ "embed"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

//go:embed test_data/refunds.json
var refundsJSON []byte

var refundIgnoreOpts = []cmp.Option{
	// ignore unstable / API-generated fields
	cmpopts.IgnoreFields(Refund{}),
}

func TestRefundsCreate(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/orders/10/refunds")
		writeJSON(w, &Refund{ID: 1, Amount: "15.00", Reason: "Customer request"})
	})

	refund, _, err := client.Refunds.Create(context.Background(), "10", &Refund{Amount: "15.00", Reason: "Customer request"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if refund.ID != 1 {
		t.Errorf("ID: got %d, want 1", refund.ID)
	}
	if refund.Amount != "15.00" {
		t.Errorf("Amount: got %s, want 15.00", refund.Amount)
	}
}

func TestRefundsGet(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/orders/10/refunds/1")
		writeJSON(w, &Refund{ID: 1, Amount: "15.00"})
	})

	refund, _, err := client.Refunds.Get(context.Background(), "10", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if refund.ID != 1 {
		t.Errorf("ID: got %d, want 1", refund.ID)
	}
}

func TestRefundsList(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/orders/10/refunds")
		if r.URL.Query().Get("per_page") != "5" {
			t.Errorf("expected per_page=5, got %q", r.URL.Query().Get("per_page"))
		}
		writeJSON(w, &[]Refund{{ID: 1}, {ID: 2}})
	})

	refunds, _, err := client.Refunds.List(context.Background(), "10", &ListRefundParams{PerPage: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(refunds) != 2 {
		t.Errorf("len: got %d, want 2", len(refunds))
	}
}

func TestRefundsDelete(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodDelete)
		assertPathSuffix(t, r, "/orders/10/refunds/1")
		writeJSON(w, &Refund{ID: 1})
	})

	refund, _, err := client.Refunds.Delete(context.Background(), "10", "1", &DeleteRefundParams{Force: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if refund.ID != 1 {
		t.Errorf("ID: got %d, want 1", refund.ID)
	}
}

func TestRefundsError(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusNotFound, "Refund not found")
	})

	_, _, err := client.Refunds.Get(context.Background(), "10", "999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if apiErr, ok := err.(*ErrorResponse); !ok {
		t.Errorf("expected *ErrorResponse, got %T", err)
	} else if apiErr.Message != "Refund not found" {
		t.Errorf("message: got %q, want %q", apiErr.Message, "Refund not found")
	}
}

func TestRefundsList_RealJSON(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/orders/723/refunds")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(refundsJSON)))
	})

	refunds, _, err := client.Refunds.List(context.Background(), "723", &ListRefundParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(refunds) != 2 {
		t.Fatalf("len: got %d, want 2", len(refunds))
	}

	wantRefunds := loadRefundsFixture(t)

	for i := range refunds {
		assertRefund(t, &refunds[i], &wantRefunds[i])
	}
}

func assertRefund(t *testing.T, got *Refund, want *Refund) {
	t.Helper()

	if diff := cmp.Diff(want, got, refundIgnoreOpts...); diff != "" {
		t.Fatalf("refund mismatch (-want +got):\n%s", diff)
	}
}

func loadRefundsFixture(t *testing.T) []Refund {
	t.Helper()

	var refunds []Refund
	if err := json.Unmarshal(refundsJSON, &refunds); err != nil {
		t.Fatalf("invalid fixture: %v", err)
	}

	return refunds
}
