package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendErrorResponse sends a JSON error response
func SendErrorResponse(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(ErrorResponse{
		Success: false,
		Message: message,
	})
}

// SendSuccessResponse sends a JSON success response
func SendSuccessResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ValidatorErrors(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", e.Field()))
			case "min":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param()))
			case "max":
				errorMessages = append(errorMessages, fmt.Sprintf("%s must not be longer than %s characters", e.Field(), e.Param()))
			default:
				errorMessages = append(errorMessages, fmt.Sprintf("%s is not valid", e.Field()))
			}
		}
		return fmt.Sprintf("Validation error: %s", errorMessages)
	}
	return "Validation error occurred"
}
