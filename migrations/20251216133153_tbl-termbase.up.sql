CREATE TABLE "termbase_tbl" (
    "id" TEXT PRIMARY KEY NOT NULL,
    
    "name" TEXT NOT NULL UNIQUE,
    "description" TEXT,
    
    "creator_id" TEXT NOT NULL REFERENCES "user_tbl"("id"),
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_trgm_termbase_name ON "termbase_tbl" USING GIN ("name" gin_trgm_ops);

CREATE INDEX idx_termbase_updated_at_desc ON "termbase_tbl" ("updated_at" DESC);
