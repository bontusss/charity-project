package handlers

import (
	"charity/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService *services.ProjectService
}

func NewProjectHandler(projectService *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// CreateProject creates a new project (admin only)
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	project, err := h.projectService.CreateProject(c.Request.Context(), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Project created successfully",
		"project": project,
	})
}

// GetProject retrieves a project by ID
func (h *ProjectHandler) GetProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := h.projectService.GetProject(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

// ListProjects retrieves all projects
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	projects, err := h.projectService.ListProjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// UpdateProject updates a project (admin only)
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	project, err := h.projectService.UpdateProject(c.Request.Context(), int32(id), req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project updated successfully",
		"project": project,
	})
}

// UpdateProjectStatus updates a project's status (admin only)
func (h *ProjectHandler) UpdateProjectStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=ongoing completed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	project, err := h.projectService.UpdateProjectStatus(c.Request.Context(), int32(id), req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Project status updated successfully",
		"project": project,
	})
}

// DeleteProject deletes a project (admin only)
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := h.projectService.DeleteProject(c.Request.Context(), int32(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// CreateProjectBefore creates or updates project "before" section (admin only)
func (h *ProjectHandler) CreateProjectBefore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		Body            string `json:"body" binding:"required"`
		EstimatedTarget string `json:"estimated_target" binding:"required"`
		CurrentFunds    string `json:"current_funds" binding:"required"`
		VideoLink       string `json:"video_link"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	projectBefore, err := h.projectService.CreateProjectBefore(c.Request.Context(), int32(id), req.Body, req.EstimatedTarget, req.CurrentFunds, req.VideoLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Project before section saved successfully",
		"project_before": projectBefore,
	})
}

// UpdateProjectBefore updates only specific fields of project "before" section (admin only)
func (h *ProjectHandler) UpdateProjectBefore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		EstimatedTarget string `json:"estimated_target"`
		CurrentFunds    string `json:"current_funds"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Get existing project before data
	existingBefore, err := h.projectService.GetProjectBefore(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project before section not found"})
		return
	}

	// Update only the specified fields
	estimatedTarget := existingBefore.EstimatedTarget
	currentFunds := existingBefore.CurrentFunds

	if req.EstimatedTarget != "" {
		if target, err := strconv.Atoi(req.EstimatedTarget); err == nil {
			estimatedTarget = int32(target)
		}
	}

	if req.CurrentFunds != "" {
		if funds, err := strconv.Atoi(req.CurrentFunds); err == nil {
			currentFunds = int32(funds)
		}
	}

	// Create a minimal body if it doesn't exist (required field)
	body := existingBefore.Body
	if body == "" {
		body = "Project details"
	}

	projectBefore, err := h.projectService.CreateProjectBefore(c.Request.Context(), int32(id), body, strconv.Itoa(int(estimatedTarget)), strconv.Itoa(int(currentFunds)), "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Project before section updated successfully",
		"project_before": projectBefore,
	})
}

// GetProjectBefore retrieves project "before" section
func (h *ProjectHandler) GetProjectBefore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	projectBefore, err := h.projectService.GetProjectBefore(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project before section not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project_before": projectBefore})
}

// CreateProjectAfter creates or updates project "after" section (admin only)
func (h *ProjectHandler) CreateProjectAfter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req struct {
		Body        string `json:"body" binding:"required"`
		ProjectCost string `json:"project_cost"`
		VideoLink   string `json:"video_link"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	projectAfter, err := h.projectService.CreateProjectAfter(c.Request.Context(), int32(id), req.Body, req.ProjectCost, req.VideoLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Project after section saved successfully",
		"project_after": projectAfter,
	})
}

// GetProjectAfter retrieves project "after" section
func (h *ProjectHandler) GetProjectAfter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	projectAfter, err := h.projectService.GetProjectAfter(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project after section not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project_after": projectAfter})
}

func (h *ProjectHandler) UploadProjectImage(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
        return
    }

    phase := c.PostForm("phase")
    if phase == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Phase is required"})
        return
    }

    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
        return
    }

    files := form.File["image"]
    if len(files) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No image files provided"})
        return
    }

    var uploadedImages []interface{}
    for _, file := range files {
        projectImage, err := h.projectService.UploadProjectImage(c.Request.Context(), int32(id), phase, file)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        uploadedImages = append(uploadedImages, projectImage)
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Images uploaded successfully",
        "images":  uploadedImages,
    })
}

// UploadProjectImage uploads an image for a project (admin only)
// func (h *ProjectHandler) UploadProjectImage(c *gin.Context) {
// 	idStr := c.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
// 		return
// 	}

// 	phase := c.PostForm("phase")
// 	if phase == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Phase is required"})
// 		return
// 	}

// 	file, err := c.FormFile("image")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
// 		return
// 	}

// 	projectImage, err := h.projectService.UploadProjectImage(c.Request.Context(), int32(id), phase, file)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Image uploaded successfully",
// 		"image":   projectImage,
// 	})
// }

// ListProjectImages retrieves all images for a project
func (h *ProjectHandler) ListProjectImages(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	images, err := h.projectService.ListProjectImages(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

// ListProjectImagesByPhase retrieves images for a project by phase
func (h *ProjectHandler) ListProjectImagesByPhase(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	phase := c.Query("phase")
	if phase == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phase parameter is required"})
		return
	}

	images, err := h.projectService.ListProjectImagesByPhase(c.Request.Context(), int32(id), phase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

// DeleteProjectImage deletes a project image (admin only)
func (h *ProjectHandler) DeleteProjectImage(c *gin.Context) {
	imageIDStr := c.Param("image_id")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	if err := h.projectService.DeleteProjectImage(c.Request.Context(), int32(imageID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image deleted successfully"})
}
