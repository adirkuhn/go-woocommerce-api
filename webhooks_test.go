package woocommerce

import (
	"context"
	"net/http"
	"testing"
)

func TestWebhooksCreate(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/webhooks")
		writeJSON(w, &Webhook{ID: 1, Topic: "order.created", Status: "active"})
	})

	webhook, _, err := client.Webhooks.Create(context.Background(), &Webhook{
		Topic:       "order.created",
		DeliveryURL: "https://example.com/webhook",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if webhook.ID != 1 {
		t.Errorf("ID: got %d, want 1", webhook.ID)
	}
	if webhook.Topic != "order.created" {
		t.Errorf("Topic: got %s, want order.created", webhook.Topic)
	}
}

func TestWebhooksGet(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/webhooks/1")
		writeJSON(w, &Webhook{ID: 1, Status: "active", DeliveryURL: "https://example.com/webhook"})
	})

	webhook, _, err := client.Webhooks.Get(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if webhook.Status != "active" {
		t.Errorf("Status: got %s, want active", webhook.Status)
	}
	if webhook.DeliveryURL != "https://example.com/webhook" {
		t.Errorf("DeliveryURL: got %s", webhook.DeliveryURL)
	}
}

func TestWebhooksList(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/webhooks")
		if r.URL.Query().Get("status") != "active" {
			t.Errorf("expected status=active, got %q", r.URL.Query().Get("status"))
		}
		writeJSON(w, &[]Webhook{{ID: 1}, {ID: 2}})
	})

	webhooks, _, err := client.Webhooks.List(context.Background(), &ListWebhooksParams{Status: "active"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(webhooks) != 2 {
		t.Errorf("len: got %d, want 2", len(webhooks))
	}
}

func TestWebhooksUpdate(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPut)
		assertPathSuffix(t, r, "/webhooks/1")
		writeJSON(w, &Webhook{ID: 1, Status: "paused"})
	})

	webhook, _, err := client.Webhooks.Update(context.Background(), "1", &Webhook{Status: "paused"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if webhook.Status != "paused" {
		t.Errorf("Status: got %s, want paused", webhook.Status)
	}
}

func TestWebhooksDelete(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodDelete)
		assertPathSuffix(t, r, "/webhooks/1")
		writeJSON(w, &Webhook{ID: 1})
	})

	webhook, _, err := client.Webhooks.Delete(context.Background(), "1", &DeleteWebhookParams{Force: "true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if webhook.ID != 1 {
		t.Errorf("ID: got %d, want 1", webhook.ID)
	}
}

func TestWebhooksBatch(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/webhooks/batch")
		writeJSON(w, &BatchWebhookUpdateResponse{
			Create: &[]Webhook{{ID: 3}},
		})
	})

	result, _, err := client.Webhooks.Batch(context.Background(), &BatchWebhookUpdate{
		Create: &[]Webhook{{Topic: "order.created", DeliveryURL: "https://example.com/wh"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*result.Create) != 1 {
		t.Errorf("create len: got %d, want 1", len(*result.Create))
	}
}

func TestWebhooksError(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusUnprocessableEntity, "Invalid topic")
	})

	_, _, err := client.Webhooks.Create(context.Background(), &Webhook{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if apiErr, ok := err.(*ErrorResponse); !ok {
		t.Errorf("expected *ErrorResponse, got %T", err)
	} else if apiErr.Data.Status != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", apiErr.Data.Status)
	}
}
