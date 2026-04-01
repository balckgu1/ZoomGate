package store

import (
	"time"

	"zoomgate/internal/model"

	"gorm.io/gorm"
)

// --- UserStore ---

func (s *SQLiteStore) CreateUser(user *model.User) error {
	return s.db.Create(user).Error
}

func (s *SQLiteStore) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := s.db.First(&user, id).Error
	return &user, err
}

func (s *SQLiteStore) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := s.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (s *SQLiteStore) GetUserByAPIKeyHash(hash string) (*model.User, error) {
	var user model.User
	err := s.db.Where("api_key_hash = ?", hash).First(&user).Error
	return &user, err
}

func (s *SQLiteStore) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	s.db.Model(&model.User{}).Count(&total)
	err := s.db.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("id DESC").Find(&users).Error
	return users, total, err
}

func (s *SQLiteStore) UpdateUser(user *model.User) error {
	return s.db.Save(user).Error
}

func (s *SQLiteStore) DeleteUser(id uint) error {
	return s.db.Delete(&model.User{}, id).Error
}

// --- ProviderStore ---

func (s *SQLiteStore) CreateProvider(p *model.ProviderConfig) error {
	return s.db.Create(p).Error
}

func (s *SQLiteStore) GetProviderByID(id uint) (*model.ProviderConfig, error) {
	var p model.ProviderConfig
	err := s.db.Preload("Models").First(&p, id).Error
	return &p, err
}

func (s *SQLiteStore) GetProviderByName(name string) (*model.ProviderConfig, error) {
	var p model.ProviderConfig
	err := s.db.Preload("Models").Where("name = ?", name).First(&p).Error
	return &p, err
}

func (s *SQLiteStore) ListProviders() ([]model.ProviderConfig, error) {
	var providers []model.ProviderConfig
	err := s.db.Preload("Models").Order("priority DESC").Find(&providers).Error
	return providers, err
}

func (s *SQLiteStore) ListEnabledProviders() ([]model.ProviderConfig, error) {
	var providers []model.ProviderConfig
	err := s.db.Preload("Models", "enabled = ?", true).
		Where("enabled = ?", true).
		Order("priority DESC").Find(&providers).Error
	return providers, err
}

func (s *SQLiteStore) UpdateProvider(p *model.ProviderConfig) error {
	return s.db.Save(p).Error
}

func (s *SQLiteStore) DeleteProvider(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("provider_id = ?", id).Delete(&model.ModelConfig{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ProviderConfig{}, id).Error
	})
}

func (s *SQLiteStore) CreateModel(m *model.ModelConfig) error {
	return s.db.Create(m).Error
}

func (s *SQLiteStore) ListModelsByProvider(providerID uint) ([]model.ModelConfig, error) {
	var models []model.ModelConfig
	err := s.db.Where("provider_id = ?", providerID).Find(&models).Error
	return models, err
}

func (s *SQLiteStore) UpdateModel(m *model.ModelConfig) error {
	return s.db.Save(m).Error
}

func (s *SQLiteStore) DeleteModel(id uint) error {
	return s.db.Delete(&model.ModelConfig{}, id).Error
}

// --- AuditStore ---

func (s *SQLiteStore) CreateAuditLog(log *model.AuditLog) error {
	return s.db.Create(log).Error
}

func (s *SQLiteStore) BatchCreateAuditLogs(logs []model.AuditLog) error {
	if len(logs) == 0 {
		return nil
	}
	return s.db.CreateInBatches(logs, 100).Error
}

func (s *SQLiteStore) SearchAuditLogs(filter model.AuditFilter) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := s.db.Model(&model.AuditLog{})

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
		if t, err := time.Parse(time.RFC3339, filter.StartTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if filter.EndTime != "" {
		if t, err := time.Parse(time.RFC3339, filter.EndTime); err == nil {
			query = query.Where("created_at <= ?", t)
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
		Order("id DESC").Find(&logs).Error
	return logs, total, err
}

// --- PolicyStore ---

func (s *SQLiteStore) CreatePolicy(p *model.SecurityPolicy) error {
	return s.db.Create(p).Error
}

func (s *SQLiteStore) GetPolicyByID(id uint) (*model.SecurityPolicy, error) {
	var p model.SecurityPolicy
	err := s.db.First(&p, id).Error
	return &p, err
}

func (s *SQLiteStore) ListPolicies() ([]model.SecurityPolicy, error) {
	var policies []model.SecurityPolicy
	err := s.db.Order("id ASC").Find(&policies).Error
	return policies, err
}

func (s *SQLiteStore) ListEnabledPolicies(policyType model.PolicyType) ([]model.SecurityPolicy, error) {
	var policies []model.SecurityPolicy
	err := s.db.Where("type = ? AND enabled = ?", policyType, true).Find(&policies).Error
	return policies, err
}

func (s *SQLiteStore) UpdatePolicy(p *model.SecurityPolicy) error {
	return s.db.Save(p).Error
}

func (s *SQLiteStore) DeletePolicy(id uint) error {
	return s.db.Delete(&model.SecurityPolicy{}, id).Error
}
