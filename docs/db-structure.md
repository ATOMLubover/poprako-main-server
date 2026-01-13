# 原 TS 下数据库 schema 设计

## Drizzle schema

```ts
import { sql } from "drizzle-orm";
import {
  boolean,
  doublePrecision,
  index,
  integer,
  pgTable,
  primaryKey,
  text,
  TIMESTAMPTZ,
  uniqueIndex,
} from "drizzle-orm/pg-core";

export const userSchema = pgTable(
  "user_tbl",
  {
    id: text("id").primaryKey().notNull(),

    qq_id: text("qq_id").unique().notNull(),
    nickname: text("nickname").unique().notNull(),

    passwordHash: text("password_hash").notNull(),

    isAdmin: boolean("is_admin").default(false).notNull(),
    assignedTranslatorAt: timestamp("assigned_translator_at"),
    assignedProoverAt: timestamp("assigned_proover_at"),
    assignedTypesetterAt: timestamp("assigned_typesetter_at"),
    assignedRedrawerAt: timestamp("assigned_redrawer_at"),
    assignedReviewerAt: timestamp("assigned_reviewer_at"),
    assignedUploaderAt: timestamp("assigned_uploader_at"),

    createdAt: timestamp("created_at").defaultNow().notNull(),
    updatedAt: timestamp("updated_at").defaultNow().notNull(),
  },
  (schema) => [
    index("idx_trgm_user_nickname").using(
      "gin",
      schema.nickname.op("gin_trgm_ops")
    ),
    uniqueIndex("uidx_user_qq_id").on(schema.qq_id),
    index("idx_user_updated_at_desc").on(schema.updatedAt.desc()),

    index("idx_user_assigned_translator_at")
      .on(schema.assignedTranslatorAt)
      .where(sql`${schema.assignedTranslatorAt} IS NOT NULL`),
    index("idx_user_assigned_proover_at")
      .on(schema.assignedProoverAt)
      .where(sql`${schema.assignedProoverAt} IS NOT NULL`),
    index("idx_user_assigned_typesetter_at")
      .on(schema.assignedTypesetterAt)
      .where(sql`${schema.assignedTypesetterAt} IS NOT NULL`),
    index("idx_user_assigned_redrawer_at")
      .on(schema.assignedRedrawerAt)
      .where(sql`${schema.assignedRedrawerAt} IS NOT NULL`),
    index("idx_user_assigned_reviewer_at")
      .on(schema.assignedReviewerAt)
      .where(sql`${schema.assignedReviewerAt} IS NOT NULL`),
    index("idx_user_assigned_uploader_at")
      .on(schema.assignedUploaderAt)
      .where(sql`${schema.assignedUploaderAt} IS NOT NULL`),
  ]
);

export type BasicUser = typeof userSchema.$inferSelect;
export type NewUser = typeof userSchema.$inferInsert;

export const tagSchema = pgTable("tag_tbl", {
  id: text("id").primaryKey().notNull(),

  name: text("name").unique().notNull(),

  picaCandidates: text("pica_candidates").array().notNull().default([]),
  ehentaiCandidates: text("ehentai_candidates").array().notNull().default([]),

  creatorId: text("creator_id")
    .references(() => userSchema.id)
    .notNull(),

  createdAt: timestamp("created_at").defaultNow().notNull(),
  updatedAt: timestamp("updated_at").defaultNow().notNull(),
});

export type BasicTag = typeof tagSchema.$inferSelect;
export type NewTag = typeof tagSchema.$inferInsert;

export const userTagSchema = pgTable(
  "user_tag_tbl",
  {
    userId: text("user_id")
      .references(() => userSchema.id)
      .notNull(),
    tagId: text("tag_id")
      .references(() => tagSchema.id)
      .notNull(),
  },
  (schema) => [primaryKey({ columns: [schema.userId, schema.tagId] })]
);

export type NewUserTag = typeof userTagSchema.$inferInsert;

export const worksetSchema = pgTable("workset_tbl", {
  id: text("id").primaryKey().notNull(),

  index: integer("index").unique().notNull(),

  comicCount: integer("comic_count").default(0).notNull(),
  description: text("description"),

  creatorId: text("creator_id"),

  createdAt: timestamp("created_at").defaultNow().notNull(),
  updatedAt: timestamp("updated_at").defaultNow().notNull(),
});

export type BasicWorkset = typeof worksetSchema.$inferSelect;
export type NewWorkset = typeof worksetSchema.$inferInsert;

export const comicSchema = pgTable(
  "comic_tbl",
  {
    id: text("id").primaryKey().notNull(),

    worksetId: text("workset_id")
      .references(() => worksetSchema.id)
      .notNull(),
    index: integer("index").notNull(),

    creatorId: text("creator_id")
      .references(() => userSchema.id)
      .notNull(),

    author: text("author").notNull(),
    title: text("title").notNull(),
    comment: text("comment"),
    description: text("description"),

    pageCount: integer("page_count").default(0).notNull(),

    likesCount: integer("likes_count").default(0).notNull(),

    translatingStartedAt: timestamp("translating_started_at"),
    translatingCompletedAt: timestamp("translating_completed_at"),

    proofreadingStartedAt: timestamp("proofreading_started_at"),
    proofreadingCompletedAt: timestamp("proofreading_completed_at"),

    typesettingStartedAt: timestamp("typesetting_started_at"),
    typesettingCompletedAt: timestamp("typesetting_completed_at"),

    reviewingCompletedAt: timestamp("reviewing_completed_at"),
    uploadingCompletedAt: timestamp("uploading_completed_at"),

    createdAt: timestamp("created_at").defaultNow().notNull(),
    updatedAt: timestamp("updated_at").defaultNow().notNull(),
  },
  (schema) => [
    uniqueIndex("uidx_comic_workset_id_index").on(
      schema.worksetId,
      schema.index
    ),
    index("idx_comic_likes_count_desc_updated_at_desc").on(
      schema.likesCount.desc(),
      schema.updatedAt.desc()
    ),
    index("idx_trgm_comic_author").using(
      "gin",
      schema.author.op("gin_trgm_ops")
    ),
    index("idx_trgm_comic_title").using("gin", schema.title.op("gin_trgm_ops")),

    index("idx_comic_translating_started_at")
      .on(schema.translatingStartedAt)
      .where(sql`${schema.translatingStartedAt} IS NOT NULL`),
    index("idx_comic_translating_completedat")
      .on(schema.translatingCompletedAt)
      .where(sql`${schema.translatingCompletedAt} IS NOT NULL`),
    index("idx_comic_proofreading_started_at")
      .on(schema.proofreadingStartedAt)
      .where(sql`${schema.proofreadingStartedAt} IS NOT NULL`),
    index("idx_comic_proofreading_completed_at")
      .on(schema.proofreadingCompletedAt)
      .where(sql`${schema.proofreadingCompletedAt} IS NOT NULL`),
    index("idx_comic_typesetting_started_at")
      .on(schema.typesettingStartedAt)
      .where(sql`${schema.typesettingStartedAt} IS NOT NULL`),
    index("idx_comic_typesetting_completed_at")
      .on(schema.typesettingCompletedAt)
      .where(sql`${schema.typesettingCompletedAt} IS NOT NULL`),
    index("idx_comic_reviewing_completed_at")
      .on(schema.reviewingCompletedAt)
      .where(sql`${schema.reviewingCompletedAt} IS NOT NULL`),
    index("idx_comic_uploading_completed_at")
      .on(schema.uploadingCompletedAt)
      .where(sql`${schema.uploadingCompletedAt} IS NOT NULL`),
  ]
);

export type BasicComic = typeof comicSchema.$inferSelect;
export type NewComic = typeof comicSchema.$inferInsert;

export const comicTagSchema = pgTable(
  "comic_tag_tbl",
  {
    comicId: text("comic_id")
      .references(() => comicSchema.id)
      .notNull(),
    tagId: text("tag_id")
      .references(() => tagSchema.id)
      .notNull(),
  },
  (schema) => [
    primaryKey({ columns: [schema.comicId, schema.tagId] }),
    index("idx_comic_tag_tag_id").on(schema.tagId),
    index("idx_comic_tag_comic_id").on(schema.comicId),
  ]
);

export type NewComicTag = typeof comicTagSchema.$inferInsert;

export const comicAssignmentSchema = pgTable(
  "comic_assignment_tbl",
  {
    id: text("id").primaryKey().notNull(),

    comicId: text("comic_id")
      .references(() => comicSchema.id)
      .notNull(),
    userId: text("user_id")
      .references(() => userSchema.id)
      .notNull(),

    // Uploader can upload any finished comics, so no need to track assignment time
    assignedTranslatorAt: timestamp("assigned_translator_at"),
    assignedProofreaderAt: timestamp("assigned_proofreader_at"),
    assignedTypesetterAt: timestamp("assigned_typesetter_at"),
    assignedRedrawerAt: timestamp("assigned_redrawer_at"),
    assignedReviewerAt: timestamp("assigned_reviewer_at"),

    createdAt: timestamp("created_at").defaultNow().notNull(),
    updatedAt: timestamp("updated_at").defaultNow().notNull(),
  },
  (schema) => [
    uniqueIndex("uidx_comic_assignment_comic_id_user_id").on(
      schema.comicId,
      schema.userId
    ),
    index("idx_comic_assignment_comic_id").on(schema.comicId),
    index("idx_comic_assignment_user_id").on(schema.userId),
  ]
);

export type BasicComicAssignment = typeof comicAssignmentSchema.$inferSelect;
export type NewComicAssignment = typeof comicAssignmentSchema.$inferInsert;

export const comicPageSchema = pgTable(
  "comic_page_tbl",
  {
    id: text("id").primaryKey().notNull(),

    comicId: text("comic_id")
      .references(() => comicSchema.id)
      .notNull(),
    index: integer("index").notNull(),

    ossKey: text("oss_key").notNull(),
    sizeBytes: integer("size_bytes").notNull(),

    uploaded: boolean("uploaded").default(true).notNull(),

    createdAt: timestamp("created_at").defaultNow().notNull(),
    updatedAt: timestamp("updated_at").defaultNow().notNull(),
  },
  (schema) => [
    uniqueIndex("uidx_comic_page_comic_id_index").on(
      schema.comicId,
      schema.index
    ),
  ]
);

export type BasicComicPage = typeof comicPageSchema.$inferSelect;
export type NewComicPage = typeof comicPageSchema.$inferInsert;

export const comicUnitSchema = pgTable(
  "comic_unit_tbl",
  {
    id: text("id").primaryKey().notNull(),

    pageId: text("page_id")
      .references(() => comicPageSchema.id)
      .notNull(),
    index: integer("index").notNull(),

    x_coordinate: doublePrecision("x_coordinate").notNull(),
    y_coordinate: doublePrecision("y_coordinate").notNull(),

    isInBox: boolean("is_in_box").default(false).notNull(),

    translatedText: text("translated_text"),
    translatorId: text("translator_id")
      .references(() => userSchema.id)
      .notNull(),
    translatorComment: text("translator_comment"),

    provedText: text("proved_text"),
    proved: boolean("proved").default(false).notNull(),
    proofreaderId: text("proofreader_id")
      .references(() => userSchema.id)
      .notNull(),
    proofreaderComment: text("proofreader_comment"),

    creatorId: text("creator_id")
      .references(() => userSchema.id)
      .notNull(),

    createdAt: timestamp("created_at").defaultNow().notNull(),
    updatedAt: timestamp("updated_at").defaultNow().notNull(),
  },
  (schema) => [
    uniqueIndex("uidx_comic_unit_page_id_index").on(
      schema.pageId,
      schema.index
    ),
  ]
);

export type BasicComicUnit = typeof comicUnitSchema.$inferSelect;
export type NewComicUnit = typeof comicUnitSchema.$inferInsert;

export const termbaseSchama = pgTable(
  "termbase_tbl",
  {
    id: text("id").primaryKey().notNull(),

    name: text("name").unique().notNull(),
    description: text("description"),

    creatorId: text("creator_id")
      .references(() => userSchema.id)
      .notNull(),

    createdAt: timestamp("created_at").defaultNow().notNull(),
    updatedAt: timestamp("updated_at").defaultNow().notNull(),
  },
  (schema) => [
    index("idx_trgm_termbase_name").using(
      "gin",
      schema.name.op("gin_trgm_ops")
    ),
    index("idx_termbase_updated_at_desc").on(schema.updatedAt.desc()),
  ]
);

export type BasicTermbase = typeof termbaseSchama.$inferSelect;
export type NewTermbase = typeof termbaseSchama.$inferInsert;

export const termSchema = pgTable(
  "term_tbl",
  {
    id: text("id").primaryKey().notNull(),

    termbaseId: text("termbase_id")
      .references(() => termbaseSchama.id)
      .notNull(),

    sourceText: text("source_text").notNull(),
    targetText: text("target_text").notNull(),

    creatorId: text("creator_id")
      .references(() => userSchema.id)
      .notNull(),

    createdAt: timestamp("created_at").defaultNow().notNull(),
    updatedAt: timestamp("updated_at").defaultNow().notNull(),
  },
  (schema) => [
    index("idx_trgm_term_source_text").using(
      "gin",
      schema.sourceText.op("gin_trgm_ops")
    ),
    index("idx_term_updated_at_desc").on(schema.updatedAt.desc()),
  ]
);

export type BasicTerm = typeof termSchema.$inferSelect;
export type NewTerm = typeof termSchema.$inferInsert;
```
