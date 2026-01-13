package po

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

type BasicComic struct {
	ID string `gorm:"id;primaryKey"`

	WorksetID   string  `gorm:"workset_id"`
	Index       int64   `gorm:"index"`
	CreatorID   string  `gorm:"creator_id"`
	Author      string  `gorm:"author"`
	Title       string  `gorm:"title"`
	Comment     *string `gorm:"comment"`
	Description *string `gorm:"description"`

	PageCount  int64 `gorm:"page_count"`
	LikesCount int64 `gorm:"likes_count"`

	TranslatingStartedAt    *int64 `gorm:"translating_started_at"`
	TranslatingCompletedAt  *int64 `gorm:"translating_completed_at"`
	ProofreadingStartedAt   *int64 `gorm:"proofreading_started_at"`
	ProofreadingCompletedAt *int64 `gorm:"proofreading_completed_at"`
	TypesettingStartedAt    *int64 `gorm:"typesetting_started_at"`
	TypesettingCompletedAt  *int64 `gorm:"typesetting_completed_at"`
	ReviewingCompletedAt    *int64 `gorm:"reviewing_completed_at"`
	UploadingCompletedAt    *int64 `gorm:"uploading_completed_at"`

	CreatedAt int64 `gorm:"created_at"`
	UpdatedAt int64 `gorm:"updated_at"`
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

	TranslatingStartedAt    *int64 `gorm:"translating_started_at"`
	TranslatingCompletedAt  *int64 `gorm:"translating_completed_at"`
	ProofreadingStartedAt   *int64 `gorm:"proofreading_started_at"`
	ProofreadingCompletedAt *int64 `gorm:"proofreading_completed_at"`
	TypesettingStartedAt    *int64 `gorm:"typesetting_started_at"`
	TypesettingCompletedAt  *int64 `gorm:"typesetting_completed_at"`
	ReviewingCompletedAt    *int64 `gorm:"reviewing_completed_at"`
	UploadingCompletedAt    *int64 `gorm:"uploading_completed_at"`
}

func (*NewComic) TableName() string { return COMIC_TABLE }

func (*BasicComic) TableName() string { return COMIC_TABLE }

func (*PatchComic) TableName() string { return COMIC_TABLE }
