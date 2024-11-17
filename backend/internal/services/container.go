package services

import (
	"database/sql"
)

type ServiceContainer struct {
	UserService     *UserService
	TemplateService *TemplateService
}

func NewServiceContainer(db *sql.DB) *ServiceContainer {
	return &ServiceContainer{
		UserService:     NewUserService(db),
		TemplateService: NewTemplateService(db),
	}
}