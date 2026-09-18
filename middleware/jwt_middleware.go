package middleware

import (
	"net/http"
	"os"
	"strings"

	"bicycle-rent-api/helper"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "missing authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "invalid authorization header",
			})
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(
			tokenString,
			&helper.JWTClaims{},
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}

				return []byte(os.Getenv("JWT_SECRET")), nil
			},
		)

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "invalid or expired token",
			})
		}

		claims, ok := token.Claims.(*helper.JWTClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "invalid token claims",
			})
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		return next(c)
	}
}
