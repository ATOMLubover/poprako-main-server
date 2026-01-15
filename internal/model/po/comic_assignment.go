package po

import "time"

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

	AssignedTranslatorAt  *time.Time `gorm:"assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"assigned_reviewer_at"`

	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

type PatchComicAsgn struct {
	ID      string  `gorm:"id;primaryKey"`
	ComicID *string `gorm:"comic_id"`
	UserID  *string `gorm:"user_id"`

	// If fields below are set to zero time, that means to erase the assignment time.
	AssignedTranslatorAt  *time.Time `gorm:"assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"assigned_reviewer_at"`
}

func (*NewComicAsgn) TableName() string { return COMIC_ASSIGNMENT_TABLE }

func (*BasicComicAsgn) TableName() string { return COMIC_ASSIGNMENT_TABLE }

func (*PatchComicAsgn) TableName() string { return COMIC_ASSIGNMENT_TABLE }
