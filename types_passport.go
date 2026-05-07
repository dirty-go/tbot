package tbot

// PassportData contains information about Telegram Passport data shared with
// the bot by the user.
type PassportData struct {
	Data        []EncryptedPassportElement `json:"data"`
	Credentials EncryptedCredentials       `json:"credentials"`
}

// PassportFile represents a file uploaded to Telegram Passport.
type PassportFile struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size"`
	FileDate     int64  `json:"file_date"`
}

// EncryptedPassportElement contains information about documents or other
// Telegram Passport elements shared with the bot.
type EncryptedPassportElement struct {
	Type        string         `json:"type"`
	Hash        string         `json:"hash"`
	Data        string         `json:"data,omitempty"`
	PhoneNumber string         `json:"phone_number,omitempty"`
	Email       string         `json:"email,omitempty"`
	Files       []PassportFile `json:"files,omitempty"`
	FrontSide   *PassportFile  `json:"front_side,omitempty"`
	ReverseSide *PassportFile  `json:"reverse_side,omitempty"`
	Selfie      *PassportFile  `json:"selfie,omitempty"`
	Translation []PassportFile `json:"translation,omitempty"`
}

// EncryptedCredentials contains the data required to decrypt and authenticate
// a EncryptedPassportElement.
type EncryptedCredentials struct {
	Data   string `json:"data"`
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
}

// PassportElementError describes an error in a Telegram Passport element to
// be reported via setPassportDataErrors.
//
// Variant is given by Source — one of "data", "front_side", "reverse_side",
// "selfie", "file", "files", "translation_file", "translation_files",
// "unspecified".
type PassportElementError struct {
	Source      string   `json:"source"`
	Type        string   `json:"type"`
	Message     string   `json:"message"`
	FieldName   string   `json:"field_name,omitempty"`
	DataHash    string   `json:"data_hash,omitempty"`
	FileHash    string   `json:"file_hash,omitempty"`
	FileHashes  []string `json:"file_hashes,omitempty"`
	ElementHash string   `json:"element_hash,omitempty"`
}

// PassportElementError Source values.
const (
	PassportElementErrorSourceData             = "data"
	PassportElementErrorSourceFrontSide        = "front_side"
	PassportElementErrorSourceReverseSide      = "reverse_side"
	PassportElementErrorSourceSelfie           = "selfie"
	PassportElementErrorSourceFile             = "file"
	PassportElementErrorSourceFiles            = "files"
	PassportElementErrorSourceTranslationFile  = "translation_file"
	PassportElementErrorSourceTranslationFiles = "translation_files"
	PassportElementErrorSourceUnspecified      = "unspecified"
)
