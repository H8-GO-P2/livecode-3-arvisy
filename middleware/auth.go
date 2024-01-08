package middleware

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo"
)

var jwtSecret = []byte("your-secret-key")

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		if tokenString == "" {
			return c.JSON(401, echo.Map{
				"message": "unauthorized",
			})
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(401, echo.Map{
				"message": "unauthorized",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(401, echo.Map{
				"message": "unauthorized",
			})
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			return c.JSON(401, echo.Map{
				"message": "unauthorized",
			})
		}

		c.Set("user_id", int(userID))

		return next(c)
	}
}
