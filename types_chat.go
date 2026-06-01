package tbot

import "encoding/json"

// Chat represents a chat. This is the small Chat object used in updates and
// embedded in many other types. Use ChatFullInfo (returned by getChat) when
// you need the bio, permissions, pinned message, etc.
type Chat struct {
	ID               int64  `json:"id"`
	Type             string `json:"type"`
	Title            string `json:"title,omitempty"`
	Username         string `json:"username,omitempty"`
	FirstName        string `json:"first_name,omitempty"`
	LastName         string `json:"last_name,omitempty"`
	IsForum          bool   `json:"is_forum,omitempty"`
	IsDirectMessages bool   `json:"is_direct_messages,omitempty"`
}

// ChatFullInfo contains full information about a chat, returned by getChat.
type ChatFullInfo struct {
	ID                                 int64                 `json:"id"`
	Type                               string                `json:"type"`
	Title                              string                `json:"title,omitempty"`
	Username                           string                `json:"username,omitempty"`
	FirstName                          string                `json:"first_name,omitempty"`
	LastName                           string                `json:"last_name,omitempty"`
	IsForum                            bool                  `json:"is_forum,omitempty"`
	IsDirectMessages                   bool                  `json:"is_direct_messages,omitempty"`
	AccentColorID                      int                   `json:"accent_color_id"`
	MaxReactionCount                   int                   `json:"max_reaction_count"`
	Photo                              *ChatPhoto            `json:"photo,omitempty"`
	ActiveUsernames                    []string              `json:"active_usernames,omitempty"`
	Birthdate                          *Birthdate            `json:"birthdate,omitempty"`
	BusinessIntro                      *BusinessIntro        `json:"business_intro,omitempty"`
	BusinessLocation                   *BusinessLocation     `json:"business_location,omitempty"`
	BusinessOpeningHours               *BusinessOpeningHours `json:"business_opening_hours,omitempty"`
	PersonalChat                       *Chat                 `json:"personal_chat,omitempty"`
	ParentChat                         *Chat                 `json:"parent_chat,omitempty"`
	AvailableReactions                 []ReactionType        `json:"available_reactions,omitempty"`
	BackgroundCustomEmojiID            string                `json:"background_custom_emoji_id,omitempty"`
	ProfileAccentColorID               int                   `json:"profile_accent_color_id,omitempty"`
	ProfileBackgroundCustomEmojiID     string                `json:"profile_background_custom_emoji_id,omitempty"`
	EmojiStatusCustomEmojiID           string                `json:"emoji_status_custom_emoji_id,omitempty"`
	EmojiStatusExpirationDate          int64                 `json:"emoji_status_expiration_date,omitempty"`
	Bio                                string                `json:"bio,omitempty"`
	HasPrivateForwards                 bool                  `json:"has_private_forwards,omitempty"`
	HasRestrictedVoiceAndVideoMessages bool                  `json:"has_restricted_voice_and_video_messages,omitempty"`
	JoinToSendMessages                 bool                  `json:"join_to_send_messages,omitempty"`
	JoinByRequest                      bool                  `json:"join_by_request,omitempty"`
	Description                        string                `json:"description,omitempty"`
	InviteLink                         string                `json:"invite_link,omitempty"`
	PinnedMessage                      *Message              `json:"pinned_message,omitempty"`
	Permissions                        *ChatPermissions      `json:"permissions,omitempty"`
	AcceptedGiftTypes                  AcceptedGiftTypes     `json:"accepted_gift_types"`
	CanSendPaidMedia                   bool                  `json:"can_send_paid_media,omitempty"`
	SlowModeDelay                      int                   `json:"slow_mode_delay,omitempty"`
	UnrestrictBoostCount               int                   `json:"unrestrict_boost_count,omitempty"`
	MessageAutoDeleteTime              int                   `json:"message_auto_delete_time,omitempty"`
	HasAggressiveAntiSpamEnabled       bool                  `json:"has_aggressive_anti_spam_enabled,omitempty"`
	HasHiddenMembers                   bool                  `json:"has_hidden_members,omitempty"`
	HasProtectedContent                bool                  `json:"has_protected_content,omitempty"`
	HasVisibleHistory                  bool                  `json:"has_visible_history,omitempty"`
	StickerSetName                     string                `json:"sticker_set_name,omitempty"`
	CanSetStickerSet                   bool                  `json:"can_set_sticker_set,omitempty"`
	CustomEmojiStickerSetName          string                `json:"custom_emoji_sticker_set_name,omitempty"`
	LinkedChatID                       int64                 `json:"linked_chat_id,omitempty"`
	Location                           *ChatLocation         `json:"location,omitempty"`
	Rating                             *UserRating           `json:"rating,omitempty"`
	FirstProfileAudio                  *Audio                `json:"first_profile_audio,omitempty"`
	UniqueGiftColors                   *UniqueGiftColors     `json:"unique_gift_colors,omitempty"`
	PaidMessageStarCount               int                   `json:"paid_message_star_count,omitempty"`
}

// ChatPhoto represents a chat photo.
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

// ChatLocation represents a location to which a chat is connected.
type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}

// ChatPermissions describes actions that a non-administrator user is allowed
// to take in a chat.
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`
	CanSendAudios         bool `json:"can_send_audios,omitempty"`
	CanSendDocuments      bool `json:"can_send_documents,omitempty"`
	CanSendPhotos         bool `json:"can_send_photos,omitempty"`
	CanSendVideos         bool `json:"can_send_videos,omitempty"`
	CanSendVideoNotes     bool `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes     bool `json:"can_send_voice_notes,omitempty"`
	CanSendPolls          bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	CanEditTag            bool `json:"can_edit_tag,omitempty"`
	CanChangeInfo         bool `json:"can_change_info,omitempty"`
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`
	CanManageTopics       bool `json:"can_manage_topics,omitempty"`
}

// ChatAdministratorRights describes actions that an administrator can take.
type ChatAdministratorRights struct {
	IsAnonymous           bool `json:"is_anonymous"`
	CanManageChat         bool `json:"can_manage_chat"`
	CanDeleteMessages     bool `json:"can_delete_messages"`
	CanManageVideoChats   bool `json:"can_manage_video_chats"`
	CanRestrictMembers    bool `json:"can_restrict_members"`
	CanPromoteMembers     bool `json:"can_promote_members"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPostStories        bool `json:"can_post_stories"`
	CanEditStories        bool `json:"can_edit_stories"`
	CanDeleteStories      bool `json:"can_delete_stories"`
	CanPostMessages       bool `json:"can_post_messages,omitempty"`
	CanEditMessages       bool `json:"can_edit_messages,omitempty"`
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`
	CanManageTopics       bool `json:"can_manage_topics,omitempty"`
	CanManageDirectMessages bool `json:"can_manage_direct_messages,omitempty"`
	CanManageTags         bool `json:"can_manage_tags,omitempty"`
}

// ChatInviteLink represents an invite link for a chat.
type ChatInviteLink struct {
	InviteLink              string `json:"invite_link"`
	Creator                 User   `json:"creator"`
	CreatesJoinRequest      bool   `json:"creates_join_request"`
	IsPrimary               bool   `json:"is_primary"`
	IsRevoked               bool   `json:"is_revoked"`
	Name                    string `json:"name,omitempty"`
	ExpireDate              int64  `json:"expire_date,omitempty"`
	MemberLimit             int    `json:"member_limit,omitempty"`
	PendingJoinRequestCount int    `json:"pending_join_request_count,omitempty"`
	SubscriptionPeriod      int    `json:"subscription_period,omitempty"`
	SubscriptionPrice       int    `json:"subscription_price,omitempty"`
}

// ChatMember is the union representation of a chat member's status.
//
// Telegram returns one of "creator", "administrator", "member", "restricted",
// "left" or "kicked" via the Status field; the other fields are populated
// according to the variant. To send to the API, set Status and only the fields
// that apply to that status.
type ChatMember struct {
	Status      string `json:"status"`
	User        User   `json:"user"`
	IsAnonymous bool   `json:"is_anonymous,omitempty"`
	CustomTitle string `json:"custom_title,omitempty"`
	Tag         string `json:"tag,omitempty"`

	// Administrator-only fields.
	CanBeEdited             bool `json:"can_be_edited,omitempty"`
	CanManageChat           bool `json:"can_manage_chat,omitempty"`
	CanDeleteMessages       bool `json:"can_delete_messages,omitempty"`
	CanManageVideoChats     bool `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers      bool `json:"can_restrict_members,omitempty"`
	CanPromoteMembers       bool `json:"can_promote_members,omitempty"`
	CanChangeInfo           bool `json:"can_change_info,omitempty"`
	CanInviteUsers          bool `json:"can_invite_users,omitempty"`
	CanPostStories          bool `json:"can_post_stories,omitempty"`
	CanEditStories          bool `json:"can_edit_stories,omitempty"`
	CanDeleteStories        bool `json:"can_delete_stories,omitempty"`
	CanPostMessages         bool `json:"can_post_messages,omitempty"`
	CanEditMessages         bool `json:"can_edit_messages,omitempty"`
	CanPinMessages          bool `json:"can_pin_messages,omitempty"`
	CanManageTopics         bool `json:"can_manage_topics,omitempty"`
	CanManageDirectMessages bool `json:"can_manage_direct_messages,omitempty"`
	CanManageTags           bool `json:"can_manage_tags,omitempty"`

	// Restricted-only fields.
	IsMember              bool  `json:"is_member,omitempty"`
	CanSendMessages       bool  `json:"can_send_messages,omitempty"`
	CanSendAudios         bool  `json:"can_send_audios,omitempty"`
	CanSendDocuments      bool  `json:"can_send_documents,omitempty"`
	CanSendPhotos         bool  `json:"can_send_photos,omitempty"`
	CanSendVideos         bool  `json:"can_send_videos,omitempty"`
	CanSendVideoNotes     bool  `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes     bool  `json:"can_send_voice_notes,omitempty"`
	CanSendPolls          bool  `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool  `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool  `json:"can_add_web_page_previews,omitempty"`
	CanEditTag            bool  `json:"can_edit_tag,omitempty"`
	UntilDate             int64 `json:"until_date,omitempty"`
}

// MarshalJSON emits only fields valid for the selected chat-member status.
func (m ChatMember) MarshalJSON() ([]byte, error) {
	type payload struct {
		Status      string `json:"status"`
		User        User   `json:"user"`
		IsAnonymous bool   `json:"is_anonymous,omitempty"`
		CustomTitle string `json:"custom_title,omitempty"`
		Tag         string `json:"tag,omitempty"`

		CanBeEdited             bool `json:"can_be_edited,omitempty"`
		CanManageChat           bool `json:"can_manage_chat,omitempty"`
		CanDeleteMessages       bool `json:"can_delete_messages,omitempty"`
		CanManageVideoChats     bool `json:"can_manage_video_chats,omitempty"`
		CanRestrictMembers      bool `json:"can_restrict_members,omitempty"`
		CanPromoteMembers       bool `json:"can_promote_members,omitempty"`
		CanChangeInfo           bool `json:"can_change_info,omitempty"`
		CanInviteUsers          bool `json:"can_invite_users,omitempty"`
		CanPostStories          bool `json:"can_post_stories,omitempty"`
		CanEditStories          bool `json:"can_edit_stories,omitempty"`
		CanDeleteStories        bool `json:"can_delete_stories,omitempty"`
		CanPostMessages         bool `json:"can_post_messages,omitempty"`
		CanEditMessages         bool `json:"can_edit_messages,omitempty"`
		CanPinMessages          bool `json:"can_pin_messages,omitempty"`
		CanManageTopics         bool `json:"can_manage_topics,omitempty"`
		CanManageDirectMessages bool `json:"can_manage_direct_messages,omitempty"`
		CanManageTags           bool `json:"can_manage_tags,omitempty"`

		IsMember              bool  `json:"is_member,omitempty"`
		CanSendMessages       bool  `json:"can_send_messages,omitempty"`
		CanSendAudios         bool  `json:"can_send_audios,omitempty"`
		CanSendDocuments      bool  `json:"can_send_documents,omitempty"`
		CanSendPhotos         bool  `json:"can_send_photos,omitempty"`
		CanSendVideos         bool  `json:"can_send_videos,omitempty"`
		CanSendVideoNotes     bool  `json:"can_send_video_notes,omitempty"`
		CanSendVoiceNotes     bool  `json:"can_send_voice_notes,omitempty"`
		CanSendPolls          bool  `json:"can_send_polls,omitempty"`
		CanSendOtherMessages  bool  `json:"can_send_other_messages,omitempty"`
		CanAddWebPagePreviews bool  `json:"can_add_web_page_previews,omitempty"`
		CanEditTag            bool  `json:"can_edit_tag,omitempty"`
		UntilDate             int64 `json:"until_date,omitempty"`
	}
	out := payload{
		Status: m.Status,
		User:   m.User,
	}
	switch m.Status {
	case ChatMemberStatusCreator:
		out.IsAnonymous = m.IsAnonymous
		out.CustomTitle = m.CustomTitle
		out.Tag = m.Tag
	case ChatMemberStatusAdministrator:
		out.IsAnonymous = m.IsAnonymous
		out.CustomTitle = m.CustomTitle
		out.Tag = m.Tag
		out.CanBeEdited = m.CanBeEdited
		out.CanManageChat = m.CanManageChat
		out.CanDeleteMessages = m.CanDeleteMessages
		out.CanManageVideoChats = m.CanManageVideoChats
		out.CanRestrictMembers = m.CanRestrictMembers
		out.CanPromoteMembers = m.CanPromoteMembers
		out.CanChangeInfo = m.CanChangeInfo
		out.CanInviteUsers = m.CanInviteUsers
		out.CanPostStories = m.CanPostStories
		out.CanEditStories = m.CanEditStories
		out.CanDeleteStories = m.CanDeleteStories
		out.CanPostMessages = m.CanPostMessages
		out.CanEditMessages = m.CanEditMessages
		out.CanPinMessages = m.CanPinMessages
		out.CanManageTopics = m.CanManageTopics
		out.CanManageDirectMessages = m.CanManageDirectMessages
		out.CanManageTags = m.CanManageTags
	case ChatMemberStatusRestricted:
		out.IsMember = m.IsMember
		out.CanSendMessages = m.CanSendMessages
		out.CanSendAudios = m.CanSendAudios
		out.CanSendDocuments = m.CanSendDocuments
		out.CanSendPhotos = m.CanSendPhotos
		out.CanSendVideos = m.CanSendVideos
		out.CanSendVideoNotes = m.CanSendVideoNotes
		out.CanSendVoiceNotes = m.CanSendVoiceNotes
		out.CanSendPolls = m.CanSendPolls
		out.CanSendOtherMessages = m.CanSendOtherMessages
		out.CanAddWebPagePreviews = m.CanAddWebPagePreviews
		out.CanEditTag = m.CanEditTag
		out.UntilDate = m.UntilDate
	}
	return json.Marshal(out)
}

// ChatMember status constants.
const (
	ChatMemberStatusCreator       = "creator"
	ChatMemberStatusAdministrator = "administrator"
	ChatMemberStatusMember        = "member"
	ChatMemberStatusRestricted    = "restricted"
	ChatMemberStatusLeft          = "left"
	ChatMemberStatusBanned        = "kicked"
)

// IsCreator reports whether the member is the chat owner.
func (m ChatMember) IsCreator() bool { return m.Status == ChatMemberStatusCreator }

// IsAdministrator reports whether the member is an administrator.
func (m ChatMember) IsAdministrator() bool { return m.Status == ChatMemberStatusAdministrator }

// IsMemberStatus reports whether the member has the plain "member" status.
// Named IsMemberStatus to avoid collision with the IsMember bool field.
func (m ChatMember) IsMemberStatus() bool { return m.Status == ChatMemberStatusMember }

// IsRestricted reports whether the member is restricted.
func (m ChatMember) IsRestricted() bool { return m.Status == ChatMemberStatusRestricted }

// HasLeft reports whether the member has left the chat.
func (m ChatMember) HasLeft() bool { return m.Status == ChatMemberStatusLeft }

// IsBanned reports whether the member has been banned (kicked).
func (m ChatMember) IsBanned() bool { return m.Status == ChatMemberStatusBanned }

// ChatMemberUpdated represents changes in the status of a chat member.
type ChatMemberUpdated struct {
	Chat                    Chat            `json:"chat"`
	From                    User            `json:"from"`
	Date                    int64           `json:"date"`
	OldChatMember           ChatMember      `json:"old_chat_member"`
	NewChatMember           ChatMember      `json:"new_chat_member"`
	InviteLink              *ChatInviteLink `json:"invite_link,omitempty"`
	ViaJoinRequest          bool            `json:"via_join_request,omitempty"`
	ViaChatFolderInviteLink bool            `json:"via_chat_folder_invite_link,omitempty"`
}

// ChatJoinRequest represents a join request sent to a chat.
type ChatJoinRequest struct {
	Chat       Chat            `json:"chat"`
	From       User            `json:"from"`
	UserChatID int64           `json:"user_chat_id"`
	Date       int64           `json:"date"`
	Bio        string          `json:"bio,omitempty"`
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

// Birthdate is a user's date of birth (year is optional).
type Birthdate struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year,omitempty"`
}
