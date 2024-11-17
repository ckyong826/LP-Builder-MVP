package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt sql.NullTime   `json:"deletedAt,omitempty"`
}

// ScanRow implements the Scanner interface for a single row
func (u *User) ScanRow(row *sql.Row) error {
	return row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
}

// ScanRows implements the Scanner interface for multiple rows
func (u *User) ScanRows(rows *sql.Rows) error {
	return rows.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.DeletedAt,
	)
}

// TableName returns the database table name for the user model
func (User) TableName() string {
	return "users"
}


// Validate performs basic validation on the user model
func (u *User) Validate() error {
	if u.Name == "" {
		return ErrEmptyName
	}
	if u.Email == "" {
		return ErrEmptyEmail
	}
	// Add more validation as needed
	return nil
}

// BeforeCreate sets up timestamps before creation
func (u *User) BeforeCreate() {
	now := time.Now().UTC()
	u.CreatedAt = now
	u.UpdatedAt = now
}

// BeforeUpdate updates the updated_at timestamp
func (u *User) BeforeUpdate() {
	u.UpdatedAt = time.Now().UTC()
}

// Custom errors for validation
var (
	ErrEmptyName  = Error("name cannot be empty")
	ErrEmptyEmail = Error("email cannot be empty")
)

