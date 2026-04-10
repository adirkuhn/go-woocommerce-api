package woocommerce

// WooCommerce is the top-level client callers interact with.
// All fields are interfaces, so they can be swapped for mocks in tests.
//
// Usage:
//
//	woo, err := woocommerce.New(woocommerce.Config{...})
//	products, _, err := woo.Products.List(ctx, nil)
type WooCommerce struct {
	Products   ProductsServiceInterface
	Orders     OrdersServiceInterface
	Customers  CustomersServiceInterface
	Coupons    CouponsServiceInterface
	OrderNotes OrderNotesServiceInterface
	Refunds    RefundsServiceInterface
	TaxRates   TaxRatesServiceInterface
	Webhooks   WebhookServiceInterface
}

// New creates a fully wired WooCommerce client from the given config.
func New(cfg Config) (*WooCommerce, error) {
	hc, err := newHTTPClient(cfg)
	if err != nil {
		return nil, err
	}

	return NewWithHTTPClient(hc), nil
}

// NewWithHTTPClient wires up services against any HTTPClient implementation.
// Use this in tests to inject a mock transport without going through New().
//
//	woo := woocommerce.NewWithHTTPClient(&myMockHTTPClient{})
func NewWithHTTPClient(hc HTTPClient) *WooCommerce {
	return &WooCommerce{
		Products:   &ProductsService{client: hc},
		Orders:     &OrdersService{client: hc},
		Customers:  &CustomersService{client: hc},
		Coupons:    &CouponsService{client: hc},
		OrderNotes: &OrderNotesService{client: hc},
		Refunds:    &RefundsService{client: hc},
		TaxRates:   &TaxRatesService{client: hc},
		Webhooks:   &WebhookService{client: hc},
	}
}
