package middleware

import (
	"fmt"
	"strings"

	"github.com/Ra1nz0r/zero_agency/internal/logger"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware возвращает middleware, который проверяет JWT-токен в заголовке Authorization.
func JWTMiddleware(jwtSecret string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Получаем значение заголовка "Authorization" из входящего запроса.
		authHeader := c.Get("Authorization")

		// Если заголовок пуст, логируем ошибку и возвращаем статус 401 Unauthorized с сообщением об ошибке.
		if authHeader == "" {
			logger.Zap.Error("Missing authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Убираем префикс "Bearer " из заголовка Authorization, чтобы получить чистый JWT-токен.
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Если после удаления префикса строка не изменилась значит формат токена неверный.
		if tokenString == authHeader {
			logger.Zap.Error("Invalid token format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}

		// Парсим и проверием JWT-токен, используя переданный секретный ключ.
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Проверяем, что токен был подписан с использованием HMAC-алгоритма (например, HS256).
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				// Если метод подписи не HMAC, возвращаем ошибку.
				logger.Zap.Error("unexpected signing method")
				return nil, fmt.Errorf("unexpected signing method")
			}
			// Возвращаем секретный ключ для проверки подписи токена.
			return []byte(jwtSecret), nil
		})

		// Если произошла ошибка парсинга токена или он оказался недействительным, логируем ошибку и возвращаем статус 401 Unauthorized.
		if err != nil || !token.Valid {
			logger.Zap.Error(err.Error())
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		return c.Next()
	}
}
