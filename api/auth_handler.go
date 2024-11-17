package api

import (
	"errors"
	"fmt"
	"github.com/Mohammadmohebi33/hotel-reservation/db"
	"github.com/Mohammadmohebi33/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
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

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResponse struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentialResponse(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResponse{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

func (h *AuthHandler) HandleAuthentication(c *fiber.Ctx) error {
	var AuthParam AuthParam
	if err := c.BodyParser(&AuthParam); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByEmail(c.Context(), AuthParam.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentialResponse(c)
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(AuthParam.Password))
	if err != nil {
		return invalidCredentialResponse(c)
	}

	resp := AuthResponse{
		User:  user,
		Token: CreateTokenFromUser(user),
	}
	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expireTime := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   expireTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	tokenSrt, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenSrt
}
