# go-woocommerce-api
A Woocommerce API Golang Wrapper for the [Woocommerce Rest API (v3)](https://woocommerce.github.io/woocommerce-rest-api-docs/)

# Install

```console
$ go get github.com/adirkuhn/go-woocommerce-api
```

# Tests

```console
$ go test ./...
```

# Usage

Create a new API client and authenticate with your REST API keys. You should specify the URL protocol by prefixing your domain name with `https://` or `http://`. Follow the [Woocommerce documentation](https://woocommerce.github.io/woocommerce-rest-api-docs/#authentication) to create your API keys through the [Wordpress admin interface](https://woocommerce.github.io/woocommerce-rest-api-docs/#generating-api-keys-in-the-wordpress-admin-interface) or using the [authentication endpoint](https://woocommerce.github.io/woocommerce-rest-api-docs/#creating-an-authentication-endpoint-url). 

```go
import (
  "github.com/adirkuhn/go-woocommerce-api"
)

func main(){
  shopURL := "https://example.com"
  key := "ck_xxxxxxx"
  secret := "cs_xxxxxxx"

  client, err := woocommerce.New(shopURL)

  if err != nil {
    // handle error

    return
  }

  client.Authenticate(key, secret)
}

```

The API routes are broken down into services, the supported services are: 
* Coupons `(Create, Get, List, Update, Delete, Batch)`
* Customers `(Create, Get, List, Update, Delete, Batch, GetDownloads)`
* Orders `(Create, Get, List, Update, Delete, Batch)`
* OrderNotes `(Create, Get, List, Delete)`
* Refunds `(Create, Get, List, Delete)`
* Products `(Create, Get, List, Update, Delete, Batch)`
* Webhooks `(Create, Get, List, Update, Delete, Batch)`
* TaxRates `(Create, Get, List, Update, Delete, Batch)`

List Orders by customer ID and page number.

```go
func (client *woocommerce.Client) listOrders() {
  customerID := 3
  pageNumber := 1
  
  opts := woocommerce.ListOrdersParams{
    Customer: customerID,
    Page: pageNumber,
  }

  ctx := context.Background()
  orders, resp, err := client.Orders.List(ctx, opts)

  if err != nil {
    // Handle errors

    return
  }

  // Pagination headers.
  totalPages := resp.Header.Get("X-Wp-Totalpages")
  totalItems := resp.Header.Get("X-Wp-Total")

  // ....

}

```

