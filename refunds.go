package woocommerce

import (
	"context"
	"net/http"
)

type RefundsServiceInterface interface {
	Create(ctx context.Context, orderID string, refund *Refund) (*Refund, *http.Response, error)
	Get(ctx context.Context, orderID string, refundID string) (*Refund, *http.Response, error)
	List(ctx context.Context, orderID string, opts *ListRefundParams) (*[]Refund, *http.Response, error)
	Delete(ctx context.Context, orderID string, refundID string, opts *DeleteRefundParams) (*Refund, *http.Response, error)
}

// Refunds service
type RefundsService service

// Refund object. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#order-refund-properties
type Refund struct {
	ID              int               `json:"id,omitempty"`
	DateCreated     string            `json:"date_created,omitempty"`
	DateCreatedGmt  string            `json:"date_created_gmt,omitempty"`
	Amount          string            `json:"amount,omitempty"`
	Reason          string            `json:"reason,omitempty"`
	RefundedBy      int               `json:"refunded_by,omitempty"`
	RefundedPayment bool              `json:"refunded_payment,omitempty"`
	ApiRefund       bool              `json:"api_refund,omitempty"`
	MetaData        *[]MetaData       `json:"meta_data,omitempty"`
	LineItems       *[]RefundLineItem `json:"line_items,omitempty"`
}

type RefundLineItem struct {
	ID          int          `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	ProductID   int          `json:"product_id,omitempty"`
	VariationID int          `json:"variation_id,omitempty"`
	Quantity    int          `json:"quantity,omitempty"`
	TaxClass    int          `json:"tax_class,omitempty"`
	Subtotal    string       `json:"subtotal,omitempty"`
	SubtotalTax string       `json:"subtotal_tax,omitempty"`
	Total       string       `json:"total,omitempty"`
	TotalTax    string       `json:"total_tax,omitempty"`
	Sku         string       `json:"sku,omitempty"`
	Price       string       `json:"price,omitempty"`
	RefundTotal float64      `json:"refund_total,omitempty"`
	Taxes       *[]RefundTax `json:"taxes,omitempty"`
	MetaData    *[]MetaData  `json:"meta_data,omitempty"`
}

type RefundTax struct {
	ID          int     `json:"id,omitempty"`
	Total       string  `json:"total,omitempty"`
	Subtotal    string  `json:"subtotal,omitempty"`
	RefundTotal float64 `json:"refund_total,omitempty"`
}

type ListRefundParams struct {
	Context       string `url:"context,omitempty"`
	Page          int    `url:"page,omitempty"`
	PerPage       int    `url:"per_page,omitempty"`
	Search        string `url:"search,omitempty"`
	Exclude       *[]int `url:"exclude,omitempty"`
	Include       *[]int `url:"include,omitempty"`
	Offset        int    `url:"offset,omitempty"`
	Order         string `url:"order,omitempty"`
	OrderBy       string `url:"orderby,omitempty"`
	After         string `url:"after,omitempty"`
	Before        string `url:"before,omitempty"`
	Parent        *[]int `url:"parent,omitempty"`
	ParentExclude *[]int `url:"parent_exclude,omitempty"`
	Dp            int    `url:"dp,omitempty"`
}

type DeleteRefundParams struct {
	Force bool `url:"force"`
}

// Create a refund. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#create-a-refund
func (service *RefundsService) Create(ctx context.Context, orderID string, refund *Refund) (*Refund, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "POST", "/orders/"+orderID+"/refunds", nil, refund)
	if err != nil {
		return nil, nil, err
	}

	createdRefund := new(Refund)
	response, err := service.client.Do(req, createdRefund)
	if err != nil {
		return nil, response, err
	}

	return createdRefund, response, nil
}

// Get a refund. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#retrieve-a-refund
func (service *RefundsService) Get(ctx context.Context, orderID string, refundID string) (*Refund, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/orders/"+orderID+"/refunds/"+refundID, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	refund := new(Refund)
	response, err := service.client.Do(req, refund)
	if err != nil {
		return nil, response, err
	}

	return refund, response, nil
}

// List refunds. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-refunds
func (service *RefundsService) List(ctx context.Context, orderID string, opts *ListRefundParams) (*[]Refund, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/orders/"+orderID+"/refunds", opts, nil)
	if err != nil {
		return nil, nil, err
	}

	refunds := new([]Refund)
	response, err := service.client.Do(req, refunds)
	if err != nil {
		return nil, response, err
	}

	return refunds, response, nil
}

// Delete a refund. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#delete-a-refund
func (service *RefundsService) Delete(ctx context.Context, orderID string, refundID string, opts *DeleteRefundParams) (*Refund, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "DELETE", "/orders/"+orderID+"/refunds/"+refundID, opts, nil)
	if err != nil {
		return nil, nil, err
	}

	refund := new(Refund)
	response, err := service.client.Do(req, refund)
	if err != nil {
		return nil, response, err
	}

	return refund, response, nil
}
