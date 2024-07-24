package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/darrenjon/restaurant-ordering-system/internal/config"
	"github.com/darrenjon/restaurant-ordering-system/internal/logger"
	"github.com/darrenjon/restaurant-ordering-system/internal/models"
)

type Manager struct {
	db     *gorm.DB
	Config *config.DatabaseConfig
}

func NewManager(cfg *config.DatabaseConfig) (*Manager, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Taipei",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	// Initially set log level to Silent
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.GetGormLogger(gormlogger.Silent),
		NowFunc: func() time.Time {
			return time.Now().In(time.FixedZone("Asia/Taipei", 8*60*60))
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.InfoLogger.Println("Connected to database successfully")
	return &Manager{
		db:     db,
		Config: cfg,
	}, nil
}

func (m *Manager) GetDB() *gorm.DB {
	return m.db
}

func (m *Manager) SetLogMode(logMode gormlogger.LogLevel) {
	m.db.Logger = logger.GetGormLogger(logMode)
}

func (m *Manager) AutoMigrate() error {
	return m.db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.MenuItem{},
		&models.AddOn{},
		&models.Order{},
		&models.OrderDetail{},
		&models.SelectedAddOn{},
		&models.RestaurantInfo{},
	)
}
