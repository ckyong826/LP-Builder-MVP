package models

import (
	"database/sql"
	"time"
)

type Template struct {
    ID           int64          `json:"id"`
    OriginalURL  string         `json:"original_url"`
    HTMLPath     string         `json:"html_path"`
    FilePaths    string         `json:"file_paths"`
    Status       string         `json:"status"`
    ErrorMessage sql.NullString `json:"error_message,omitempty"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    sql.NullTime   `json:"deleted_at,omitempty"`
}

type FileContent struct {
	HTML   string            `json:"html"`
	CSS    map[string]string `json:"css"`
	JS     map[string]string `json:"js"`
	Images map[string][]byte `json:"images"`
}


// ConvertUrlToFile represents the request payload for URL conversion
type ConvertUrlToFile struct {
	URL string `json:"url" binding:"required"`
}

// TableName returns the database table name for the template model
func (Template) TableName() string {
	return "templates"
}

// Template status constants
const (
	StatusPending   = "pending"
	StatusComplete  = "complete"
	StatusFailed    = "failed"
	StatusProgress  = "in_progress"
)


// Validate performs basic validation on the template model
func (t *Template) Validate() error {
	if t.OriginalURL == "" {
		return ErrEmptyOriginalURL
	}
	if t.Status == "" {
		t.Status = StatusPending
	}
	if !isValidStatus(t.Status) {
		return ErrInvalidStatus
	}
	return nil
}

// isValidStatus checks if the status is valid
func isValidStatus(status string) bool {
	validStatuses := map[string]bool{
		StatusPending:   true,
		StatusComplete:  true,
		StatusFailed:    true,
		StatusProgress:  true,
	}
	return validStatuses[status]
}

// BeforeCreate sets up timestamps before creation
func (t *Template) BeforeCreate() {
	now := time.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now
	if t.Status == "" {
		t.Status = StatusPending
	}
}

// BeforeUpdate updates the updated_at timestamp
func (t *Template) BeforeUpdate() {
	t.UpdatedAt = time.Now().UTC()
}

// SetError sets the error message and updates the status
func (t *Template) SetError(err error) {
	t.Status = StatusFailed
	t.ErrorMessage = sql.NullString{
		String: err.Error(),
		Valid:  true,
	}
}

// SetComplete marks the template as complete
func (t *Template) SetComplete() {
	t.Status = StatusComplete
	t.ErrorMessage = sql.NullString{Valid: false}
}

// Custom errors for validation
var (
	ErrEmptyOriginalURL = Error("original URL cannot be empty")
	ErrInvalidStatus    = Error("invalid status")
)



