package services

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) FindAll(ctx context.Context, page, pageSize int, orderBy, sort string) ([]models.User, int64, error) {
	// Count total records
	var total int64
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count error: %w", err)
	}

	// Build query with pagination
	query := `SELECT id, name, email, created_at, updated_at, deleted_at 
			  FROM users 
			  WHERE deleted_at IS NULL`

	if orderBy != "" {
		direction := "ASC"
		if strings.ToLower(sort) == "desc" {
			direction = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", orderBy, direction)
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
	}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (s *UserService) FindOneById(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, email, created_at, updated_at, deleted_at 
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	return user, nil
}

func (s *UserService) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := s.db.QueryRowContext(ctx,
		`INSERT INTO users (name, email, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id`,
		user.Name, user.Email, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("create error: %w", err)
	}

	return nil
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	result, err := s.db.ExecContext(ctx,
		`UPDATE users 
		 SET name = $1, email = $2, updated_at = $3 
		 WHERE id = $4 AND deleted_at IS NULL`,
		user.Name, user.Email, user.UpdatedAt, user.ID,
	)

	if err != nil {
		return fmt.Errorf("update error: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE users 
		 SET deleted_at = $1 
		 WHERE id = $2 AND deleted_at IS NULL`,
		time.Now(), id,
	)

	if err != nil {
		return fmt.Errorf("delete error: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected error: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// WithTx executes operations within a transaction
func (s *UserService) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit()
}