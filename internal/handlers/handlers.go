package handlers

import (
	"context"
	"database/sql"
	"strconv"

	db "github.com/Ra1nz0r/zero_agency/db/sqlc"
	"github.com/gofiber/fiber/v3"
)

type HandleQueries struct {
	*db.Queries
}

func NewHandlerQueries(queries *sql.DB) *HandleQueries {
	return &HandleQueries{
		db.New(queries),
	}
}

// GET /list - список новостей с пагинацией
func (hq *HandleQueries) ListNews(c fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid limit"})
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid offset"})
	}

	newsList, err := hq.Queries.ListNews(context.Background(), db.ListNewsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"Success": true,
		"News":    newsList,
	})
}

// POST /edit/:id - изменение новости по Id
func (hq *HandleQueries) EditNews(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid news id"})
	}

	var input struct {
		Title      string  `json:"Title"`
		Content    string  `json:"Content"`
		Categories []int64 `json:"Categories"`
	}

	if err := c.Bind().Body(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	// Обновление новости
	err = hq.UpdateNews(context.Background(), db.UpdateNewsParams{
		Id:      id,
		Column2: input.Title,
		Column3: input.Content,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Удаление старых категорий и вставка новых
	if err := hq.DeleteNewsCategories(context.Background(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if len(input.Categories) > 0 {
		err = hq.InsertNewsCategories(context.Background(), db.InsertNewsCategoriesParams{
			NewsId:  id,
			Column2: input.Categories,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return c.JSON(fiber.Map{"Success": true})
}
