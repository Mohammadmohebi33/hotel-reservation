package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("-- JWT auth")

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("token not found")
	}

	if err := parseJWTToken(token[0]); err != nil {
		return err
	}
	return nil
}

func parseJWTToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthrize")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		return fmt.Errorf("unauthrize")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	}
	return fmt.Errorf("unauthrize")
}
