package model

import "time"

type ProviderType string

const (
	ProviderOpenAI    ProviderType = "openai"
	ProviderAnthropic ProviderType = "anthropic"
	ProviderGemini    ProviderType = "gemini"
	ProviderDeepSeek  ProviderType = "deepseek"
	ProviderQwen      ProviderType = "qwen"
	ProviderOllama    ProviderType = "ollama"
)

type ProviderConfig struct {
	ID             uint          `gorm:"primarykey" json:"id"`
	Name           string        `gorm:"uniqueIndex;not null;size:64" json:"name"`
	Type           ProviderType  `gorm:"not null;size:32" json:"type"`
	APIKeyEnc      string        `gorm:"size:512" json:"-"`
	BaseURL        string        `gorm:"not null;size:256" json:"base_url"`
	Enabled        bool          `gorm:"default:true" json:"enabled"`
	Priority       int           `gorm:"default:0" json:"priority"`
	Weight         int           `gorm:"default:1" json:"weight"`
	HealthCheckURL string        `gorm:"size:256" json:"health_check_url"`
	Models         []ModelConfig `gorm:"foreignKey:ProviderID" json:"models,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type ModelConfig struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	ProviderID   uint      `gorm:"index;not null" json:"provider_id"`
	ModelName    string    `gorm:"not null;size:128" json:"model_name"`
	DisplayName  string    `gorm:"size:128" json:"display_name"`
	CostInput1K  float64   `gorm:"default:0" json:"cost_input_1k"`
	CostOutput1K float64   `gorm:"default:0" json:"cost_output_1k"`
	MaxContext   int       `gorm:"default:4096" json:"max_context"`
	Enabled      bool      `gorm:"default:true" json:"enabled"`
	CreatedAt    time.Time `json:"created_at"`
}
