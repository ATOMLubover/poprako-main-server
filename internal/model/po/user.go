package po

// Model objects (persistence objects) for user table.

const (
	USER_TABLE = "user_tbl"
)

// Used when creating a new user.
type NewUser struct {
	ID           string `gorm:"id;primaryKey"`
	QQ           string `gorm:"qq"`
	Nickname     string `gorm:"nickname"`
	PasswordHash string `gorm:"password_hash"`
}

// Used when retrieving basic user info.
type BasicUser struct {
	ID       string `gorm:"id;primaryKey"`
	QQ       string `gorm:"qq"`
	Nickname string `gorm:"nickname"`

	IsAdmin               bool  `gorm:"is_admin"`
	AssignedTranslatorAt  int64 `gorm:"assigned_translator_at"`
	AssignedProofreaderAt int64 `gorm:"assigned_proofreader_at"`
	AssignedTypesetterAt  int64 `gorm:"assigned_typesetter_at"`
	AssignedRedrawerAt    int64 `gorm:"assigned_redrawer_at"`
	AssignedReviewerAt    int64 `gorm:"assigned_reviewer_at"`
	AssignedUploaderAt    int64 `gorm:"assigned_uploader_at"`

	CreatedAt int64 `gorm:"created_at"`
	UpdatedAt int64 `gorm:"updated_at"`
}

// Used when login.
type SecretUser struct {
	ID      string `gorm:"id;primaryKey"`
	PwdHash string `gorm:"password_hash"`
}

// Used when updating user info.
// Any fields with default zero values will not be updated.
type PatchUser struct {
	ID       string  `gorm:"id;primaryKey"`
	QQ       *string `gorm:"qq"`
	Nickname *string `gorm:"nickname"`
	IsAdmin  *bool   `gorm:"is_admin"`

	// If fields below are set 0, that means to erase the assignment time.
	AssignedTranslatorAt  *int64 `gorm:"assigned_translator_at"`
	AssignedProofreaderAt *int64 `gorm:"assigned_proofreader_at"`
	AssignedTypesetterAt  *int64 `gorm:"assigned_typesetter_at"`
	AssignedRedrawerAt    *int64 `gorm:"assigned_redrawer_at"`
	AssignedReviewerAt    *int64 `gorm:"assigned_reviewer_at"`
	AssignedUploaderAt    *int64 `gorm:"assigned_uploader_at"`
}

func (*NewUser) TableName() string { return USER_TABLE }

func (*BasicUser) TableName() string { return USER_TABLE }

func (*SecretUser) TableName() string { return USER_TABLE }

func (*PatchUser) TableName() string { return USER_TABLE }
