package po

import "time"

const (
	COMIC_TABLE = "comic_tbl"
)

type NewComic struct {
	ID string `gorm:"id;primaryKey"`

	WorksetID   string  `gorm:"workset_id"`
	Index       int64   `gorm:"index"`
	CreatorID   string  `gorm:"creator_id"`
	Author      string  `gorm:"author"`
	Title       string  `gorm:"title"`
	Comment     *string `gorm:"comment"`
	Description *string `gorm:"description"`
}

type BriefComic struct {
	ID string `gorm:"id;primaryKey"`

	WorksetID    string `gorm:"workset_id"`
	WorksetIndex string `gorm:"workset_index"`
	Index        int64  `gorm:"index"`

	Author string `gorm:"author"`
	Title  string `gorm:"title"`

	TranslatingStartedAt    *time.Time `gorm:"translating_started_at"`
	TranslatingCompletedAt  *time.Time `gorm:"translating_completed_at"`
	ProofreadingStartedAt   *time.Time `gorm:"proofreading_started_at"`
	ProofreadingCompletedAt *time.Time `gorm:"proofreading_completed_at"`
	TypesettingStartedAt    *time.Time `gorm:"typesetting_started_at"`
	TypesettingCompletedAt  *time.Time `gorm:"typesetting_completed_at"`
	ReviewingCompletedAt    *time.Time `gorm:"reviewing_completed_at"`
	UploadingCompletedAt    *time.Time `gorm:"uploading_completed_at"`
}

type BasicComic struct {
	ID string `gorm:"id;primaryKey"`

	WorksetID    string `gorm:"workset_id"`
	WorksetIndex string `gorm:"workset_index"`
	Index        int64  `gorm:"index"`

	CreatorID       string `gorm:"creator_id"`
	CreatorNickname string

	Author      string  `gorm:"author"`
	Title       string  `gorm:"title"`
	Comment     *string `gorm:"comment"`
	Description *string `gorm:"description"`

	PageCount  int64 `gorm:"page_count"`
	LikesCount int64 `gorm:"likes_count"`

	TranslatingStartedAt    *time.Time `gorm:"translating_started_at"`
	TranslatingCompletedAt  *time.Time `gorm:"translating_completed_at"`
	ProofreadingStartedAt   *time.Time `gorm:"proofreading_started_at"`
	ProofreadingCompletedAt *time.Time `gorm:"proofreading_completed_at"`
	TypesettingStartedAt    *time.Time `gorm:"typesetting_started_at"`
	TypesettingCompletedAt  *time.Time `gorm:"typesetting_completed_at"`
	ReviewingCompletedAt    *time.Time `gorm:"reviewing_completed_at"`
	UploadingCompletedAt    *time.Time `gorm:"uploading_completed_at"`

	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

type PatchComic struct {
	ID string `gorm:"id;primaryKey"`

	WorksetID   *string `gorm:"workset_id"`
	Index       *int64  `gorm:"index"`
	CreatorID   *string `gorm:"creator_id"`
	Author      *string `gorm:"author"`
	Title       *string `gorm:"title"`
	Comment     *string `gorm:"comment"`
	Description *string `gorm:"description"`

	PageCount  *int64 `gorm:"page_count"`
	LikesCount *int64 `gorm:"likes_count"`

	TranslatingStartedAt    *time.Time `gorm:"translating_started_at"`
	TranslatingCompletedAt  *time.Time `gorm:"translating_completed_at"`
	ProofreadingStartedAt   *time.Time `gorm:"proofreading_started_at"`
	ProofreadingCompletedAt *time.Time `gorm:"proofreading_completed_at"`
	TypesettingStartedAt    *time.Time `gorm:"typesetting_started_at"`
	TypesettingCompletedAt  *time.Time `gorm:"typesetting_completed_at"`
	ReviewingCompletedAt    *time.Time `gorm:"reviewing_completed_at"`
	UploadingCompletedAt    *time.Time `gorm:"uploading_completed_at"`
}

func (*NewComic) TableName() string { return COMIC_TABLE }

func (*BriefComic) TableName() string { return COMIC_TABLE }

func (*BasicComic) TableName() string { return COMIC_TABLE }

func (*PatchComic) TableName() string { return COMIC_TABLE }
