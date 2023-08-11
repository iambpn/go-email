package controller

import (
	"github.com/gofiber/fiber/v2"
	imapsession "github.com/iambpn/go-email/src/imapSession"
)

func MailBoxRouter(router fiber.Router) {
	router.Get("/", getMailboxes)
}

// @Summary Get list of Mailboxes
// @Description get List of Mailboxes
// @Tags Mailbox
// @Accept json
// @Produce json
// @Router /mailbox [get]
func getMailboxes(c *fiber.Ctx) error {
	mailboxes, err := imapsession.IW.GetMailBoxes()
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": mailboxes,
	})
}
