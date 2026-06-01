package tbot

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendInvoice_PricesSerialised(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/sendInvoice") {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, okMessage)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	prices := []LabeledPrice{
		{Label: "Item", Amount: 1000},
		{Label: "Tax", Amount: 100},
	}
	_, err := c.SendInvoice(context.Background(), Int64(42), "Test", "A test invoice", "payload123", "USD", prices)
	if err != nil {
		t.Fatalf("SendInvoice() error: %v", err)
	}
	for _, want := range []string{
		"chat_id=42",
		"title=Test",
		"description=A+test+invoice",
		"payload=payload123",
		"currency=USD",
		"prices=",
		"Item",
		"1000",
		"Tax",
		"100",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q; body=%q", want, body)
		}
	}
}

func TestAnswerShippingQuery_OkTrue(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/answerShippingQuery") {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	opts := []ShippingOption{
		{ID: "express", Title: "Express", Prices: []LabeledPrice{{Label: "Express delivery", Amount: 500}}},
	}
	ok, err := c.AnswerShippingQuery(context.Background(), "sq-1", true, opts, "")
	if err != nil {
		t.Fatalf("AnswerShippingQuery() error: %v", err)
	}
	if !ok {
		t.Fatalf("AnswerShippingQuery() = false, want true")
	}
	for _, want := range []string{"ok=true", "shipping_options="} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q; body=%q", want, body)
		}
	}
	if strings.Contains(body, "error_message") {
		t.Errorf("body should not contain error_message; body=%q", body)
	}
}

func TestAnswerShippingQuery_OkFalse(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/answerShippingQuery") {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	ok, err := c.AnswerShippingQuery(context.Background(), "sq-2", false, nil, "Cannot ship to this address")
	if err != nil {
		t.Fatalf("AnswerShippingQuery() error: %v", err)
	}
	if !ok {
		t.Fatalf("AnswerShippingQuery() = false, want true")
	}
	for _, want := range []string{"ok=false", "error_message="} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q; body=%q", want, body)
		}
	}
	if strings.Contains(body, "shipping_options") {
		t.Errorf("body should not contain shipping_options; body=%q", body)
	}
}

func TestAnswerPreCheckoutQuery_OkFalse(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/answerPreCheckoutQuery") {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	ok, err := c.AnswerPreCheckoutQuery(context.Background(), "pcq-1", false, "Out of stock")
	if err != nil {
		t.Fatalf("AnswerPreCheckoutQuery() error: %v", err)
	}
	if !ok {
		t.Fatalf("AnswerPreCheckoutQuery() = false, want true")
	}
	for _, want := range []string{"ok=false", "error_message=Out+of+stock"} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q; body=%q", want, body)
		}
	}
}

func TestEditUserStarSubscription_Canceled(t *testing.T) {
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/editUserStarSubscription") {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		fmt.Fprint(w, `{"ok":true,"result":true}`)
	}))
	defer srv.Close()

	c, _ := newTestClient(t, srv)
	ok, err := c.EditUserStarSubscription(context.Background(), 12345, "charge-abc", true)
	if err != nil {
		t.Fatalf("EditUserStarSubscription() error: %v", err)
	}
	if !ok {
		t.Fatalf("EditUserStarSubscription() = false, want true")
	}
	for _, want := range []string{"user_id=12345", "telegram_payment_charge_id=charge-abc", "is_canceled=true"} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing %q; body=%q", want, body)
		}
	}
}
