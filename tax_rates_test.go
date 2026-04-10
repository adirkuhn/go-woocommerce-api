package woocommerce

import (
	"context"
	"net/http"
	"testing"
)

func TestTaxRatesCreate(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/taxes")
		writeJSON(w, &TaxRate{ID: 1, Country: "US", State: "CA", Rate: "8.250", Name: "CA State Tax"})
	})

	taxRate, _, err := client.TaxRates.Create(context.Background(), &TaxRate{
		Country: "US",
		State:   "CA",
		Rate:    "8.250",
		Name:    "CA State Tax",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if taxRate.ID != 1 {
		t.Errorf("ID: got %d, want 1", taxRate.ID)
	}
	if taxRate.Rate != "8.250" {
		t.Errorf("Rate: got %s, want 8.250", taxRate.Rate)
	}
	if taxRate.Country != "US" {
		t.Errorf("Country: got %s, want US", taxRate.Country)
	}
}

func TestTaxRatesGet(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/taxes/1")
		writeJSON(w, &TaxRate{ID: 1, Name: "CA State Tax", Class: "standard", Compound: false, Shipping: true})
	})

	taxRate, _, err := client.TaxRates.Get(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if taxRate.Class != "standard" {
		t.Errorf("Class: got %s, want standard", taxRate.Class)
	}
	if !taxRate.Shipping {
		t.Error("expected Shipping=true")
	}
}

func TestTaxRatesList(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/taxes")
		if r.URL.Query().Get("class") != "standard" {
			t.Errorf("expected class=standard, got %q", r.URL.Query().Get("class"))
		}
		if r.URL.Query().Get("per_page") != "10" {
			t.Errorf("expected per_page=10, got %q", r.URL.Query().Get("per_page"))
		}
		writeJSON(w, &[]TaxRate{{ID: 1}, {ID: 2}})
	})

	taxRates, _, err := client.TaxRates.List(context.Background(), &ListTaxRatesParams{
		Class:   "standard",
		PerPage: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(taxRates) != 2 {
		t.Errorf("len: got %d, want 2", len(taxRates))
	}
}

func TestTaxRatesUpdate(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPut)
		assertPathSuffix(t, r, "/taxes/1")
		writeJSON(w, &TaxRate{ID: 1, Rate: "9.000"})
	})

	taxRate, _, err := client.TaxRates.Update(context.Background(), "1", &TaxRate{Rate: "9.000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if taxRate.Rate != "9.000" {
		t.Errorf("Rate: got %s, want 9.000", taxRate.Rate)
	}
}

func TestTaxRatesDelete(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodDelete)
		assertPathSuffix(t, r, "/taxes/1")
		if r.URL.Query().Get("force") != "true" {
			t.Errorf("expected force=true, got %q", r.URL.Query().Get("force"))
		}
		writeJSON(w, &TaxRate{ID: 1})
	})

	taxRate, _, err := client.TaxRates.Delete(context.Background(), "1", &DeleteTaxRateParams{Force: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if taxRate.ID != 1 {
		t.Errorf("ID: got %d, want 1", taxRate.ID)
	}
}

func TestTaxRatesBatch(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/taxes/batch")
		writeJSON(w, &BatchTaxRateUpdateResponse{
			Create: &[]TaxRate{{ID: 3, Rate: "5.000"}, {ID: 4, Rate: "10.000"}},
		})
	})

	result, _, err := client.TaxRates.Batch(context.Background(), &BatchTaxRateUpdate{
		Create: &[]TaxRate{
			{Country: "US", State: "NY", Rate: "5.000"},
			{Country: "US", State: "TX", Rate: "10.000"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*result.Create) != 2 {
		t.Errorf("create len: got %d, want 2", len(*result.Create))
	}
}

func TestTaxRatesError(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusNotFound, "Tax rate not found")
	})

	_, _, err := client.TaxRates.Get(context.Background(), "999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if apiErr, ok := err.(*ErrorResponse); !ok {
		t.Errorf("expected *ErrorResponse, got %T", err)
	} else if apiErr.Data.Status != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", apiErr.Data.Status)
	}
}
