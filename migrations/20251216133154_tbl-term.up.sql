CREATE TABLE "term_tbl" (
    "id" TEXT PRIMARY KEY NOT NULL,
    
    "termbase_id" TEXT NOT NULL REFERENCES "termbase_tbl"("id"),
    
    "source_text" TEXT NOT NULL,
    "target_text" TEXT NOT NULL,
    
    "creator_id" TEXT NOT NULL REFERENCES "user_tbl"("id"),
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_trgm_term_source_text ON "term_tbl" USING GIN ("source_text" gin_trgm_ops);

CREATE INDEX idx_term_updated_at_desc ON "term_tbl" ("updated_at" DESC);
