CREATE TABLE IF NOT EXISTS "News" (
    "Id" BIGSERIAL PRIMARY KEY,
    "Title" TEXT NOT NULL,
    "Content" TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS "NewsCategories" (
    "NewsId" BIGINT NOT NULL,
    "CategoryId" BIGINT NOT NULL,
    PRIMARY KEY ("NewsId", "CategoryId"),
    FOREIGN KEY ("NewsId") REFERENCES "News"("Id") ON DELETE CASCADE
);