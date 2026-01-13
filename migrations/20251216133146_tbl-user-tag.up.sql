CREATE TABLE "user_tag_tbl" (
    "user_id" TEXT NOT NULL REFERENCES "user_tbl"("id"),
    "tag_id" TEXT NOT NULL REFERENCES "tag_tbl"("id"),
    
    PRIMARY KEY ("user_id", "tag_id")
);
