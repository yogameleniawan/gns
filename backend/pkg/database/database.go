package database

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/base-go/backend/pkg/config"
)

type Database interface {
	GetDB() *gorm.DB
	Close() error
	Ping() error
}

type database struct {
	db *gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase() (Database, error) {
	cfg := config.GetConfig()
	dbCfg := cfg.Database

	// Build PostgreSQL DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Name,
		dbCfg.SSLMode,
	)

	// Open database connection with GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbCfg.ConnMaxLifetime) * time.Minute)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.Info("Successfully connected to PostgreSQL database")

	return &database{
		db: db,
	}, nil
}

// GetDB returns the GORM DB instance
func (d *database) GetDB() *gorm.DB {
	return d.db
}

// Close closes the database connection
func (d *database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// Ping checks if the database connection is alive
func (d *database) Ping() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Ping()
}
