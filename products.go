package woocommerce

import (
	"context"
	"net/http"
)

type ProductsServiceInterface interface {
	Create(ctx context.Context, product *Product) (*Product, *http.Response, error)
	Get(ctx context.Context, productID string) (*Product, *http.Response, error)
	List(ctx context.Context, opts *ListProductParams) ([]Product, *http.Response, error)
	Update(ctx context.Context, productID string, product *Product) (*Product, *http.Response, error)
	Delete(ctx context.Context, productID string, opts *DeleteProductParams) (*Product, *http.Response, error)
	Batch(ctx context.Context, opts *BatchProductUpdate) (*BatchProductUpdateResponse, *http.Response, error)
}

// Product service
type ProductsService struct {
	client HTTPClient
}

// Product object. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#product-properties
type Product struct {
	ID                int                 `json:"id,omitempty"`
	Name              string              `json:"name,omitempty"`
	Slug              string              `json:"slug,omitempty"`
	Permalink         string              `json:"permalink,omitempty"`
	DateCreated       string              `json:"date_created,omitempty"`
	DateCreatedGmt    string              `json:"date_created_gmt,omitempty"`
	DateModified      string              `json:"date_modified,omitempty"`
	DateModifiedGmt   string              `json:"date_modified_gmt,omitempty"`
	Type              string              `json:"type,omitempty"`
	Status            string              `json:"status,omitempty"`
	Featured          bool                `json:"featured,omitempty"`
	CatalogVisibility string              `json:"catalog_visibility,omitempty"`
	Description       string              `json:"description,omitempty"`
	ShortDescription  string              `json:"short_description,omitempty"`
	Sku               string              `json:"sku,omitempty"`
	Price             string              `json:"price,omitempty"`
	RegularPrice      string              `json:"regular_price,omitempty"`
	SalePrice         string              `json:"sale_price,omitempty"`
	DateOnSaleFrom    string              `json:"date_on_sale_from,omitempty"`
	DateOnSaleFromGmt string              `json:"date_on_sale_from_gmt,omitempty"`
	DateOnSaleTo      string              `json:"date_on_sale_to,omitempty"`
	DateOnSaleToGmt   string              `json:"date_on_sale_to_gmt,omitempty"`
	PriceHtml         string              `json:"price_html,omitempty"`
	OnSale            bool                `json:"on_sale,omitempty"`
	Purchasable       bool                `json:"purchasable,omitempty"`
	TotalSales        int                 `json:"total_sales,omitempty"`
	Virtual           bool                `json:"virtual,omitempty"`
	Downloadable      bool                `json:"downloadable,omitempty"`
	DownloadLimit     int                 `json:"download_limit,omitempty"`
	DownloadExpiry    int                 `json:"download_expiry,omitempty"`
	ExternalUrl       string              `json:"external_url,omitempty"`
	ButtonText        string              `json:"button_text,omitempty"`
	TaxStatus         string              `json:"tax_status,omitempty"`
	TaxClass          string              `json:"tax_class,omitempty"`
	ManageStock       bool                `json:"manage_stock,omitempty"`
	StockQuantity     int                 `json:"stock_quantity,omitempty"`
	StockStatus       string              `json:"stock_status,omitempty"`
	LowStockAmount    *int                `json:"low_stock_amount,omitempty"`
	Backorders        string              `json:"backorders,omitempty"`
	BackordersAllowed bool                `json:"backorders_allowed,omitempty"`
	Backordered       bool                `json:"backordered,omitempty"`
	SoldIndividually  bool                `json:"sold_individually,omitempty"`
	Weight            string              `json:"weight,omitempty"`
	ShippingRequired  bool                `json:"shipping_required,omitempty"`
	ShippingTaxable   bool                `json:"shipping_taxable,omitempty"`
	ShippingClass     string              `json:"shipping_class,omitempty"`
	ShippingClassID   int                 `json:"shipping_class_id,omitempty"`
	ReviewsAllowed    bool                `json:"reviews_allowed,omitempty"`
	AverageRating     string              `json:"average_rating,omitempty"`
	RatingCount       int                 `json:"rating_count,omitempty"`
	ParentID          int                 `json:"parent_id,omitempty"`
	PurchaseNote      string              `json:"purchase_note,omitempty"`
	MenuOrder         int                 `json:"menu_order,omitempty"`
	Variations        []int               `json:"variations,omitempty"`
	GroupedProducts   []int               `json:"grouped_products,omitempty"`
	MetaData          []MetaData          `json:"meta_data,omitempty"`
	RelatedIds        []int               `json:"related_ids,omitempty"`
	CrossSellIds      []int               `json:"cross_sell_ids,omitempty"`
	UpsellIds         []int               `json:"upsell_ids,omitempty"`
	Images            []Image             `json:"images,omitempty"`
	Dimensions        *ProductDimensions  `json:"dimensions,omitempty"`
	Downloads         []ProductDownloads  `json:"downloads,omitempty"`
	Categories        []ProductCategory   `json:"categories,omitempty"`
	Tags              []ProductTag        `json:"tags,omitempty"`
	Attributes        []ProductAttributes `json:"attributes,omitempty"`
	DefaultAttributes []DefaultAttributes `json:"default_attributes,omitempty"`
}

type ProductDownloads struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	File string `json:"file,omitempty"`
}

type ProductDimensions struct {
	Length string `json:"length,omitempty"`
	Width  string `json:"width,omitempty"`
	Height string `json:"height,omitempty"`
}

type ProductCategory struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

type ProductTag struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

type Image struct {
	ID              any    `json:"id,omitempty"`
	DateCreated     string `json:"date_created,omitempty"`
	DateCreatedGmt  string `json:"date_created_gmt,omitempty"`
	DateModified    string `json:"date_modified,omitempty"`
	DateModifiedGmt string `json:"date_modified_gmt,omitempty"`
	Source          string `json:"src,omitempty"`
	Name            string `json:"name,omitempty"`
	Alt             string `json:"alt,omitempty"`
}

type ProductAttributes struct {
	ID        int      `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Position  int      `json:"position,omitempty"`
	Visible   bool     `json:"visible,omitempty"`
	Variation bool     `json:"variation,omitempty"`
	Options   []string `json:"options,omitempty"`
}

type DefaultAttributes struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Option string `json:"option,omitempty"`
}

type ListProductParams struct {
	Context        string   `url:"context,omitempty"`
	Page           int      `url:"page,omitempty"`
	PerPage        int      `url:"per_page,omitempty"`
	Search         string   `url:"search,omitempty"`
	Exclude        []int    `url:"exclude,omitempty"`
	Include        []int    `url:"include,omitempty"`
	Offset         int      `url:"offset,omitempty"`
	Order          string   `url:"order,omitempty"`
	OrderBy        string   `url:"orderby,omitempty"`
	After          string   `url:"after,omitempty"`
	Before         string   `url:"before,omitempty"`
	ModifiedAfter  string   `url:"modified_after,omitempty"`
	ModifiedBefore string   `url:"modified_before,omitempty"`
	DatesAreGmt    bool     `url:"dates_are_gmt,omitempty"`
	Slug           string   `url:"slug,omitempty"`
	Status         string   `url:"status,omitempty"`
	Type           string   `url:"type,omitempty"`
	Sku            string   `url:"sku,omitempty"`
	Featured       bool     `url:"featured,omitempty"`
	Category       string   `url:"category,omitempty"`
	Tag            string   `url:"tag,omitempty"`
	ShippingClass  string   `url:"shipping_class,omitempty"`
	Attribute      string   `url:"attribute,omitempty"`
	AttributeTerm  string   `url:"attribute_term,omitempty"`
	TaxClass       string   `url:"tax_class,omitempty"`
	OnSale         bool     `url:"on_sale,omitempty"`
	MinPrice       string   `url:"min_price,omitempty"`
	MaxPrice       string   `url:"max_price,omitempty"`
	StockStatus    string   `url:"stock_status,omitempty"`
	Parent         []int    `url:"parent,omitempty"`
	ParentExclude  []int    `url:"parent_exclude,omitempty"`
	Fields         []string `url:"_fields,omitempty"`
}

type DeleteProductParams struct {
	Force string `json:"force,omitempty"`
}

type BatchProductUpdate struct {
	Create []Product `json:"create,omitempty"`
	Update []Product `json:"update,omitempty"`
	Delete []int     `json:"delete,omitempty"`
}

type BatchProductUpdateResponse struct {
	Create []Product `json:"create,omitempty"`
	Update []Product `json:"update,omitempty"`
	Delete []Product `json:"delete,omitempty"`
}

// Create a product. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#create-a-product
func (service *ProductsService) Create(ctx context.Context, product *Product) (*Product, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "POST", "/products", nil, product)
	if err != nil {
		return nil, nil, err
	}

	createdProduct := new(Product)
	response, err := service.client.Do(req, createdProduct)
	if err != nil {
		return nil, response, err
	}

	return createdProduct, response, nil
}

// Get a product. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#retrieve-a-product
func (service *ProductsService) Get(ctx context.Context, productID string) (*Product, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/products/"+productID, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	product := new(Product)
	response, err := service.client.Do(req, product)
	if err != nil {
		return nil, response, err
	}

	return product, response, nil
}

// List products. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-products
func (service *ProductsService) List(ctx context.Context, opts *ListProductParams) ([]Product, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/products", opts, nil)
	if err != nil {
		return nil, nil, err
	}

	products := new([]Product)
	response, err := service.client.Do(req, products)
	if err != nil {
		return nil, response, err
	}

	return *products, response, nil
}

// Update a product. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#update-a-product
func (service *ProductsService) Update(ctx context.Context, productID string, product *Product) (*Product, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "PUT", "/products/"+productID, nil, product)
	if err != nil {
		return nil, nil, err
	}

	updatedProduct := new(Product)
	response, err := service.client.Do(req, updatedProduct)
	if err != nil {
		return nil, response, err
	}

	return updatedProduct, response, nil
}

// Delete a product. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#delete-a-product
func (service *ProductsService) Delete(ctx context.Context, productID string, opts *DeleteProductParams) (*Product, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "DELETE", "/products/"+productID, opts, nil)
	if err != nil {
		return nil, nil, err
	}

	product := new(Product)
	response, err := service.client.Do(req, product)
	if err != nil {
		return nil, response, err
	}

	return product, response, nil
}

// Batch update products. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#batch-update-products
func (service *ProductsService) Batch(ctx context.Context, opts *BatchProductUpdate) (*BatchProductUpdateResponse, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "POST", "/products/batch", nil, opts)
	if err != nil {
		return nil, nil, err
	}

	products := new(BatchProductUpdateResponse)
	response, err := service.client.Do(req, products)
	if err != nil {
		return nil, response, err
	}

	return products, response, nil
}
