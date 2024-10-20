package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	db "github.com/Ra1nz0r/zero_agency/db/sqlc"
	cfg "github.com/Ra1nz0r/zero_agency/internal/config"
	"github.com/Ra1nz0r/zero_agency/internal/logger"
	"github.com/Ra1nz0r/zero_agency/internal/models"
	srvs "github.com/Ra1nz0r/zero_agency/internal/services"
	"github.com/gofiber/fiber/v3"
)

type HandleQueries struct {
	*sql.DB
	*db.Queries
	cfg.Config
}

func NewHandlerQueries(queries *sql.DB, cfg cfg.Config) *HandleQueries {
	return &HandleQueries{
		queries,
		db.New(queries),
		cfg,
	}
}

// POST /edit/:id - изменение новости по Id
func (hq *HandleQueries) EditNews(c fiber.Ctx) error {
	logger.Zap.Debug("-> `EditNews` - calling handler.")

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		logger.Zap.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid news id",
		})
	}

	logger.Zap.Debug("Getting data from JSON.")

	// Записываем данные из JSON в структуру.
	var input models.InputEditNews
	if err := c.Bind().Body(&input); err != nil {
		logger.Zap.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid input",
		})
	}

	// Сравниваем ID из URL с переданным в JSON
	if input.ID != id {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID in URL does not match ID in JSON",
		})
	}

	logger.Zap.Debug("Checking ID in database.")

	// Проверяем существует ли переданное ID новости в базе данных.
	if _, err := hq.GetNews(c.Context(), id); err != nil {
		logger.Zap.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Entered ID does not exist in the database.",
		})
	}

	logger.Zap.Debug("+ Beginning transaction.")

	// Начинаем выполнение транзакции.
	tx, err := hq.Begin()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	qtx := hq.WithTx(tx)

	logger.Zap.Debug("- Updating data in the database.")

	// Обновление новости
	err = qtx.Update(c.Context(), db.UpdateParams{
		Id:      id,
		Column2: input.Title,
		Column3: input.Content,
	})
	if err != nil {
		logger.Zap.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Удаление старых категорий и вставка новых, если категории переданы в запросе.
	if len(input.Categories) > 0 {
		logger.Zap.Debug("- Removing categories from the database.")

		if err := qtx.DeleteCategories(c.Context(), id); err != nil {
			logger.Zap.Error(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		logger.Zap.Debug("- Adding new categories to the database.")

		err = qtx.InsertCategories(c.Context(), db.InsertCategoriesParams{
			NewsId:  id,
			Column2: input.Categories,
		})
		if err != nil {
			logger.Zap.Error(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Завершаем выполнение транзакции.
	if err = tx.Commit(); err != nil {
		logger.Zap.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Zap.Debug("+ Transaction committed.")

	logger.Zap.Debug("-> `EditNews` - successful called.")

	return c.JSON(&models.WriteResponse{
		Success: true,
	})
}

// GET /list - список новостей с пагинацией
func (hq *HandleQueries) ListNews(c fiber.Ctx) error {
	logger.Zap.Debug("-> `ListNews` - calling handler.")

	limit, err := srvs.StringToInt32WithOverflowCheck(c.Query("limit", hq.DefaultPaginationLimit))
	if err != nil {
		logger.Zap.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	offset, err := srvs.StringToInt32WithOverflowCheck(c.Query("offset", hq.DefaultOffset))
	if err != nil {
		logger.Zap.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Zap.Debug("Getting a list of news from the database.")

	// Получаем список новостей из базы данных.
	newsList, err := hq.List(c.Context(), db.ListParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Zap.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	logger.Zap.Debug("-> `ListNews` - successful called.")

	return c.JSON(&models.WriteResponse{
		Success: true,
		News:    newsList,
	})
}

// Ручка для аутентификации (логин)
func (hq *HandleQueries) Login(c fiber.Ctx) error {
	logger.Zap.Debug("-> `Login` - calling handler.")

	// Записываем данные из JSON в структуру.
	var lr models.LoginRequest
	if err := c.Bind().Body(&lr); err != nil {
		logger.Zap.Error(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Генерация JWT токена
	token, err := srvs.GenerateJWT(lr, hq.SecretKeyJWT, hq.JwtExpiresIn)
	if err != nil {
		logger.Zap.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	logger.Zap.Debug("-> `Login` - successful called.")

	// Возвращаем токен клиенту
	return c.JSON(fiber.Map{
		"token": token,
	})
}
