package models

import (
	"time"

	"gorm.io/gorm"
)

type Template struct {
  ID        uint           `gorm:"primaryKey" json:"id"`
  Name      string         `json:"name"`
  Type      string         `json:"type"`
  CreatedAt time.Time      `json:"createdAt"`
  UpdatedAt time.Time      `json:"updatedAt"`
  DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ConvertUrlToFile struct {
  URL string `json:"url" binding:"required"`
}