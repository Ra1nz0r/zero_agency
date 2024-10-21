package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fmt"

	cfg "github.com/Ra1nz0r/zero_agency/internal/config"
	hd "github.com/Ra1nz0r/zero_agency/internal/handlers"
	"github.com/Ra1nz0r/zero_agency/internal/logger"
	"github.com/Ra1nz0r/zero_agency/internal/middleware"
	"github.com/Ra1nz0r/zero_agency/internal/models"
	"github.com/Ra1nz0r/zero_agency/internal/services"
	"github.com/go-playground/validator/v10/non-standard/validators"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

// Run запускает сервер.
func Run() {
	logger.Zap.Info("Loading config.")

	// Загружаем переменные окружения из '.env' файла.
	cfg, errLoad := cfg.LoadConfig(".")
	if errLoad != nil {
		log.Fatal(fmt.Errorf("unable to load config: %w", errLoad))
	}

	logger.Zap.Info("Initialize logger.")

	// Инициализируем логгер, с указанием уровня логирования.
	if errLog := logger.Initialize(cfg.LogLevel); errLog != nil {
		log.Fatal(fmt.Errorf("failed to initialize the logger: %w", errLog))
	}

	logger.Zap.Info("Connecting to the database.")

	// Создаём подключение к базе данных с пуллом.
	connect, errConn := Connect(&cfg)
	if errConn != nil {
		logger.Zap.Fatal(fmt.Errorf("unable to create connection to database: %w", errConn))
	}

	logger.Zap.Info("Running migrations.")

	// Запускаем и проверяем миграции, если есть не применённые, то они будут выполнены.
	if errRunMigr := services.RunMigrations(connect, cfg); errRunMigr != nil {
		logger.Zap.Fatal(fmt.Errorf("failed to run migrations: %w", errRunMigr))
	}

	logger.Zap.Info("Configuring and starting the server.")

	// Конфигурируем и запускаем сервер.
	srv := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		AppName:       "News App v1.0.0",
		ReadTimeout:   5 * time.Second,
		WriteTimeout:  10 * time.Second,
		IdleTimeout:   120 * time.Second,
	})

	logger.Zap.Info("Sending connection to handlers.")

	// Передаём подключение и настройки приложения нашим обработчикам.
	queries := hd.NewHandlerQueries(connect, cfg)

	logger.Zap.Info("Creating and initialize the validator.")

	// Создаём и конфигурируем валидатор.
	valid := models.InitValidator()

	// Создаём notblank для проверки, что строка только из пробелов
	if err := valid.RegisterValidation("notblank", validators.NotBlank); err != nil {
		logger.Zap.Error(err.Error())

	}

	logger.Zap.Info("Running handlers.")

	// Создаём и запускаем эндпойнты.
	srv.Use(swagger.New(swagger.Config{
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}))
	srv.Use("/edit/:id", middleware.JWTMiddleware(cfg.SecretKeyJWT))

	srv.Get("/list", queries.ListNews)

	srv.Post("/edit/:id", queries.EditNews)
	srv.Post("/login", queries.Login)

	// Запускаем сервер.
	host := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)

	logger.Zap.Info(fmt.Sprintf("Server is running on: '%s'", host))

	go func() {
		if errListn := srv.Listen(host); errListn != nil {
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
