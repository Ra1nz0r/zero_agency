package server

import (
	"database/sql"
	"time"

	"fmt"

	cfg "github.com/Ra1nz0r/zero_agency/internal/config"
	"github.com/Ra1nz0r/zero_agency/internal/logger"
)

// Connect открывает подключение к базе данных и настраивает пул соединений.
// Параметры подключения берутся из структуры конфигурации cfg.
// Функция возвращает объект sql.DB для работы с базой данных или ошибку, если подключение не удалось.
func Connect(cfg *cfg.Config) (*sql.DB, error) {
	logger.Zap.Debug("-> `Connect` - launching function.")

	// Формируем URL для подключения к базе данных PostgreSQL с параметрами из конфигурации.
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)

	logger.Zap.Debug(fmt.Sprintf("The link to connect to the DB generated: %s", dbURL))
	logger.Zap.Debug("Opening connection to the database.")

	// Открываем подключение к базе данных с помощью драйвера, указанного в конфигурации.
	db, err := sql.Open(cfg.DatabaseDriver, dbURL)
	if err != nil {
		return nil, err
	}

	logger.Zap.Debug("Setting up a connection pool.")

	// Настраиваем пул соединений: максимальное количество открытых соединений, максимальное количество неактивных и время жизни соединений.
	db.SetMaxOpenConns(cfg.DatabaseMaxOpenConns)
	db.SetMaxIdleConns(cfg.DatabaseMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DatabaseMaxLifetimeInMin * time.Minute)

	logger.Zap.Debug("Let's ping connection to the database.")

	// Пингуем базу данных, чтобы убедиться, что подключение успешно установлено.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Zap.Debug(fmt.Sprintf("-> `Connect` - successful via `%s`.", cfg.DatabaseDriver))

	return db, nil
}
