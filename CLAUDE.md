# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
# Run all tests
go test ./...

# Run a single test
go test -run TestName ./...

# Build
go build ./...

# Vet
go vet ./...
```

## Architecture

This is a single-package Go library (`package woocommerce`) that wraps the WooCommerce REST API v3.

**Core wiring (`woocommerce.go`):**
- `Client` is the root struct. It holds the HTTP client, base URL, auth, and a pointer to each service.
- `New(shopURL)` / `NewWithConfig(ClientConfig)` construct the client and wire all services.
- `Authenticate(key, secret)` encodes Basic auth from WooCommerce consumer key/secret.
- `NewRequest` builds HTTP requests, serializing query params via `go-querystring` and bodies via `encoding/json`.
- `Do` executes requests with retry logic (2 attempts, 1s hold, retries on 5xx or transport errors).

**Service pattern:**
Each resource (Orders, Customers, Products, etc.) is its own file and type alias of the shared `service` struct:
```go
type OrdersService service  // embeds *Client via service.client
```
Every service exposes methods matching the WooCommerce API: `Create`, `Get`, `List`, `Update`, `Delete`, `Batch`. Methods return `(ResourceType, *http.Response, error)` — the raw response is always passed back so callers can read pagination headers (`X-Wp-Total`, `X-Wp-Totalpages`).

**Types:**
- Resource structs (e.g. `Order`, `Product`) and their list/batch param structs live in the same file as their service.
- Shared types used across resources (`Billing`, `Shipping`, `MetaData`, `Links`, `Image`) live in `global_types.go`.
- Query param structs use `url:` struct tags (consumed by `go-querystring`); JSON structs use `json:` tags.

**Adding a new resource:** create a new file, define `type XxxService service`, declare the resource struct with `json` tags, declare `ListXxxParams` with `url` tags, implement CRUD methods following the existing pattern, and register `client.Xxx = &XxxService{client: client}` in `NewWithConfig`.
