# go-woocommerce-api

A Go client library for the [WooCommerce REST API v3](https://woocommerce.github.io/woocommerce-rest-api-docs/).

## Install

```console
go get github.com/adirkuhn/go-woocommerce-api
```

Requires Go 1.18+.

## Usage

Create a client by passing a `Config` struct to `New`. The `ConsumerKey` and `ConsumerSecret` are your WooCommerce REST API credentials — see the [WooCommerce authentication docs](https://woocommerce.github.io/woocommerce-rest-api-docs/#authentication) for how to generate them.

```go
import (
    "context"

    woocommerce "github.com/adirkuhn/go-woocommerce-api"
)

func main() {
    woo, err := woocommerce.New(woocommerce.Config{
        ShopURL:        "https://example.com",
        ConsumerKey:    "ck_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        ConsumerSecret: "cs_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
    })
    if err != nil {
        // handle error
        return
    }

    ctx := context.Background()

    // List orders for a specific customer
    orders, resp, err := woo.Orders.List(ctx, &woocommerce.ListOrdersParams{
        Customer: 3,
        Page:     1,
    })
    if err != nil {
        // handle error
        return
    }

    // Pagination is exposed via response headers
    totalPages := resp.Header.Get("X-Wp-Totalpages")
    totalItems := resp.Header.Get("X-Wp-Total")
    _ = orders
    _ = totalPages
    _ = totalItems
}
```

### Config options

| Field            | Type            | Default | Description                                  |
|------------------|-----------------|---------|----------------------------------------------|
| `ShopURL`        | `string`        | —       | Required. Base URL of the WooCommerce store. |
| `ConsumerKey`    | `string`        | —       | Required. WooCommerce REST API consumer key. |
| `ConsumerSecret` | `string`        | —       | Required. WooCommerce REST API consumer secret. |
| `Version`        | `string`        | `"v3"`  | API version path segment.                    |
| `HTTPClient`     | `*http.Client`  | 10s timeout | Custom HTTP client.                     |

## Supported services

All services are exposed as fields on the `*WooCommerce` value returned by `New`.

| Service      | Methods                                             |
|--------------|-----------------------------------------------------|
| `Orders`     | `Create`, `Get`, `List`, `Update`, `Delete`, `Batch` |
| `OrderNotes` | `Create`, `Get`, `List`, `Delete`                   |
| `Refunds`    | `Create`, `Get`, `List`, `Delete`                   |
| `Customers`  | `Create`, `Get`, `List`, `Update`, `Delete`, `Batch`, `GetDownloads` |
| `Products`   | `Create`, `Get`, `List`, `Update`, `Delete`, `Batch` |
| `Coupons`    | `Create`, `Get`, `List`, `Update`, `Delete`, `Batch` |
| `TaxRates`   | `Create`, `Get`, `List`, `Update`, `Delete`, `Batch` |
| `Webhooks`   | `Create`, `Get`, `List`, `Update`, `Delete`, `Batch` |

Every method returns `(Resource, *http.Response, error)`. The raw `*http.Response` is always passed back so callers can read WooCommerce pagination headers (`X-Wp-Total`, `X-Wp-Totalpages`).

## Examples

### Create a product

```go
product, _, err := woo.Products.Create(ctx, &woocommerce.Product{
    Name:          "My Product",
    Type:          "simple",
    RegularPrice:  "9.99",
    Status:        "publish",
})
```

### Get a customer

```go
customer, _, err := woo.Customers.Get(ctx, "42")
```

### Batch update orders

```go
result, _, err := woo.Orders.Batch(ctx, &woocommerce.BatchOrderUpdate{
    Update: &[]woocommerce.Order{
        {ID: 101, Status: "completed"},
        {ID: 102, Status: "completed"},
    },
})
```

## Testing

Each service is backed by an interface (`OrdersServiceInterface`, `ProductsServiceInterface`, etc.), making it straightforward to substitute a mock in your own tests.

For injecting a custom transport without going through `New`, use `NewWithHTTPClient`:

```go
woo := woocommerce.NewWithHTTPClient(myMockHTTPClient)
```

### Running the library's own tests

```console
go test ./...
```

## Error handling

API errors are returned as `*ErrorResponse`, which carries the HTTP status code, WooCommerce error code, and message:

```go
orders, _, err := woo.Orders.List(ctx, nil)
if err != nil {
    if apiErr, ok := err.(*woocommerce.ErrorResponse); ok {
        fmt.Println(apiErr.Code, apiErr.Message, apiErr.Data.Status)
    }
}
```

The client automatically retries once on 5xx responses and transient transport errors before returning an error to the caller.

## License

See [LICENSE](LICENSE).
