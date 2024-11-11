package services_test

import (
	"testing"

	"github.com/toffysoft/go-hexagonal-example/internal/core/domain"
	"github.com/toffysoft/go-hexagonal-example/internal/core/services"
	"github.com/toffysoft/go-hexagonal-example/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBlogRepository is a mock type for the BlogRepository
type MockBlogRepository struct {
	mock.Mock
}

func (m *MockBlogRepository) Create(blog *domain.Blog) error {
	args := m.Called(blog)
	return args.Error(0)
}

func (m *MockBlogRepository) GetByID(id uint) (*domain.Blog, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Blog), args.Error(1)
}

func (m *MockBlogRepository) Update(blog *domain.Blog) error {
	args := m.Called(blog)
	return args.Error(0)
}

func (m *MockBlogRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBlogRepository) List() ([]*domain.Blog, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Blog), args.Error(1)
}

func TestCreateBlog(t *testing.T) {
	mockRepo := new(MockBlogRepository)
	blogService := services.NewBlogService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		blog := &domain.Blog{Title: "Test Blog", Content: "Test Content", Author: "Test Author"}
		mockRepo.On("Create", blog).Return(nil).Once()

		err := blogService.CreateBlog(blog)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyTitle", func(t *testing.T) {
		blog := &domain.Blog{Content: "Test Content", Author: "Test Author"}

		err := blogService.CreateBlog(blog)

		assert.Error(t, err)
		assert.IsType(t, errors.AppError{}, err)
		assert.Equal(t, errors.InvalidInput, err.(errors.AppError).Type)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		blog := &domain.Blog{Title: "Test Blog", Content: "Test Content", Author: "Test Author"}
		mockRepo.On("Create", blog).Return(errors.NewInternalServerError("Database error")).Once()

		err := blogService.CreateBlog(blog)

		assert.Error(t, err)
		assert.IsType(t, errors.AppError{}, err)
		assert.Equal(t, errors.InternalServer, err.(errors.AppError).Type)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetBlog(t *testing.T) {
	mockRepo := new(MockBlogRepository)
	blogService := services.NewBlogService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		blog := &domain.Blog{ID: 1, Title: "Test Blog", Content: "Test Content", Author: "Test Author"}
		mockRepo.On("GetByID", uint(1)).Return(blog, nil).Once()

		result, err := blogService.GetBlog(1)

		assert.NoError(t, err)
		assert.Equal(t, blog, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return((*domain.Blog)(nil), errors.NewNotFoundError("Blog not found")).Once()

		result, err := blogService.GetBlog(999)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.IsType(t, errors.AppError{}, err)
		assert.Equal(t, errors.NotFound, err.(errors.AppError).Type)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateBlog(t *testing.T) {
	mockRepo := new(MockBlogRepository)
	blogService := services.NewBlogService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		blog := &domain.Blog{ID: 1, Title: "Updated Blog", Content: "Updated Content", Author: "Updated Author"}
		mockRepo.On("GetByID", uint(1)).Return(blog, nil).Once()
		mockRepo.On("Update", blog).Return(nil).Once()

		err := blogService.UpdateBlog(blog)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		blog := &domain.Blog{ID: 999, Title: "Non-existent Blog"}
		mockRepo.On("GetByID", uint(999)).Return((*domain.Blog)(nil), errors.NewNotFoundError("Blog not found")).Once()

		err := blogService.UpdateBlog(blog)

		assert.Error(t, err)
		assert.IsType(t, errors.AppError{}, err)
		assert.Equal(t, errors.NotFound, err.(errors.AppError).Type)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		blog := &domain.Blog{ID: 0, Title: "Invalid Blog"}

		err := blogService.UpdateBlog(blog)

		assert.Error(t, err)
		assert.IsType(t, errors.AppError{}, err)
		assert.Equal(t, errors.InvalidInput, err.(errors.AppError).Type)
	})
}

func TestDeleteBlog(t *testing.T) {
	mockRepo := new(MockBlogRepository)
	blogService := services.NewBlogService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("GetByID", uint(1)).Return(&domain.Blog{ID: 1}, nil).Once()
		mockRepo.On("Delete", uint(1)).Return(nil).Once()

		err := blogService.DeleteBlog(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.On("GetByID", uint(999)).Return((*domain.Blog)(nil), errors.NewNotFoundError("Blog not found")).Once()

		err := blogService.DeleteBlog(999)

		assert.Error(t, err)
		assert.IsType(t, errors.AppError{}, err)
		assert.Equal(t, errors.NotFound, err.(errors.AppError).Type)
		mockRepo.AssertExpectations(t)
	})
}

func TestListBlogs(t *testing.T) {
	mockRepo := new(MockBlogRepository)
	blogService := services.NewBlogService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		blogs := []*domain.Blog{
			{ID: 1, Title: "Blog 1"},
			{ID: 2, Title: "Blog 2"},
		}
		mockRepo.On("List").Return(blogs, nil).Once()

		result, err := blogService.ListBlogs()

		assert.NoError(t, err)
		assert.Equal(t, blogs, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyList", func(t *testing.T) {
		mockRepo.On("List").Return([]*domain.Blog{}, nil).Once()

		result, err := blogService.ListBlogs()

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo.On("List").Return(([]*domain.Blog)(nil), errors.NewInternalServerError("Database error")).Once()

		result, err := blogService.ListBlogs()

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.IsType(t, errors.AppError{}, err)
		assert.Equal(t, errors.InternalServer, err.(errors.AppError).Type)
		mockRepo.AssertExpectations(t)
	})
}
