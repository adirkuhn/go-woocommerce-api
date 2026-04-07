package woocommerce

import (
	"context"
	"net/http"
)

type TaxRatesServiceInterface interface {
	Create(ctx context.Context, taxRate *TaxRate) (*TaxRate, *http.Response, error)
	Get(ctx context.Context, taxRateID string) (*TaxRate, *http.Response, error)
	List(ctx context.Context, opts *ListTaxRatesParams) ([]TaxRate, *http.Response, error)
	Update(ctx context.Context, taxRateID string, taxRate *TaxRate) (*TaxRate, *http.Response, error)
	Delete(ctx context.Context, taxRateID string, opts *DeleteTaxRateParams) (*TaxRate, *http.Response, error)
	Batch(ctx context.Context, opts *BatchTaxRateUpdate) (*BatchTaxRateUpdateResponse, *http.Response, error)
}

// Tax Rates service
type TaxRatesService service

// TaxRate object. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#tax-rate-properties
type TaxRate struct {
	ID       int    `json:"id,omitempty"`
	Country  string `json:"country,omitempty"`
	State    string `json:"state,omitempty"`
	Postcode string `json:"postcode,omitempty"`
	City     string `json:"city,omitempty"`
	Rate     string `json:"rate,omitempty"`
	Name     string `json:"name,omitempty"`
	Priority int    `json:"priority,omitempty"`
	Compound bool   `json:"compound,omitempty"`
	Shipping bool   `json:"shipping,omitempty"`
	Order    int    `json:"order,omitempty"`
	Class    string `json:"class,omitempty"`
}

type ListTaxRatesParams struct {
	Context string `url:"context,omitempty"`
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Search  string `url:"search,omitempty"`
	Exclude *[]int `url:"exclude,omitempty"`
	Include *[]int `url:"include,omitempty"`
	Offset  int    `url:"offset,omitempty"`
	Order   string `url:"order,omitempty"`
	OrderBy string `url:"orderby,omitempty"`
	Class   string `url:"class,omitempty"`
}

type DeleteTaxRateParams struct {
	Force bool `url:"force"`
}

type BatchTaxRateUpdate struct {
	Create *[]TaxRate `json:"create,omitempty"`
	Update *[]TaxRate `json:"update,omitempty"`
	Delete *[]int     `json:"delete,omitempty"`
}

type BatchTaxRateUpdateResponse struct {
	Create *[]TaxRate `json:"create,omitempty"`
	Update *[]TaxRate `json:"update,omitempty"`
	Delete *[]TaxRate `json:"delete,omitempty"`
}

// Create a tax rate. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#create-a-tax-rate
func (service *TaxRatesService) Create(ctx context.Context, taxRate *TaxRate) (*TaxRate, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "POST", "/taxes", nil, taxRate)
	if err != nil {
		return nil, nil, err
	}

	created := new(TaxRate)
	response, err := service.client.Do(req, created)
	if err != nil {
		return nil, response, err
	}

	return created, response, nil
}

// Get a tax rate. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#retrieve-a-tax-rate
func (service *TaxRatesService) Get(ctx context.Context, taxRateID string) (*TaxRate, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/taxes/"+taxRateID, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	taxRate := new(TaxRate)
	response, err := service.client.Do(req, taxRate)
	if err != nil {
		return nil, response, err
	}

	return taxRate, response, nil
}

// List tax rates. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-tax-rates
func (service *TaxRatesService) List(ctx context.Context, opts *ListTaxRatesParams) ([]TaxRate, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/taxes", opts, nil)
	if err != nil {
		return nil, nil, err
	}

	taxRates := new([]TaxRate)
	response, err := service.client.Do(req, taxRates)
	if err != nil {
		return nil, response, err
	}

	return *taxRates, response, nil
}

// Update a tax rate. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#update-a-tax-rate
func (service *TaxRatesService) Update(ctx context.Context, taxRateID string, taxRate *TaxRate) (*TaxRate, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "PUT", "/taxes/"+taxRateID, nil, taxRate)
	if err != nil {
		return nil, nil, err
	}

	updated := new(TaxRate)
	response, err := service.client.Do(req, updated)
	if err != nil {
		return nil, response, err
	}

	return updated, response, nil
}

// Delete a tax rate. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#delete-a-tax-rate
func (service *TaxRatesService) Delete(ctx context.Context, taxRateID string, opts *DeleteTaxRateParams) (*TaxRate, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "DELETE", "/taxes/"+taxRateID, opts, nil)
	if err != nil {
		return nil, nil, err
	}

	taxRate := new(TaxRate)
	response, err := service.client.Do(req, taxRate)
	if err != nil {
		return nil, response, err
	}

	return taxRate, response, nil
}

// Batch update tax rates. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#batch-update-tax-rates
func (service *TaxRatesService) Batch(ctx context.Context, opts *BatchTaxRateUpdate) (*BatchTaxRateUpdateResponse, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "POST", "/taxes/batch", nil, opts)
	if err != nil {
		return nil, nil, err
	}

	result := new(BatchTaxRateUpdateResponse)
	response, err := service.client.Do(req, result)
	if err != nil {
		return nil, response, err
	}

	return result, response, nil
}
