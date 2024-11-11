package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func JWTAuthentication(c *fiber.Ctx) error {

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("token not found")
	}

	clams, err := validateJWTToken(token[0])
	if err != nil {
		return err
	}
	fmt.Println(clams)

	expireFloat := clams["exp"].(float64)
	expire := int64(expireFloat)
	if time.Now().Unix() > expire {
		return fmt.Errorf("token expired")
	}
	return c.Next()
}

func validateJWTToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthrize")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("unauthrize")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
