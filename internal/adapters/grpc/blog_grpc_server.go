package grpc

import (
	"context"

	"github.com/toffysoft/go-hexagonal-example/internal/adapters/grpc/proto"
	"github.com/toffysoft/go-hexagonal-example/internal/core/domain"
	"github.com/toffysoft/go-hexagonal-example/internal/core/ports"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BlogServer struct {
	proto.UnimplementedBlogServiceServer
	blogService ports.BlogService
}

func NewBlogServer(blogService ports.BlogService) *BlogServer {
	return &BlogServer{blogService: blogService}
}

func (s *BlogServer) CreateBlog(ctx context.Context, req *proto.CreateBlogRequest) (*proto.BlogResponse, error) {
	blog := &domain.Blog{
		Title:   req.Title,
		Content: req.Content,
		Author:  req.Author,
	}

	err := s.blogService.CreateBlog(blog)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create blog: %v", err)
	}

	return &proto.BlogResponse{
		Blog: &proto.Blog{
			Id:      uint64(blog.ID),
			Title:   blog.Title,
			Content: blog.Content,
			Author:  blog.Author,
		},
	}, nil
}

// Implement other methods (GetBlog, UpdateBlog, DeleteBlog, ListBlogs) similarly

func (s *BlogServer) ListBlogs(ctx context.Context, req *proto.ListBlogsRequest) (*proto.ListBlogsResponse, error) {
	blogs, err := s.blogService.ListBlogs()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list blogs: %v", err)
	}

	var blogResponses []*proto.Blog
	for _, blog := range blogs {
		blogResponses = append(blogResponses, &proto.Blog{
			Id:      uint64(blog.ID),
			Title:   blog.Title,
			Content: blog.Content,
		})
	}

	return &proto.ListBlogsResponse{Blogs: blogResponses}, nil
}

// Implement other methods (GetBlog, UpdateBlog, DeleteBlog) similarly

func (s *BlogServer) GetBlog(ctx context.Context, req *proto.GetBlogRequest) (*proto.BlogResponse, error) {
	blog, err := s.blogService.GetBlog(uint(req.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Blog not found: %v", err)
	}

	return &proto.BlogResponse{
		Blog: &proto.Blog{
			Id:      uint64(blog.ID),
			Title:   blog.Title,
			Content: blog.Content,
		},
	}, nil
}

func (s *BlogServer) UpdateBlog(ctx context.Context, req *proto.UpdateBlogRequest) (*proto.BlogResponse, error) {
	blog := &domain.Blog{
		ID:      uint(req.Id),
		Title:   req.Title,
		Content: req.Content,
	}

	err := s.blogService.UpdateBlog(blog)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update blog: %v", err)
	}

	return &proto.BlogResponse{
		Blog: &proto.Blog{
			Id:      uint64(blog.ID),
			Title:   blog.Title,
			Content: blog.Content,
		},
	}, nil
}

func (s *BlogServer) DeleteBlog(ctx context.Context, req *proto.DeleteBlogRequest) (*proto.DeleteBlogResponse, error) {
	err := s.blogService.DeleteBlog(uint(req.Id))
	if err != nil {
		return &proto.DeleteBlogResponse{
			Success: false,
		}, status.Errorf(codes.Internal, "Failed to delete blog: %v", err)
	}

	return &proto.DeleteBlogResponse{
		Success: true,
	}, nil
}
