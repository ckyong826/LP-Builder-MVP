package migrations

import (
	"backend/src/api/models"
	"backend/src/database"
)

func Migrate() {
    database.DB.AutoMigrate(&models.User{})
}
