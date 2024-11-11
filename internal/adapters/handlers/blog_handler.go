package handlers

import (
	"strconv"

	"github.com/toffysoft/go-hexagonal-example/internal/core/domain"
	"github.com/toffysoft/go-hexagonal-example/internal/core/ports"
	"github.com/toffysoft/go-hexagonal-example/pkg/errors"
	"github.com/toffysoft/go-hexagonal-example/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BlogHandler struct {
	blogService ports.BlogService
	validate    *validator.Validate
}

func NewBlogHandler(blogService ports.BlogService) *BlogHandler {
	return &BlogHandler{
		blogService: blogService,
		validate:    validator.New(),
	}
}

type CreateBlogRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=100"`
	Content string `json:"content" validate:"required,min=10"`
	Author  string `json:"author" validate:"required,min=2,max=50"`
}

func (h *BlogHandler) CreateBlog(c *fiber.Ctx) error {
	var req CreateBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	blog := &domain.Blog{
		Title:   req.Title,
		Content: req.Content,
		Author:  req.Author,
	}

	if err := h.blogService.CreateBlog(blog); err != nil {
		if appErr, ok := err.(errors.AppError); ok {
			return utils.SendErrorResponse(c, appErr.StatusCode(), appErr.Error())
		}
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to create blog")
	}

	return utils.SendSuccessResponse(c, fiber.StatusCreated, "Blog created successfully", blog)
}

type UpdateBlogRequest struct {
	Title   string `json:"title" validate:"omitempty,min=3,max=100"`
	Content string `json:"content" validate:"omitempty,min=10"`
	Author  string `json:"author" validate:"omitempty,min=2,max=50"`
}

func (h *BlogHandler) UpdateBlog(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid blog ID")
	}

	var req UpdateBlogRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.validate.Struct(req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	blog, err := h.blogService.GetBlog(uint(id))
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "Blog not found")
	}

	if req.Title != "" {
		blog.Title = req.Title
	}
	if req.Content != "" {
		blog.Content = req.Content
	}
	if req.Author != "" {
		blog.Author = req.Author
	}

	if err := h.blogService.UpdateBlog(blog); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to update blog")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, "Blog updated successfully", blog)
}

func (h *BlogHandler) GetBlog(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid blog ID")
	}

	blog, err := h.blogService.GetBlog(uint(id))
	if err != nil {
		if appErr, ok := err.(errors.AppError); ok {
			return utils.SendErrorResponse(c, appErr.StatusCode(), appErr.Error())
		}
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve blog")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, "Blog retrieved successfully", blog)
}

func (h *BlogHandler) DeleteBlog(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "Invalid blog ID")
	}

	if err := h.blogService.DeleteBlog(uint(id)); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete blog")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, "Blog deleted successfully", nil)
}

func (h *BlogHandler) ListBlogs(c *fiber.Ctx) error {
	blogs, err := h.blogService.ListBlogs()
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve blogs")
	}

	return utils.SendSuccessResponse(c, fiber.StatusOK, "Blogs retrieved successfully", blogs)
}
