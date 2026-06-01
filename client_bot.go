package tbot

import (
	"context"
	"net/url"
)

// SetMyName changes the bot's name. Pass an empty languageCode for the
// default name; BCP-47 codes set language-specific names.
func (c *Client) SetMyName(ctx context.Context, name, languageCode string) (bool, error) {
	req := url.Values{}
	if name != "" {
		req.Set("name", name)
	}
	if languageCode != "" {
		req.Set("language_code", languageCode)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setMyName", req, &ok)
	return ok, err
}

// GetMyName returns the bot's name for the given language. Pass an empty
// languageCode to get the default name.
func (c *Client) GetMyName(ctx context.Context, languageCode string) (*BotName, error) {
	req := url.Values{}
	if languageCode != "" {
		req.Set("language_code", languageCode)
	}
	name := &BotName{}
	err := c.sendRequest(ctx, "/getMyName", req, name)
	return name, err
}

// SetMyDescription changes the bot's description — the text shown on the bot's
// profile page and sent to users who start a chat with the bot via a link.
func (c *Client) SetMyDescription(ctx context.Context, description, languageCode string) (bool, error) {
	req := url.Values{}
	if description != "" {
		req.Set("description", description)
	}
	if languageCode != "" {
		req.Set("language_code", languageCode)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setMyDescription", req, &ok)
	return ok, err
}

// GetMyDescription returns the bot's description for the given language.
func (c *Client) GetMyDescription(ctx context.Context, languageCode string) (*BotDescription, error) {
	req := url.Values{}
	if languageCode != "" {
		req.Set("language_code", languageCode)
	}
	desc := &BotDescription{}
	err := c.sendRequest(ctx, "/getMyDescription", req, desc)
	return desc, err
}

// SetMyShortDescription changes the bot's short description shown on the
// bot's profile page and sent to users who start a chat.
func (c *Client) SetMyShortDescription(ctx context.Context, shortDescription, languageCode string) (bool, error) {
	req := url.Values{}
	if shortDescription != "" {
		req.Set("short_description", shortDescription)
	}
	if languageCode != "" {
		req.Set("language_code", languageCode)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setMyShortDescription", req, &ok)
	return ok, err
}

// GetMyShortDescription returns the bot's short description for the given
// language.
func (c *Client) GetMyShortDescription(ctx context.Context, languageCode string) (*BotShortDescription, error) {
	req := url.Values{}
	if languageCode != "" {
		req.Set("language_code", languageCode)
	}
	desc := &BotShortDescription{}
	err := c.sendRequest(ctx, "/getMyShortDescription", req, desc)
	return desc, err
}

// SetMyCommands registers the list of bot commands. Scope and languageCode
// narrow which users see these commands; omit both to set the global default.
func (c *Client) SetMyCommands(ctx context.Context, commands []BotCommand, scope *BotCommandScope, languageCode string) (bool, error) {
	req := url.Values{}
	req.Set("commands", structString(commands))
	if scope != nil {
		req.Set("scope", structString(scope))
	}
	if languageCode != "" {
		req.Set("language_code", languageCode)
	}
	var ok bool
	err := c.sendRequest(ctx, "/setMyCommands", req, &ok)
	return ok, err
}

// DeleteMyCommands deletes the list of bot commands for the given scope and
// language. Pass nil scope and empty languageCode to clear the global default.
func (c *Client) DeleteMyCommands(ctx context.Context, scope *BotCommandScope, languageCode string) (bool, error) {
	req := url.Values{}
	if scope != nil {
		req.Set("scope", structString(scope))
	}
	if languageCode != "" {
		req.Set("language_code", languageCode)
	}
	var ok bool
	err := c.sendRequest(ctx, "/deleteMyCommands", req, &ok)
	return ok, err
}

// GetMyCommands returns the list of bot commands for the given scope and
// language. Pass nil scope and empty languageCode to get the global default.
func (c *Client) GetMyCommands(ctx context.Context, scope *BotCommandScope, languageCode string) ([]BotCommand, error) {
	req := url.Values{}
	if scope != nil {
		req.Set("scope", structString(scope))
	}
	if languageCode != "" {
		req.Set("language_code", languageCode)
	}
	var commands []BotCommand
	err := c.sendRequest(ctx, "/getMyCommands", req, &commands)
	return commands, err
}

// SetChatMenuButton changes the menu button for a private chat or sets the
// default menu button. Pass chatID 0 to set the default for all users.
func (c *Client) SetChatMenuButton(ctx context.Context, chatID ChatID, menuButton *MenuButton) (bool, error) {
	req := url.Values{}
	if !chatID.IsZero() {
		req.Set("chat_id", chatID.String())
	}
	if menuButton != nil {
		req.Set("menu_button", structString(menuButton))
	}
	var ok bool
	err := c.sendRequest(ctx, "/setChatMenuButton", req, &ok)
	return ok, err
}

// GetChatMenuButton returns the menu button for a private chat or the default
// menu button. Pass chatID 0 to get the default.
func (c *Client) GetChatMenuButton(ctx context.Context, chatID ChatID) (*MenuButton, error) {
	req := url.Values{}
	if !chatID.IsZero() {
		req.Set("chat_id", chatID.String())
	}
	btn := &MenuButton{}
	err := c.sendRequest(ctx, "/getChatMenuButton", req, btn)
	return btn, err
}

// SetMyDefaultAdministratorRights changes the default administrator rights
// requested by the bot when added as an administrator. forChannel=true sets
// the rights for channels; false sets them for groups and supergroups.
func (c *Client) SetMyDefaultAdministratorRights(ctx context.Context, rights *ChatAdministratorRights, forChannel bool) (bool, error) {
	req := url.Values{}
	if rights != nil {
		req.Set("rights", structString(rights))
	}
	if forChannel {
		req.Set("for_channels", "true")
	}
	var ok bool
	err := c.sendRequest(ctx, "/setMyDefaultAdministratorRights", req, &ok)
	return ok, err
}

// GetMyDefaultAdministratorRights returns the current default administrator
// rights for the bot. forChannel=true returns channel rights; false returns
// group/supergroup rights.
func (c *Client) GetMyDefaultAdministratorRights(ctx context.Context, forChannel bool) (*ChatAdministratorRights, error) {
	req := url.Values{}
	if forChannel {
		req.Set("for_channels", "true")
	}
	rights := &ChatAdministratorRights{}
	err := c.sendRequest(ctx, "/getMyDefaultAdministratorRights", req, rights)
	return rights, err
}
