package services

import (
	"database/sql"
	"math"
	"strconv"
	"time"

	"fmt"

	cfg "github.com/Ra1nz0r/zero_agency/internal/config"
	"github.com/Ra1nz0r/zero_agency/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

// GenerateJWT генерирует JWT (JSON Web Token) для предоставленного логина и пароля.
func GenerateJWT(username, jwtSecret string, hours time.Duration) (string, error) {
	logger.Zap.Debug("-> `GenerateJWT` - calling function.")

	// Создаем объект claims, который будет содержать полезную информацию (Payload) внутри JWT
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(hours).Unix(),
	}

	logger.Zap.Debug("Creating new token.")

	// Создаем новый токен с использованием алгоритма HMAC (HS256) и переданных claims.
	// jwt.SigningMethodHS256 - алгоритм для подписи токена.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	logger.Zap.Debug("Sigining the token.")

	// Подписываем токен с использованием секретного ключа (jwtSecret).
	// Метод SignedString берет строковое представление секретного ключа (преобразованного в байты).
	res, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	logger.Zap.Debug("-> `GenerateJWT` - successful called.")

	return res, nil
}

// RunMigrations запускает миграции, при повторном запуске программы, функция
// будет проверять, были ли применены миграции ранее, и если все миграции уже выполнены,
// она не будет применять их снова. Если в директории migration появились новые файлы миграций
// с более высокими номерами (например, 000002_add_column.up.sql), библиотека применит их.
func RunMigrations(db *sql.DB, cfg cfg.Config) error {
	logger.Zap.Debug("-> `RunMigrations` - calling function.")

	// Создаем драйвер миграции
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	logger.Zap.Debug("Crating new migration instance.")

	// Настраиваем новый экземпляр миграций из исходного URL и существующего экземпляра базы данных.
	m, err := migrate.NewWithDatabaseInstance(
		cfg.MigrationPath,
		cfg.DatabaseName,
		driver,
	)
	if err != nil {
		return err
	}

	logger.Zap.Debug("Running migration UP.")

	// Запускаем миграции Up
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	logger.Zap.Debug("-> `RunMigrations` - successful called.")

	return nil
}

// StringToInt32WithOverflowCheck преобразует строку в int32 с проверкой переполнения
func StringToInt32WithOverflowCheck(s string) (int32, error) {
	logger.Zap.Debug("-> `StringToInt32WithOverflowCheck` - calling function.")

	logger.Zap.Debug("Parsing string to int.")

	// Преобразуем строку в int64
	id64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	logger.Zap.Debug("Checking on range and overflow.")

	// Проверяем, не выходит ли значение за пределы диапазона int32
	if id64 > math.MaxInt32 || id64 < math.MinInt32 {
		return 0, fmt.Errorf("number out of range for int32")
	}

	logger.Zap.Debug("-> `StringToInt32WithOverflowCheck` - successful called.")

	return int32(id64), nil //nolint:gosec
}
