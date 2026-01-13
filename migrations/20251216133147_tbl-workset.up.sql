CREATE TABLE "workset_tbl" (
    "id" TEXT PRIMARY KEY NOT NULL,
    
    "index" INTEGER UNIQUE NOT NULL,
    "comic_count" INTEGER DEFAULT 0 NOT NULL,
    
    "description" TEXT,
    
    "creator_id" TEXT,
    
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_workset_updated_at_desc ON "workset_tbl" ("updated_at" DESC);
