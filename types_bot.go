package tbot

// BotCommand represents a registered bot command.
type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

// BotCommandScope describes the scope to which a set of commands applies.
//
// Variant is given by Type:
//   - "default", "all_private_chats", "all_group_chats",
//     "all_chat_administrators" — no other fields
//   - "chat", "chat_administrators"                       — ChatID
//   - "chat_member"                                       — ChatID, UserID
type BotCommandScope struct {
	Type   string `json:"type"`
	ChatID any    `json:"chat_id,omitempty"`
	UserID int64  `json:"user_id,omitempty"`
}

// BotCommandScope Type values.
const (
	BotCommandScopeTypeDefault               = "default"
	BotCommandScopeTypeAllPrivateChats       = "all_private_chats"
	BotCommandScopeTypeAllGroupChats         = "all_group_chats"
	BotCommandScopeTypeAllChatAdministrators = "all_chat_administrators"
	BotCommandScopeTypeChat                  = "chat"
	BotCommandScopeTypeChatAdministrators    = "chat_administrators"
	BotCommandScopeTypeChatMember            = "chat_member"
)

// BotName is the response shape of getMyName.
type BotName struct {
	Name string `json:"name"`
}

// BotDescription is the response shape of getMyDescription.
type BotDescription struct {
	Description string `json:"description"`
}

// BotShortDescription is the response shape of getMyShortDescription.
type BotShortDescription struct {
	ShortDescription string `json:"short_description"`
}
