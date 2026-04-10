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
	defaultVersion               = "v3"
	defaultHeaderName            = "Authorization"
	acceptedContentType          = "application/json"
	userAgent                    = "go-woocommerce-api/1.1"
	clientRequestRetryAttempts   = 2
	clientRequestRetryHoldMillis = 1000
	defaultClientTimeout         = 10 * time.Second
)

var (
	errAllAttemptsExhausted = errors.New("all request attempts were exhausted")
	errNilRequest           = errors.New("request could not be constructed")
)

// Config holds all configuration needed to create a WooCommerce client.
type Config struct {
	ShopURL        string
	ConsumerKey    string
	ConsumerSecret string
	// Version defaults to "v3" if empty.
	Version string
	// HTTPClient defaults to a client with a 10s timeout if nil.
	HTTPClient *http.Client
}

// HTTPClient is the interface that service implementations depend on.
// It can be replaced with a mock in tests.
type HTTPClient interface {
	NewRequest(ctx context.Context, method, path string, opts, body any) (*http.Request, error)
	Do(req *http.Request, v any) (*http.Response, error)
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

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		e.Response.Request.Method, e.Response.Request.URL,
		e.Response.StatusCode, e.Message)
}

// httpClient is the concrete implementation of HTTPClient.
type httpClient struct {
	baseURL    *url.URL
	authHeader string
	version    string
	http       *http.Client
}

func newHTTPClient(cfg Config) (*httpClient, error) {
	if cfg.ShopURL == "" {
		return nil, errors.New("ShopURL is required")
	}
	if cfg.ConsumerKey == "" || cfg.ConsumerSecret == "" {
		return nil, errors.New("ConsumerKey and ConsumerSecret are required")
	}

	version := cfg.Version
	if version == "" {
		version = defaultVersion
	}

	httpC := cfg.HTTPClient
	if httpC == nil {
		httpC = &http.Client{Timeout: defaultClientTimeout}
	}

	baseURL, err := url.Parse(cfg.ShopURL + "/wp-json/wc/" + version + "/")
	if err != nil {
		return nil, fmt.Errorf("invalid ShopURL: %w", err)
	}

	token := base64.StdEncoding.EncodeToString(
		[]byte(strings.Join([]string{cfg.ConsumerKey, cfg.ConsumerSecret}, ":")),
	)

	return &httpClient{
		baseURL:    baseURL,
		authHeader: "Basic " + token,
		version:    version,
		http:       httpC,
	}, nil
}

// NewRequest builds an *http.Request for the given method and path.
// opts is encoded as query parameters; body is JSON-encoded as the request body.
func (c *httpClient) NewRequest(ctx context.Context, method, path string, opts, body any) (*http.Request, error) {
	// Strip leading slash — baseURL already has a trailing slash.
	path = strings.TrimLeft(path, "/")

	if opts != nil {
		queryParams, err := query.Values(opts)
		if err != nil {
			return nil, err
		}
		if raw := queryParams.Encode(); raw != "" {
			path += "?" + raw
		}
	}

	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

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

	req.Header.Set(defaultHeaderName, c.authHeader)
	req.Header.Set("Accept", acceptedContentType)
	req.Header.Set("Content-Type", acceptedContentType)
	req.Header.Set("User-Agent", userAgent)

	return req, nil
}

// Do executes the request, retrying on transient errors, and decodes the
// response body into v.
func (c *httpClient) Do(req *http.Request, v any) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt < clientRequestRetryAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(clientRequestRetryHoldMillis * time.Millisecond)
		}

		resp, shouldRetry, err := c.doAttempt(req, v)
		if !shouldRetry {
			return resp, err
		}

		lastErr = err
	}

	if lastErr == nil {
		lastErr = errAllAttemptsExhausted
	}

	return nil, lastErr
}

func (c *httpClient) doAttempt(req *http.Request, v any) (*http.Response, bool, error) {
	if req == nil {
		return nil, false, errNilRequest
	}

	resp, err := c.http.Do(req)
	if err != nil || resp.StatusCode >= 500 {
		return nil, true, err
	}

	defer resp.Body.Close()

	if err := checkResponse(resp); err != nil {
		return resp, false, err
	}

	if v != nil {
		switch w := v.(type) {
		case io.Writer:
			_, _ = io.Copy(w, resp.Body)
		default:
			if err := json.NewDecoder(resp.Body).Decode(v); err != nil && err != io.EOF {
				return resp, false, err
			}
		}
	}

	return resp, false, nil
}

func checkResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}

	errResp := &ErrorResponse{Response: resp}

	if data, err := io.ReadAll(resp.Body); err == nil {
		_ = json.Unmarshal(data, errResp)
	}

	return errResp
}
