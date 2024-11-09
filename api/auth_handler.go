package api

import (
	"errors"
	"fmt"
	"github.com/Mohammadmohebi33/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthentication(c *fiber.Ctx) error {
	var AuthParam AuthParam
	if err := c.BodyParser(&AuthParam); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), AuthParam.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid credentials")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(AuthParam.Password))
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}

	fmt.Println("authenticated ->", user)
	return nil
}
