package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"com.spruce.flapi/ext"
	"github.com/gofiber/fiber/v2"
)



func main() {
	app := fiber.New()
	gen := ext.Gen{}

	// Fiber router
	app.All("/gen", func(c *fiber.Ctx) error {
		// Parse query parameters
		amountStr := c.Query("amount", "1") // Default to 1
		amount, err := strconv.Atoi(amountStr)
		if err != nil || amount < 1 || amount > 100 {
			return c.Status(400).JSON(fiber.Map{"error": "Amount must be between 1 and 100"})
		}

		// Parse schema from query or body
		schema := c.Query("schema", "")
		if len(schema) == 0 {
			body := c.Body()
			if len(body) > 0 {
				schema = string(body)
			}
		}

		if len(schema) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Schema is required"})
		}

		// Parse schema as JSON
		var schemaMap map[string]interface{}
		if err := json.Unmarshal([]byte(schema), &schemaMap); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON schema"})
		}
		fmt.Println(schemaMap)

		// Generate mock data
		data := gen.GenerateObject(schemaMap, amount)
		return c.JSON(data)
	})

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
