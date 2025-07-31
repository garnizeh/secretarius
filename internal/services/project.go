package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/garnizeh/englog/internal/database"
	"github.com/garnizeh/englog/internal/logging"
	"github.com/garnizeh/englog/internal/models"
	"github.com/garnizeh/englog/internal/store"
	"github.com/google/uuid"
)

// ProjectService handles all business logic for projects
type ProjectService struct {
	db     *database.DB
	logger *logging.Logger
}

// NewProjectService creates a new ProjectService instance
func NewProjectService(db *database.DB, logger *logging.Logger) *ProjectService {
	return &ProjectService{
		db:     db,
		logger: logger.WithComponent("project_service"),
	}
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(ctx context.Context, userID string, req *models.ProjectRequest) (*models.Project, error) {
	// Validate request
	if err := s.validateProjectRequest(req); err != nil {
		s.logger.LogError(ctx, err, "Project validation failed", "user_id", userID, "project_name", req.Name)
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Creating new project", "user_id", userID, "project_name", req.Name, "status", req.Status)

	var project *models.Project

	// Start write transaction to create project
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		sqlcProject, err := qtx.CreateProject(ctx, store.CreateProjectParams{
			Name:        req.Name,
			Description: stringToPgText(req.Description),
			Color:       stringToPgText(&req.Color),
			Status:      stringToPgText((*string)(&req.Status)),
			StartDate:   timeToPgDate(req.StartDate),
			EndDate:     timeToPgDate(req.EndDate),
			CreatedBy:   userUUID,
			IsDefault:   boolToPgBool(req.IsDefault),
		})
		if err != nil {
			s.logger.LogError(ctx, err, "Database error creating project", "user_id", userID, "project_name", req.Name)
			return fmt.Errorf("failed to create project: %w", err)
		}

		project = s.sqlcToModel(sqlcProject)
		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Transaction failed for project creation", "user_id", userID, "project_name", req.Name)
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	s.logger.Info("Project created successfully", "user_id", userID, "project_id", project.ID, "project_name", project.Name)
	return project, nil
}

// GetProject retrieves a single project by ID
func (s *ProjectService) GetProject(ctx context.Context, userID, projectID string) (*models.Project, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetProject", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid project ID format", "user_id", userID, "project_id", projectID)
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	s.logger.Info("Getting project", "user_id", userID, "project_id", projectID)

	var project *models.Project

	// Read operation to get project
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcProject, err := qtx.GetProjectByID(ctx, projectUUID)
		if err != nil {
			return fmt.Errorf("failed to get project: %w", err)
		}

		// Verify ownership
		if sqlcProject.CreatedBy != userUUID {
			s.logger.Warn("Unauthorized access attempt to project", "user_id", userID, "project_id", projectID, "owner_id", sqlcProject.CreatedBy)
			return fmt.Errorf("project not found")
		}

		project = s.sqlcToModel(sqlcProject)
		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get project", "user_id", userID, "project_id", projectID)
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	s.logger.Info("Project retrieved successfully", "user_id", userID, "project_id", projectID, "project_name", project.Name)
	return project, nil
}

// GetUserProjects retrieves all projects for a user
func (s *ProjectService) GetUserProjects(ctx context.Context, userID string) ([]*models.Project, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetUserProjects", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting user projects", "user_id", userID)

	var projects []*models.Project

	// Read operation to get projects
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcProjects, err := qtx.GetProjectsByUser(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("failed to get projects: %w", err)
		}

		projects = make([]*models.Project, len(sqlcProjects))
		for i, sqlcProject := range sqlcProjects {
			projects[i] = s.sqlcToModel(sqlcProject)
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get user projects", "user_id", userID)
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	s.logger.Info("User projects retrieved successfully", "user_id", userID, "projects_count", len(projects))
	return projects, nil
}

// GetActiveProjects retrieves all active projects for a user
func (s *ProjectService) GetActiveProjects(ctx context.Context, userID string) ([]*models.Project, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format in GetActiveProjects", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting active projects", "user_id", userID)

	var projects []*models.Project

	// Read operation to get active projects
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcProjects, err := qtx.GetActiveProjectsByUser(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("failed to get active projects: %w", err)
		}

		projects = make([]*models.Project, len(sqlcProjects))
		for i, sqlcProject := range sqlcProjects {
			projects[i] = s.sqlcToModel(sqlcProject)
		}

		return nil
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to get active projects", "user_id", userID)
		return nil, fmt.Errorf("failed to get active projects: %w", err)
	}

	s.logger.Info("Active projects retrieved successfully", "user_id", userID, "active_projects_count", len(projects))
	return projects, nil
}

// GetDefaultProject retrieves the default project for a user
func (s *ProjectService) GetDefaultProject(ctx context.Context, userID string) (*models.Project, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	s.logger.Info("Getting default project", "user_id", userID)

	var project *models.Project

	// Read operation to get default project
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcProject, err := qtx.GetUserDefaultProject(ctx, userUUID)
		if err != nil {
			return fmt.Errorf("failed to get default project: %w", err)
		}

		project = s.sqlcToModel(sqlcProject)
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Info("No default project found for user", "user_id", userID)
		} else {
			s.logger.LogError(ctx, err, "Failed to get default project", "user_id", userID)
		}
		return nil, fmt.Errorf("failed to get default project: %w", err)
	}

	s.logger.Info("Default project retrieved successfully", "user_id", userID, "project_id", project.ID, "project_name", project.Name)
	return project, nil
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(ctx context.Context, userID, projectID string, req *models.ProjectRequest) (*models.Project, error) {
	// Validate request
	if err := s.validateProjectRequest(req); err != nil {
		s.logger.Warn("Invalid project update request", "user_id", userID, "project_id", projectID, "error", err.Error())
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format", "user_id", userID)
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid project ID format", "project_id", projectID)
		return nil, fmt.Errorf("invalid project ID: %w", err)
	}

	s.logger.Info("Updating project", "user_id", userID, "project_id", projectID, "project_name", req.Name)

	var project *models.Project

	// Start write transaction to update project
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		sqlcProject, err := qtx.UpdateProject(ctx, store.UpdateProjectParams{
			ID:          projectUUID,
			Name:        req.Name,
			Description: stringToPgText(req.Description),
			Color:       stringToPgText(&req.Color),
			Status:      stringToPgText((*string)(&req.Status)),
			StartDate:   timeToPgDate(req.StartDate),
			EndDate:     timeToPgDate(req.EndDate),
			IsDefault:   boolToPgBool(req.IsDefault),
			CreatedBy:   userUUID,
		})
		if err != nil {
			return fmt.Errorf("failed to update project: %w", err)
		}

		project = s.sqlcToModel(sqlcProject)
		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Warn("Unauthorized project update attempt", "user_id", userID, "project_id", projectID)
		} else {
			s.logger.LogError(ctx, err, "Failed to update project", "user_id", userID, "project_id", projectID)
		}
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	s.logger.Info("Project updated successfully", "user_id", userID, "project_id", projectID, "project_name", project.Name)
	return project, nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(ctx context.Context, userID, projectID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid user ID format", "user_id", userID)
		return fmt.Errorf("invalid user ID: %w", err)
	}

	projectUUID, err := uuid.Parse(projectID)
	if err != nil {
		s.logger.LogError(ctx, err, "Invalid project ID format", "project_id", projectID)
		return fmt.Errorf("invalid project ID: %w", err)
	}

	s.logger.Info("Deleting project", "user_id", userID, "project_id", projectID)

	// First verify the project exists and user has ownership
	if err := s.db.Read(ctx, func(qtx *store.Queries) error {
		sqlcProject, err := qtx.GetProjectByID(ctx, projectUUID)
		if err != nil {
			return fmt.Errorf("project not found: %w", err)
		}

		// Verify ownership
		if sqlcProject.CreatedBy != userUUID {
			return fmt.Errorf("project not found")
		}

		return nil
	}); err != nil {
		if strings.Contains(err.Error(), "no rows") {
			s.logger.Warn("Project not found for deletion", "user_id", userID, "project_id", projectID)
		} else if strings.Contains(err.Error(), "project not found") && !strings.Contains(err.Error(), "no rows") {
			s.logger.Warn("Unauthorized project deletion attempt", "user_id", userID, "project_id", projectID)
		} else {
			s.logger.LogError(ctx, err, "Failed to verify project for deletion", "user_id", userID, "project_id", projectID)
		}
		return fmt.Errorf("failed to verify project: %w", err)
	}

	// Now delete the project
	if err := s.db.Write(ctx, func(qtx *store.Queries) error {
		return qtx.DeleteProject(ctx, store.DeleteProjectParams{
			ID:        projectUUID,
			CreatedBy: userUUID,
		})
	}); err != nil {
		s.logger.LogError(ctx, err, "Failed to delete project", "user_id", userID, "project_id", projectID)
		return fmt.Errorf("failed to delete project: %w", err)
	}

	s.logger.Info("Project deleted successfully", "user_id", userID, "project_id", projectID)
	return nil
}

// validateProjectRequest validates the project request
func (s *ProjectService) validateProjectRequest(req *models.ProjectRequest) error {
	if req.Name == "" {
		return fmt.Errorf("project name is required")
	}

	// Check for whitespace-only name
	if strings.TrimSpace(req.Name) == "" {
		return fmt.Errorf("project name cannot be whitespace only")
	}

	if len(req.Name) > 200 {
		return fmt.Errorf("project name must be at most 200 characters")
	}

	if !req.Status.IsValid() {
		return fmt.Errorf("invalid project status: %s", req.Status)
	}

	if err := models.ValidateHexColor(req.Color); err != nil {
		return fmt.Errorf("invalid color format: %w", err)
	}

	// Validate description length if provided
	if req.Description != nil && len(*req.Description) > 1000 {
		return fmt.Errorf("project description must be at most 1000 characters")
	}

	// Validate date range
	if req.StartDate != nil && req.EndDate != nil {
		if req.EndDate.Before(*req.StartDate) {
			return fmt.Errorf("end date must be after start date")
		}
	}

	return nil
}

// sqlcToModel converts SQLC Project to models.Project
func (s *ProjectService) sqlcToModel(sqlcProject store.Project) *models.Project {
	return &models.Project{
		ID:          sqlcProject.ID,
		Name:        sqlcProject.Name,
		Description: pgTextToString(sqlcProject.Description),
		Color:       pgTextToStringRequired(sqlcProject.Color),
		Status:      models.ProjectStatus(pgTextToStringRequired(sqlcProject.Status)),
		StartDate:   pgDateToTime(sqlcProject.StartDate),
		EndDate:     pgDateToTime(sqlcProject.EndDate),
		CreatedBy:   sqlcProject.CreatedBy,
		IsDefault:   pgBoolToBool(sqlcProject.IsDefault),
		CreatedAt:   pgTimestamptzToTime(sqlcProject.CreatedAt),
		UpdatedAt:   pgTimestamptzToTime(sqlcProject.UpdatedAt),
	}
}
