package handlers

import (
	"net/http"
	"strconv"

	"trithong.com/task-golang/internal/dto"
	"trithong.com/task-golang/internal/mappers"
	"trithong.com/task-golang/internal/responses"
	"trithong.com/task-golang/internal/services"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projectService services.ProjectService
}

func NewProjectHandler(projectService services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

func (h *ProjectHandler) Create(c *gin.Context) {

	var request dto.CreateProjectRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	project := mappers.ToProjectEntity(request)

	userID := c.GetInt("user_id")

	project.OwnerID = userID

	result, err := h.projectService.Create(project)

	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := mappers.ToProjectResponse(result)

	responses.Success(c, http.StatusCreated, "Create project successfully", response)
}

func (h *ProjectHandler) GetAll(c *gin.Context) {
	userID := c.GetInt("user_id")

	var params dto.ProjectQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid query params")
		return
	}

	// Nếu không có pagination thì trả về tất cả như cũ
	if params.Page == 0 && params.Limit == 0 {
		projects, err := h.projectService.GetAllByOwnerID(userID)
		if err != nil {
			responses.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		var projectResponses []dto.ProjectResponse
		for _, project := range projects {
			projectResponses = append(projectResponses, mappers.ToProjectResponse(project))
		}

		responses.Success(c, http.StatusOK, "Get projects successfully", projectResponses)
		return
	}

	projects, total, err := h.projectService.GetAllWithPagination(userID, params)
	if err != nil {
		responses.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	totalPages := (total + params.Limit - 1) / params.Limit

	var projectResponses []dto.ProjectResponse
	for _, project := range projects {
		projectResponses = append(projectResponses, mappers.ToProjectResponse(project))
	}

	responses.Success(c, http.StatusOK, "Get projects successfully", dto.PaginationResponse{
		Items:      projectResponses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	})
}

func (h *ProjectHandler) GetByID(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid project id")
		return
	}

	userID := c.GetInt("user_id")

	project, err := h.projectService.GetByIDAndOwnerID(id, userID)

	if err != nil {
		responses.Error(c, http.StatusNotFound, "Project not found")
		return
	}

	response := mappers.ToProjectResponse(project)

	responses.Success(
		c,
		http.StatusOK,
		"Get project successfully",
		response,
	)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid project id")
		return
	}

	var request dto.UpdateProjectRequest

	err = c.ShouldBindJSON(&request)

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID := c.GetInt("user_id")

	project := mappers.ToUpdateProjectEntity(request)

	result, err := h.projectService.Update(id, userID, project)

	if err != nil {
		responses.Error(c, http.StatusNotFound, "Project not found")
		return
	}

	response := mappers.ToProjectResponse(result)

	responses.Success(c, http.StatusOK, "Update project successfully", response)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		responses.Error(c, http.StatusBadRequest, "Invalid project id")
		return
	}

	userID := c.GetInt("user_id")

	err = h.projectService.Delete(id, userID)

	if err != nil {
		responses.Error(c, http.StatusNotFound, "Project not found")
		return
	}

	responses.Success(c, http.StatusOK, "Delete project successfully", nil)
}
