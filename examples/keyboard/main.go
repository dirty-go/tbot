// keyboard demonstrates the three reply-markup types: InlineKeyboardMarkup,
// ReplyKeyboardMarkup, and ReplyKeyboardRemove.
//
// Required environment variables:
//
//	TBOT_TOKEN      — the bot token from @BotFather
//	TBOT_CHAT_ID    — numeric chat ID to send messages to (e.g. your own user ID)
//
// The program sends three messages and exits — no polling loop is needed.
package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/dirty-go/tbot"
)

func main() {
	token := os.Getenv("TBOT_TOKEN")
	if token == "" {
		slog.Error("TBOT_TOKEN environment variable is not set")
		os.Exit(1)
	}
	rawChatID := os.Getenv("TBOT_CHAT_ID")
	if rawChatID == "" {
		slog.Error("TBOT_CHAT_ID environment variable is not set")
		os.Exit(1)
	}
	chatIDInt, err := strconv.ParseInt(rawChatID, 10, 64)
	if err != nil {
		slog.Error("TBOT_CHAT_ID is not a valid integer", "value", rawChatID, "err", err)
		os.Exit(1)
	}

	bot := tbot.NewClient(token, "")
	ctx := context.Background()
	chatID := tbot.Int64(chatIDInt)

	if err := sendInlineKeyboard(ctx, bot, chatID); err != nil {
		slog.Error("sendInlineKeyboard failed", "err", err)
		os.Exit(1)
	}

	if err := sendReplyKeyboard(ctx, bot, chatID); err != nil {
		slog.Error("sendReplyKeyboard failed", "err", err)
		os.Exit(1)
	}

	if err := sendKeyboardRemove(ctx, bot, chatID); err != nil {
		slog.Error("sendKeyboardRemove failed", "err", err)
		os.Exit(1)
	}

	slog.Info("all keyboard messages sent")
}

// sendInlineKeyboard sends a message with an InlineKeyboardMarkup containing
// three buttons: "Option A", "Option B", and "Help".
func sendInlineKeyboard(ctx context.Context, bot *tbot.Client, chatID tbot.ChatID) error {
	markup := &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			{
				{Text: "Option A", CallbackData: "option_a"},
				{Text: "Option B", CallbackData: "option_b"},
				{Text: "Help", CallbackData: "help"},
			},
		},
	}
	msg, err := bot.SendMessage(ctx, chatID,
		"Choose an option:",
		tbot.OptInlineKeyboardMarkup(markup),
	)
	if err != nil {
		return err
	}
	slog.Info("inline keyboard sent", "message_id", msg.MessageID)
	return nil
}

// sendReplyKeyboard sends a message with a two-row ReplyKeyboardMarkup:
// row 1: "Yes", "No" — row 2: "Cancel".
func sendReplyKeyboard(ctx context.Context, bot *tbot.Client, chatID tbot.ChatID) error {
	markup := &tbot.ReplyKeyboardMarkup{
		Keyboard: [][]tbot.KeyboardButton{
			{
				{Text: "Yes"},
				{Text: "No"},
			},
			{
				{Text: "Cancel"},
			},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
	}
	msg, err := bot.SendMessage(ctx, chatID,
		"Confirm your choice:",
		tbot.OptReplyKeyboardMarkup(markup),
	)
	if err != nil {
		return err
	}
	slog.Info("reply keyboard sent", "message_id", msg.MessageID)
	return nil
}

// sendKeyboardRemove sends a message that removes any active reply keyboard.
func sendKeyboardRemove(ctx context.Context, bot *tbot.Client, chatID tbot.ChatID) error {
	msg, err := bot.SendMessage(ctx, chatID,
		"Keyboard removed.",
		tbot.OptReplyKeyboardRemove,
	)
	if err != nil {
		return err
	}
	slog.Info("keyboard remove sent", "message_id", msg.MessageID)
	return nil
}
