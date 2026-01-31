CREATE TABLE "comic_page_tbl" (
    "id" TEXT PRIMARY KEY NOT NULL,
    
    "comic_id" TEXT NOT NULL REFERENCES "comic_tbl"("id"),
    "index" INTEGER NOT NULL,
    
    "uploaded" BOOLEAN DEFAULT false NOT NULL,
    
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    
    UNIQUE ("comic_id", "index")
);

CREATE INDEX idx_comic_page_comic_id ON "comic_page_tbl" ("comic_id");
