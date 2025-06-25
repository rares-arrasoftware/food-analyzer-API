package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type dbSqlite[T any] struct {
	db *gorm.DB
}

// newSqlite creates a new SQLite-backed generic database.
func newSqlite[T any](dsn string, model T) (Database[T], error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// AutoMigrate ensures the table exists for T
	if err := db.AutoMigrate(&model); err != nil {
		return nil, fmt.Errorf("failed to migrate model: %w", err)
	}

	return &dbSqlite[T]{db}, nil
}

func (d *dbSqlite[T]) Create(item T) error {
	return d.db.Create(&item).Error
}

func (d *dbSqlite[T]) Update(item T) error {
	return d.db.Save(&item).Error
}

func (d *dbSqlite[T]) GetByID(id uint) (*T, error) {
	var item T
	if err := d.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (d *dbSqlite[T]) Delete(id uint) error {
	var item T
	return d.db.Delete(&item, id).Error
}

func (d *dbSqlite[T]) GetByField(field string, value any) (*T, error) {
	var item T
	// Note: If no record is found, GORM will log "record not found" internally.
	// This is expected behavior for lookups like login or register checks.
	// If needed, GORM's logger can be configured later to suppress these logs
	if err := d.db.Where("? = ?", gorm.Expr(field), value).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
