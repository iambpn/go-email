package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	imapwrapper "github.com/iambpn/go-email/pkg/imap_wrapper"
	imapsession "github.com/iambpn/go-email/src/imapSession"
	"github.com/iambpn/go-email/src/middleware"
)

func AuthRouter(router fiber.Router) {
	router.Post("/login", login)
	router.Get("/logout", middleware.IsImapWrapperDefined, logout)
}

type LoginBody struct {
	Username string `json:"username" example:"test@test.com" validate:"required" format:"string"`
	Password string `json:"password" example:"password1122" validate:"required" format:"string"`
	Host     string `json:"host" example:"outlook.office365.com" validate:"required"`
	Port     string `json:"port" example:"993" validate:"required"`
}

// @summary		Login
// @Description	Login API
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			RequestBody	body LoginBody	false	"Login Body"
// @Success		200
// @Failure		500,404
// @Router			/auth/login [post]
func login(c *fiber.Ctx) error {
	loginBody := LoginBody{}
	err := c.BodyParser(&loginBody)

	if err != nil {
		return err
	}

	imapsession.IW = imapwrapper.ImapWrapper{
		Host:     loginBody.Host,
		Port:     loginBody.Port,
		Username: loginBody.Username,
		Password: loginBody.Password,
	}

	err = imapsession.IW.Connect()

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
	if imapsession.IW.Host != "" {
		fmt.Print(imapsession.IW)
		imapsession.IW.Logout()
		return c.JSON(fiber.Map{
			"message": "Logout Successful",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Unable to logout. Imap Wrapper is not initialized.",
	})
}
