CREATE TABLE "tag_tbl" (
    "id" TEXT PRIMARY KEY NOT NULL,

    "name" TEXT NOT NULL UNIQUE,
    
    "pica_candidates" TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
    "ehentai_candidates" TEXT[] NOT NULL DEFAULT '{}'::TEXT[],
    
    "creator_id" TEXT NOT NULL REFERENCES "user_tbl"("id"),
    
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_trgm_tag_name ON "tag_tbl" USING GIN ("name" gin_trgm_ops);

CREATE INDEX idx_tag_updated_at_desc ON "tag_tbl" ("updated_at" DESC);
