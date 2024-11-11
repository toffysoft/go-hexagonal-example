package integration_test

import (
	"context"
	"log"
	"net"
	"testing"

	bloggrpc "github.com/toffysoft/go-hexagonal-example/internal/adapters/grpc"
	"github.com/toffysoft/go-hexagonal-example/internal/adapters/grpc/proto"
	"github.com/toffysoft/go-hexagonal-example/internal/adapters/repositories"
	"github.com/toffysoft/go-hexagonal-example/internal/core/services"
	"github.com/toffysoft/go-hexagonal-example/internal/infrastructure/database"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	// Setup your actual dependencies here
	db, err := database.InitTestDB()
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	blogRepo := repositories.NewBlogRepository(db)
	blogService := services.NewBlogService(blogRepo)
	blogServer := bloggrpc.NewBlogServer(blogService)

	proto.RegisterBlogServiceServer(s, blogServer)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestCreateBlogIntegration(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewBlogServiceClient(conn)

	req := &proto.CreateBlogRequest{
		Title:   "Integration Test Blog",
		Content: "This is an integration test",
		Author:  "Test Author",
	}

	resp, err := client.CreateBlog(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Title, resp.Blog.Title)
	assert.Equal(t, req.Content, resp.Blog.Content)
	assert.Equal(t, req.Author, resp.Blog.Author)
	assert.NotZero(t, resp.Blog.Id)
}

// Implement other integration test cases...

func TestListBlogsIntegration(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewBlogServiceClient(conn)

	req := &proto.ListBlogsRequest{}

	resp, err := client.ListBlogs(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Blogs)
}

// Implement other integration test cases...

func TestGetBlogIntegration(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewBlogServiceClient(conn)

	createBlogReq := &proto.CreateBlogRequest{
		Title:   "Integration Test Blog",
		Content: "This is an integration test",
		Author:  "Test Author",
	}

	createBlogResp, _ := client.CreateBlog(ctx, createBlogReq)

	req := &proto.GetBlogRequest{
		Id: createBlogResp.Blog.Id,
	}

	resp, err := client.GetBlog(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotZero(t, resp.Blog.Id)
}

func TestUpdateBlogIntegration(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewBlogServiceClient(conn)

	createBlogReq := &proto.CreateBlogRequest{
		Title:   "Integration Test Blog",
		Content: "This is an integration test",
		Author:  "Test Author",
	}

	createBlogResp, _ := client.CreateBlog(ctx, createBlogReq)

	req := &proto.UpdateBlogRequest{
		Id:      createBlogResp.Blog.Id,
		Title:   "Updated Integration Test Blog",
		Content: "This is an updated integration test",
	}

	resp, err := client.UpdateBlog(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Title, resp.Blog.Title)
	assert.Equal(t, req.Content, resp.Blog.Content)
	assert.NotZero(t, resp.Blog.Id)
}

func TestDeleteBlogIntegration(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := proto.NewBlogServiceClient(conn)

	createBlogReq := &proto.CreateBlogRequest{
		Title:   "Integration Test Blog",
		Content: "This is an integration test",
		Author:  "Test Author",
	}

	createBlogResp, _ := client.CreateBlog(ctx, createBlogReq)

	req := &proto.DeleteBlogRequest{
		Id: createBlogResp.Blog.Id,
	}

	resp, err := client.DeleteBlog(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, true, resp.Success)
}
