package handlers

import (
	"service/internal/apperrors"
	"service/internal/models"
	"service/internal/service"
	"service/internal/validators"
	"strconv"

	"service/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type NewsHandler struct {
	service service.INewsService
	log     *logger.Logger
}

func NewNewsHandler(service service.INewsService, log *logger.Logger) NewsHandler {
	return NewsHandler{
		service: service,
		log:     log,
	}
}

type ErrorResponse struct {
	Success bool   `json:"Success" example:"false"`
	Error   string `json:"Error"`
}

type SuccessResponse struct {
	Success bool `json:"Success" example:"true"`
}

type SuccessResponseCreate struct {
	Success bool  `json:"Success" example:"true"`
	Id      int64 `json:"Id" example:"1"`
}

type NewsListsResponse struct {
	Success bool                        `json:"Success" example:"true"`
	News    []models.NewsWithCategories `json:"News"`
}

// CreateNews godoc
// @Summary Create news
// @Description Create news with title, content and categories(optional). Categories must be positive integers, example: [1, 2, 3]
// @Tags news
// @Accept json
// @Produce json
// @Param request body models.NewsCreateForm true "News data"
// @Success 201 {object} SuccessResponseCreate "News created successful"
// @Failure 400 {object} ErrorResponse "Validation error"
// @Failure 401 {object} ErrorResponse "No authorization"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /create [post]
func (h *NewsHandler) CreateNews(c *fiber.Ctx) error {
	if err := validators.ValidateCreateNewsRequest(c.Body()); err != nil {
		return apperrors.NewValidation(err.Error())
	}

	var reqForm models.NewsCreateForm
	if err := c.BodyParser(&reqForm); err != nil {
		return apperrors.NewBadRequest("Failed to parse request body")
	}

	reqForm.Normalize()

	if err := reqForm.Validate(); err != nil {
		return apperrors.NewValidation(err.Error())
	}

	id, err := h.service.CreateNews(reqForm)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(SuccessResponseCreate{
		Success: true,
		Id:      id,
	})
}

// EditNews godoc
// @Summary Edit news
// @Description Edit news fields (title, content, categories). Categories must be positive integers, example: [1, 2, 3]
// @Tags news
// @Accept json
// @Produce json
// @Param id path int true "ID news"
// @Param request body models.NewsEditForm true "News updated data"
// @Success 200 {object} SuccessResponse "Success updated"
// @Failure 400 {object} ErrorResponse "Error validation"
// @Failure 401 {object} ErrorResponse "Not authorized"
// @Failure 404 {object} ErrorResponse "News not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /edit/{id} [post]
func (h *NewsHandler) EditNews(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return apperrors.NewBadRequest("Invalid ID format")
	}

	if err = validators.ValidateEditNewsRequest(c.Body()); err != nil {
		return apperrors.NewValidation(err.Error())
	}

	var editForm models.NewsEditForm
	if err = c.BodyParser(&editForm); err != nil {
		return apperrors.NewBadRequest("Failed to parse request body")
	}

	editForm.Normalize()
	if err = editForm.Validate(); err != nil {
		return apperrors.NewValidation(err.Error())
	}

	if err = h.service.EditNews(int64(id), editForm); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Success: true,
	})
}

// ListNews godoc
// @Summary Get news
// @Tags news
// @Accept json
// @Produce json
// @Param limit query int false "default=10, max=100"
// @Param offset query int false "default=0"
// @Success 200 {object} NewsListsResponse "List news"
// @Failure 400 {object} ErrorResponse "Error validation params"
// @Failure 401 {object} ErrorResponse "Not authorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /list [get]
func (h *NewsHandler) ListNews(c *fiber.Ctx) error {
	limit, err := strconv.ParseInt(c.Query("limit", "10"), 10, 64)
	if err != nil {
		return apperrors.NewBadRequest("limit must be a valid number")
	}

	offset, err := strconv.ParseInt(c.Query("offset", "0"), 10, 64)
	if err != nil {
		return apperrors.NewBadRequest("offset must be a valid number")
	}

	if err = validators.ValidatePaginationParams(limit, offset); err != nil {
		return err
	}

	newsList, err := h.service.ListNews(limit, offset)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(NewsListsResponse{Success: true, News: newsList})
}
