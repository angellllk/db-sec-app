package core

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func StartApp(db *sql.DB) error {
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Minute * 2,
		WriteTimeout: time.Minute * 2,
	})

	app.Use(cors.New())

	api := app.Group("/dealership-api/v1/")
	{
		api.Get("/status", func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		api.Post("/register-buyer", func(c *fiber.Ctx) error {
			return registerBuyer(c, db)
		})

		api.Post("/import-car", func(c *fiber.Ctx) error {
			return importCar()
		})

		api.Post("/register-sale", func(c *fiber.Ctx) error {
			return createSale()
		})

		api.Get("/sales-details/:name", func(c *fiber.Ctx) error {
			return salesDetails(c.Params("name"))
		})
	}

	return app.Listen(":80")
}
