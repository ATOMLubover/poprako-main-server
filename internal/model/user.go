package model

type UserInfo struct {
	UserID    string `json:"user_id"`
	QQ        string `json:"qq"`
	Nickname  string `json:"nickname"`
	CreatedAt int64  `json:"created_at"`
}

type LoginArgs struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginToken struct {
	Token string `json:"token"`
}

type UpdateUserArgs struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email,omitempty"`
	Nickname string `json:"nickname,omitempty"`
}
