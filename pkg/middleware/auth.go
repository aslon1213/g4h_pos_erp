package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(c *fiber.Ctx) {
	token := c.Get("Authorization")
	if token == "" {
		c.Status(http.StatusUnauthorized)
		return
	}

	claims, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	fmt.Println(claims)

	c.Locals("user", "test")

}
