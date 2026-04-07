package woocommerce

import (
	"context"
	"net/http"
	"testing"
)

func TestCouponsCreate(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/coupons")
		writeJSON(w, &Coupon{ID: 1, Code: "SAVE10"})
	})

	coupon, _, err := client.Coupons.Create(context.Background(), &Coupon{Code: "SAVE10"})
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
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/coupons/1")
		writeJSON(w, &Coupon{ID: 1, Code: "SAVE10", DiscountType: "percent"})
	})

	coupon, _, err := client.Coupons.Get(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coupon.DiscountType != "percent" {
		t.Errorf("DiscountType: got %s, want percent", coupon.DiscountType)
	}
}

func TestCouponsList(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/coupons")
		if r.URL.Query().Get("code") != "SAVE10" {
			t.Errorf("expected code=SAVE10, got %q", r.URL.Query().Get("code"))
		}
		writeJSON(w, &[]Coupon{{ID: 1}, {ID: 2}})
	})

	coupons, _, err := client.Coupons.List(context.Background(), &ListCouponParams{Code: "SAVE10"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*coupons) != 2 {
		t.Errorf("len: got %d, want 2", len(*coupons))
	}
}

func TestCouponsUpdate(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPut)
		assertPathSuffix(t, r, "/coupons/1")
		writeJSON(w, &Coupon{ID: 1, Amount: "20"})
	})

	coupon, _, err := client.Coupons.Update(context.Background(), "1", &Coupon{Amount: "20"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coupon.Amount != "20" {
		t.Errorf("Amount: got %s, want 20", coupon.Amount)
	}
}

func TestCouponsDelete(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodDelete)
		assertPathSuffix(t, r, "/coupons/1")
		writeJSON(w, &Coupon{ID: 1})
	})

	coupon, _, err := client.Coupons.Delete(context.Background(), "1", &DeleteCouponParams{Force: "true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if coupon.ID != 1 {
		t.Errorf("ID: got %d, want 1", coupon.ID)
	}
}

func TestCouponsBatch(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/coupons/batch")
		writeJSON(w, &BatchCouponUpdateResponse{
			Create: &[]Coupon{{ID: 5}},
		})
	})

	result, _, err := client.Coupons.Batch(context.Background(), &BatchCouponUpdate{
		Create: &[]Coupon{{Code: "NEW"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*result.Create) != 1 {
		t.Errorf("create len: got %d, want 1", len(*result.Create))
	}
}

func TestCouponsError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusUnprocessableEntity, "Invalid coupon code")
	})

	_, _, err := client.Coupons.Create(context.Background(), &Coupon{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*ErrorResponse); !ok {
		t.Errorf("expected *ErrorResponse, got %T", err)
	}
}
