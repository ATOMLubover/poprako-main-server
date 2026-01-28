-- Add workset_index column to comic_tbl
ALTER TABLE "comic_tbl" ADD COLUMN "workset_index" INTEGER;

-- Backfill workset_index from workset_tbl.index
UPDATE "comic_tbl"
SET "workset_index" = "workset_tbl"."index"
FROM "workset_tbl"
WHERE "comic_tbl"."workset_id" = "workset_tbl"."id";

-- Set NOT NULL constraint after backfilling
ALTER TABLE "comic_tbl" ALTER COLUMN "workset_index" SET NOT NULL;

-- Add index for potential queries
CREATE INDEX idx_comic_workset_index ON "comic_tbl" ("workset_index");
