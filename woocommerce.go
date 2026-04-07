package woocommerce

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultRestEndpointVersion   = "v3"
	defaultHeaderName            = "Authorization"
	acceptedContentType          = "application/json"
	userAgent                    = "go-woocommerce-api/1.1"
	clientRequestRetryAttempts   = 2
	clientRequestRetryHoldMillis = 1000
	clientTimeout                = 10
)

var errorDoAllAttemptsExhausted = errors.New("all request attempts were exhausted")
var errorDoAttemptNilRequest = errors.New("request could not be constructed")

type ClientConfig struct {
	HttpClient          *http.Client
	RestEndpointURL     string
	RestEndpointVersion string
}

type auth struct {
	HeaderName string
	ApiKey     string
}

type Client struct {
	config  *ClientConfig
	client  *http.Client
	auth    *auth
	baseURL *url.URL

	Coupons    CouponsServiceInterface
	Customers  CustomersServiceInterface
	Orders     OrdersServiceInterface
	OrderNotes OrderNotesServiceInterface
	Refunds    RefundsServiceInterface
	Products   ProductsServiceInterface
	TaxRates   TaxRatesServiceInterface
	Webhooks   WebhookServiceInterface
}

type service struct {
	client *Client
}

// ErrorResponse represents an error returned by the WooCommerce API.
type ErrorResponse struct {
	Response *http.Response

	Code    string    `json:"code"`
	Message string    `json:"message"`
	Data    ErrorData `json:"data"`
}

type ErrorData struct {
	Status int `json:"status"`
}

func (response *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		response.Response.Request.Method, response.Response.Request.URL,
		response.Response.StatusCode, response.Message)
}

func New(shopURL string) (*Client, error) {
	if shopURL == "" {
		return nil, errors.New("store url is required")
	}

	return NewWithConfig(ClientConfig{
		RestEndpointURL: shopURL,
	})
}

func NewWithConfig(config ClientConfig) (*Client, error) {
	if config.RestEndpointURL == "" {
		return nil, errors.New("rest endpoint url is required")
	}

	if config.HttpClient == nil {
		config.HttpClient = &http.Client{
			Timeout: time.Duration(clientTimeout * time.Second),
		}
	}

	if config.RestEndpointVersion == "" {
		config.RestEndpointVersion = defaultRestEndpointVersion
	}

	baseURL, err := url.Parse(config.RestEndpointURL + "/wp-json/wc/" + defaultRestEndpointVersion)
	if err != nil {
		return nil, err
	}

	client := &Client{config: &config, client: config.HttpClient, auth: &auth{}, baseURL: baseURL}

	client.Coupons = &CouponsService{client: client}
	client.Customers = &CustomersService{client: client}
	client.Orders = &OrdersService{client: client}
	client.OrderNotes = &OrderNotesService{client: client}
	client.Refunds = &RefundsService{client: client}
	client.Products = &ProductsService{client: client}
	client.TaxRates = &TaxRatesService{client: client}
	client.Webhooks = &WebhookService{client: client}

	return client, nil
}

// Authenticate saves authentication parameters for the client.
func (client *Client) Authenticate(consumer_key string, consumer_secret string) {
	client.auth.HeaderName = defaultHeaderName
	client.auth.ApiKey = "Basic " + base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{consumer_key, consumer_secret}, ":")))
}

// NewRequest creates an API request. ctx must be non-nil.
func (client *Client) NewRequest(ctx context.Context, method, urlStr string, opts any, body any) (*http.Request, error) {
	if opts != nil {
		queryParams, err := query.Values(opts)
		if err != nil {
			return nil, err
		}

		rawQuery := queryParams.Encode()
		if rawQuery != "" {
			urlStr += "?" + rawQuery
		}
	}

	rel, err := url.Parse(client.config.RestEndpointVersion + urlStr)
	if err != nil {
		return nil, err
	}

	u := client.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add(client.auth.HeaderName, client.auth.ApiKey)
	req.Header.Add("Accept", acceptedContentType)
	req.Header.Add("Content-type", acceptedContentType)
	req.Header.Add("User-Agent", userAgent)

	return req, nil
}

// Do sends an API request and decodes the response into v.
func (client *Client) Do(req *http.Request, v any) (*http.Response, error) {
	var lastErr error

	attempts := 0

	for attempts < clientRequestRetryAttempts {
		if attempts > 0 {
			time.Sleep(clientRequestRetryHoldMillis * time.Millisecond)
		}

		attempts++
		resp, shouldRetry, err := client.doAttempt(req, v)

		if !shouldRetry {
			return resp, err
		}

		lastErr = err
	}

	if lastErr == nil {
		lastErr = errorDoAllAttemptsExhausted
	}

	return nil, lastErr
}

func (client *Client) doAttempt(req *http.Request, v any) (*http.Response, bool, error) {
	if req == nil {
		return nil, false, errorDoAttemptNilRequest
	}

	resp, err := client.client.Do(req)

	if checkRequestRetry(resp, err) {
		return nil, true, err
	}

	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return resp, false, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, _ = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil
			}
		}
	}

	return resp, false, err
}

func checkRequestRetry(response *http.Response, err error) bool {
	if err != nil || response.StatusCode >= 500 {
		return true
	}
	return false
}

func checkResponse(response *http.Response) error {
	if code := response.StatusCode; 200 <= code && code <= 299 {
		return nil
	}

	errResp := &ErrorResponse{Response: response}

	data, err := io.ReadAll(response.Body)
	if err == nil && data != nil {
		_ = json.Unmarshal(data, errResp)
	}

	return errResp
}
