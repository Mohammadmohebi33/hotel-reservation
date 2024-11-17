package api

import (
	"fmt"
	"github.com/Mohammadmohebi33/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("not authorization")
	}
	if !user.IsAdmin {
		return fmt.Errorf("not admin")
	}
	return c.Next()
}
