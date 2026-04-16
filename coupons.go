package woocommerce

import (
	"context"
	"net/http"
)

type CouponsServiceInterface interface {
	Create(ctx context.Context, coupon *Coupon) (*Coupon, *http.Response, error)
	Get(ctx context.Context, couponID string) (*Coupon, *http.Response, error)
	List(ctx context.Context, opts *ListCouponParams) ([]Coupon, *http.Response, error)
	Update(ctx context.Context, couponID string, coupon *Coupon) (*Coupon, *http.Response, error)
	Delete(ctx context.Context, couponID string, opts *DeleteCouponParams) (*Coupon, *http.Response, error)
	Batch(ctx context.Context, opts *BatchCouponUpdate) (*BatchCouponUpdateResponse, *http.Response, error)
}

// Coupon service
type CouponsService struct {
	client HTTPClient
}

// Coupon object. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#coupon-properties
type Coupon struct {
	ID                        int        `json:"id,omitempty"`
	Code                      string     `json:"code,omitempty"`
	Amount                    string     `json:"amount,omitempty"`
	DateCreated               string     `json:"date_created,omitempty"`
	DateCreatedGmt            string     `json:"date_created_gmt,omitempty"`
	DateModified              string     `json:"date_modified,omitempty"`
	DateModifiedGmt           string     `json:"date_modified_gmt,omitempty"`
	DiscountType              string     `json:"discount_type,omitempty"`
	Description               string     `json:"description,omitempty"`
	DateExpires               string     `json:"date_expires,omitempty"`
	DateExpiresGmt            string     `json:"date_expires_gmt,omitempty"`
	UsageCount                int        `json:"usage_count,omitempty"`
	IndividualUse             bool       `json:"individual_use,omitempty"`
	UsageLimit                int        `json:"usage_limit,omitempty"`
	UsageLimitPerUser         int        `json:"usage_limit_per_user,omitempty"`
	LimitUsageToXItems        int        `json:"limit_usage_to_x_items,omitempty"`
	FreeShipping              bool       `json:"free_shipping,omitempty"`
	ExcludeSaleItems          bool       `json:"exclude_sale_items,omitempty"`
	MinimumAmount             string     `json:"minimum_amount,omitempty"`
	MaximumAmount             string     `json:"maximum_amount,omitempty"`
	EmailRestrictions         any        `json:"email_restrictions,omitempty"`
	UsedBy                    any        `json:"used_by,omitempty"`
	ProductIds                []int      `json:"product_ids,omitempty"`
	ExcludedProductIds        []int      `json:"excluded_product_ids,omitempty"`
	ProductCategories         []int      `json:"product_categories,omitempty"`
	ExcludedProductCategories []int      `json:"excluded_product_categories,omitempty"`
	MetaData                  []MetaData `json:"meta_data,omitempty"`
}

type ListCouponParams struct {
	Context        string `url:"context,omitempty"`
	Page           int    `url:"page,omitempty"`
	PerPage        int    `url:"per_page,omitempty"`
	Search         string `url:"search,omitempty"`
	Exclude        []int  `url:"exclude,omitempty,comma"`
	Include        []int  `url:"include,omitempty,comma"`
	Offset         int    `url:"offset,omitempty"`
	Order          string `url:"order,omitempty"`
	OrderBy        string `url:"orderby,omitempty"`
	After          string `url:"after,omitempty"`
	Before         string `url:"before,omitempty"`
	ModifiedAfter  string `url:"modified_after,omitempty"`
	ModifiedBefore string `url:"modified_before,omitempty"`
	DatesAreGmt    bool   `url:"dates_are_gmt,omitempty"`
	Code           string `url:"code,omitempty"`
}

type DeleteCouponParams struct {
	Force string `json:"force,omitempty"`
}

type BatchCouponUpdate struct {
	Create []Coupon `json:"create,omitempty"`
	Update []Coupon `json:"update,omitempty"`
	Delete []int    `json:"delete,omitempty"`
}

type BatchCouponUpdateResponse struct {
	Create []Coupon `json:"create,omitempty"`
	Update []Coupon `json:"update,omitempty"`
	Delete []Coupon `json:"delete,omitempty"`
}

// Create a coupon. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#create-a-coupon
func (service *CouponsService) Create(ctx context.Context, coupon *Coupon) (*Coupon, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "POST", "/coupons", nil, coupon)
	if err != nil {
		return nil, nil, err
	}

	createdCoupon := new(Coupon)
	response, err := service.client.Do(req, createdCoupon)
	if err != nil {
		return nil, response, err
	}

	return createdCoupon, response, nil
}

// Get a coupon. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#retrieve-a-coupon
func (service *CouponsService) Get(ctx context.Context, couponID string) (*Coupon, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/coupons/"+couponID, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	coupon := new(Coupon)
	response, err := service.client.Do(req, coupon)
	if err != nil {
		return nil, response, err
	}

	return coupon, response, nil
}

// List coupons. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-coupons
func (service *CouponsService) List(ctx context.Context, opts *ListCouponParams) ([]Coupon, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/coupons", opts, nil)
	if err != nil {
		return nil, nil, err
	}

	coupons := new([]Coupon)
	response, err := service.client.Do(req, coupons)
	if err != nil {
		return nil, response, err
	}

	return *coupons, response, nil
}

// Update a coupon. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#update-a-coupon
func (service *CouponsService) Update(ctx context.Context, couponID string, coupon *Coupon) (*Coupon, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "PUT", "/coupons/"+couponID, nil, coupon)
	if err != nil {
		return nil, nil, err
	}

	updatedCoupon := new(Coupon)
	response, err := service.client.Do(req, updatedCoupon)
	if err != nil {
		return nil, response, err
	}

	return updatedCoupon, response, nil
}

// Delete a coupon. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#delete-a-coupon
func (service *CouponsService) Delete(ctx context.Context, couponID string, opts *DeleteCouponParams) (*Coupon, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "DELETE", "/coupons/"+couponID, opts, nil)
	if err != nil {
		return nil, nil, err
	}

	coupon := new(Coupon)
	response, err := service.client.Do(req, coupon)
	if err != nil {
		return nil, response, err
	}

	return coupon, response, nil
}

// Batch update coupons. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#batch-update-coupons
func (service *CouponsService) Batch(ctx context.Context, opts *BatchCouponUpdate) (*BatchCouponUpdateResponse, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "POST", "/coupons/batch", nil, opts)
	if err != nil {
		return nil, nil, err
	}

	coupons := new(BatchCouponUpdateResponse)
	response, err := service.client.Do(req, coupons)
	if err != nil {
		return nil, response, err
	}

	return coupons, response, nil
}
