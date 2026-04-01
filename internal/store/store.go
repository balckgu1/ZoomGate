package store

import "zoomgate/internal/model"

type Store interface {
	UserStore
	ProviderStore
	AuditStore
	PolicyStore
	AutoMigrate() error
	Close() error
}

type UserStore interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByAPIKeyHash(hash string) (*model.User, error)
	ListUsers(page, pageSize int) ([]model.User, int64, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
}

type ProviderStore interface {
	CreateProvider(p *model.ProviderConfig) error
	GetProviderByID(id uint) (*model.ProviderConfig, error)
	GetProviderByName(name string) (*model.ProviderConfig, error)
	ListProviders() ([]model.ProviderConfig, error)
	ListEnabledProviders() ([]model.ProviderConfig, error)
	UpdateProvider(p *model.ProviderConfig) error
	DeleteProvider(id uint) error

	CreateModel(m *model.ModelConfig) error
	ListModelsByProvider(providerID uint) ([]model.ModelConfig, error)
	UpdateModel(m *model.ModelConfig) error
	DeleteModel(id uint) error
}

type AuditStore interface {
	CreateAuditLog(log *model.AuditLog) error
	BatchCreateAuditLogs(logs []model.AuditLog) error
	SearchAuditLogs(filter model.AuditFilter) ([]model.AuditLog, int64, error)
}

type PolicyStore interface {
	CreatePolicy(p *model.SecurityPolicy) error
	GetPolicyByID(id uint) (*model.SecurityPolicy, error)
	ListPolicies() ([]model.SecurityPolicy, error)
	ListEnabledPolicies(policyType model.PolicyType) ([]model.SecurityPolicy, error)
	UpdatePolicy(p *model.SecurityPolicy) error
	DeletePolicy(id uint) error
}
