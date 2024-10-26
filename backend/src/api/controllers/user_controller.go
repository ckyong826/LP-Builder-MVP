package controllers

import (
	"backend/src/api/models"
	"backend/src/api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(s *services.UserService) *UserController {
	return &UserController{userService: s}
}

func (ctrl *UserController) FindAll(c *gin.Context) {
	users, err := ctrl.userService.FindAll()
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
	}
	c.JSON(http.StatusOK, users)
}

func (ctrl *UserController) FindOneById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
	}
	user, err := ctrl.userService.FindOneById(uint(id))
	if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
	}
	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) Create(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}
	if err := ctrl.userService.Create(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
	}
	c.JSON(http.StatusCreated, user)
}

func (ctrl *UserController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
	}
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}
	user.ID = uint(id)
	if err := ctrl.userService.Update(user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
	}
	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
	}
	if err := ctrl.userService.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
