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

//go:embed test_data/coupons.json
var couponsJSON []byte

var couponIgnoreOpts = []cmp.Option{
	// ignore unstable / API-generated fields
	cmpopts.IgnoreFields(Coupon{}),
}

func TestCouponsCreate(t *testing.T) {
	woo := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/coupons")
		writeJSON(w, &Coupon{ID: 1, Code: "SAVE10"})
	})

	coupon, _, err := woo.Coupons.Create(context.Background(), &Coupon{Code: "SAVE10"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coupon.ID != 1 {
		t.Errorf("ID: got %d, want 1", coupon.ID)
	}
	if coupon.Code != "SAVE10" {
		t.Errorf("Code: got %s, want SAVE10", coupon.Code)
	}
}

func TestCouponsGet(t *testing.T) {
	woo := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/coupons/1")
		writeJSON(w, &Coupon{ID: 1, Code: "SAVE10", DiscountType: "percent"})
	})

	coupon, _, err := woo.Coupons.Get(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coupon.DiscountType != "percent" {
		t.Errorf("DiscountType: got %s, want percent", coupon.DiscountType)
	}
}

func TestCouponsList(t *testing.T) {
	woo := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/coupons")
		if r.URL.Query().Get("code") != "SAVE10" {
			t.Errorf("expected code=SAVE10, got %q", r.URL.Query().Get("code"))
		}
		writeJSON(w, &[]Coupon{{ID: 1}, {ID: 2}})
	})

	coupons, _, err := woo.Coupons.List(context.Background(), &ListCouponParams{Code: "SAVE10"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(coupons) != 2 {
		t.Errorf("len: got %d, want 2", len(coupons))
	}
}

func TestCouponsUpdate(t *testing.T) {
	woo := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPut)
		assertPathSuffix(t, r, "/coupons/1")
		writeJSON(w, &Coupon{ID: 1, Amount: "20"})
	})

	coupon, _, err := woo.Coupons.Update(context.Background(), "1", &Coupon{Amount: "20"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coupon.Amount != "20" {
		t.Errorf("Amount: got %s, want 20", coupon.Amount)
	}
}

func TestCouponsDelete(t *testing.T) {
	woo := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodDelete)
		assertPathSuffix(t, r, "/coupons/1")
		writeJSON(w, &Coupon{ID: 1})
	})

	coupon, _, err := woo.Coupons.Delete(context.Background(), "1", &DeleteCouponParams{Force: "true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coupon.ID != 1 {
		t.Errorf("ID: got %d, want 1", coupon.ID)
	}
}

func TestCouponsBatch(t *testing.T) {
	woo := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/coupons/batch")
		writeJSON(w, &BatchCouponUpdateResponse{
			Create: []Coupon{{ID: 5}},
		})
	})

	result, _, err := woo.Coupons.Batch(context.Background(), &BatchCouponUpdate{
		Create: []Coupon{{Code: "NEW"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Create) != 1 {
		t.Errorf("create len: got %d, want 1", len(result.Create))
	}
}

func TestCouponsError(t *testing.T) {
	woo := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusUnprocessableEntity, "Invalid coupon code")
	})

	_, _, err := woo.Coupons.Create(context.Background(), &Coupon{})
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
func TestCouponsList_RealJSON(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/coupons")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(couponsJSON)))
	})

	coupons, _, err := client.Coupons.List(context.Background(), &ListCouponParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(coupons) != 2 {
		t.Fatalf("len: got %d, want 2", len(coupons))
	}

	wantCoupons := loadCouponsFixture(t)

	for i := range coupons {
		assertCoupon(t, &coupons[i], &wantCoupons[i])
	}
}

func assertCoupon(t *testing.T, got *Coupon, want *Coupon) {
	t.Helper()

	if diff := cmp.Diff(want, got, couponIgnoreOpts...); diff != "" {
		t.Fatalf("coupon mismatch (-want +got):\n%s", diff)
	}
}

func loadCouponsFixture(t *testing.T) []Coupon {
	t.Helper()

	var coupons []Coupon
	if err := json.Unmarshal(couponsJSON, &coupons); err != nil {
		t.Fatalf("invalid fixture: %v", err)
	}

	return coupons
}
