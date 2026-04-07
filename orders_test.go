package woocommerce

import (
	"context"
	"net/http"
	"testing"
)

func TestOrdersCreate(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/orders")
		writeJSON(w, &Order{ID: 42, Status: "pending"})
	})

	order, resp, err := client.Orders.Create(context.Background(), &Order{Status: "pending"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status: got %d, want 200", resp.StatusCode)
	}
	if order.ID != 42 {
		t.Errorf("ID: got %d, want 42", order.ID)
	}
	if order.Status != "pending" {
		t.Errorf("Status: got %s, want pending", order.Status)
	}
}

func TestOrdersGet(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/orders/42")
		writeJSON(w, &Order{ID: 42, Status: "processing"})
	})

	order, _, err := client.Orders.Get(context.Background(), "42", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.ID != 42 {
		t.Errorf("ID: got %d, want 42", order.ID)
	}
}

func TestOrdersGet_WithParams(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("dp") != "2" {
			t.Errorf("expected dp=2, got %q", r.URL.Query().Get("dp"))
		}
		writeJSON(w, &Order{ID: 1})
	})

	client.Orders.Get(context.Background(), "1", &GetOrderParams{DecimalPoints: 2}) //nolint
}

func TestOrdersList(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/orders")
		if r.URL.Query().Get("customer") != "5" {
			t.Errorf("expected customer=5, got %q", r.URL.Query().Get("customer"))
		}
		writeJSON(w, &[]Order{{ID: 1}, {ID: 2}})
	})

	orders, _, err := client.Orders.List(context.Background(), &ListOrdersParams{Customer: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*orders) != 2 {
		t.Errorf("len: got %d, want 2", len(*orders))
	}
}

func TestOrdersUpdate(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPut)
		assertPathSuffix(t, r, "/orders/42")
		writeJSON(w, &Order{ID: 42, Status: "completed"})
	})

	order, _, err := client.Orders.Update(context.Background(), "42", &Order{Status: "completed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.Status != "completed" {
		t.Errorf("Status: got %s, want completed", order.Status)
	}
}

func TestOrdersDelete(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodDelete)
		assertPathSuffix(t, r, "/orders/42")
		writeJSON(w, &Order{ID: 42})
	})

	order, _, err := client.Orders.Delete(context.Background(), "42", &DeleteOrderParams{Force: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.ID != 42 {
		t.Errorf("ID: got %d, want 42", order.ID)
	}
}

func TestOrdersBatch(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/orders/batch")
		writeJSON(w, &BatchOrderUpdateResponse{
			Create: &[]Order{{ID: 10}},
		})
	})

	result, _, err := client.Orders.Batch(context.Background(), &BatchOrderUpdate{
		Create: &[]Order{{Status: "pending"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*result.Create) != 1 {
		t.Errorf("create len: got %d, want 1", len(*result.Create))
	}
}

func TestOrdersError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusNotFound, "Order not found")
	})

	_, _, err := client.Orders.Get(context.Background(), "999", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if apiErr, ok := err.(*ErrorResponse); !ok {
		t.Errorf("expected *ErrorResponse, got %T", err)
	} else if apiErr.Data.Status != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", apiErr.Data.Status)
	}
}
