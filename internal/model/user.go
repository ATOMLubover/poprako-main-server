package model

type UserInfo struct {
	UserID    string `json:"user_id"`
	QQ        string `json:"qq"`
	Nickname  string `json:"nickname"`
	CreatedAt int64  `json:"created_at"`
}

type LoginArgs struct {
	QQ       string `json:"qq"`
	Password string `json:"password"`
	Nickname string `json:"nickname,omitempty"`
	InvCode  string `json:"inv_code,omitempty"`
}

type LoginToken struct {
	Token string `json:"token"`
}

type UpdateUserArgs struct {
	UserID            string  `json:"user_id"`
	QQ                *string `json:"qq,omitempty"`
	Nickname          *string `json:"nickname,omitempty"`
	IsAdmin           *bool   `json:"is_admin,omitempty"`
	AssignTranslator  *bool   `json:"assign_translator,omitempty"`
	AssignProofreader *bool   `json:"assign_proofreader,omitempty"`
	AssignTypesetter  *bool   `json:"assign_typesetter,omitempty"`
	AssignRedrawer    *bool   `json:"assign_redrawer,omitempty"`
	AssignReviewer    *bool   `json:"assign_reviewer,omitempty"`
	AssignUploader    *bool   `json:"assign_uploader,omitempty"`
}
