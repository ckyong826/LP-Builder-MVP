package repositories

import (
	"backend/src/api/models"

	"gorm.io/gorm"
)

type UserRepository struct {
    Repository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{Repository{DB: db}}
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
    var users []models.User
    if err := r.DB.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

func (r *UserRepository) CreateUser(user models.User) error {
    return r.DB.Create(&user).Error
}
