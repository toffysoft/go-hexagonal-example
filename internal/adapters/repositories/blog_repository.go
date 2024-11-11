package repositories

import (
	"github.com/toffysoft/go-hexagonal-example/internal/core/domain"
	"github.com/toffysoft/go-hexagonal-example/internal/core/ports"

	"gorm.io/gorm"
)

type blogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) ports.BlogRepository {
	return &blogRepository{db: db}
}

func (r *blogRepository) Create(blog *domain.Blog) error {
	return r.db.Create(blog).Error
}

func (r *blogRepository) GetByID(id uint) (*domain.Blog, error) {
	var blog domain.Blog
	err := r.db.First(&blog, id).Error
	return &blog, err
}

func (r *blogRepository) Update(blog *domain.Blog) error {
	return r.db.Save(blog).Error
}

func (r *blogRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Blog{}, id).Error
}

func (r *blogRepository) List() ([]*domain.Blog, error) {
	var blogs []*domain.Blog
	err := r.db.Find(&blogs).Error
	return blogs, err
}
