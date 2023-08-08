package controller

import "github.com/gofiber/fiber/v2"

func AuthRouter(router fiber.Router) {
	router.Get("/logout", logout)
}

func login(c *fiber.Ctx) {
	return
}

//	@Summary		Logout
//	@Description	Logout
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@failure		500,404
//	@Router			/logout [get]
func logout(c *fiber.Ctx) error {
	return nil
}
