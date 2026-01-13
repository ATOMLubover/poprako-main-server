package po

// Model objects (persistence objects) for user table.

const (
	USER_TABLE = "user_tbl"
)

// Used when creating a new user.
type NewUser struct {
	ID           string `gorm:"id;primaryKey"`
	Email        string `gorm:"email"`
	PasswordHash string `gorm:"password_hash"`
}

// Used when retrieving basic user info.
type UserBasic struct {
	ID        string `gorm:"id;primaryKey"`
	Email     string `gorm:"email"`
	Nickname  string `gorm:"nickname"`
	CreatedAt int64  `gorm:"created_at"`
}

// Used when updating user info.
// Any fields with default zero values will not be updated.
type UpdateUser struct {
	ID       string `gorm:"id;primaryKey"`
	Email    string `gorm:"email"`
	Nickname string `gorm:"nickname"`
}

func (*NewUser) TableName() string   { return USER_TABLE }
func (*UserBasic) TableName() string { return USER_TABLE }
