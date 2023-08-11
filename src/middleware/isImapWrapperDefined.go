package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	imapsession "github.com/iambpn/go-email/src/imapSession"
)

func IsImapWrapperDefined(c *fiber.Ctx) error {
	if imapsession.IW.Host == "" {
		return errors.New("imap wrapper is not initialized")
	}

	return c.Next()
}
