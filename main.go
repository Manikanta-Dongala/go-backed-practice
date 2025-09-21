package main

import (
	"github.com/gofiber/fiber/v2"
)


func main() {
	app := fiber.New()
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg":"Hello world!"})
	})
	app.Listen(":3000")
}
