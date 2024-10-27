package services

import (
    "gorm.io/gorm"
)

type ServiceContainer struct {
    UserService *UserService

}

func NewServiceContainer(db *gorm.DB) *ServiceContainer {
    return &ServiceContainer{
        UserService: NewUserService(db),

    }
}
