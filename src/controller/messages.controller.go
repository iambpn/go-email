package controller

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	imapsession "github.com/iambpn/go-email/src/imapSession"
)

func MessageRouter(router fiber.Router) {
	router.Get("/preview", getAllMessages)
}

// @Summary Get All Preview Message
// @Description Get All Preview Message
// @Tags Message
// @Accept json
// @Produce json
// @Param mailboxName query string true "Mailbox Name"
// @Param page query number false "Page Number" default(1)
// @Param pageSize query number false "Page Size" default(20)
// @Router /message/preview [get]
func getAllMessages(c *fiber.Ctx) error {
	mailbox := c.Query("mailboxName", "")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)

	if mailbox == "" {
		return errors.New("mailboxName is required")
	}

	prevMsg, err := imapsession.IW.GetPreviewMessages(mailbox, page, pageSize)

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"data": prevMsg})
}
