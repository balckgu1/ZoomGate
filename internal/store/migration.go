package store

import (
	"zoomgate/internal/model"

	"golang.org/x/crypto/bcrypt"
)

// AutoMigrate runs GORM auto-migration for all registered entity models,
// creating or updating database tables as needed.
func (sqliteStore *SQLiteStore) AutoMigrate() error {
	return sqliteStore.db.AutoMigrate(
		&model.User{},
		&model.ProviderConfig{},
		&model.ModelConfig{},
		&model.AuditLog{},
		&model.SecurityPolicy{},
	)
}

// Seed creates the initial admin user if no users exist in the database.
// This should be called once during application startup.
func (sqliteStore *SQLiteStore) Seed(adminUsername, adminPassword string) error {
	var count int64
	sqliteStore.db.Model(&model.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := model.User{
		Username:     adminUsername,
		PasswordHash: string(passwordHash),
		Role:         model.RoleAdmin,
		RateLimit:    0, // unlimited for admin
	}
	return sqliteStore.db.Create(&admin).Error
}
