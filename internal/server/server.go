package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	"github.com/Ra1nz0r/zero_agency/db"
	"github.com/Ra1nz0r/zero_agency/internal/config"
	hd "github.com/Ra1nz0r/zero_agency/internal/handlers"
	"github.com/Ra1nz0r/zero_agency/internal/logger"
	srv "github.com/Ra1nz0r/zero_agency/internal/services"
	"github.com/gofiber/fiber/v3"
)

// Run запускает сервер.
func Run() {
	logger.Zap.Debug()
	// Загружаем переменные окружения из '.env' файла.
	cfg, errLoad := config.LoadConfig(".")
	if errLoad != nil {
		log.Fatal(fmt.Errorf("unable to load config: %w", errLoad))
	}

	// Инициализируем логгер, с указанием уровня логирования.
	if errLog := logger.Initialize(cfg.LogLevel); errLog != nil {
		log.Fatal(fmt.Errorf("failed to initialize the logger: %w", errLog))
	}

	logger.Zap.Debug("Connecting to the database.")
	connect, dbURL, errConn := db.Connect(&cfg)
	if errConn != nil {
		logger.Zap.Fatal(fmt.Errorf("unable to create connection to database: %w", errConn))
	}

	logger.Zap.Debug("Checking the existence of a TABLE in the database.")
	// Проверяем существование TABLE в базе данных.
	exists, errExs := srv.TableExists(connect, cfg.DatabaseName)
	if errExs != nil {
		logger.Zap.Fatal(fmt.Errorf("failed to check if TABLE exists: %w", errExs))
	}

	// Создаём TABLE, если он не существует.
	if !exists {
		logger.Zap.Debug(fmt.Sprintf("Creating TABLE in '%s' database.", cfg.DatabaseName))
		if errRunMigr := srv.RunMigrations(dbURL, cfg.MigrationPath); errRunMigr != nil {
			logger.Zap.Fatal(fmt.Errorf("failed to run migrations: %w", errConn))
		}
	}

	logger.Zap.Debug("Configuring and starting the server.")
	// Конфигурируем и запускаем сервер.
	srv := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		AppName:       "News App v1.0.0",
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  10 * time.Second,
		IdleTimeout:   120 * time.Second,
	})

	// Передаём подключение и настройки приложения нашим обработчикам.
	queries := hd.NewHandlerQueries(connect, cfg)

	//srv.Use(swagger.New())
	logger.Zap.Debug("Running handlers.")

	// Ручки
	srv.Get("/list", queries.ListNews)
	srv.Post("/edit/:id", queries.EditNews)

	logger.Zap.Info(fmt.Sprintf("Server is running on: '%s'", cfg.ServerHost))

	go func() {
		if errListn := srv.Listen(cfg.ServerHost); errListn != nil {
			logger.Zap.Fatal(fmt.Errorf("HTTP server error: %w", errListn))
		}
		logger.Zap.Info("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if errShut := srv.ShutdownWithContext(shutdownCtx); errShut != nil {
		logger.Zap.Fatal(fmt.Errorf("HTTP shutdown error: %w", errShut))
	}
	logger.Zap.Info("Graceful shutdown complete.")
}
