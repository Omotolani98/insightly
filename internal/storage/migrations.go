package storage

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&LogEntry{},
		&Summary{},
	)
}
