// payment demonstrates the Telegram Payments flow: send an invoice, handle a
// ShippingQuery and PreCheckoutQuery via GetUpdates, then log a SuccessfulPayment.
//
// Required environment variables:
//
//	TBOT_TOKEN                  — the bot token from @BotFather
//	TBOT_CHAT_ID                — numeric chat ID of the recipient
//	TBOT_PAYMENT_PROVIDER_TOKEN — the payment provider token (from @BotFather)
//
// The program sends one invoice, then polls GetUpdates until it receives a
// SuccessfulPayment or times out after 5 minutes.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dirty-go/tbot"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	token := os.Getenv("TBOT_TOKEN")
	if token == "" {
		return errors.New("TBOT_TOKEN environment variable is not set")
	}
	rawChatID := os.Getenv("TBOT_CHAT_ID")
	if rawChatID == "" {
		return errors.New("TBOT_CHAT_ID environment variable is not set")
	}
	providerToken := os.Getenv("TBOT_PAYMENT_PROVIDER_TOKEN")
	if providerToken == "" {
		return errors.New("TBOT_PAYMENT_PROVIDER_TOKEN environment variable is not set")
	}

	chatIDInt, err := strconv.ParseInt(rawChatID, 10, 64)
	if err != nil {
		return err
	}

	bot := tbot.NewClient(token, "", tbot.WithMaxRetries(3))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	chatID := tbot.Int64(chatIDInt)

	prices := []tbot.LabeledPrice{
		{Label: "30 Days", Amount: 500},
	}

	// providerTokenOpt injects the provider_token field which is not a
	// named parameter in SendInvoice.
	providerTokenOpt := func(v url.Values) {
		v.Set("provider_token", providerToken)
	}

	msg, err := bot.SendInvoice(ctx, chatID,
		"Premium Access",
		"30-day premium subscription",
		"sub_30d",
		"USD",
		prices,
		providerTokenOpt,
	)
	if err != nil {
		return err
	}
	slog.Info("invoice sent", "message_id", msg.MessageID)

	if err := pollPaymentUpdates(ctx, bot); err != nil {
		return err
	}
	return nil
}

// pollPaymentUpdates polls GetUpdates and handles payment-related update types
// until a SuccessfulPayment arrives or the context is cancelled.
func pollPaymentUpdates(ctx context.Context, bot *tbot.Client) error {
	for {
		updates, err := bot.GetUpdates(ctx)
		if err != nil {
			// Context cancellation and deadline expiry are expected stop conditions,
			// not errors from the caller's perspective.
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		}

		for _, u := range updates {
			switch {
			case u.ShippingQuery != nil:
				handleShippingQuery(ctx, bot, u.ShippingQuery)

			case u.PreCheckoutQuery != nil:
				handlePreCheckoutQuery(ctx, bot, u.PreCheckoutQuery)

			case u.Message != nil && u.Message.SuccessfulPayment != nil:
				sp := u.Message.SuccessfulPayment
				slog.Info("payment successful",
					"currency", sp.Currency,
					"total_amount", sp.TotalAmount,
					"invoice_payload", sp.InvoicePayload,
					"telegram_charge_id", sp.TelegramPaymentChargeID,
				)
				// Payment complete — exit the poll loop.
				return nil
			}
		}

		// No updates were returned; wait one second before the next poll so we
		// do not busy-loop against the Telegram API on an empty update stream.
		if len(updates) == 0 {
			select {
			case <-ctx.Done():
				return nil
			case <-time.After(time.Second):
			}
		}
	}
}

// handleShippingQuery answers a shipping query with a single flat-rate option.
func handleShippingQuery(ctx context.Context, bot *tbot.Client, q *tbot.ShippingQuery) {
	slog.Info("shipping query received", "query_id", q.ID)

	options := []tbot.ShippingOption{
		{
			ID:    "standard",
			Title: "Standard",
			Prices: []tbot.LabeledPrice{
				{Label: "Shipping", Amount: 0},
			},
		},
	}

	_, err := bot.AnswerShippingQuery(ctx, q.ID, true, options, "")
	if err != nil {
		slog.Error("AnswerShippingQuery failed", "query_id", q.ID, "err", err)
	}
}

// handlePreCheckoutQuery approves a pre-checkout query.
func handlePreCheckoutQuery(ctx context.Context, bot *tbot.Client, q *tbot.PreCheckoutQuery) {
	slog.Info("pre-checkout query received",
		"query_id", q.ID,
		"currency", q.Currency,
		"total_amount", q.TotalAmount,
	)

	_, err := bot.AnswerPreCheckoutQuery(ctx, q.ID, true, "")
	if err != nil {
		slog.Error("AnswerPreCheckoutQuery failed", "query_id", q.ID, "err", err)
	}
}
