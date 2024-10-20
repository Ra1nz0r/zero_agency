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
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"

	"github.com/gofiber/fiber/v3"
)

// Run запускает сервер.
func Run() {
	logger.Zap.Debug()
	// Загружаем переменные окружения из '.env' файла.
	cfg, errLoad := cfg.LoadConfig(".")
	if errLoad != nil {
		log.Fatal(fmt.Errorf("unable to load config: %w", errLoad))
	}

	// Инициализируем логгер, с указанием уровня логирования.
	if errLog := logger.Initialize(cfg.LogLevel); errLog != nil {
		log.Fatal(fmt.Errorf("failed to initialize the logger: %w", errLog))
	}

	logger.Zap.Debug("Connecting to the database.")
	connect, errConn := Connect(&cfg)
	if errConn != nil {
		logger.Zap.Fatal(fmt.Errorf("unable to create connection to database: %w", errConn))
	}

	if errRunMigr := services.RunMigrations(connect, cfg); errRunMigr != nil {
		logger.Zap.Fatal(fmt.Errorf("failed to run migrations: %w", errRunMigr))
	}

	// Создаём и конфигурируем валидатор.
	v := validator.New()
	v.RegisterValidation("notblank", validators.NotBlank) // для проверки, что строка только из пробелов

	logger.Zap.Debug("Configuring and starting the server.")

	// Конфигурируем и запускаем сервер.
	srv := fiber.New(fiber.Config{
		CaseSensitive:   true,
		StrictRouting:   true,
		AppName:         "News App v1.0.0",
		ReadTimeout:     5 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     120 * time.Second,
		StructValidator: models.NewValidator(v),
	})

	// Передаём подключение и настройки приложения нашим обработчикам.
	queries := hd.NewHandlerQueries(connect, cfg)

	//srv.Use(swagger.New())
	logger.Zap.Debug("Running handlers.")

	srv.Post("/login", queries.Login)

	srv.Use("/list", middleware.JWTMiddleware(cfg.SecretKeyJWT))
	srv.Use("/edit/:id", middleware.JWTMiddleware(cfg.SecretKeyJWT))

	// Ручки
	srv.Get("/list", queries.ListNews)
	srv.Post("/edit/:id", queries.EditNews)

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
