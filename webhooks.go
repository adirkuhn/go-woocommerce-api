package woocommerce

import (
	"context"
	"net/http"
)

// Webhooks service
type WebhookService service

type Webhook struct {
	ID              int      `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Status          string   `json:"status,omitempty"`
	Topic           string   `json:"topic,omitempty"`
	Resource        string   `json:"resource,omitempty"`
	Event           string   `json:"event,omitempty"`
	Hooks           []string `json:"hooks,omitempty"`
	DeliveryURL     string   `json:"delivery_url,omitempty"`
	Secret          string   `json:"secret,omitempty"`
	DateCreated     string   `json:"date_created,omitempty"`
	DateCreatedGmt  string   `json:"date_created_gmt,omitempty"`
	DateModified    string   `json:"date_modified,omitempty"`
	DateModifiedGmt string   `json:"date_modified_gmt,omitempty"`
	Links           *Links   `json:"links,omitempty"`
}

type ListWebhooksParams struct {
	Context string `url:"context,omitempty"`
	Page    int    `url:"page,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
	Search  string `url:"search,omitempty"`
	Exclude *[]int `url:"exclude,omitempty"`
	Include *[]int `url:"include,omitempty"`
	Offset  int    `url:"offset,omitempty"`
	Order   string `url:"order,omitempty"`
	OrderBy string `url:"orderby,omitempty"`
	After   string `url:"after,omitempty"`
	Before  string `url:"before,omitempty"`
	Status  string `url:"status,omitempty"`
}

type DeleteWebhookParams struct {
	Force string `json:"force,omitempty"`
}

type BatchWebhookUpdate struct {
	Create *[]Webhook `json:"create,omitempty"`
	Delete *[]int     `json:"delete,omitempty"`
}

type BatchWebhookUpdateResponse struct {
	Create *[]Webhook `json:"create,omitempty"`
	Delete *[]Webhook `json:"delete,omitempty"`
}

// Create a webhook. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#create-a-webhook
func (service *WebhookService) Create(ctx context.Context, webhook *Webhook) (*Webhook, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "POST", "/webhooks", nil, webhook)
	if err != nil {
		return nil, nil, err
	}

	createdWebhook := new(Webhook)
	response, err := service.client.Do(req, createdWebhook)
	if err != nil {
		return nil, response, err
	}

	return createdWebhook, response, nil
}

// Get a webhook. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#retrieve-a-webhook
func (service *WebhookService) Get(ctx context.Context, webhookID string) (*Webhook, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/webhooks/"+webhookID, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	webhook := new(Webhook)
	response, err := service.client.Do(req, webhook)
	if err != nil {
		return nil, response, err
	}

	return webhook, response, nil
}

// List webhooks. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-webhooks
func (service *WebhookService) List(ctx context.Context, opts *ListWebhooksParams) (*[]Webhook, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/webhooks", opts, nil)
	if err != nil {
		return nil, nil, err
	}

	webhooks := new([]Webhook)
	response, err := service.client.Do(req, webhooks)
	if err != nil {
		return nil, response, err
	}

	return webhooks, response, nil
}

// Update a webhook. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#update-a-webhook
func (service *WebhookService) Update(ctx context.Context, webhookID string, webhook *Webhook) (*Webhook, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "PUT", "/webhooks/"+webhookID, nil, webhook)
	if err != nil {
		return nil, nil, err
	}

	updatedWebhook := new(Webhook)
	response, err := service.client.Do(req, updatedWebhook)
	if err != nil {
		return nil, response, err
	}

	return updatedWebhook, response, nil
}

// Delete a webhook. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#delete-a-webhook
func (service *WebhookService) Delete(ctx context.Context, webhookID string, opts *DeleteWebhookParams) (*Webhook, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "DELETE", "/webhooks/"+webhookID, opts, nil)
	if err != nil {
		return nil, nil, err
	}

	webhook := new(Webhook)
	response, err := service.client.Do(req, webhook)
	if err != nil {
		return nil, response, err
	}

	return webhook, response, nil
}

// Batch update webhooks. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#batch-update-webhooks
func (service *WebhookService) Batch(ctx context.Context, opts *BatchWebhookUpdate) (*BatchWebhookUpdateResponse, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "POST", "/webhooks/batch", nil, opts)
	if err != nil {
		return nil, nil, err
	}

	webhooks := new(BatchWebhookUpdateResponse)
	response, err := service.client.Do(req, webhooks)
	if err != nil {
		return nil, response, err
	}

	return webhooks, response, nil
}
