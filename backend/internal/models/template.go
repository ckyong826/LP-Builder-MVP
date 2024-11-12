package models

import (
	"time"

	"gorm.io/gorm"
)

type Template struct {
  ID          uint      `gorm:"primaryKey"`
  OriginalURL string    `gorm:"type:text;not null"`
  HTMLPath    string    `gorm:"type:text"`
  FilePaths   string    `gorm:"type:jsonb"` 
  Status      string    `gorm:"type:varchar(20)"`
  ErrorMessage string   `gorm:"type:text"`
  CreatedAt time.Time      `json:"createdAt"`
  UpdatedAt time.Time      `json:"updatedAt"`
  DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ConvertUrlToFile struct {
  URL string `json:"url" binding:"required"`
}