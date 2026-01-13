CREATE TABLE "comic_unit_tbl" (
    "id" TEXT PRIMARY KEY NOT NULL,
    
    "page_id" TEXT NOT NULL REFERENCES "comic_page_tbl"("id"),
    "index" INTEGER NOT NULL,
    
    "x_coordinate" DOUBLE PRECISION NOT NULL,
    "y_coordinate" DOUBLE PRECISION NOT NULL,
    
    "is_in_box" BOOLEAN DEFAULT FALSE NOT NULL,
    
    "translated_text" TEXT,
    "translator_id" TEXT REFERENCES "user_tbl"("id"),
    "translator_comment" TEXT,
    
    "proved_text" TEXT,
    "proved" BOOLEAN DEFAULT FALSE NOT NULL,
    "proofreader_id" TEXT REFERENCES "user_tbl"("id"),
    "proofreader_comment" TEXT,
    
    "creator_id" TEXT REFERENCES "user_tbl"("id"),
    
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    
    UNIQUE ("page_id", "index")
);

CREATE INDEX idx_comic_unit_page_id ON "comic_unit_tbl" ("page_id");
