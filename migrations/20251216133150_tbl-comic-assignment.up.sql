CREATE TABLE "comic_assignment_tbl" (
    "id" TEXT PRIMARY KEY NOT NULL,

    "comic_id" TEXT NOT NULL REFERENCES "comic_tbl"("id") ON DELETE CASCADE,
    "user_id" TEXT NOT NULL REFERENCES "user_tbl"("id") ON DELETE CASCADE,

    "assigned_translator_at" TIMESTAMPTZ,
    "assigned_proofreader_at" TIMESTAMPTZ,
    "assigned_typesetter_at" TIMESTAMPTZ,
    "assigned_redrawer_at" TIMESTAMPTZ,
    "assigned_reviewer_at" TIMESTAMPTZ,

    "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,

    UNIQUE ("comic_id", "user_id")
);

CREATE INDEX idx_comic_assignment_comic_id ON "comic_assignment_tbl" ("comic_id");

CREATE INDEX idx_comic_assignment_user_id ON "comic_assignment_tbl" ("user_id");
