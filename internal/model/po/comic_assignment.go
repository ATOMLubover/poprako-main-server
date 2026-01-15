package po

const (
	COMIC_ASSIGNMENT_TABLE = "comic_assignment_tbl"
)

type NewComicAsgn struct {
	ID      string `gorm:"id;primaryKey"`
	ComicID string `gorm:"comic_id"`
	UserID  string `gorm:"user_id"`
}

type BasicComicAsgn struct {
	ID      string `gorm:"id;primaryKey"`
	ComicID string `gorm:"comic_id"`
	UserID  string `gorm:"user_id"`

	AssignedTranslatorAt  *int64 `gorm:"assigned_translator_at"`
	AssignedProofreaderAt *int64 `gorm:"assigned_proofreader_at"`
	AssignedTypesetterAt  *int64 `gorm:"assigned_typesetter_at"`
	AssignedRedrawerAt    *int64 `gorm:"assigned_redrawer_at"`
	AssignedReviewerAt    *int64 `gorm:"assigned_reviewer_at"`

	CreatedAt int64 `gorm:"created_at"`
	UpdatedAt int64 `gorm:"updated_at"`
}

type PatchComicAsgn struct {
	ID      string  `gorm:"id;primaryKey"`
	ComicID *string `gorm:"comic_id"`
	UserID  *string `gorm:"user_id"`

	// If fields below are set 0, that means to erase the assignment time.
	AssignedTranslatorAt  *int64 `gorm:"assigned_translator_at"`
	AssignedProofreaderAt *int64 `gorm:"assigned_proofreader_at"`
	AssignedTypesetterAt  *int64 `gorm:"assigned_typesetter_at"`
	AssignedRedrawerAt    *int64 `gorm:"assigned_redrawer_at"`
	AssignedReviewerAt    *int64 `gorm:"assigned_reviewer_at"`
}

func (*NewComicAsgn) TableName() string { return COMIC_ASSIGNMENT_TABLE }

func (*BasicComicAsgn) TableName() string { return COMIC_ASSIGNMENT_TABLE }

func (*PatchComicAsgn) TableName() string { return COMIC_ASSIGNMENT_TABLE }
