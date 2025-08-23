package services

import (
	db "charity/db/sqlc"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"database/sql"
)

type ProjectService struct {
	queries *db.Queries
}

func NewProjectService(queries *db.Queries) *ProjectService {
	return &ProjectService{
		queries: queries,
	}
}

// CreateProject creates a new project with basic information
func (s *ProjectService) CreateProject(ctx context.Context, name string) (*db.Project, error) {
	project, err := s.queries.CreateProject(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	return &project, nil
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(ctx context.Context, id int32) (*db.Project, error) {
	project, err := s.queries.GetProject(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	return &project, nil
}

// ListProjects retrieves all projects ordered by creation date
func (s *ProjectService) ListProjects(ctx context.Context) ([]db.Project, error) {
	projects, err := s.queries.ListProjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	return projects, nil
}

// ListProjectsByStatus retrieves projects filtered by status
func (s *ProjectService) ListProjectsByStatus(ctx context.Context, status string) ([]db.Project, error) {
	projects, err := s.queries.ListProjectsByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects by status: %w", err)
	}
	return projects, nil
}

// UpdateProject updates a project's basic information
func (s *ProjectService) UpdateProject(ctx context.Context, id int32, name string) (*db.Project, error) {
	project, err := s.queries.UpdateProject(ctx, db.UpdateProjectParams{
		ID:   id,
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	return &project, nil
}

// UpdateProjectStatus updates a project's status
func (s *ProjectService) UpdateProjectStatus(ctx context.Context, id int32, status string) (*db.Project, error) {
	project, err := s.queries.UpdateProjectStatus(ctx, db.UpdateProjectStatusParams{ID: id, Status: status})
	if err != nil {
		return nil, fmt.Errorf("failed to update project status: %w", err)
	}
	return &project, nil
}

// DeleteProject deletes a project and all its related data
func (s *ProjectService) DeleteProject(ctx context.Context, id int32) error {
	// Delete project images first
	if err := s.queries.DeleteAllProjectImages(ctx, id); err != nil {
		log.Printf("Warning: failed to delete project images: %v", err)
	}

	// Delete project before/after data
	if err := s.queries.DeleteProjectBefore(ctx, id); err != nil {
		log.Printf("Warning: failed to delete project before data: %v", err)
	}

	if err := s.queries.DeleteProjectAfter(ctx, id); err != nil {
		log.Printf("Warning: failed to delete project after data: %v", err)
	}

	// Delete the project
	if err := s.queries.DeleteProject(ctx, id); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

// CreateProjectBefore creates or updates the "before" section of a project
func (s *ProjectService) CreateProjectBefore(ctx context.Context, projectID int32, body, estimatedTarget, currentFunds, videoLink string) (*db.ProjectBefore, error) {
	// Convert estimatedTarget string to int32
	estimatedTargetInt, err := strconv.Atoi(estimatedTarget)
	if err != nil {
		return nil, fmt.Errorf("invalid estimated target: must be a number")
	}

	// Convert currentFunds string to int32
	currentFundsInt, err := strconv.Atoi(currentFunds)
	if err != nil {
		return nil, fmt.Errorf("invalid current funds: must be a number")
	}

	// Try to update first, if it doesn't exist, create new
	_, err = s.queries.GetProjectBefore(ctx, projectID)
	if err != nil {
		// Project before doesn't exist, create new
		projectBefore, err := s.queries.CreateProjectBefore(ctx, db.CreateProjectBeforeParams{
			ProjectID:       projectID,
			Body:            body,
			EstimatedTarget: int32(estimatedTargetInt),
			CurrentFunds:    int32(currentFundsInt),
			VideoLink:       sql.NullString{String: videoLink, Valid: videoLink != ""},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create project before: %w", err)
		}
		return &projectBefore, nil
	}

	// Update existing
	projectBefore, err := s.queries.UpdateProjectBefore(ctx, db.UpdateProjectBeforeParams{
		ProjectID:       projectID,
		Body:            body,
		EstimatedTarget: int32(estimatedTargetInt),
		CurrentFunds:    int32(currentFundsInt),
		VideoLink:       sql.NullString{String: videoLink, Valid: videoLink != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update project before: %w", err)
	}
	return &projectBefore, nil
}

// GetProjectBefore retrieves the "before" section of a project
func (s *ProjectService) GetProjectBefore(ctx context.Context, projectID int32) (*db.ProjectBefore, error) {
	projectBefore, err := s.queries.GetProjectBefore(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project before: %w", err)
	}
	return &projectBefore, nil
}

// CreateProjectAfter creates or updates the "after" section of a project
func (s *ProjectService) CreateProjectAfter(ctx context.Context, projectID int32, body, projectCost, videoLink string) (*db.ProjectAfter, error) {
	// Convert projectCost string to int32
	var projectCostInt sql.NullInt32
	if projectCost != "" {
		cost, err := strconv.Atoi(projectCost)
		if err != nil {
			return nil, fmt.Errorf("invalid project cost: must be a number")
		}
		projectCostInt = sql.NullInt32{Int32: int32(cost), Valid: true}
	}

	// Try to update first, if it doesn't exist, create new
	_, err := s.queries.GetProjectAfter(ctx, projectID)
	if err != nil {
		// Project after doesn't exist, create new
		projectAfter, err := s.queries.CreateProjectAfter(ctx, db.CreateProjectAfterParams{
			ProjectID:   projectID,
			Body:        sql.NullString{String: body, Valid: body != ""},
			ProjectCost: projectCostInt,
			VideoLink:   sql.NullString{String: videoLink, Valid: videoLink != ""},
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create project after: %w", err)
		}
		return &projectAfter, nil
	}

	// Update existing
	projectAfter, err := s.queries.UpdateProjectAfter(ctx, db.UpdateProjectAfterParams{
		ProjectID:   projectID,
		Body:        sql.NullString{String: body, Valid: body != ""},
		ProjectCost: projectCostInt,
		VideoLink:   sql.NullString{String: videoLink, Valid: videoLink != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update project after: %w", err)
	}
	return &projectAfter, nil
}

// GetProjectAfter retrieves the "after" section of a project
func (s *ProjectService) GetProjectAfter(ctx context.Context, projectID int32) (*db.ProjectAfter, error) {
	projectAfter, err := s.queries.GetProjectAfter(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project after: %w", err)
	}
	return &projectAfter, nil
}

// UploadProjectImage uploads an image for a project
func (s *ProjectService) UploadProjectImage(ctx context.Context, projectID int32, phase string, file *multipart.FileHeader) (*db.ProjectImage, error) {
	// Validate file type
	if !isValidImageFile(file.Filename) {
		return nil, fmt.Errorf("invalid file type, only jpg, jpeg, png, gif are allowed")
	}

	// Create unique filename
	timestamp := time.Now().Unix()
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("project_%d_%s_%d%s", projectID, phase, timestamp, ext)

	// Ensure images directory exists
	imagesDir := "static/images"
	if err := os.MkdirAll(imagesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create images directory: %w", err)
	}

	// Create file path
	filePath := filepath.Join(imagesDir, filename)

	// Open source file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// Save to database
	imageLink := "/static/images/" + filename
	projectImage, err := s.queries.AddProjectImage(ctx, db.AddProjectImageParams{
		ProjectID: projectID,
		Phase:     phase,
		ImageLink: imageLink,
	})
	if err != nil {
		// Clean up file if database save fails
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to save image to database: %w", err)
	}

	return &projectImage, nil
}

// ListProjectImages retrieves all images for a project
func (s *ProjectService) ListProjectImages(ctx context.Context, projectID int32) ([]db.ProjectImage, error) {
	images, err := s.queries.ListProjectImages(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to list project images: %w", err)
	}
	return images, nil
}

// ListProjectImagesByPhase retrieves images for a project by phase
func (s *ProjectService) ListProjectImagesByPhase(ctx context.Context, projectID int32, phase string) ([]db.ProjectImage, error) {
	images, err := s.queries.ListProjectImagesByPhase(ctx, db.ListProjectImagesByPhaseParams{
		ProjectID: projectID,
		Phase:     phase,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list project images by phase: %w", err)
	}
	return images, nil
}

// DeleteProjectImage deletes a project image
func (s *ProjectService) DeleteProjectImage(ctx context.Context, imageID int32) error {
	// Get image info first to delete file
	images, err := s.queries.ListProjectImages(ctx, 0) // We'll need to get the specific image
	if err != nil {
		return fmt.Errorf("failed to get image info: %w", err)
	}

	// Find the image to delete
	var imageToDelete *db.ProjectImage
	for _, img := range images {
		if img.ID == imageID {
			imageToDelete = &img
			break
		}
	}

	if imageToDelete == nil {
		return fmt.Errorf("image not found")
	}

	// Delete from database first
	if err := s.queries.DeleteProjectImage(ctx, imageID); err != nil {
		return fmt.Errorf("failed to delete image from database: %w", err)
	}

	// Delete file from filesystem
	filePath := strings.TrimPrefix(imageToDelete.ImageLink, "/static")
	filePath = "static" + filePath
	if err := os.Remove(filePath); err != nil {
		log.Printf("Warning: failed to delete image file %s: %v", filePath, err)
	}

	return nil
}

// isValidImageFile checks if the uploaded file is a valid image
func isValidImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif"}
	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}
