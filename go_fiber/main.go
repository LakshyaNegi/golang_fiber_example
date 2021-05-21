package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	app := fiber.New()
	initDB()
	//defer DBConn.Close()

	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!!")
	})

	app.Post("/login", login)
	app.Post("/register", register)

	app.Get("/shoes", GetShoes)
	app.Get("/shoes/:id", GetShoe)
	app.Post("/shoes", Protected(), CreateShoe)
	app.Delete("/shoes/:id", Protected(), DeleteShoe)
	app.Put("/shoes/:id", Protected(), UpdateShoe)

	err := app.Listen(":4000")
	if err != nil {
		panic(err)
	}

}
