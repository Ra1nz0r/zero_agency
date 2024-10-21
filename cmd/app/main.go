package main

import (
	_ "github.com/Ra1nz0r/zero_agency/docs"
	"github.com/Ra1nz0r/zero_agency/internal/server"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"

	_ "github.com/lib/pq"
)

// @title News API
// @version 1.0
// @description REST API для управления новостями. Включает функции изменения новостей и получения списка с поддержкой категорий.
// @termsOfService http://swagger.io/terms/

// @contact.name Artem Rylskii
// @contact.url https://t.me/Rainz0r
// @contact.email n52rus@gmail.com

// @host localhost:7654
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description !!! ВАЖНО !!! Введите токен в формате: Bearer <токен> !!! ВАЖНО !!!

// @security BearerAuth

func main() {
	server.Run()
}
