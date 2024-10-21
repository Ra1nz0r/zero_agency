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

// EditNews обновляет существующую новость по её ID. Обрабатывает PUT запрос в формате
// JSON {"Id": 64, "Title": "Lorem ipsum", "Content": "Dolor sit amet <b>foo</b>", "Categories": [1,2,3]}.
// Данные новости обновляются в базе данных. Если указаны категории, они также обновляются.
// Если переданный ID новости не существует, возвращается ошибка. В случае успешного завершения транзакции
// возвращается ответ с подтверждением обновления.
//
// @Summary Обновляет существующую новость.
// @Description Обновляет новость и её категории в базе данных по переданному ID. Если новость не найдена или данные некорректны, возвращается ошибка.
// @Tags news
// @Accept  json
// @Produce json
// @Param id path int true "ID новости для обновления."
// @Param models.InputEditNews body models.InputEditNews true "Данные для обновления новости."
// @Success 200 {object} models.WriteResponse "Успешное обновление новости."
// @Failure 400 {object} map[string]string "Некорректный запрос. Например, если ID в URL не совпадает с ID в JSON или если ID не существует."
// @Failure 500 {object} map[string]string "Ошибка сервера при обновлении данных."
// @Router /news/edit/{id} [post]
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

// ListNews возвращает список новостей с поддержкой пагинации.
// Обрабатывает GET запросы с параметрами "limit" и "offset" для ограничения количества возвращаемых записей и их смещения.
// Если параметры отсутствуют, используются значения по умолчанию.
// В случае ошибки в обработке параметров или при работе с базой данных, возвращается соответствующий статус ошибки и сообщение.
// При успешном завершении запроса возвращается список новостей.
//
// @Summary Получение списка новостей с пагинацией.
// @Description Возвращает список новостей с возможностью ограничения числа результатов (limit) и смещения (offset). Если параметры не указаны, используются значения по умолчанию. При некорректных параметрах или ошибке базы данных возвращается соответствующее сообщение об ошибке.
// @Tags news
// @Accept  json
// @Produce json
// @Param limit query int false "Максимальное количество новостей для получения. По умолчанию — 10."
// @Param offset query int false "Смещение для пагинации. По умолчанию — 0."
// @Success 200 {object} models.WriteResponse "Список новостей."
// @Failure 400 {object} map[string]string "Некорректный запрос: ошибка обработки параметров."
// @Failure 500 {object} map[string]string "Ошибка сервера при получении новостей."
// @Router /news/list [get]
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

// Login аутентифицирует пользователя и возвращает JWT токен.
// Обрабатывает POST-запрос с JSON телом, содержащим имя пользователя и пароль.
// В случае успешной аутентификации генерируется JWT токен, который возвращается в ответе.
// Если параметры запроса некорректны, возвращается ошибка 400.
// Если при генерации токена произошла ошибка, возвращается ошибка 500.
//
// @Summary Аутентификация пользователя
// @Description Аутентифицирует пользователя на основе имени и пароля, и возвращает JWT токен при успешной аутентификации.
// @Tags auth
// @Accept  json
// @Produce json
// @Param request body models.LoginRequest true "Данные для входа (имя пользователя и пароль)"
// @Success 200 {object} map[string]string "JWT токен"
// @Failure 400 {object} map[string]string "Ошибка запроса: не удалось распарсить JSON или некорректные данные"
// @Failure 500 {object} map[string]string "Ошибка сервера при генерации токена"
// @Router /auth/login [post]
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
	token, err := srvs.GenerateJWT(lr.Username, hq.SecretKeyJWT, hq.JwtExpiresIn)
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
