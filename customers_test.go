package woocommerce

import (
	"context"
	"net/http"
	"testing"
)

func TestCustomersCreate(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/customers")
		writeJSON(w, &Customer{ID: 10, Email: "test@example.com"})
	})

	customer, _, err := client.Customers.Create(context.Background(), &Customer{Email: "test@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if customer.ID != 10 {
		t.Errorf("ID: got %d, want 10", customer.ID)
	}
	if customer.Email != "test@example.com" {
		t.Errorf("Email: got %s, want test@example.com", customer.Email)
	}
}

func TestCustomersGet(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/customers/10")
		writeJSON(w, &Customer{ID: 10, FirstName: "Jane", LastName: "Doe"})
	})

	customer, _, err := client.Customers.Get(context.Background(), "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if customer.FirstName != "Jane" {
		t.Errorf("FirstName: got %s, want Jane", customer.FirstName)
	}
}

func TestCustomersList(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/customers")
		if r.URL.Query().Get("role") != "customer" {
			t.Errorf("expected role=customer, got %q", r.URL.Query().Get("role"))
		}
		writeJSON(w, &[]Customer{{ID: 1}, {ID: 2}, {ID: 3}})
	})

	customers, _, err := client.Customers.List(context.Background(), &ListCustomerParams{Role: "customer"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*customers) != 3 {
		t.Errorf("len: got %d, want 3", len(*customers))
	}
}

func TestCustomersUpdate(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPut)
		assertPathSuffix(t, r, "/customers/10")
		writeJSON(w, &Customer{ID: 10, FirstName: "Updated"})
	})

	customer, _, err := client.Customers.Update(context.Background(), "10", &Customer{FirstName: "Updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if customer.FirstName != "Updated" {
		t.Errorf("FirstName: got %s, want Updated", customer.FirstName)
	}
}

func TestCustomersDelete(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodDelete)
		assertPathSuffix(t, r, "/customers/10")
		writeJSON(w, &Customer{ID: 10})
	})

	customer, _, err := client.Customers.Delete(context.Background(), "10", &DeleteCustomerParams{Force: "true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if customer.ID != 10 {
		t.Errorf("ID: got %d, want 10", customer.ID)
	}
}

func TestCustomersBatch(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/customers/batch")
		writeJSON(w, &BatchCustomerUpdateResponse{
			Create: &[]Customer{{ID: 11}, {ID: 12}},
		})
	})

	result, _, err := client.Customers.Batch(context.Background(), &BatchCustomerUpdate{
		Create: &[]Customer{{Email: "a@example.com"}, {Email: "b@example.com"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*result.Create) != 2 {
		t.Errorf("create len: got %d, want 2", len(*result.Create))
	}
}

func TestCustomersGetDownloads(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/customers/10/downloads")
		writeJSON(w, &[]CustomerDownload{
			{DownloadID: "abc", ProductName: "Product 1"},
		})
	})

	downloads, _, err := client.Customers.GetDownloads(context.Background(), "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*downloads) != 1 {
		t.Errorf("len: got %d, want 1", len(*downloads))
	}
	if (*downloads)[0].DownloadID != "abc" {
		t.Errorf("DownloadID: got %s, want abc", (*downloads)[0].DownloadID)
	}
}

func TestCustomersError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusNotFound, "Customer not found")
	})

	_, _, err := client.Customers.Get(context.Background(), "999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*ErrorResponse); !ok {
		t.Errorf("expected *ErrorResponse, got %T", err)
	}
}
