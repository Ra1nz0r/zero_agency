CREATE TABLE IF NOT EXISTS "News" (
    "Id" BIGSERIAL PRIMARY KEY,
    "Title" VARCHAR NOT NULL,
    "Content" TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS "NewsCategories" (
    "NewsId" BIGINT NOT NULL,
    "CategoryId" BIGINT NOT NULL,
    PRIMARY KEY ("NewsId", "CategoryId")
);
CREATE INDEX ON "News" ("Title");
CREATE INDEX ON "NewsCategories" ("CategoryId");