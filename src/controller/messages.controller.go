package controller

import "github.com/gofiber/fiber/v2"

func MessageRouter(router fiber.Router) {
	router.Get("/preview", getAllMessages)
}

// @Summary Get All Preview Message
// @Description Get All Preview Message
// @Tags Message
// @Accept json
// @Produce json
// @Router /message/preview [get]
func getAllMessages(c *fiber.Ctx) error {
	prevMsg, err := iw.GetPreviewMessages("INBOX", 1, 10)

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"data": prevMsg})
}
