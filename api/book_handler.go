package api

import (
	"github.com/Mohammadmohebi33/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type BookHandler struct {
	store *db.Store
}

func NewBookHandler(store *db.Store) *BookHandler {
	return &BookHandler{
		store: store,
	}
}

func (h BookHandler) HandleCancelBook(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}

	user, err := getAuthUser(c)
	if err != nil {
		return err
	}

	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResponse{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericResponse{
		Type: "Msg",
		Msg:  "updated",
	})
}

func (h *BookHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(bookings)
}

func (h *BookHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(booking)
}
