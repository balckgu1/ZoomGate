package model

import "time"

// PolicyType categorizes the kind of security policy.
type PolicyType string

const (
	// PolicyPromptFilter detects and blocks prompt injection attacks.
	PolicyPromptFilter PolicyType = "prompt_filter"
	// PolicyRateLimit controls request rate per user.
	PolicyRateLimit PolicyType = "rate_limit"
	// PolicyDataMask masks sensitive PII data in requests and responses.
	PolicyDataMask PolicyType = "data_mask"
)

// SecurityPolicy stores a configurable security rule managed via the admin UI.
type SecurityPolicy struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	Name      string     `gorm:"uniqueIndex;not null;size:64" json:"name"`
	Type      PolicyType `gorm:"not null;size:32" json:"type"`
	Config    string     `gorm:"type:text" json:"config"`
	Enabled   bool       `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
