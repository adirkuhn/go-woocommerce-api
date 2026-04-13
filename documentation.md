# WooCommerce API Client Documentation

This repository contains a Go implementation for interacting with a WooCommerce API, providing services for managing products, orders, customers, coupons, refunds, tax rates, and webhooks.

## 1. Overview

The project implements a comprehensive client (`woocommerce.go`) that abstracts the HTTP communication and business logic for various WooCommerce endpoints. It relies on a custom `HTTPClient` interface for all external communication.

## 2. Core Components

### 2.1 HTTP Client (`client.go`)

The `httpClient` struct and its associated methods handle all external HTTP requests, including retry logic and error checking.

*   **`HTTPClient` Interface:** Defines methods for making requests (`NewRequest`, `Do`) and handles retry mechanisms.
*   **`Config`:** Stores configuration details like `ShopURL`, `ConsumerKey`, `ConsumerSecret`, and default timeouts.
*   **Error Handling:** The client includes logic to check HTTP status codes and parse error responses (`ErrorResponse`).

### 2.2 Service Interfaces and Implementations

The system is structured around several service interfaces, each handling a specific domain of the WooCommerce API. All services utilize the shared `HTTPClient`.

#### 2.2.1 Products Service (`products.go`)
Handles CRUD operations and retrieval of product data.
*   **`ProductsServiceInterface`**: Defines methods for creating, getting, listing, updating, deleting, and batch operations for products.
*   **`ProductsService`**: Implements the interface, managing product data retrieval and manipulation.
*   **Key Data Models:** `Product`, `ProductDownloads`, `ProductDimensions`, `ProductTag`, `ProductAttributes`, `DefaultAttributes`.

#### 2.2.2 Orders Service (`orders.go`)
Manages order creation, retrieval, updates, and refunds.
*   **`OrdersServiceInterface`**: Defines methods for order management.
*   **`OrdersService`**: Implements the order service logic.
*   **Key Data Models:** `Order`, `OrderRefund`, `CouponLine`, `LineItems`, `TaxLines`, `ShippingLines`, `ListOrdersParams`, `GetOrderParams`, `BatchOrderUpdate`.

#### 2.2.3 Customers Service (`customers.go`)
Handles customer data management.
*   **`CustomersServiceInterface`**: Defines methods for customer management.
*   **`CustomersService`**: Implements the customer service logic.
*   **Key Data Models:** `Customer`, `ListCustomerParams`, `DeleteCustomerParams`, `BatchCustomerUpdate`.

#### 2.2.4 Coupons Service (`coupons.go`)
Manages coupon data.
*   **`CouponsServiceInterface`**: Defines methods for coupon management (Create, Get, List, Update, Delete, Batch).
*   **`CouponsService`**: Implements the coupon service logic.
*   **Key Data Models:** `Coupon`, `ListCouponParams`, `DeleteCouponParams`, `BatchCouponUpdate`.

#### 2.2.5 Order Notes Service (`order_notes.go`)
Manages order notes associated with orders.
*   **`OrderNotesServiceInterface`**: Defines methods for creating, getting, listing, and deleting order notes.
*   **`OrderNotesService`**: Implements the order notes service logic.
*   **Key Data Models:** `OrderNote`, `ListOrderNotesParams`, `DeleteOrderNoteParams`.

#### 2.2.6 Refunds Service (`refunds.go`)
Manages refund operations.
*   **`RefundsServiceInterface`**: Defines methods for creating, getting, listing, and deleting refunds.
*   **`RefundsService`**: Implements the refund service logic.
*   **Key Data Models:** `Refund`, `RefundLineItem`, `RefundTax`, `ListRefundParams`, `DeleteRefundParams`.

#### 2.2.7 Tax Rates Service (`tax_rates.go`)
Manages tax rate configurations.
*   **`TaxRatesServiceInterface`**: Defines methods for managing tax rates.
*   **`TaxRatesService`**: Implements the tax rate service logic.
*   **Key Data Models:** `TaxRate`, `ListTaxRatesParams`, `DeleteTaxRateParams`, `BatchTaxRateUpdate`.

#### 2.2.8 Webhooks Service (`webhooks.go`)
Manages webhook configurations.
*   **`WebhookServiceInterface`**: Defines methods for webhook management.
*   **`WebhookService`**: Implements the webhook service logic.
*   **Key Data Models:** `Webhook`, `ListWebhooksParams`, `DeleteWebhookParams`, `BatchWebhookUpdate`.

## 3. Summary of Data Models

The system utilizes various structs to represent data fetched from or sent to the WooCommerce API, including:

*   **Customer Data:** `Customer`, `Billing`, `Shipping` (from `global_types.go`).
*   **Product Data:** `Product`, `ProductAttributes`, `ProductTag`, `ProductDimensions` (from `products.go`).
*   **Order Data:** `Order`, `LineItems`, `TaxLines`, `ShippingLines`, `OrderRefund` (from `orders.go`).
*   **Coupon Data:** `Coupon` (from `coupons.go`).
*   **Refund Data:** `Refund`, `RefundLineItem`, `RefundTax` (from `refunds.go`).
*   **Tax Data:** `TaxRate` (from `tax_rates.go`).
*   **Webhook Data:** `Webhook` (from `webhooks.go`).
