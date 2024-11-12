package controllers

import (
	"backend/internal/models"
	"backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TemplateController struct {
	templateService *services.TemplateService
}

func NewTemplateController(s *services.TemplateService) *TemplateController {
	return &TemplateController{templateService: s}
}


//////////////////////////////////////////////////
///////////////////// GET ////////////////////////
//////////////////////////////////////////////////
func (ctrl *TemplateController) FindAll(c *gin.Context) {
	templates, err := ctrl.templateService.FindAll()
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve templates"})
			return
	}
	c.JSON(http.StatusOK, templates)
}

func (ctrl *TemplateController) FindOneById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
	}
	template, err := ctrl.templateService.FindOneById(uint(id))
	if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "template not found"})
			return
	}
	c.JSON(http.StatusOK, template)
}

//////////////////////////////////////////////////
//////////////////// POST ////////////////////////
//////////////////////////////////////////////////
func (ctrl *TemplateController) Create(c *gin.Context) {
	var template models.Template
	if err := c.ShouldBindJSON(&template); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}
	if err := ctrl.templateService.Create(&template); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template"})
			return
	}
	c.JSON(http.StatusCreated, template)
}

func (ctrl *TemplateController) ConvertUrlToFile(c *gin.Context) {
	var request struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	conversion := models.Template{}

	if err := ctrl.templateService.ConvertUrlToFile(&conversion, request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert URL", "details": err.Error()})
		return
	}

	// Respond with success and metadata
	c.JSON(http.StatusOK, gin.H{
		"message":     "URL converted successfully",
		"conversion":  conversion,
		"html_path":   conversion.HTMLPath,
		"file_paths":  conversion.FilePaths,
	})
}

//////////////////////////////////////////////////
/////////////////// PATCH ////////////////////////
//////////////////////////////////////////////////
func (ctrl *TemplateController) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
	}
	var template models.Template
	if err := c.ShouldBindJSON(&template); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
	}
	template.ID = uint(id)
	if err := ctrl.templateService.Update(template); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template"})
			return
	}
	c.JSON(http.StatusOK, template)
}

//////////////////////////////////////////////////
////////////////// DELETE ////////////////////////
//////////////////////////////////////////////////
func (ctrl *TemplateController) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
	}
	if err := ctrl.templateService.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template"})
			return
	}
	c.JSON(http.StatusOK, gin.H{"message": "template deleted successfully"})
}
