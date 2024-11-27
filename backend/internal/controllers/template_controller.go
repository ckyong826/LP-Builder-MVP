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

//////////////
// GET Methods
//////////////

func (ctrl *TemplateController) FindAll(c *gin.Context) {
	var query models.PaginationQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pagination parameters"})
		return
	}

	templates, total, err := ctrl.templateService.FindAll(
		c.Request.Context(),
		query.Page,
		query.PageSize,
		query.OrderBy,
		query.Sort,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  templates,
		"total": total,
		"page":  query.Page,
		"size":  query.PageSize,
	})
}

func (ctrl *TemplateController) FindOneById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	template, err := ctrl.templateService.FindOneById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	c.JSON(http.StatusOK, template)
}

///////////////
// POST Methods
///////////////

func (ctrl *TemplateController) Create(c *gin.Context) {
	var template models.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.templateService.Create(c.Request.Context(), &template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

func (ctrl *TemplateController) ConvertUrlToFile(c *gin.Context) {
	var request models.ConvertUrlToFile
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	template := &models.Template{}
	if err := ctrl.templateService.ConvertUrlToFile(c.Request.Context(), template, request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert URL", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "URL converted successfully",
		"conversion":  template,
		"html_path":   template.HTMLPath,
		"file_paths":  template.FilePaths,
	})
}

///////////////
// PUT Methods
///////////////

func (ctrl *TemplateController) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var template models.Template
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template.ID = id
	if err := ctrl.templateService.Update(c.Request.Context(), &template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update template"})
		return
	}

	c.JSON(http.StatusOK, template)
}

//////////////////
// DELETE Methods
//////////////////

func (ctrl *TemplateController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := ctrl.templateService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}