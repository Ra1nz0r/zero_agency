package services

import (
	"database/sql"
	"math"
	"strconv"

	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

// RunMigrations запускает миграцию Up по указанному пути.
func RunMigrations(databaseURL, migrationPath string) error {
	m, err := migrate.New(migrationPath, databaseURL)
	if err != nil {
		return err
	}

	// Применение миграций
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// TableExists проверяет существование table в базе данных.
func TableExists(db *sql.DB, tableName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT FROM pg_tables
			WHERE schemaname = 'public' OR schemaname = 'private'
			AND tablename = $1
		);`
	// Используем параметризованный запрос, где $1 — это плейсхолдер для tableName
	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// StringToInt32WithOverflowCheck преобразует строку в int32 с проверкой переполнения
func StringToInt32WithOverflowCheck(s string) (int32, error) {
	// Преобразуем строку в int64
	id64, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid string to number conversion: %w", err)
	}

	// Проверяем, не выходит ли значение за пределы диапазона int32
	if id64 > math.MaxInt32 || id64 < math.MinInt32 {
		return 0, fmt.Errorf("number out of range for int32")
	}

	// Возвращаем преобразованное значение
	return int32(id64), nil
}
