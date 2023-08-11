package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	imapwrapper "github.com/iambpn/go-email/pkg/imap_wrapper"
)

func AuthRouter(router fiber.Router) {
	router.Post("/login", login)
	router.Get("/logout", logout)
}

type LoginBody struct {
	Username string `json:"username" example:"test@test.com" validate:"required" format:"string"`
	Password string `json:"password" example:"password1122" validate:"required" format:"string"`
	Host     string `json:"host" example:"outlook.office365.com" validate:"required"`
	Port     string `json:"port" example:"993" validate:"required"`
}

var iw imapwrapper.ImapWrapper

// @summary		Login
// @Description	Login API
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			RequestBody	body	LoginBody	false	"Login Body"
// @Success		200
// @Failure		500,404
// @Router			/auth/login [post]
func login(c *fiber.Ctx) error {
	loginBody := LoginBody{}
	err := c.BodyParser(&loginBody)

	if err != nil {
		return err
	}

	iw = imapwrapper.ImapWrapper{
		Host:     loginBody.Host,
		Port:     loginBody.Port,
		Username: loginBody.Username,
		Password: loginBody.Password,
	}

	err = iw.Connect()

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Login Successful",
	})
}

// @Summary		Logout
// @Description	Logout API
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		500,403
// @Router			/auth/logout [get]
// //@Security		ApiKeyAuth
func logout(c *fiber.Ctx) error {
	if iw.Host != "" {
		fmt.Print(iw)
		iw.Logout()
		return c.JSON(fiber.Map{
			"message": "Logout Successful",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Unable to logout. Imap Wrapper is not initialized.",
	})
}
