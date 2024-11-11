package services

import (
	"fmt"

	"github.com/toffysoft/go-hexagonal-example/internal/core/domain"
	"github.com/toffysoft/go-hexagonal-example/internal/core/ports"
	"github.com/toffysoft/go-hexagonal-example/pkg/errors"
)

type blogService struct {
	repo ports.BlogRepository
}

func NewBlogService(repo ports.BlogRepository) ports.BlogService {
	return &blogService{repo: repo}
}

func (s *blogService) CreateBlog(blog *domain.Blog) error {
	if blog.Title == "" || blog.Content == "" || blog.Author == "" {
		return errors.NewInvalidInputError("All fields are required")
	}
	return s.repo.Create(blog)
}

func (s *blogService) GetBlog(id uint) (*domain.Blog, error) {
	blog, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("Blog with ID %d not found", id))
	}
	return blog, nil
}

func (s *blogService) UpdateBlog(blog *domain.Blog) error {
	if blog.ID == 0 {
		return errors.NewInvalidInputError("Blog ID is required")
	}
	_, err := s.GetBlog(blog.ID)
	if err != nil {
		return err
	}
	return s.repo.Update(blog)
}

func (s *blogService) DeleteBlog(id uint) error {
	_, err := s.GetBlog(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *blogService) ListBlogs() ([]*domain.Blog, error) {
	return s.repo.List()
}
