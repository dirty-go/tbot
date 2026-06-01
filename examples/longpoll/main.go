// longpoll is a minimal echo bot that receives updates via long polling.
//
// Required environment variable:
//
//	TBOT_TOKEN — the bot token from @BotFather
//
// The bot echoes every text message it receives back to the same chat.
// Interrupt with SIGINT or SIGTERM.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
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

	// Long-poll timeout is 30 s; HTTP client deadline is 30 s + 15 s headroom
	// so the transport never cuts off a legitimate server-side wait.
	bot := tbot.NewClient(token, "",
		tbot.WithHTTPClient(&http.Client{
			Timeout: 45 * time.Second,
		}),
		tbot.WithPollTimeout(30),
		tbot.WithMaxRetries(3),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	slog.Info("long-poll bot starting")

	if err := poll(ctx, bot); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	slog.Info("bot stopped")
	return nil
}

// poll runs the GetUpdates loop until ctx is cancelled.
func poll(ctx context.Context, bot *tbot.Client) error {
	for {
		updates, err := bot.GetUpdates(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return err
			}
			slog.Error("GetUpdates failed", "err", err)
			// Back off briefly before retrying so we do not hammer the API on
			// transient failures.
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second):
			}
			continue
		}

		for _, u := range updates {
			handleUpdate(ctx, bot, u)
		}
	}
}

// handleUpdate dispatches a single update and echoes text messages.
func handleUpdate(ctx context.Context, bot *tbot.Client, u tbot.Update) {
	if u.Message == nil || u.Message.Text == "" {
		return
	}

	msg := u.Message
	username := ""
	if msg.From != nil {
		username = msg.From.Username
	}
	slog.Info("message received",
		"chat_id", msg.Chat.ID,
		"from", username,
		"text", msg.Text,
	)

	_, err := bot.SendMessage(ctx, tbot.Int64(msg.Chat.ID), msg.Text)
	if err != nil {
		slog.Error("SendMessage failed", "chat_id", msg.Chat.ID, "err", err)
	}
}
