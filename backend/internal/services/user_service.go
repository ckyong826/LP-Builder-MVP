package services

import (
	"backend/internal/models"
	"backend/internal/repositories"

	"gorm.io/gorm"
)

type UserService struct {
    repo *repositories.Repository[models.User]
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{
        repo: repositories.NewRepository[models.User](db),
    }
}

func (s *UserService) FindAll() ([]models.User, error) {
    return s.repo.FindAll()
}

func (s *UserService) FindOneById(id uint) (models.User, error) {
    user, err := s.repo.FindOneByID(id)
    if err != nil {
        return models.User{}, err
    }
    return *user, nil
}

func (s *UserService) Create(user models.User) error {
    return s.repo.Create(&user)
}

func (s *UserService) Update(user models.User) error {
    return s.repo.Update(&user)
}

func (s *UserService) Delete(id uint) error {
    return s.repo.Delete(id)
}
