package middlewear

import (
	"bookmark_service/pkg/secretFuncs"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Извлекаем заголовок Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "отсутствует заголовок Authorization")
			}

			// 2. Проверяем формат "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "неверный формат токена")
			}

			tokenString := parts[1]

			// 3. Парсим и валидируем токен
			claims := &secretFuncs.MyCustomClaims{} // Твоя структура из Login
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Важно: проверяем метод подписи
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "неверный метод подписи")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "недействительный или просроченный токен")
			}

			// 4. Сохраняем данные пользователя в контекст Echo
			// Теперь в любом обработчике можно достать ID через c.Get("user_id")
			c.Set("user_id", claims.UserID)

			return next(c)
		}
	}
}
