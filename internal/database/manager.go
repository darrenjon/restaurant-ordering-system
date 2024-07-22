package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/darrenjon/restaurant-ordering-system/internal/config"
	"github.com/darrenjon/restaurant-ordering-system/internal/models"
)

type Manager struct {
	db *gorm.DB
}

func NewManager(cfg *config.DatabaseConfig) (*Manager, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Taipei",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().In(time.FixedZone("Asia/Taipei", 8*60*60))
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Manager{db: db}, nil
}

func (m *Manager) GetDB() *gorm.DB {
	return m.db
}

func (m *Manager) AutoMigrate() error {
	return m.db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.MenuItem{},
		&models.AddOn{},
		&models.MenuItemAddOn{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderItemAddOn{},
		&models.RestaurantInfo{},
	)
}
