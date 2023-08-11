package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/iambpn/go-email/docs"
	"github.com/iambpn/go-email/src/config"
	"github.com/iambpn/go-email/src/controller"
)

// @title						Go Email Service API
// @version					1.0
// @description				This is Go Email Service Build with Go-Fiber and Go-Imap
// @contact.name				Bipin Maharjan
// @contact.email				bipinmhr10@gmail.com
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath					/
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description				Bearer token authorization
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

	app.Get("/swagger/*", swagger.HandlerDefault)

	controller.AuthRouter(app.Group("/auth"))

	addr := fmt.Sprintf("%s:%s", config.GetConfig("host", ""), config.GetConfig("port", "3000"))
	err := app.Listen(addr)

	if err != nil {
		log.Fatal(err)
	}
}
