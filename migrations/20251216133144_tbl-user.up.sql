CREATE TABLE "user_tbl" (
    "id" TEXT PRIMARY KEY,

    "qq" TEXT UNIQUE NOT NULL,
    "nickname" TEXT UNIQUE NOT NULL,

    "password_hash" TEXT NOT NULL,

    "is_admin" BOOLEAN DEFAULT FALSE NOT NULL,

    "assigned_translator_at" TIMESTAMPTZ,
    "assigned_proofreader_at" TIMESTAMPTZ,
    "assigned_typesetter_at" TIMESTAMPTZ,
    "assigned_redrawer_at" TIMESTAMPTZ,
    "assigned_reviewer_at" TIMESTAMPTZ,
    "assigned_uploader_at" TIMESTAMPTZ,

    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE INDEX idx_trgm_user_nickname ON "user_tbl" USING GIN ("nickname" gin_trgm_ops);
CREATE UNIQUE INDEX uidx_user_qq_id ON "user_tbl" ("qq");
CREATE UNIQUE INDEX uidx_user_nickname ON "user_tbl" ("nickname");
CREATE INDEX idx_user_updated_at_desc ON "user_tbl" ("updated_at" DESC);

CREATE INDEX idx_user_assigned_translator_at ON "user_tbl" ("assigned_translator_at") WHERE "assigned_translator_at" IS NOT NULL;
CREATE INDEX idx_user_assigned_proofreader_at ON "user_tbl" ("assigned_proofreader_at") WHERE "assigned_proofreader_at" IS NOT NULL;
CREATE INDEX idx_user_assigned_typesetter_at ON "user_tbl" ("assigned_typesetter_at") WHERE "assigned_typesetter_at" IS NOT NULL;
CREATE INDEX idx_user_assigned_redrawer_at ON "user_tbl" ("assigned_redrawer_at") WHERE "assigned_redrawer_at" IS NOT NULL;
CREATE INDEX idx_user_assigned_reviewer_at ON "user_tbl" ("assigned_reviewer_at") WHERE "assigned_reviewer_at" IS NOT NULL;
CREATE INDEX idx_user_assigned_uploader_at ON "user_tbl" ("assigned_uploader_at") WHERE "assigned_uploader_at" IS NOT NULL;
