package handlers

import (
	"net/http"

	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/services"
	"github.com/gin-gonic/gin"
)

// ProjectHandler handles HTTP requests for projects
type ProjectHandler struct {
	projectService *services.ProjectService
}

// NewProjectHandler creates a new ProjectHandler instance
func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// CreateProject handles POST /v1/projects
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	var req models.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	project, err := h.projectService.CreateProject(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create project",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": project,
	})
}

// GetProject handles GET /v1/projects/:id
func (h *ProjectHandler) GetProject(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		RespondWithError(c, 401, "Unauthorized")
		return
	}

	projectID := c.Param("id")

	project, err := h.projectService.GetProject(c.Request.Context(), userID.(string), projectID)
	if err != nil {
		RespondWithError(c, 404, "Project not found")
		return
	}

	RespondWithSuccess(c, 200, project, "Project retrieved successfully")
}

// GetProjects handles GET /v1/projects
func (h *ProjectHandler) GetProjects(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	projects, err := h.projectService.GetUserProjects(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get projects",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  projects,
		"total": len(projects),
	})
}

// UpdateProject handles PUT /v1/projects/:id
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		RespondWithError(c, 401, "Unauthorized")
		return
	}

	projectID := c.Param("id")

	var req models.ProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondWithError(c, 400, "Invalid request format")
		return
	}

	project, err := h.projectService.UpdateProject(c.Request.Context(), userID.(string), projectID, &req)
	if err != nil {
		RespondWithError(c, 400, "Failed to update project")
		return
	}

	RespondWithSuccess(c, 200, project, "Project updated successfully")
}

// DeleteProject handles DELETE /v1/projects/:id
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		RespondWithError(c, 401, "Unauthorized")
		return
	}

	projectID := c.Param("id")

	err := h.projectService.DeleteProject(c.Request.Context(), userID.(string), projectID)
	if err != nil {
		RespondWithError(c, 500, "Failed to delete project")
		return
	}

	RespondWithSuccess(c, 200, nil, "Project deleted successfully")
}
