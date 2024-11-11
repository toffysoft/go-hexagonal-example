package database

import (
	"github.com/toffysoft/go-hexagonal-example/internal/core/domain"
	"github.com/toffysoft/go-hexagonal-example/internal/infrastructure/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DBSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.Blog{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("host=localhost user=test password=test dbname=test_db port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.Blog{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
