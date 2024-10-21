// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"

	"github.com/lib/pq"
)

const deleteCategories = `-- name: DeleteCategories :exec
DELETE FROM "NewsCategories"
WHERE "NewsId" = $1
`

func (q *Queries) DeleteCategories(ctx context.Context, newsid int64) error {
	_, err := q.db.ExecContext(ctx, deleteCategories, newsid)
	return err
}

const getNews = `-- name: GetNews :one
SELECT "Id", "Title", "Content"
FROM "News"
WHERE "Id" = $1
LIMIT 1
`

func (q *Queries) GetNews(ctx context.Context, id int64) (News, error) {
	row := q.db.QueryRowContext(ctx, getNews, id)
	var i News
	err := row.Scan(&i.Id, &i.Title, &i.Content)
	return i, err
}

const insertCategories = `-- name: InsertCategories :exec
INSERT INTO "NewsCategories" ("NewsId", "CategoryId")
VALUES ($1, unnest($2::BIGINT []))
`

type InsertCategoriesParams struct {
	NewsId  int64   `json:"NewsId"`
	Column2 []int64 `json:"column_2"`
}

func (q *Queries) InsertCategories(ctx context.Context, arg InsertCategoriesParams) error {
	_, err := q.db.ExecContext(ctx, insertCategories, arg.NewsId, pq.Array(arg.Column2))
	return err
}

const list = `-- name: List :many
SELECT n."Id",
    n."Title",
    n."Content",
    COALESCE(array_agg(nc."CategoryId"), '{}') AS "Categories"
FROM "News" n
    JOIN "NewsCategories" nc ON n."Id" = nc."NewsId"
GROUP BY n."Id",
    n."Title",
    n."Content"
ORDER BY n."Id" DESC
LIMIT $1 OFFSET $2
`

type ListParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListRow struct {
	Id         int64       `json:"Id"`
	Title      string      `json:"Title"`
	Content    string      `json:"Content"`
	Categories interface{} `json:"Categories"`
}

func (q *Queries) List(ctx context.Context, arg ListParams) ([]ListRow, error) {
	rows, err := q.db.QueryContext(ctx, list, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRow
	for rows.Next() {
		var i ListRow
		if err := rows.Scan(
			&i.Id,
			&i.Title,
			&i.Content,
			&i.Categories,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const update = `-- name: Update :exec
UPDATE "News"
SET "Title" = COALESCE(NULLIF($2, ''), "Title"),
    "Content" = COALESCE(NULLIF($3, ''), "Content")
WHERE "Id" = $1
`

type UpdateParams struct {
	Id      int64       `json:"Id"`
	Column2 interface{} `json:"column_2"`
	Column3 interface{} `json:"column_3"`
}

func (q *Queries) Update(ctx context.Context, arg UpdateParams) error {
	_, err := q.db.ExecContext(ctx, update, arg.Id, arg.Column2, arg.Column3)
	return err
}
