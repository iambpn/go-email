package controller

import (
	"github.com/gofiber/fiber/v2"
	imapsession "github.com/iambpn/go-email/src/imapSession"
)

func MessageRouter(router fiber.Router) {
	router.Get("/:mailbox", getAllPreviewMessages)
	router.Get("/:mailbox/:uid", getMessage)
	router.Put("/:mailbox/:uid", updateFlags)
}

// @Summary Get All Preview Message in Desc Order
// @Description Get All Preview Message in Desc Order
// @Tags Message
// @Produce json
// @Param mailbox path string true "Mailbox Name"
// @Param page query number false "Page Number" default(1)
// @Param pageSize query number false "Page Size" default(20)
// @Router /message/{mailbox} [get]
func getAllPreviewMessages(c *fiber.Ctx) error {
	mailbox := c.Params("mailbox", "")
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 20)

	prevMsg, err := imapsession.IW.GetPreviewMessages(mailbox, page, pageSize)

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"data": prevMsg})
}

// @Summary Get Message Details
// @Description Get Message Details
// @Tags Message
// @Produce json
// @Param mailbox path string true "Mailbox Name"
// @Param uid path number true "Uid Number"
// @Router /message/{mailbox}/{uid} [get]
func getMessage(c *fiber.Ctx) error {
	mailbox := c.Params("mailbox")
	uid, err := c.ParamsInt("uid")

	if err != nil {
		return err
	}

	message, err := imapsession.IW.GetMessage(mailbox, uint32(uid))

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": message,
	})
}

type UpdateFlagsBody struct {
	AddFlags    []string `json:"addFlags" example:"\\Seen,\\Answered" validate:"required" enums:"\\Seen,\\Answered,\\Flagged,\\Deleted,\\Draft,\\Recent"`
	RemoveFlags []string `json:"removeFlags" example:"\\Seen,\\Answered" validate:"required" enums:"\\Seen,\\Answered,\\Flagged,\\Deleted,\\Draft,\\Recent"`
}

// @Summary Update Message flags
// @Description Update Message flags
// @Tags Message
// @Produce json
// @Param mailbox path string true "Mailbox Name"
// @Param uid path number true "Uid Number"
// @Param updateBody body UpdateFlagsBody true "Update Body"
// @Router /message/{mailbox}/{uid} [put]
func updateFlags(c *fiber.Ctx) error {
	mailbox := c.Params("mailbox", "")
	uid, err := c.ParamsInt("uid")

	if err != nil {
		return err
	}

	updateBody := &UpdateFlagsBody{}
	err = c.BodyParser(updateBody)

	if err != nil {
		return err
	}

	err = imapsession.IW.UpdateMessage(mailbox, uint32(uid), updateBody.AddFlags, updateBody.RemoveFlags)

	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"message": "success",
		},
	})
}
