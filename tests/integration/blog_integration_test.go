package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/toffysoft/go-hexagonal-example/internal/adapters/handlers"
	"github.com/toffysoft/go-hexagonal-example/internal/adapters/repositories"
	"github.com/toffysoft/go-hexagonal-example/internal/core/domain"
	"github.com/toffysoft/go-hexagonal-example/internal/core/services"
	"github.com/toffysoft/go-hexagonal-example/internal/infrastructure/database"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	db, _ := database.InitTestDB()
	blogRepo := repositories.NewBlogRepository(db)
	blogService := services.NewBlogService(blogRepo)
	blogHandler := handlers.NewBlogHandler(blogService)

	app := fiber.New()
	api := app.Group("/api")
	v1 := api.Group("/v1")

	blogs := v1.Group("/blogs")
	blogs.Post("/", blogHandler.CreateBlog)
	blogs.Get("/:id", blogHandler.GetBlog)
	blogs.Put("/:id", blogHandler.UpdateBlog)
	blogs.Delete("/:id", blogHandler.DeleteBlog)
	blogs.Get("/", blogHandler.ListBlogs)

	return app
}

func TestCreateBlog(t *testing.T) {
	app := setupTestApp()

	blog := domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog post",
		Author:  "Test Author",
	}

	payload, _ := json.Marshal(blog)

	req := httptest.NewRequest("POST", "/api/v1/blogs", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Equal(t, "Blog created successfully", response["message"])
	assert.NotNil(t, response["data"])

}

func TestGetBlog(t *testing.T) {
	app := setupTestApp()

	// First, create a blog
	blog := domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog post",
		Author:  "Test Author",
	}

	payload, _ := json.Marshal(blog)

	req := httptest.NewRequest("POST", "/api/v1/blogs", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	var createResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&createResponse)

	createdBlogID := uint(createResponse["data"].(map[string]interface{})["id"].(float64))

	// Now, test getting the blog
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%d", createdBlogID), nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var getResponse map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&getResponse)

	assert.Equal(t, "Blog retrieved successfully", getResponse["message"])
	assert.Equal(t, "Test Blog", getResponse["data"].(map[string]interface{})["title"])
}

// Implement similar tests for UpdateBlog, DeleteBlog, and ListBlogs
func TestUpdateBlog(t *testing.T) {
	app := setupTestApp()

	// First, create a blog
	blog := domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog post",
		Author:  "Test Author",
	}

	payload, _ := json.Marshal(blog)

	req := httptest.NewRequest("POST", "/api/v1/blogs", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	var createResponse map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&createResponse)

	createdBlogID := uint(createResponse["data"].(map[string]interface{})["id"].(float64))

	// Now, test updating the blog
	blog.Title = "Updated Test Blog"

	payload, _ = json.Marshal(blog)

	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/blogs/%d", createdBlogID), bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)

	var updateResponse map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&updateResponse)

	assert.Equal(t, "Blog updated successfully", updateResponse["message"])

	// Now, test getting the blog to verify the update

	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%d", createdBlogID), nil)

	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var getResponse map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&getResponse)

	assert.Equal(t, "Blog retrieved successfully", getResponse["message"])

	assert.Equal(t, "Updated Test Blog", getResponse["data"].(map[string]interface{})["title"])

}

func TestDeleteBlog(t *testing.T) {
	app := setupTestApp()

	// First, create a blog
	blog := domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog post",
		Author:  "Test Author",
	}

	payload, _ := json.Marshal(blog)

	req := httptest.NewRequest("POST", "/api/v1/blogs", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	var createResponse map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&createResponse)

	createdBlogID := uint(createResponse["data"].(map[string]interface{})["id"].(float64))

	// Now, test deleting the blog

	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/blogs/%d", createdBlogID), nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)

	var deleteResponse map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&deleteResponse)

	assert.Equal(t, "Blog deleted successfully", deleteResponse["message"])

	// Now, test getting the blog to verify that it was deleted

	req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/blogs/%d", createdBlogID), nil)

	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var getResponse map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&getResponse)

	assert.Equal(t, fmt.Sprintf("Blog with ID %d not found", createdBlogID), getResponse["message"])
}

func TestListBlogs(t *testing.T) {
	app := setupTestApp()

	// First, create a blog
	blog := domain.Blog{
		Title:   "Test Blog",
		Content: "This is a test blog post",
	}

	payload, _ := json.Marshal(blog)

	req := httptest.NewRequest("POST", "/api/v1/blogs", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	var createResponse map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&createResponse)

	// Now, test listing all blogs

	req = httptest.NewRequest("GET", "/api/v1/blogs", nil)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var listResponse map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&listResponse)

	assert.Equal(t, "Blogs retrieved successfully", listResponse["message"])
	assert.NotNil(t, listResponse["data"])
}
