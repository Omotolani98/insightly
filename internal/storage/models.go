package storage

import (
	"time"

	"gorm.io/datatypes"
)

type LogEntry struct {
	ID        uint      `gorm:"primaryKey"`
	Stream    string    `gorm:"index"`
	Timestamp time.Time `gorm:"index"`
	Level     string    `gorm:"size:10"`
	Message   string    `gorm:"type:text"`
	Metadata  datatypes.JSON
	CreatedAt time.Time
}

type Summary struct {
	ID          uint      `gorm:"primaryKey"`
	Stream      string    `gorm:"index"`
	WindowStart time.Time `gorm:"index"`
	WindowEnd   time.Time `gorm:"index"`
	Text        string    `gorm:"type:text"`
	CreatedAt   time.Time
}
