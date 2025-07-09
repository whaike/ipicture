package models

import (
	"gorm.io/gorm"
	"time"
)

type File struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	Filename  string         `json:"filename"`
	Path      string         `json:"path"`
	Type      string         `json:"type"` // image/video
	Size      int64          `json:"size"`
	MD5       string         `gorm:"index" json:"md5"`
	ThumbPath string         `json:"thumb_path"`
	Tags      string         `json:"tags"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
} 