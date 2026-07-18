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

//go:embed test_data/products.json
var productsJSON []byte

var productIgnoreOpts = []cmp.Option{
	// ignore unstable / API-generated fields
	cmpopts.IgnoreFields(Product{}),
}

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
	if len(products) != 3 {
		t.Errorf("len: got %d, want 3", len(products))
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
			Delete: []Product{{ID: 21}, {ID: 22}},
		})
		_ = ids
	})

	ids := []int{21, 22}
	result, _, err := client.Products.Batch(context.Background(), &BatchProductUpdate{
		Delete: ids,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Delete) != 2 {
		t.Errorf("delete len: got %d, want 2", len(result.Delete))
	}
}

func TestFilterParams(t *testing.T) {
	filter := &ListProductParams{
		After:   "2022-01-01T00:00:00Z",
		PerPage: 20,
	}

	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/products")
		if r.URL.Query().Get("after") != "2022-01-01T00:00:00Z" {
			t.Errorf("expected after=2022-01-01T00:00:00Z, got %q", r.URL.Query().Get("after"))
		}
		if r.URL.Query().Get("per_page") != "20" {
			t.Errorf("expected per_page=20, got %q", r.URL.Query().Get("per_page"))
		}
		writeJSON(w, &[]Product{{ID: 1}, {ID: 2}, {ID: 3}})
	})

	products, _, err := client.Products.List(context.Background(), filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 3 {
		t.Errorf("len: got %d, want 3", len(products))
	}
}

func TestFilterParamsWithFields(t *testing.T) {
	filter := &ListProductParams{
		Fields: []string{"id", "name"},
	}

	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/products")
		if r.URL.Query().Get("_fields") != "id,name" {
			t.Errorf("expected _fields=id,name, got %q", r.URL.Query().Get("_fields"))
		}
		writeJSON(w, &[]Product{{ID: 1}, {ID: 2}, {ID: 3}})
	})

	products, _, err := client.Products.List(context.Background(), filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 3 {
		t.Errorf("len: got %d, want 3", len(products))
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

// TestProductsList_RealJSON verifies the List endpoint correctly deserialises
// a response containing a real product payload (including meta_data as an
// empty array, which previously caused an unmarshal panic).
func TestProductsList_RealJSON(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/products")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(productsJSON)))
	})

	products, _, err := client.Products.List(context.Background(), &ListProductParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 2 {
		t.Fatalf("len: got %d, want 2", len(products))
	}

	wantProducts := loadProductsFixture(t)

	for i := range products {
		assertProduct(t, &products[i], &wantProducts[i])
	}
}

func TestProduct_SetStockQuantity_ToZero_NotIgnored(t *testing.T) {
	stockQuantity := 0
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPut)
		assertPathSuffix(t, r, "/products/1")

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}

		// Verify field exists
		value, ok := payload["stock_quantity"]
		if !ok {
			t.Fatal("stock_quantity was omitted from payload")
		}

		// JSON numbers decode as float64
		if value.(float64) != 0 {
			t.Fatalf("stock_quantity: got %v, want 0", value)
		}

		writeJSON(w, &Product{
			ID:            1,
			StockQuantity: &stockQuantity,
		})
	})

	_, _, err := client.Products.Update(
		context.Background(),
		"1",
		&Product{
			ID:            1,
			StockQuantity: &stockQuantity,
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestProductsGet_NumericSalePrice_UnmarshalsAsString reproduces a real
// payload from a store running a bundle-product plugin ("woosb" type) that
// recalculates sale_price and re-serialises it as a bare JSON number instead
// of WooCommerce's documented string — this used to fail the whole decode
// with "cannot unmarshal number into Go struct field Product.sale_price of
// type string".
func TestProductsGet_NumericSalePrice_UnmarshalsAsString(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"id": 3263,
			"name": "PROMOÇÃO - Pacote com 2 Tapiocas da Terrinha",
			"type": "woosb",
			"price": "6.99",
			"regular_price": "7.98",
			"sale_price": 6.99
		}`))
	})

	product, _, err := client.Products.Get(context.Background(), "3263")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.Price != "6.99" {
		t.Errorf("Price: got %q, want %q", product.Price, "6.99")
	}
	if product.RegularPrice != "7.98" {
		t.Errorf("RegularPrice: got %q, want %q", product.RegularPrice, "7.98")
	}
	if product.SalePrice != "6.99" {
		t.Errorf("SalePrice: got %q, want %q", product.SalePrice, "6.99")
	}
}

func TestFlexibleString_NumericAndStringAndNull(t *testing.T) {
	cases := []struct {
		name string
		json string
		want FlexibleString
	}{
		{"quoted string", `"6.99"`, "6.99"},
		{"bare integer", `7`, "7"},
		{"bare decimal", `6.99`, "6.99"},
		{"null", `null`, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got FlexibleString
			if err := json.Unmarshal([]byte(tc.json), &got); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestFlexibleString_MarshalsAsPlainJSONString(t *testing.T) {
	out, err := json.Marshal(FlexibleString("6.99"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != `"6.99"` {
		t.Errorf("got %s, want %s", out, `"6.99"`)
	}
}

func assertProduct(t *testing.T, got *Product, want *Product) {
	t.Helper()

	if diff := cmp.Diff(want, got, productIgnoreOpts...); diff != "" {
		t.Fatalf("product mismatch (-want +got):\n%s", diff)
	}
}

func loadProductsFixture(t *testing.T) []Product {
	t.Helper()

	var products []Product
	if err := json.Unmarshal(productsJSON, &products); err != nil {
		t.Fatalf("invalid fixture: %v", err)
	}

	return products
}
