package tbot

// StarAmount is an amount of Telegram Stars, optionally with sub-star
// precision (1/1_000_000_000 Stars).
type StarAmount struct {
	Amount         int `json:"amount"`
	NanostarAmount int `json:"nanostar_amount,omitempty"`
}

// StarTransaction describes a single Star transaction. Source is set for
// incoming transactions; Receiver is set for outgoing ones. Some pseudo
// transactions (gifts, refunds) populate neither.
type StarTransaction struct {
	ID             string              `json:"id"`
	Amount         int                 `json:"amount"`
	NanostarAmount int                 `json:"nanostar_amount,omitempty"`
	Date           int64               `json:"date"`
	Source         *TransactionPartner `json:"source,omitempty"`
	Receiver       *TransactionPartner `json:"receiver,omitempty"`
}

// StarTransactions is the response shape of getStarTransactions.
type StarTransactions struct {
	Transactions []StarTransaction `json:"transactions"`
}

// TransactionPartner describes the other side of a Star transaction.
//
// Variant is given by Type:
//   - "user"              — User, TransactionType, plus optional fields
//   - "chat"              — Chat, optional Gift
//   - "affiliate_program" — SponsorUser, CommissionPerMille
//   - "fragment"          — optional WithdrawalState
//   - "telegram_ads"      — no extra fields
//   - "telegram_api"      — RequestCount
//   - "other"             — no extra fields
type TransactionPartner struct {
	Type                        string                  `json:"type"`
	TransactionType             string                  `json:"transaction_type,omitempty"`
	User                        *User                   `json:"user,omitempty"`
	Chat                        *Chat                   `json:"chat,omitempty"`
	Affiliate                   *AffiliateInfo          `json:"affiliate,omitempty"`
	InvoicePayload              string                  `json:"invoice_payload,omitempty"`
	SubscriptionPeriod          int                     `json:"subscription_period,omitempty"`
	PaidMedia                   []PaidMedia             `json:"paid_media,omitempty"`
	PaidMediaPayload            string                  `json:"paid_media_payload,omitempty"`
	Gift                        *Gift                   `json:"gift,omitempty"`
	PremiumSubscriptionDuration int                     `json:"premium_subscription_duration,omitempty"`
	SponsorUser                 *User                   `json:"sponsor_user,omitempty"`
	CommissionPerMille          int                     `json:"commission_per_mille,omitempty"`
	WithdrawalState             *RevenueWithdrawalState `json:"withdrawal_state,omitempty"`
	RequestCount                int                     `json:"request_count,omitempty"`
}

// TransactionPartner Type values.
const (
	TransactionPartnerTypeUser             = "user"
	TransactionPartnerTypeChat             = "chat"
	TransactionPartnerTypeAffiliateProgram = "affiliate_program"
	TransactionPartnerTypeFragment         = "fragment"
	TransactionPartnerTypeTelegramAds      = "telegram_ads"
	TransactionPartnerTypeTelegramAPI      = "telegram_api"
	TransactionPartnerTypeOther            = "other"
)

// RevenueWithdrawalState describes the state of a revenue withdrawal. Variant
// is given by Type: "pending", "succeeded" (Date, URL set) or "failed".
type RevenueWithdrawalState struct {
	Type string `json:"type"`
	Date int64  `json:"date,omitempty"`
	URL  string `json:"url,omitempty"`
}

// RevenueWithdrawalState Type values.
const (
	RevenueWithdrawalStatePending   = "pending"
	RevenueWithdrawalStateSucceeded = "succeeded"
	RevenueWithdrawalStateFailed    = "failed"
)

// AffiliateInfo describes the affiliate that received a commission via a Star
// transaction.
type AffiliateInfo struct {
	AffiliateUser      *User `json:"affiliate_user,omitempty"`
	AffiliateChat      *Chat `json:"affiliate_chat,omitempty"`
	CommissionPerMille int   `json:"commission_per_mille"`
	Amount             int   `json:"amount"`
	NanostarAmount     int   `json:"nanostar_amount,omitempty"`
}
