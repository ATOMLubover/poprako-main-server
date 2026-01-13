CREATE TABLE "comic_tag_tbl" (
    "comic_id" TEXT NOT NULL REFERENCES "comic_tbl"("id"),
    "tag_id" TEXT NOT NULL REFERENCES "tag_tbl"("id"),

    PRIMARY KEY ("comic_id", "tag_id")
);

CREATE INDEX idx_comic_tag_tag_id ON "comic_tag_tbl" ("tag_id");

CREATE INDEX idx_comic_tag_comic_id ON "comic_tag_tbl" ("comic_id");
