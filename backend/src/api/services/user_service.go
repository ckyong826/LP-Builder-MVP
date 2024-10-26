package services

import (
	"backend/src/api/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// FindAll retrieves all users from the database.
func (s *UserService) FindAll() ([]models.User, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// FindOneById retrieves a user by ID.
func (s *UserService) FindOneById(id uint) (models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

// Create adds a new user to the database.
func (s *UserService) Create(user models.User) error {
	if err := s.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// Update modifies an existing user.
func (s *UserService) Update(user models.User) error {
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a user by ID.
func (s *UserService) Delete(id uint) error {
	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}
