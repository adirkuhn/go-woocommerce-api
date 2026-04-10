package woocommerce

import (
	"context"
	"net/http"
	"testing"
)

func TestNew_MissingShopURL(t *testing.T) {
	_, err := New(Config{ConsumerKey: "ck", ConsumerSecret: "cs"})
	if err == nil {
		t.Fatal("expected error for missing ShopURL")
	}
}

func TestNew_MissingCredentials(t *testing.T) {
	_, err := New(Config{ShopURL: "http://example.com"})
	if err == nil {
		t.Fatal("expected error for missing credentials")
	}
}

func TestNew_WiresAllServices(t *testing.T) {
	woo, err := New(Config{
		ShopURL:        "http://example.com",
		ConsumerKey:    "ck_test",
		ConsumerSecret: "cs_test",
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if woo.Products == nil {
		t.Error("Products is nil")
	}
	if woo.Orders == nil {
		t.Error("Orders is nil")
	}
	if woo.Customers == nil {
		t.Error("Customers is nil")
	}
	if woo.Coupons == nil {
		t.Error("Coupons is nil")
	}
	if woo.OrderNotes == nil {
		t.Error("OrderNotes is nil")
	}
	if woo.Refunds == nil {
		t.Error("Refunds is nil")
	}
	if woo.TaxRates == nil {
		t.Error("TaxRates is nil")
	}
	if woo.Webhooks == nil {
		t.Error("Webhooks is nil")
	}
}

func TestNewWithHTTPClient_WiresAllServices(t *testing.T) {
	woo := NewWithHTTPClient(&noopHTTPClient{})

	if woo.Products == nil {
		t.Error("Products is nil")
	}
	if woo.Orders == nil {
		t.Error("Orders is nil")
	}
	if woo.Customers == nil {
		t.Error("Customers is nil")
	}
	if woo.Coupons == nil {
		t.Error("Coupons is nil")
	}
	if woo.OrderNotes == nil {
		t.Error("OrderNotes is nil")
	}
	if woo.Refunds == nil {
		t.Error("Refunds is nil")
	}
	if woo.TaxRates == nil {
		t.Error("TaxRates is nil")
	}
	if woo.Webhooks == nil {
		t.Error("Webhooks is nil")
	}
}

// noopHTTPClient satisfies HTTPClient without doing anything.
// Used when we only want to test wiring, not actual requests.
type noopHTTPClient struct{}

func (n *noopHTTPClient) NewRequest(_ context.Context, _, _ string, _, _ any) (*http.Request, error) {
	return nil, nil
}

func (n *noopHTTPClient) Do(_ *http.Request, _ any) (*http.Response, error) {
	return nil, nil
}
