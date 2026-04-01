package store

import "zoomgate/internal/model"

// Store is the top-level data access interface that composes all sub-store interfaces.
type Store interface {
	UserStore
	ProviderStore
	AuditStore
	PolicyStore

	// AutoMigrate runs database schema auto-migration for all registered models.
	AutoMigrate() error

	// Close gracefully closes the underlying database connection.
	Close() error
}

// UserStore defines data access operations for the User entity.
type UserStore interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByAPIKeyHash(hash string) (*model.User, error)
	ListUsers(page, pageSize int) ([]model.User, int64, error)
	UpdateUser(user *model.User) error
	DeleteUser(id uint) error
}

// ProviderStore defines data access operations for ProviderConfig and ModelConfig entities.
type ProviderStore interface {
	CreateProvider(provider *model.ProviderConfig) error
	GetProviderByID(id uint) (*model.ProviderConfig, error)
	GetProviderByName(name string) (*model.ProviderConfig, error)
	ListProviders() ([]model.ProviderConfig, error)
	ListEnabledProviders() ([]model.ProviderConfig, error)
	UpdateProvider(provider *model.ProviderConfig) error
	DeleteProvider(id uint) error

	CreateModel(modelConfig *model.ModelConfig) error
	ListModelsByProvider(providerID uint) ([]model.ModelConfig, error)
	UpdateModel(modelConfig *model.ModelConfig) error
	DeleteModel(id uint) error
}

// AuditStore defines data access operations for audit log entries.
type AuditStore interface {
	CreateAuditLog(auditLog *model.AuditLog) error
	BatchCreateAuditLogs(auditLogs []model.AuditLog) error
	SearchAuditLogs(filter model.AuditFilter) ([]model.AuditLog, int64, error)
}

// PolicyStore defines data access operations for security policies.
type PolicyStore interface {
	CreatePolicy(policy *model.SecurityPolicy) error
	GetPolicyByID(id uint) (*model.SecurityPolicy, error)
	ListPolicies() ([]model.SecurityPolicy, error)
	ListEnabledPolicies(policyType model.PolicyType) ([]model.SecurityPolicy, error)
	UpdatePolicy(policy *model.SecurityPolicy) error
	DeletePolicy(id uint) error
}
