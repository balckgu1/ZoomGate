package store

import (
	"zoomgate/internal/model"

	"golang.org/x/crypto/bcrypt"
)

func (s *SQLiteStore) AutoMigrate() error {
	return s.db.AutoMigrate(
		&model.User{},
		&model.ProviderConfig{},
		&model.ModelConfig{},
		&model.AuditLog{},
		&model.SecurityPolicy{},
	)
}

func (s *SQLiteStore) Seed(adminUsername, adminPassword string) error {
	var count int64
	s.db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := model.User{
		Username:     adminUsername,
		PasswordHash: string(hash),
		Role:         model.RoleAdmin,
		RateLimit:    0, // unlimited for admin
	}
	return s.db.Create(&admin).Error
}
