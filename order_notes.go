package woocommerce

import (
	"context"
	"net/http"
)

// Order Notes service
type OrderNotesService service

// OrderNote object. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#order-note-properties
type OrderNote struct {
	ID             int    `json:"id,omitempty"`
	Author         string `json:"author,omitempty"`
	DateCreated    string `json:"date_created,omitempty"`
	DateCreatedGmt string `json:"date_created_gmt,omitempty"`
	Note           string `json:"note,omitempty"`
	CustomerNote   bool   `json:"customer_note,omitempty"`
	AddedByUser    bool   `json:"added_by_user,omitempty"`
}

type ListOrderNotesParams struct {
	Context string `url:"context,omitempty"`
	Type    string `url:"type,omitempty"`
}

type DeleteOrderNoteParams struct {
	Force bool `url:"force"`
}

// Create an order note. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#create-an-order-note
func (service *OrderNotesService) Create(ctx context.Context, orderID string, orderNote *OrderNote) (*OrderNote, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "POST", "/orders/"+orderID+"/notes", nil, orderNote)
	if err != nil {
		return nil, nil, err
	}

	createdNote := new(OrderNote)
	response, err := service.client.Do(req, createdNote)
	if err != nil {
		return nil, response, err
	}

	return createdNote, response, nil
}

// Get an order note. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#retrieve-an-order-note
func (service *OrderNotesService) Get(ctx context.Context, orderID string, noteID string) (*OrderNote, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/orders/"+orderID+"/notes/"+noteID, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	orderNote := new(OrderNote)
	response, err := service.client.Do(req, orderNote)
	if err != nil {
		return nil, response, err
	}

	return orderNote, response, nil
}

// List order notes. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#list-all-order-notes
func (service *OrderNotesService) List(ctx context.Context, orderID string, opts *ListOrderNotesParams) (*[]OrderNote, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "GET", "/orders/"+orderID+"/notes", opts, nil)
	if err != nil {
		return nil, nil, err
	}

	notes := new([]OrderNote)
	response, err := service.client.Do(req, notes)
	if err != nil {
		return nil, response, err
	}

	return notes, response, nil
}

// Delete an order note. Reference: https://woocommerce.github.io/woocommerce-rest-api-docs/#delete-an-order-note
func (service *OrderNotesService) Delete(ctx context.Context, orderID string, noteID string, opts *DeleteOrderNoteParams) (*OrderNote, *http.Response, error) {
	req, err := service.client.NewRequest(ctx, "DELETE", "/orders/"+orderID+"/notes/"+noteID, opts, nil)
	if err != nil {
		return nil, nil, err
	}

	orderNote := new(OrderNote)
	response, err := service.client.Do(req, orderNote)
	if err != nil {
		return nil, response, err
	}

	return orderNote, response, nil
}
