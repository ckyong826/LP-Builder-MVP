package services

import (
	"backend/internal/models"
	"log"

	"gorm.io/gorm"
)

type ServiceContainer struct {
    UserService *UserService
    TemplateService *TemplateService
}

func AutoMigrate(db *gorm.DB){
    modelsToMigrate := []interface{}{
        &models.User{},
        &models.Template{},
    }

    for _, model := range modelsToMigrate {
        if err := db.AutoMigrate(model); err != nil {
            log.Fatalf("failed to migrate model %T: %v", model, err)
        }
    }
}

func NewServiceContainer(db *gorm.DB) *ServiceContainer {

    AutoMigrate(db)
   
    return &ServiceContainer{
        UserService: NewUserService(db),
        TemplateService : NewTemplateService(db),
    }
}
