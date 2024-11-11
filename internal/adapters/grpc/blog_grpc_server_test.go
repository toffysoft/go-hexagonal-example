package grpc_test

import (
	"context"
	"testing"

	"github.com/toffysoft/go-hexagonal-example/internal/adapters/grpc"
	"github.com/toffysoft/go-hexagonal-example/internal/adapters/grpc/proto"
	"github.com/toffysoft/go-hexagonal-example/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBlogService struct {
	mock.Mock
}

func (m *MockBlogService) CreateBlog(blog *domain.Blog) error {
	args := m.Called(blog)
	return args.Error(0)
}

func (m *MockBlogService) GetBlog(id uint) (*domain.Blog, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Blog), args.Error(1)
}

func (m *MockBlogService) UpdateBlog(blog *domain.Blog) error {
	args := m.Called(blog)
	return args.Error(0)
}

func (m *MockBlogService) DeleteBlog(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBlogService) ListBlogs() ([]*domain.Blog, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Blog), args.Error(1)
}

// Implement other methods...

func TestCreateBlog(t *testing.T) {
	mockService := new(MockBlogService)
	server := grpc.NewBlogServer(mockService)

	req := &proto.CreateBlogRequest{
		Title:   "Test Blog",
		Content: "Test Content",
		Author:  "Test Author",
	}

	mockService.On("CreateBlog", mock.AnythingOfType("*domain.Blog")).Return(nil)

	resp, err := server.CreateBlog(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Title, resp.Blog.Title)
	assert.Equal(t, req.Content, resp.Blog.Content)
	assert.Equal(t, req.Author, resp.Blog.Author)

	mockService.AssertExpectations(t)
}

// Implement other test cases...

func TestListBlogs(t *testing.T) {
	mockService := new(MockBlogService)
	server := grpc.NewBlogServer(mockService)

	blogs := []*domain.Blog{
		{ID: 1, Title: "Test Blog 1", Content: "Test Content 1", Author: "Test Author 1"},
		{ID: 2, Title: "Test Blog 2", Content: "Test Content 2", Author: "Test Author 2"},
	}

	mockService.On("ListBlogs").Return(blogs, nil)

	resp, err := server.ListBlogs(context.Background(), &proto.ListBlogsRequest{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Blogs, len(blogs))

	for i, blog := range blogs {
		assert.Equal(t, blog.ID, uint(resp.Blogs[i].Id))
		assert.Equal(t, blog.Title, resp.Blogs[i].Title)
		assert.Equal(t, blog.Content, resp.Blogs[i].Content)
	}

	mockService.AssertExpectations(t)
}

// Implement other test cases...

func TestGetBlog(t *testing.T) {
	mockService := new(MockBlogService)
	server := grpc.NewBlogServer(mockService)

	blog := &domain.Blog{
		ID:      1,
		Title:   "Test Blog",
		Content: "Test Content",
	}

	mockService.On("GetBlog", uint(1)).Return(blog, nil)

	resp, err := server.GetBlog(context.Background(), &proto.GetBlogRequest{Id: 1})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, blog.ID, uint(resp.Blog.Id))
	assert.Equal(t, blog.Title, resp.Blog.Title)
	assert.Equal(t, blog.Content, resp.Blog.Content)

	mockService.AssertExpectations(t)
}

// Implement other test cases...

func TestUpdateBlog(t *testing.T) {
	mockService := new(MockBlogService)
	server := grpc.NewBlogServer(mockService)

	req := &proto.UpdateBlogRequest{
		Id:      1,
		Title:   "Test Blog",
		Content: "Test Content",
		Author:  "Test Author",
	}

	mockService.On("UpdateBlog", mock.AnythingOfType("*domain.Blog")).Return(nil)

	resp, err := server.UpdateBlog(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Id, uint64(resp.Blog.Id))

	mockService.AssertExpectations(t)
}

// Implement other test cases...

func TestDeleteBlog(t *testing.T) {
	mockService := new(MockBlogService)
	server := grpc.NewBlogServer(mockService)

	mockService.On("DeleteBlog", uint(1)).Return(nil)

	resp, err := server.DeleteBlog(context.Background(), &proto.DeleteBlogRequest{Id: 1})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)

	mockService.AssertExpectations(t)
}
