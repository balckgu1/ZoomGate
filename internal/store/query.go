package store

import (
	"time"

	"zoomgate/internal/model"

	"gorm.io/gorm"
)

// ==================== UserStore Implementation ====================

// CreateUser inserts a new user record into the database.
func (sqliteStore *SQLiteStore) CreateUser(user *model.User) error {
	return sqliteStore.db.Create(user).Error
}

// GetUserByID retrieves a user by their primary key ID.
func (sqliteStore *SQLiteStore) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := sqliteStore.db.First(&user, id).Error
	return &user, err
}

// GetUserByUsername retrieves a user by their unique username.
func (sqliteStore *SQLiteStore) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := sqliteStore.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetUserByAPIKeyHash retrieves a user by the SHA-256 hash of their API key.
func (sqliteStore *SQLiteStore) GetUserByAPIKeyHash(hash string) (*model.User, error) {
	var user model.User
	err := sqliteStore.db.Where("api_key_hash = ?", hash).First(&user).Error
	return &user, err
}

// ListUsers returns a paginated list of users, ordered by ID descending.
func (sqliteStore *SQLiteStore) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	sqliteStore.db.Model(&model.User{}).Count(&total)
	err := sqliteStore.db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id DESC").Find(&users).Error
	return users, total, err
}

// UpdateUser saves the modified user record to the database.
func (sqliteStore *SQLiteStore) UpdateUser(user *model.User) error {
	return sqliteStore.db.Save(user).Error
}

// DeleteUser removes a user record by their primary key ID.
func (sqliteStore *SQLiteStore) DeleteUser(id uint) error {
	return sqliteStore.db.Delete(&model.User{}, id).Error
}

// ==================== ProviderStore Implementation ====================

// CreateProvider inserts a new provider configuration record.
func (sqliteStore *SQLiteStore) CreateProvider(provider *model.ProviderConfig) error {
	return sqliteStore.db.Create(provider).Error
}

// GetProviderByID retrieves a provider config by ID, including its associated models.
func (sqliteStore *SQLiteStore) GetProviderByID(id uint) (*model.ProviderConfig, error) {
	var provider model.ProviderConfig
	err := sqliteStore.db.Preload("Models").First(&provider, id).Error
	return &provider, err
}

// GetProviderByName retrieves a provider config by unique name, including its models.
func (sqliteStore *SQLiteStore) GetProviderByName(name string) (*model.ProviderConfig, error) {
	var provider model.ProviderConfig
	err := sqliteStore.db.Preload("Models").Where("name = ?", name).First(&provider).Error
	return &provider, err
}

// ListProviders returns all provider configs ordered by priority descending.
func (sqliteStore *SQLiteStore) ListProviders() ([]model.ProviderConfig, error) {
	var providers []model.ProviderConfig
	err := sqliteStore.db.Preload("Models").Order("priority DESC").Find(&providers).Error
	return providers, err
}

// ListEnabledProviders returns only enabled providers with their enabled models.
func (sqliteStore *SQLiteStore) ListEnabledProviders() ([]model.ProviderConfig, error) {
	var providers []model.ProviderConfig
	err := sqliteStore.db.Preload("Models", "enabled = ?", true).
		Where("enabled = ?", true).
		Order("priority DESC").Find(&providers).Error
	return providers, err
}

// UpdateProvider saves the modified provider config to the database.
func (sqliteStore *SQLiteStore) UpdateProvider(provider *model.ProviderConfig) error {
	return sqliteStore.db.Save(provider).Error
}

// DeleteProvider removes a provider and all its associated model configs in a transaction.
func (sqliteStore *SQLiteStore) DeleteProvider(id uint) error {
	return sqliteStore.db.Transaction(func(transaction *gorm.DB) error {
		if err := transaction.Where("provider_id = ?", id).Delete(&model.ModelConfig{}).Error; err != nil {
			return err
		}
		return transaction.Delete(&model.ProviderConfig{}, id).Error
	})
}

// CreateModel inserts a new model configuration record.
func (sqliteStore *SQLiteStore) CreateModel(modelConfig *model.ModelConfig) error {
	return sqliteStore.db.Create(modelConfig).Error
}

// ListModelsByProvider returns all model configs belonging to a given provider.
func (sqliteStore *SQLiteStore) ListModelsByProvider(providerID uint) ([]model.ModelConfig, error) {
	var models []model.ModelConfig
	err := sqliteStore.db.Where("provider_id = ?", providerID).Find(&models).Error
	return models, err
}

// UpdateModel saves the modified model config to the database.
func (sqliteStore *SQLiteStore) UpdateModel(modelConfig *model.ModelConfig) error {
	return sqliteStore.db.Save(modelConfig).Error
}

// DeleteModel removes a model config by its primary key ID.
func (sqliteStore *SQLiteStore) DeleteModel(id uint) error {
	return sqliteStore.db.Delete(&model.ModelConfig{}, id).Error
}

// ==================== AuditStore Implementation ====================

// CreateAuditLog inserts a single audit log entry.
func (sqliteStore *SQLiteStore) CreateAuditLog(auditLog *model.AuditLog) error {
	return sqliteStore.db.Create(auditLog).Error
}

// BatchCreateAuditLogs inserts multiple audit log entries in batches of 100.
func (sqliteStore *SQLiteStore) BatchCreateAuditLogs(auditLogs []model.AuditLog) error {
	if len(auditLogs) == 0 {
		return nil
	}
	return sqliteStore.db.CreateInBatches(auditLogs, 100).Error
}

// SearchAuditLogs retrieves audit logs matching the provided filter criteria with pagination.
func (sqliteStore *SQLiteStore) SearchAuditLogs(filter model.AuditFilter) ([]model.AuditLog, int64, error) {
	var auditLogs []model.AuditLog
	var total int64

	query := sqliteStore.db.Model(&model.AuditLog{})

	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.Model != "" {
		query = query.Where("model = ?", filter.Model)
	}
	if filter.Provider != "" {
		query = query.Where("provider = ?", filter.Provider)
	}
	if filter.Status > 0 {
		query = query.Where("status_code = ?", filter.Status)
	}
	if filter.StartTime != "" {
		if parsedTime, err := time.Parse(time.RFC3339, filter.StartTime); err == nil {
			query = query.Where("created_at >= ?", parsedTime)
		}
	}
	if filter.EndTime != "" {
		if parsedTime, err := time.Parse(time.RFC3339, filter.EndTime); err == nil {
			query = query.Where("created_at <= ?", parsedTime)
		}
	}

	query.Count(&total)

	page := filter.Page
	pageSize := filter.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id DESC").Find(&auditLogs).Error
	return auditLogs, total, err
}

// ==================== PolicyStore Implementation ====================

// CreatePolicy inserts a new security policy record.
func (sqliteStore *SQLiteStore) CreatePolicy(policy *model.SecurityPolicy) error {
	return sqliteStore.db.Create(policy).Error
}

// GetPolicyByID retrieves a security policy by its primary key ID.
func (sqliteStore *SQLiteStore) GetPolicyByID(id uint) (*model.SecurityPolicy, error) {
	var policy model.SecurityPolicy
	err := sqliteStore.db.First(&policy, id).Error
	return &policy, err
}

// ListPolicies returns all security policies ordered by ID ascending.
func (sqliteStore *SQLiteStore) ListPolicies() ([]model.SecurityPolicy, error) {
	var policies []model.SecurityPolicy
	err := sqliteStore.db.Order("id ASC").Find(&policies).Error
	return policies, err
}

// ListEnabledPolicies returns all enabled policies of a given type.
func (sqliteStore *SQLiteStore) ListEnabledPolicies(policyType model.PolicyType) ([]model.SecurityPolicy, error) {
	var policies []model.SecurityPolicy
	err := sqliteStore.db.Where("type = ? AND enabled = ?", policyType, true).Find(&policies).Error
	return policies, err
}

// UpdatePolicy saves the modified security policy to the database.
func (sqliteStore *SQLiteStore) UpdatePolicy(policy *model.SecurityPolicy) error {
	return sqliteStore.db.Save(policy).Error
}

// DeletePolicy removes a security policy by its primary key ID.
func (sqliteStore *SQLiteStore) DeletePolicy(id uint) error {
	return sqliteStore.db.Delete(&model.SecurityPolicy{}, id).Error
}
