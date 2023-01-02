package routes

import (
	"github.com/gofiber/fiber/v2"
)

func web() {
	app := fiber.New()

	app.Group("/", func(c *fiber.Ctx) error {
		// Set หรือ Get Header
		c.Set("Version", "v1")
		return c.Next()
	})
	app.Get("/signup")

}
