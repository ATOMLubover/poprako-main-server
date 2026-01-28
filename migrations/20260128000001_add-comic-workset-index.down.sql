-- Drop the index
DROP INDEX IF EXISTS idx_comic_workset_index;

-- Remove workset_index column
ALTER TABLE "comic_tbl" DROP COLUMN "workset_index";
