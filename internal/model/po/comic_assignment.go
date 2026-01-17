package po

import (
	"time"
)

const (
	COMIC_ASSIGNMENT_TABLE = "comic_assignment_tbl"
)

// Used when creating a new comic assignment.
type NewComicAsgn struct {
	ID      string `gorm:"column:id;primaryKey"`
	ComicID string `gorm:"column:comic_id"`
	UserID  string `gorm:"column:user_id"`
}

// Used when retrieving basic comic assignment info.
type BasicComicAsgn struct {
	ID      string `gorm:"column:id;primaryKey"`
	ComicID string `gorm:"column:comic_id"`
	UserID  string `gorm:"column:user_id"`

	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// Used when updating comic assignment info.
type PatchComicAsgn struct {
	ID      string  `gorm:"column:id;primaryKey"`
	ComicID *string `gorm:"column:comic_id"`
	UserID  *string `gorm:"column:user_id"`

	// If fields below are set to zero time, that means to erase the assignment time.
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
}

func (*NewComicAsgn) TableName() string { return COMIC_ASSIGNMENT_TABLE }

func (*BasicComicAsgn) TableName() string { return COMIC_ASSIGNMENT_TABLE }

func (*PatchComicAsgn) TableName() string { return COMIC_ASSIGNMENT_TABLE }