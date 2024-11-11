package main

import (
	"log"

	"github.com/toffysoft/go-hexagonal-example/internal/adapters/handlers"
	"github.com/toffysoft/go-hexagonal-example/internal/adapters/repositories"
	"github.com/toffysoft/go-hexagonal-example/internal/core/services"
	"github.com/toffysoft/go-hexagonal-example/internal/infrastructure/config"
	"github.com/toffysoft/go-hexagonal-example/internal/infrastructure/database"
	"github.com/toffysoft/go-hexagonal-example/pkg/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repositories
	blogRepo := repositories.NewBlogRepository(db)

	// Initialize services
	blogService := services.NewBlogService(blogRepo)

	// Initialize handlers
	blogHandler := handlers.NewBlogHandler(blogService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Add middlewares
	app.Use(recover.New()) // Recover from panics and sends 500 internal server error
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Setup routes
	api := app.Group("/api")
	v1 := api.Group("/v1")

	blogs := v1.Group("/blogs")
	blogs.Post("/", blogHandler.CreateBlog)
	blogs.Get("/:id", blogHandler.GetBlog)
	blogs.Put("/:id", blogHandler.UpdateBlog)
	blogs.Delete("/:id", blogHandler.DeleteBlog)
	blogs.Get("/", blogHandler.ListBlogs)

	// Start server
	log.Printf("Starting server on %s", cfg.ServerAddress)
	log.Fatal(app.Listen(cfg.ServerAddress))
}

func customErrorHandler(c *fiber.Ctx, err error) error {

	if appErr, ok := err.(errors.AppError); ok {
		return c.Status(appErr.StatusCode()).JSON(fiber.Map{
			"error": appErr.Error(),
		})
	}

	// Default 500 statuscode
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		// Override status code if fiber.Error type
		code = e.Code
	}

	// Set Content-Type: text/plain; charset=utf-8
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	// Return statuscode with error message
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
