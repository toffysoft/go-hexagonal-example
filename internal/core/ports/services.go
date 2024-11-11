package ports

import "github.com/toffysoft/go-hexagonal-example/internal/core/domain"

type BlogService interface {
	CreateBlog(blog *domain.Blog) error
	GetBlog(id uint) (*domain.Blog, error)
	UpdateBlog(blog *domain.Blog) error
	DeleteBlog(id uint) error
	ListBlogs() ([]*domain.Blog, error)
}
