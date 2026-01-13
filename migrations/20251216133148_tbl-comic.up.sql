CREATE TABLE "comic_tbl" (
    "id" TEXT PRIMARY KEY NOT NULL,
    
    "workset_id" TEXT NOT NULL REFERENCES "workset_tbl"("id"),
    "index" INTEGER NOT NULL,
    
    "creator_id" TEXT NOT NULL REFERENCES "user_tbl"("id"),
    
    "author" TEXT NOT NULL,
    "title" TEXT NOT NULL,
    "comment" TEXT,
    "description" TEXT,
    
    "page_count" INTEGER DEFAULT 0 NOT NULL,
    "likes_count" INTEGER DEFAULT 0 NOT NULL,
    
    "translating_started_at" TIMESTAMPTZ,
    "translating_completed_at" TIMESTAMPTZ,
    "proofreading_started_at" TIMESTAMPTZ,
    "proofreading_completed_at" TIMESTAMPTZ,
    "typesetting_started_at" TIMESTAMPTZ,
    "typesetting_completed_at" TIMESTAMPTZ,
    "reviewing_completed_at" TIMESTAMPTZ,
    "uploading_completed_at" TIMESTAMPTZ,
    
    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    
    UNIQUE ("workset_id", "index")
);

CREATE INDEX idx_comic_likes_count_desc_updated_at_desc ON "comic_tbl" ("likes_count" DESC, "updated_at" DESC);
CREATE INDEX idx_trgm_comic_author ON "comic_tbl" USING GIN ("author" gin_trgm_ops);
CREATE INDEX idx_trgm_comic_title ON "comic_tbl" USING GIN ("title" gin_trgm_ops);

CREATE INDEX idx_comic_translating_started_at ON "comic_tbl" ("translating_started_at") WHERE "translating_started_at" IS NOT NULL;
CREATE INDEX idx_comic_translating_completedat ON "comic_tbl" ("translating_completed_at") WHERE "translating_completed_at" IS NOT NULL;
CREATE INDEX idx_comic_proofreading_started_at ON "comic_tbl" ("proofreading_started_at") WHERE "proofreading_started_at" IS NOT NULL;
CREATE INDEX idx_comic_proofreading_completed_at ON "comic_tbl" ("proofreading_completed_at") WHERE "proofreading_completed_at" IS NOT NULL;
CREATE INDEX idx_comic_typesetting_started_at ON "comic_tbl" ("typesetting_started_at") WHERE "typesetting_started_at" IS NOT NULL;
CREATE INDEX idx_comic_typesetting_completed_at ON "comic_tbl" ("typesetting_completed_at") WHERE "typesetting_completed_at" IS NOT NULL;
CREATE INDEX idx_comic_reviewing_completed_at ON "comic_tbl" ("reviewing_completed_at") WHERE "reviewing_completed_at" IS NOT NULL;
CREATE INDEX idx_comic_uploading_completed_at ON "comic_tbl" ("uploading_completed_at") WHERE "uploading_completed_at" IS NOT NULL;
