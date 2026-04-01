package model

import "time"

type PolicyType string

const (
	PolicyPromptFilter PolicyType = "prompt_filter"
	PolicyRateLimit    PolicyType = "rate_limit"
	PolicyDataMask     PolicyType = "data_mask"
)

type SecurityPolicy struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	Name      string     `gorm:"uniqueIndex;not null;size:64" json:"name"`
	Type      PolicyType `gorm:"not null;size:32" json:"type"`
	Config    string     `gorm:"type:text" json:"config"`
	Enabled   bool       `gorm:"default:true" json:"enabled"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
