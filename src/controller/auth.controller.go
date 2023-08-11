package controller

import "github.com/gofiber/fiber/v2"

func AuthRouter(router fiber.Router) {
	router.Post("/login", login)
	router.Get("/logout", logout)
}

type LoginBody struct {
	Username string `json:"username" example:"test@test.com" validate:"required" format:"string"`
	Password string `json:"password" example:"password123" validate:"required" format:"string"`
}

// @summary		Login
// @Description	Login API
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			RequestBody	body	LoginBody	false	"Login Body"
// @Success		200
// @Failure		500,404
// @Router			/login [post]
func login(c *fiber.Ctx) error {
	return nil
}

// @Summary		Logout
// @Description	Logout API
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		500,403
// @Router			/logout [get]
// @Security		ApiKeyAuth
func logout(c *fiber.Ctx) error {
	return nil
}
