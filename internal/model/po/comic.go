package po

import (
	"time"
)

const (
	COMIC_TABLE = "comic_tbl"
)

// Used when creating a new comic.
type NewComic struct {
	ID string `gorm:"column:id;primaryKey"`

	WorksetID    string  `gorm:"column:workset_id"`
	WorksetIndex int     `gorm:"column:workset_index"`
	Index        int64   `gorm:"column:index"`
	CreatorID    string  `gorm:"column:creator_id"`
	Author       string  `gorm:"column:author"`
	Title        string  `gorm:"column:title"`
	Comment      *string `gorm:"column:comment"`
	Description  *string `gorm:"column:description"`
}

// Used when retrieving brief comic info.
type BriefComic struct {
	ID string `gorm:"column:id;primaryKey"`

	WorksetID    string `gorm:"column:workset_id"`
	WorksetIndex int    `gorm:"column:workset_index"`
	Index        int64  `gorm:"column:index"`

	Author string `gorm:"column:author"`
	Title  string `gorm:"column:title"`

	PageCount int64 `gorm:"column:page_count"`

	TranslatingStartedAt    *time.Time `gorm:"column:translating_started_at"`
	TranslatingCompletedAt  *time.Time `gorm:"column:translating_completed_at"`
	ProofreadingStartedAt   *time.Time `gorm:"column:proofreading_started_at"`
	ProofreadingCompletedAt *time.Time `gorm:"column:proofreading_completed_at"`
	TypesettingStartedAt    *time.Time `gorm:"column:typesetting_started_at"`
	TypesettingCompletedAt  *time.Time `gorm:"column:typesetting_completed_at"`
	ReviewingCompletedAt    *time.Time `gorm:"column:reviewing_completed_at"`
	UploadingCompletedAt    *time.Time `gorm:"column:uploading_completed_at"`
}

// Used when retrieving basic comic info.
type BasicComic struct {
	ID string `gorm:"column:id;primaryKey"`

	WorksetID    string `gorm:"column:workset_id"`
	WorksetIndex int    `gorm:"column:workset_index"`
	Index        int64  `gorm:"column:index"`

	CreatorID       string `gorm:"column:creator_id"`
	CreatorNickname string `gorm:"column:creator_nickname"`

	Author      string  `gorm:"column:author"`
	Title       string  `gorm:"column:title"`
	Comment     *string `gorm:"column:comment"`
	Description *string `gorm:"column:description"`

	PageCount int64 `gorm:"column:page_count"`

	TranslatingStartedAt    *time.Time `gorm:"column:translating_started_at"`
	TranslatingCompletedAt  *time.Time `gorm:"column:translating_completed_at"`
	ProofreadingStartedAt   *time.Time `gorm:"column:proofreading_started_at"`
	ProofreadingCompletedAt *time.Time `gorm:"column:proofreading_completed_at"`
	TypesettingStartedAt    *time.Time `gorm:"column:typesetting_started_at"`
	TypesettingCompletedAt  *time.Time `gorm:"column:typesetting_completed_at"`
	ReviewingCompletedAt    *time.Time `gorm:"column:reviewing_completed_at"`
	UploadingCompletedAt    *time.Time `gorm:"column:uploading_completed_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// Used when updating comic info.
// Any fields with default zero values (nil) will not be updated.
type PatchComic struct {
	ID string `gorm:"column:id;primaryKey"`

	Author      *string `gorm:"column:author"`
	Title       *string `gorm:"column:title"`
	Comment     *string `gorm:"column:comment"`
	Description *string `gorm:"column:description"`

	TranslatingStartedAt    *time.Time `gorm:"column:translating_started_at"`
	TranslatingCompletedAt  *time.Time `gorm:"column:translating_completed_at"`
	ProofreadingStartedAt   *time.Time `gorm:"column:proofreading_started_at"`
	ProofreadingCompletedAt *time.Time `gorm:"column:proofreading_completed_at"`
	TypesettingStartedAt    *time.Time `gorm:"column:typesetting_started_at"`
	TypesettingCompletedAt  *time.Time `gorm:"column:typesetting_completed_at"`
	ReviewingCompletedAt    *time.Time `gorm:"column:reviewing_completed_at"`
	UploadingCompletedAt    *time.Time `gorm:"column:uploading_completed_at"`
}

func (*NewComic) TableName() string { return COMIC_TABLE }

func (*BriefComic) TableName() string { return COMIC_TABLE }

func (*BasicComic) TableName() string { return COMIC_TABLE }

func (*PatchComic) TableName() string { return COMIC_TABLE }
