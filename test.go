package main;

import (
	"github.com/gofiber/fiber/v2"
	"github.com/brianvoe/gofakeit/v6"
)

func test() {
	// Initialize Fiber app
	app := fiber.New()

	// Route to generate mock data
	app.Get("/mock", func(c *fiber.Ctx) error {
		// Generate mock data
		mockData := map[string]interface{}{
			"name":      gofakeit.Name(),
			"email":     gofakeit.Email(),
			"address":   gofakeit.Address().Address,
			"phone":     gofakeit.Phone(),
			"company":   gofakeit.Company(),
			"job_title": gofakeit.JobTitle(),
		}

		// Send mock data as JSON response
		return c.JSON(mockData)
	})

	// Start the server on port 3000
	app.Listen(":3000")
}
