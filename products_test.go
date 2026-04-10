package woocommerce

import (
	"context"
	"net/http"
	"testing"
)

func TestProductsCreate(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/products")
		writeJSON(w, &Product{ID: 20, Name: "T-Shirt", Status: "publish"})
	})

	product, _, err := client.Products.Create(context.Background(), &Product{Name: "T-Shirt", Status: "publish"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.ID != 20 {
		t.Errorf("ID: got %d, want 20", product.ID)
	}
	if product.Name != "T-Shirt" {
		t.Errorf("Name: got %s, want T-Shirt", product.Name)
	}
}

func TestProductsGet(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/products/20")
		writeJSON(w, &Product{ID: 20, Sku: "TS-001", StockStatus: "instock"})
	})

	product, _, err := client.Products.Get(context.Background(), "20")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.Sku != "TS-001" {
		t.Errorf("Sku: got %s, want TS-001", product.Sku)
	}
	if product.StockStatus != "instock" {
		t.Errorf("StockStatus: got %s, want instock", product.StockStatus)
	}
}

func TestProductsList(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/products")
		if r.URL.Query().Get("status") != "publish" {
			t.Errorf("expected status=publish, got %q", r.URL.Query().Get("status"))
		}
		if r.URL.Query().Get("per_page") != "20" {
			t.Errorf("expected per_page=20, got %q", r.URL.Query().Get("per_page"))
		}
		writeJSON(w, &[]Product{{ID: 1}, {ID: 2}, {ID: 3}})
	})

	products, _, err := client.Products.List(context.Background(), &ListProductParams{
		Status:  "publish",
		PerPage: 20,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*products) != 3 {
		t.Errorf("len: got %d, want 3", len(*products))
	}
}

func TestProductsUpdate(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPut)
		assertPathSuffix(t, r, "/products/20")
		writeJSON(w, &Product{ID: 20, RegularPrice: "29.99"})
	})

	product, _, err := client.Products.Update(context.Background(), "20", &Product{RegularPrice: "29.99"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.RegularPrice != "29.99" {
		t.Errorf("RegularPrice: got %s, want 29.99", product.RegularPrice)
	}
}

func TestProductsDelete(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodDelete)
		assertPathSuffix(t, r, "/products/20")
		writeJSON(w, &Product{ID: 20})
	})

	product, _, err := client.Products.Delete(context.Background(), "20", &DeleteProductParams{Force: "true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.ID != 20 {
		t.Errorf("ID: got %d, want 20", product.ID)
	}
}

func TestProductsBatch(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/products/batch")
		ids := []int{21, 22}
		writeJSON(w, &BatchProductUpdateResponse{
			Delete: &[]Product{{ID: 21}, {ID: 22}},
		})
		_ = ids
	})

	ids := []int{21, 22}
	result, _, err := client.Products.Batch(context.Background(), &BatchProductUpdate{
		Delete: &ids,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*result.Delete) != 2 {
		t.Errorf("delete len: got %d, want 2", len(*result.Delete))
	}
}

func TestProductsError(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusNotFound, "Product not found")
	})

	_, _, err := client.Products.Get(context.Background(), "999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*ErrorResponse); !ok {
		t.Errorf("expected *ErrorResponse, got %T", err)
	}
}
