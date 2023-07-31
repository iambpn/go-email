package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iambpn/go-email/src/config"
)

func main() {
	// load env
	config.LoadConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// un-caught error
			log.Println(err)
			c.Status(500).JSON(fiber.Map{
				"message": err.Error(),
				"path":    c.Path(),
				"time":    time.Now().UTC().String(),
			})
			return nil
		},
	})

	addr := fmt.Sprintf("%s:%s", config.GetConfig("host", ""), config.GetConfig("port", "3000"))
	err := app.Listen(addr)

	if err != nil {
		log.Fatal(err)
	}
}
