-- name: ListNews :many
SELECT n."Id",
    n."Title",
    n."Content",
    COALESCE(array_agg(nc."CategoryId"), '{}') AS "Categories"
FROM "News" n
    LEFT JOIN "NewsCategories" nc ON n."Id" = nc."NewsId"
GROUP BY n."Id",
    n."Title",
    n."Content"
LIMIT $1 OFFSET $2;
-- name: UpdateNews :exec
UPDATE "News"
SET "Title" = COALESCE(NULLIF($2, ''), "Title"),
    "Content" = COALESCE(NULLIF($3, ''), "Content")
WHERE "Id" = $1;
-- name: DeleteNewsCategories :exec
DELETE FROM "NewsCategories"
WHERE "NewsId" = $1;
-- name: InsertNewsCategories :exec
INSERT INTO "NewsCategories" ("NewsId", "CategoryId")
VALUES ($1, unnest($2::BIGINT []));