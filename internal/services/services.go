package services

import (
	"database/sql"
	"math"
	"strconv"
	"time"

	"fmt"

	cfg "github.com/Ra1nz0r/zero_agency/internal/config"
	"github.com/Ra1nz0r/zero_agency/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func GenerateJWT(lr models.LoginRequest, jwtSecret string, hours time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"username": lr.Username,
		"password": lr.Password,
		"exp":      time.Now().Add(hours).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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
	return int32(id64), nil
}
