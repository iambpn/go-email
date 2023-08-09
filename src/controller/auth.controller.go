package controller

import "github.com/gofiber/fiber/v2"

func AuthRouter(router fiber.Router) {
	router.Post("/login", login)
	router.Get("/logout", logout)
}

// @summary Login
// @Description Login API
// @Tags Auth
// @Accept json
// @Produce json
// @Param username body string true "email address" example(test@test.com) default(test@test.com)
// @Param password body string true "password" example(password123) default(password123)
// @Success 200
// @Failure 500,404
// @Router /login [post]
func login(c *fiber.Ctx) error {
	return nil
}

// @Summary		Logout
// @Description	Logout API
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Success		200
// @Failure		500,404
// @Router			/logout [get]
func logout(c *fiber.Ctx) error {
	return nil
}
