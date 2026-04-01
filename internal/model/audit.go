package model

import "time"

type AuditLog struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	RequestID    string    `gorm:"size:64;index" json:"request_id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	Username     string    `gorm:"size:64" json:"username"`
	Model        string    `gorm:"size:128" json:"model"`
	Provider     string    `gorm:"size:64" json:"provider"`
	InputTokens  int       `gorm:"default:0" json:"input_tokens"`
	OutputTokens int       `gorm:"default:0" json:"output_tokens"`
	LatencyMs    int64     `gorm:"default:0" json:"latency_ms"`
	StatusCode   int       `json:"status_code"`
	ErrorMessage string    `gorm:"size:1024" json:"error_message,omitempty"`
	ClientIP     string    `gorm:"size:45" json:"client_ip"`
	RoutedBy     string    `gorm:"size:32" json:"routed_by"`
	CreatedAt    time.Time `gorm:"index" json:"created_at"`
}

type AuditFilter struct {
	UserID    uint   `form:"user_id"`
	Username  string `form:"username"`
	Model     string `form:"model"`
	Provider  string `form:"provider"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
	Status    int    `form:"status"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}
