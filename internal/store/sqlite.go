package store

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SQLiteStore implements the Store interface using SQLite with GORM.
type SQLiteStore struct {
	db *gorm.DB
}

// NewSQLiteStore opens a SQLite database at the given path with WAL mode enabled
// and returns a ready-to-use SQLiteStore instance.
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	database, err := gorm.Open(sqlite.Open(dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// SQLite best practice: limit to 1 writer connection under WAL mode
	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying db: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)

	return &SQLiteStore{db: database}, nil
}

// Close gracefully closes the underlying database connection.
func (sqliteStore *SQLiteStore) Close() error {
	sqlDB, err := sqliteStore.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GormDB returns the underlying *gorm.DB instance for advanced queries.
func (sqliteStore *SQLiteStore) GormDB() *gorm.DB {
	return sqliteStore.db
}
