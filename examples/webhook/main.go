// webhook is an echo bot that receives updates via a Telegram webhook.
//
// Required environment variables:
//
//	TBOT_TOKEN      — the bot token from @BotFather
//	WEBHOOK_URL     — the public HTTPS URL Telegram should post updates to
//	                  (e.g. https://example.com/webhook)
//	WEBHOOK_SECRET  — secret token passed to setWebhook; Telegram echoes it in
//	                  X-Telegram-Bot-Api-Secret-Token on every incoming request
//
// The server listens on :8080.  Interrupt with SIGINT or SIGTERM; the bot
// calls DeleteWebhook before exiting so Telegram stops sending updates.
package main

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
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
	webhookURL := os.Getenv("WEBHOOK_URL")
	if webhookURL == "" {
		return errors.New("WEBHOOK_URL environment variable is not set")
	}
	secret := os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		return errors.New("WEBHOOK_SECRET environment variable is not set")
	}

	bot := tbot.NewClient(token, "", tbot.WithMaxRetries(3))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	secretOpt := func(v url.Values) { v.Set("secret_token", secret) }
	if _, err := bot.SetWebhook(ctx, webhookURL, secretOpt); err != nil {
		return err
	}
	slog.Info("webhook registered", "url", webhookURL)

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", makeHandler(bot, secret))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the HTTP server in the background.
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("http server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	// Wait for either a signal or a server error.
	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case err := <-serverErr:
		if err != nil {
			return err
		}
	}

	// Graceful shutdown: remove the webhook first, then stop the server.
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if _, err := bot.DeleteWebhook(shutdownCtx, false); err != nil {
		slog.Error("DeleteWebhook failed", "err", err)
	} else {
		slog.Info("webhook deleted")
	}

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "err", err)
	}

	slog.Info("bot stopped")
	return nil
}

// makeHandler returns an http.HandlerFunc that decodes Telegram webhook
// payloads and dispatches them. It verifies the X-Telegram-Bot-Api-Secret-Token
// header with a constant-time comparison to prevent forged-update injection.
func makeHandler(bot *tbot.Client, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Reject requests that do not carry the correct secret token.
		if subtle.ConstantTimeCompare(
			[]byte(r.Header.Get("X-Telegram-Bot-Api-Secret-Token")),
			[]byte(secret),
		) != 1 {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		// Limit request body to 1 MiB to prevent memory exhaustion.
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		var update tbot.Update
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			slog.Error("decode update failed", "err", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Always respond 200 immediately; Telegram will retry on non-2xx.
		w.WriteHeader(http.StatusOK)

		// Handle in the same goroutine — the request context is still live.
		handleWebhookUpdate(r.Context(), bot, update)
	}
}

// handleWebhookUpdate dispatches one update received via webhook.
func handleWebhookUpdate(ctx context.Context, bot *tbot.Client, u tbot.Update) {
	if u.Message == nil {
		return
	}

	msg := u.Message
	chatID := tbot.Int64(msg.Chat.ID)

	username := ""
	if msg.From != nil {
		username = msg.From.Username
	}
	slog.Info("message received",
		"chat_id", msg.Chat.ID,
		"from", username,
		"text", msg.Text,
	)

	if msg.Text == "/start" {
		keyboard := &tbot.InlineKeyboardMarkup{
			InlineKeyboard: [][]tbot.InlineKeyboardButton{
				{
					{Text: "About", CallbackData: "about"},
					{Text: "Help", CallbackData: "help"},
				},
			},
		}
		_, err := bot.SendMessage(ctx, chatID,
			"Welcome! Choose an option:",
			tbot.OptInlineKeyboardMarkup(keyboard),
		)
		if err != nil {
			slog.Error("SendMessage /start failed", "chat_id", msg.Chat.ID, "err", err)
		}
		return
	}

	if msg.Text != "" {
		_, err := bot.SendMessage(ctx, chatID, msg.Text)
		if err != nil {
			slog.Error("SendMessage echo failed", "chat_id", msg.Chat.ID, "err", err)
		}
	}
}
