package server

import (
	"database/sql"
	"time"

	"fmt"

	cfg "github.com/Ra1nz0r/zero_agency/internal/config"
	"github.com/Ra1nz0r/zero_agency/internal/logger"
)

func Connect(cfg *cfg.Config) (*sql.DB, error) {
	logger.Zap.Debug("-> `Connect` - launching function.")

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseName,
	)
	logger.Zap.Debug(fmt.Sprintf("The link to connect to the DB generated: %s", dbURL))

	logger.Zap.Debug("Opening connection to the database.")
	db, err := sql.Open(cfg.DatabaseDriver, dbURL)
	if err != nil {
		return nil, err
	}

	logger.Zap.Debug("Setting up a connection pool.")
	db.SetMaxOpenConns(cfg.DatabaseMaxOpenConns)
	db.SetMaxIdleConns(cfg.DatabaseMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DatabaseMaxLifetimeInMin * time.Minute)

	logger.Zap.Debug("Let's ping connection to the database.")
	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Zap.Debug(fmt.Sprintf("-> `Connect` - successful via `%s`.", cfg.DatabaseDriver))

	return db, nil
}
