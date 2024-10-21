package services

import (
	"database/sql"
	"math"
	"strconv"
	"time"

	"fmt"

	cfg "github.com/Ra1nz0r/zero_agency/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

// GenerateJWT генерирует JWT (JSON Web Token) для предоставленного логина и пароля.
//
// Параметры:
// - lr (models.LoginRequest): структура, содержащая логин и пароль пользователя.
// - jwtSecret (string): секретный ключ для подписания токена.
// - hours (time.Duration): количество времени, на которое токен будет действительным.
//
// Возвращает:
// - (string): сгенерированный JWT в виде строки.
// - (error): ошибка, если токен не удалось сгенерировать.
func GenerateJWT(username, jwtSecret string, hours time.Duration) (string, error) {
	// Создаем объект claims, который будет содержать полезную информацию (Payload) внутри JWT
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(hours).Unix(),
	}

	// Создаем новый токен с использованием алгоритма HMAC (HS256) и переданных claims.
	// jwt.SigningMethodHS256 - алгоритм для подписи токена.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с использованием секретного ключа (jwtSecret).
	// Метод SignedString берет строковое представление секретного ключа (преобразованного в байты).
	res, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return res, nil
}

// RunMigrations запускает миграции, при повторном запуске программы, функция
// будет проверять, были ли применены миграции ранее, и если все миграции уже выполнены,
// она не будет применять их снова. Если в директории migration появились новые файлы миграций
// с более высокими номерами (например, 000002_add_column.up.sql), библиотека применит их.
func RunMigrations(db *sql.DB, cfg cfg.Config) error {
	// Создаем драйвер миграции
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Настраиваем миграции (указываем путь к файлам миграций)
	m, err := migrate.NewWithDatabaseInstance(
		cfg.MigrationPath,
		cfg.DatabaseName,
		driver,
	)
	if err != nil {
		return err
	}

	// Запускаем миграции Up
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// StringToInt32WithOverflowCheck преобразует строку в int32 с проверкой переполнения
func StringToInt32WithOverflowCheck(s string) (int32, error) {
	// Преобразуем строку в int64
	id64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	// Проверяем, не выходит ли значение за пределы диапазона int32
	if id64 > math.MaxInt32 || id64 < math.MinInt32 {
		return 0, fmt.Errorf("number out of range for int32")
	}

	// Возвращаем преобразованное значение
	return int32(id64), nil //nolint:gosec
}
