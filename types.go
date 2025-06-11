package tbot

// Chat represents a chat
type Chat struct {
	ID                          int              `json:"id"`
	Type                        string           `json:"type"`
	Title                       string           `json:"title"`
	Username                    string           `json:"username"`
	FirstName                   string           `json:"first_name"`
	LastName                    string           `json:"last_name"`
	Photo                       *ChatPhoto       `json:"photo"`
	Bio                         string           `json:"bio,omitempty"`
	Description                 string           `json:"description"`
	InviteLink                  string           `json:"invite_link"`
	PinnedMessage               *Message         `json:"pinned_message"`
	Permissions                 *ChatPermissions `json:"permissions"`
	SlowModeDelay               int              `json:"slow_mode_delay"`
	MessageAutoDeleteTime       int64            `json:"message_auto_delete_time"`
	StickerSetName              string           `json:"sticker_set_name"`
	CanSetStickerSet            bool             `json:"can_set_sticker_set"`
	LinkedChatID                int              `json:"linked_chat_id"`
	Location                    *ChatLocation    `json:"location"`
	AllMembersAreAdministrators bool             `json:"all_members_are_administrators"`
}

// ChatLocation Represents a location to which a chat is connected.
type ChatLocation struct {
	Location *Location `json:"location"`
	Address  string    `json:"address"`
}

// MessageEntity represents one special entity in a text message.
// For example, hashtags, usernames, URLs, etc.
type MessageEntity struct {
	Type     string `json:"type"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
	URL      string `json:"url"`
	User     *User  `json:"user,omitempty"`
	Language string `json:"language"`
}

// Audio represents an audio file to bea treated as music by Telegram clients
type Audio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	Performer    string `json:"performer"`
	Title        string `json:"title"`
	MIMEType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

// PhotoSize represents one size of a photo or a file/sticker thumbnail.
type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size"`
}

// Document represents a general file
// (as opposed to photos, voice messages and audio files)
type Document struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Thumb        *PhotoSize `json:"thumb"`
	FileName     string     `json:"file_name"`
	MIMEType     string     `json:"mime_type"`
	FileSize     int        `json:"file_size"`
}

// Animation represents an animation file
// to be displayed in the message containing a game
type Animation struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Thumb        *PhotoSize `json:"thumb"`
	FileName     string     `json:"file_name"`
	MimeType     string     `json:"mime_type"`
	FileSize     int        `json:"file_size"`
}

// Game represents a game. Use BotFather to create and edit games,
// their short names will act as unique identifiers.
type Game struct {
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Photo        []*PhotoSize     `json:"photo"`
	Text         string           `json:"text"`
	TextEntities []*MessageEntity `json:"text_entities"`
	Animation    *Animation       `json:"animation"`
}

// MaskPosition describes the position on faces
// where a mask should be placed by default
type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float32 `json:"x_shift"`
	YShift float32 `json:"y_shift"`
	Scale  float32 `json:"scale"`
}

// Sticker represents a sticker
type Sticker struct {
	FileID       string        `json:"file_id"`
	FileUniqueID string        `json:"file_unique_id"`
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	IsAnimated   bool          `json:"is_animated"`
	Thumb        *PhotoSize    `json:"thumb"`
	Emoji        string        `json:"emoji"`
	MaskPosition *MaskPosition `json:"mask_position"`
	SetName      string        `json:"set_name"`
	FileSize     int           `json:"file_size"`
}

// Video represents a video file
type Video struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumb"`
	MimeType     string     `json:"mime_type"`
	FileSize     int        `json:"file_size"`
}

// Voice represents a voice note
type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

// VideoNote represents a video message
type VideoNote struct {
	FileID       string     `json:"file_id"`
	FileUniqueID string     `json:"file_unique_id"`
	Length       int        `json:"length"`
	Duration     int        `json:"duration"`
	Thumb        *PhotoSize `json:"thumb"`
	FileSize     int        `json:"file_size"`
}

// Contact represents a phone contact
type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserID      int    `json:"user_id"`
}

// Location represents a point on the map
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// Venue represents a venue
type Venue struct {
	Location     Location `json:"location"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	FoursquareID string   `json:"foursquare_id"`
}

// PollOption is an option for Poll
type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

// PollAnswer represents an answer of a user in a non-anonymous poll
type PollAnswer struct {
	PollID    int   `json:"poll_id"`
	User      User  `json:"user"`
	OptionIDs []int `json:"option_ids"`
}

// Poll represents native telegram poll
type Poll struct {
	ID                    string       `json:"id"`
	Question              string       `json:"question"`
	Options               []PollOption `json:"options"`
	TotalVoterCount       int          `json:"total_voter_count"`
	IsClosed              bool         `json:"is_closed"`
	IsAnonymous           bool         `json:"is_anonymous"`
	Type                  string       `json:"type"`
	AllowsMultipleAnswers bool         `json:"allows_multiple_answers"`
	CorrectOptionID       int          `json:"correct_option_id"`
}

// Dice represents native telegram dice
type Dice struct {
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
}

// Invoice contains basic information about an invoice
type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
}

// ShippingAddress represents a shipping address
type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

// OrderInfo represents information about an order
type OrderInfo struct {
	Name            string           `json:"name"`
	PhoneNumber     string           `json:"phone_number"`
	Email           string           `json:"email"`
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

// SuccessfulPayment contains basic information about a successful payment
type SuccessfulPayment struct {
	Currency                string     `json:"currency"`
	TotalAmount             int        `json:"total_amount"`
	InvoicePayload          string     `json:"invoice_payload"`
	ShippingOptionID        string     `json:"shipping_option_id"`
	OrderInfo               *OrderInfo `json:"order_info"`
	TelegramPaymentChargeID string     `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string     `json:"provider_payment_charge_id"`
}

// PassportFile represents a file uploaded to Telegram Passport
type PassportFile struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FileDate     int    `json:"file_date"`
}

// EncryptedPassportElement contains information about documents or other Telegram Passport elements shared with the bot by the user
type EncryptedPassportElement struct {
	Type        string         `json:"type"`
	Data        string         `json:"data"`
	PhoneNumber string         `json:"phone_number"`
	Email       string         `json:"email"`
	Files       []PassportFile `json:"files"`
	FrontSide   *PassportFile  `json:"front_side"`
	ReverseSide *PassportFile  `json:"reverse_side"`
	Selfie      *PassportFile  `json:"selfie"`
}

// EncryptedCredentials contains data required for decrypting and authenticating EncryptedPassportElement
type EncryptedCredentials struct {
	Data   string `json:"data"`
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
}

// PassportData contains information about Telegram Passport data shared with the bot by the user
type PassportData struct {
	Data        []EncryptedPassportElement `json:"data"`
	Credentials EncryptedCredentials       `json:"credentials"`
}

// MessageAutoDeleteTimerChanged This object represents a service message about a change in auto-delete timer settings.
type MessageAutoDeleteTimerChanged struct {
	MessageAutoDeleteTime int32 `json:"message_auto_delete_time"`
}

// ProximityAlertTriggered This object represents the content of a service message,
// sent whenever a user in the chat triggers a proximity alert set by another user.
type ProximityAlertTriggered struct {
	Traveler *User `json:"traveler"`
	Watched  *User `json:"watched"`
	Distance int32 `json:"distance"`
}

// VoiceChatScheduled This object represents a service message about a voice chat
// scheduled in the chat.
type VoiceChatScheduled struct {
	StartDate int32 `json:"start_date"`
}

// VoiceChatStarted This object represents a service message about a voice chat
// started in the chat. Currently holds no information.
type VoiceChatStarted struct{}

// VoiceChatEnded  This object represents a service message about a voice chat ended in the chat.
type VoiceChatEnded struct {
	Duration int32 `json:"duration"`
}

// VoiceChatParticipantsInvited This object represents a service message about new
// members invited to a voice chat.
type VoiceChatParticipantsInvited struct {
	Users []User `json:"users"`
}

// Message represents a message
type Message struct {
	MessageID                     int                            `json:"message_id"`
	From                          *User                          `json:"from,omitempty"`
	SenderChat                    *Chat                          `json:"sender_chat,omitempty"`
	Date                          int64                          `json:"date"`
	Chat                          Chat                           `json:"chat"`
	ForwardFrom                   *User                          `json:"forward_from"`
	ForwardFromChat               *Chat                          `json:"forward_from_chat"`
	ForwardFromMessageID          int                            `json:"forward_from_message_id"`
	ForwardSignature              string                         `json:"forward_signature"`
	ForwardSenderName             string                         `json:"forward_sender_name"`
	ForwardDate                   int64                          `json:"forward_date"`
	ReplyToMessage                *Message                       `json:"reply_to_message"`
	ViaBOT                        *User                          `json:"via_bot,omitempty"`
	EditDate                      int64                          `json:"edit_date"`
	MediaGroupID                  string                         `json:"media_group_id"`
	AuthorSignature               string                         `json:"author_signature"`
	Text                          string                         `json:"text"`
	Entities                      []*MessageEntity               `json:"entities"`
	Animation                     *Animation                     `json:"animation"`
	Audio                         *Audio                         `json:"audio"`
	Document                      *Document                      `json:"document"`
	Photo                         []*PhotoSize                   `json:"photo"`
	CaptionEntities               []*MessageEntity               `json:"caption_entities"`
	Game                          *Game                          `json:"game"`
	Sticker                       *Sticker                       `json:"sticker"`
	Video                         *Video                         `json:"video"`
	Voice                         *Voice                         `json:"voice"`
	VideoNote                     *VideoNote                     `json:"video_note"`
	Caption                       string                         `json:"caption"`
	Contact                       *Contact                       `json:"contact"`
	Location                      *Location                      `json:"location"`
	Venue                         *Venue                         `json:"venue"`
	Poll                          *Poll                          `json:"poll"`
	Dice                          *Dice                          `json:"dice"`
	NewChatMembers                []User                         `json:"new_chat_members"`
	LeftChatMember                *User                          `json:"left_chat_member"`
	NewChatTitle                  string                         `json:"new_chat_title"`
	NewChatPhoto                  []PhotoSize                    `json:"new_chat_photo"`
	DeleteChatPhoto               bool                           `json:"delete_chat_photo"`
	GroupChatCreated              bool                           `json:"group_chat_created"`
	SupergroupChatCreated         bool                           `json:"supergroup_chat_created"`
	ChannelChatCreated            bool                           `json:"channel_chat_created"`
	MigrateToChatID               int                            `json:"migrate_to_chat_id"`
	MigrateFromChatID             int                            `json:"migrate_from_chat_id"`
	PinnedMessage                 *Message                       `json:"pinned_message"`
	Invoice                       *Invoice                       `json:"invoice"`
	SuccessfulPayment             *SuccessfulPayment             `json:"successful_payment"`
	ConnectedWebsite              string                         `json:"connected_website"`
	PassportData                  *PassportData                  `json:"passport_data"`
	ReplyMarkup                   *InlineKeyboardMarkup          `json:"reply_markup"`
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
	ProximityAlertTriggered       *ProximityAlertTriggered       `json:"proximity_alert_triggered,omitempty"`
	VoiceChatScheduled            *VoiceChatScheduled            `json:"voice_chat_scheduled,omitempty"`
	VoiceChatStarted              *VoiceChatStarted              `json:"voice_chat_started,omitempty"`
	VoiceChatEnded                *VoiceChatEnded                `json:"voice_chat_ended,omitempty"`
	VoiceChatParticipantsInvited  *VoiceChatParticipantsInvited  `json:"voice_chat_participants_invited,omitempty"`
}

type InlineKeyboardButton struct {
	Text                         string    `json:"text"`
	URL                          string    `json:"url,omitempty"`
	LoginURL                     *LoginURL `json:"login_url,omitempty"`
	CallbackData                 string    `json:"callback_data,omitempty"`
	SwitchInlineQuery            *string   `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat *string   `json:"switch_inline_query_current_chat,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type KeyboardButtonPoolType struct {
	Type string `json:"type"`
}

// KeyboardButton represents one button of the reply keyboard
type KeyboardButton struct {
	Text            string                  `json:"text"`
	RequestContact  bool                    `json:"request_contact"`
	RequestLocation bool                    `json:"request_location"`
	RequestPool     *KeyboardButtonPoolType `json:"request_pool,omitempty"`
}

type ReplyKeyboardMarkup struct {
	Keyboard        []KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool             `json:"resize_keyboard"`
	OneTimeKeyboard bool             `json:"one_time_keyboard"`
	Selective       bool             `json:"selective"`
}
