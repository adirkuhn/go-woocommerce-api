package woocommerce

import (
	"context"
	"net/http"
	"testing"
)

func TestOrderNotesCreate(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodPost)
		assertPathSuffix(t, r, "/orders/5/notes")
		writeJSON(w, &OrderNote{ID: 1, Note: "Order confirmed"})
	})

	note, _, err := client.OrderNotes.Create(context.Background(), "5", &OrderNote{Note: "Order confirmed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.ID != 1 {
		t.Errorf("ID: got %d, want 1", note.ID)
	}
	if note.Note != "Order confirmed" {
		t.Errorf("Note: got %s, want 'Order confirmed'", note.Note)
	}
}

func TestOrderNotesGet(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/orders/5/notes/1")
		writeJSON(w, &OrderNote{ID: 1, Author: "system", CustomerNote: true})
	})

	note, _, err := client.OrderNotes.Get(context.Background(), "5", "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Author != "system" {
		t.Errorf("Author: got %s, want system", note.Author)
	}
	if !note.CustomerNote {
		t.Error("expected CustomerNote=true")
	}
}

func TestOrderNotesList(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodGet)
		assertPathSuffix(t, r, "/orders/5/notes")
		if r.URL.Query().Get("type") != "customer" {
			t.Errorf("expected type=customer, got %q", r.URL.Query().Get("type"))
		}
		writeJSON(w, &[]OrderNote{{ID: 1}, {ID: 2}})
	})

	notes, _, err := client.OrderNotes.List(context.Background(), "5", &ListOrderNotesParams{Type: "customer"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*notes) != 2 {
		t.Errorf("len: got %d, want 2", len(*notes))
	}
}

func TestOrderNotesDelete(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		assertMethod(t, r, http.MethodDelete)
		assertPathSuffix(t, r, "/orders/5/notes/1")
		writeJSON(w, &OrderNote{ID: 1})
	})

	note, _, err := client.OrderNotes.Delete(context.Background(), "5", "1", &DeleteOrderNoteParams{Force: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.ID != 1 {
		t.Errorf("ID: got %d, want 1", note.ID)
	}
}

func TestOrderNotesError(t *testing.T) {
	client := newTestServerFn(t, func(w http.ResponseWriter, r *http.Request) {
		writeAPIError(w, http.StatusNotFound, "Note not found")
	})

	_, _, err := client.OrderNotes.Get(context.Background(), "5", "999")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if _, ok := err.(*ErrorResponse); !ok {
		t.Errorf("expected *ErrorResponse, got %T", err)
	}
}
