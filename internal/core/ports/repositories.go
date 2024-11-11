package ports

import "github.com/toffysoft/go-hexagonal-example/internal/core/domain"

type BlogRepository interface {
	Create(blog *domain.Blog) error
	GetByID(id uint) (*domain.Blog, error)
	Update(blog *domain.Blog) error
	Delete(id uint) error
	List() ([]*domain.Blog, error)
}
