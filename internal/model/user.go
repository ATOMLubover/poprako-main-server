package model

type UserInfo struct {
	UserID                string `json:"user_id"`
	QQ                    string `json:"qq"`
	Nickname              string `json:"nickname"`
	AssignedTranslatorAt  int64  `json:"assigned_translator_at"`
	AssignedProofreaderAt int64  `json:"assigned_proofreader_at"`
	AssignedTypesetterAt  int64  `json:"assigned_typesetter_at"`
	AssignedRedrawerAt    int64  `json:"assigned_redrawer_at"`
	AssignedReviewerAt    int64  `json:"assigned_reviewer_at"`
	AssignedUploaderAt    int64  `json:"assigned_uploader_at"`
	IsAdmin               bool   `json:"is_admin"`
	CreatedAt             int64  `json:"created_at"`
}

type LoginArgs struct {
	QQ       string `json:"qq"`
	Password string `json:"password"`
	Nickname string `json:"nickname,omitempty"`
	InvCode  string `json:"invitation_code,omitempty"`
}

type LoginReply struct {
	Token string `json:"token"`
}

type UpdateUserArgs struct {
	UserID   string  `json:"user_id"`
	QQ       *string `json:"qq,omitempty"`
	Nickname *string `json:"nickname,omitempty"`
	// IsAdmin           *bool   `json:"is_admin,omitempty"`
	// AssignTranslator  *bool   `json:"assign_translator,omitempty"`
	// AssignProofreader *bool   `json:"assign_proofreader,omitempty"`
	// AssignTypesetter  *bool   `json:"assign_typesetter,omitempty"`
	// AssignRedrawer    *bool   `json:"assign_redrawer,omitempty"`
	// AssignReviewer    *bool   `json:"assign_reviewer,omitempty"`
	// AssignUploader    *bool   `json:"assign_uploader,omitempty"`
}

type InviteUserArgs struct {
	InviteeID string `json:"invitee_id"`
}

type InviteUserReply struct {
	InvCode string `json:"invitation_code"`
}

type RetrieveUserOpt struct {
	Nickname *string `url:"nn,omitempty"` // Fuzzy

	IsAdmin *bool `url:"ia,omitempty"`

	IsTranslator  *bool `url:"itsl,omitempty"`
	IsProofreader *bool `url:"ipr,omitempty"`
	IsTypesetter  *bool `url:"itst,omitempty"`
	IsRedrawer    *bool `url:"ird,omitempty"`
	IsReviewer    *bool `url:"irv,omitempty"`
	IsUploader    *bool `url:"iul,omitempty"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

type RoleAssignment struct {
	Role     string `json:"role"`
	Assigned bool   `json:"assigned"`
}

type AssignUserRoleArgs struct {
	UserID string           `json:"user_id"`
	Roles  []RoleAssignment `json:"roles"`
}
