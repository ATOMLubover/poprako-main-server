package po

import (
	"time"
)

// Model objects (persistence objects) for user table.

const (
	USER_TABLE = "user_tbl"
)

// Used when creating a new user.
type NewUser struct {
	ID           string `gorm:"column:id;primaryKey"`
	QQ           string `gorm:"column:qq"`
	Nickname     string `gorm:"column:nickname"`
	PasswordHash string `gorm:"column:password_hash"`

	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedUploaderAt    *time.Time `gorm:"column:assigned_uploader_at"`
}

// Used when retrieving basic user info.
type BasicUser struct {
	ID       string `gorm:"column:id;primaryKey"`
	QQ       string `gorm:"column:qq"`
	Nickname string `gorm:"column:nickname"`

	IsAdmin               bool       `gorm:"column:is_admin"`
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedUploaderAt    *time.Time `gorm:"column:assigned_uploader_at"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// Added LastAssignedAt field below

// Used when login.
type SecretUser struct {
	ID      string `gorm:"column:id;primaryKey"`
	PwdHash string `gorm:"column:password_hash"`
}

// Used when updating user info.
// Any fields with default zero values will not be updated.
type PatchUser struct {
	ID       string  `gorm:"column:id;primaryKey"`
	QQ       *string `gorm:"column:qq"`
	Nickname *string `gorm:"column:nickname"`
	IsAdmin  *bool   `gorm:"column:is_admin"`

	// If fields below are set to zero time, that means to erase the assignment time.
	AssignedTranslatorAt  *time.Time `gorm:"column:assigned_translator_at"`
	AssignedProofreaderAt *time.Time `gorm:"column:assigned_proofreader_at"`
	AssignedTypesetterAt  *time.Time `gorm:"column:assigned_typesetter_at"`
	AssignedRedrawerAt    *time.Time `gorm:"column:assigned_redrawer_at"`
	AssignedReviewerAt    *time.Time `gorm:"column:assigned_reviewer_at"`
	AssignedUploaderAt    *time.Time `gorm:"column:assigned_uploader_at"`
}

func (*NewUser) TableName() string { return USER_TABLE }

func (*BasicUser) TableName() string { return USER_TABLE }

func (*SecretUser) TableName() string { return USER_TABLE }

func (*PatchUser) TableName() string { return USER_TABLE }
