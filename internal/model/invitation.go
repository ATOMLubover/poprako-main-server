package model

type InvitationInfo struct {
	ID                string `json:"id"`
	InvitorID         string `json:"invitor_id"`
	InviteeQQ         string `json:"invitee_qq"`
	InvCode           string `json:"invitation_code"`
	AssignTranslator  bool   `json:"assign_translator"`
	AssignProofreader bool   `json:"assign_proofreader"`
	AssignTypesetter  bool   `json:"assign_typesetter"`
	AssignRedrawer    bool   `json:"assign_redrawer"`
	AssignReviewer    bool   `json:"assign_reviewer"`
	AssignUploader    bool   `json:"assign_uploader"`
	Pending           bool   `json:"pending"`
	CreatedAt         int64  `json:"created_at"`
}

type CreateInvitationArgs struct {
	InviteeQQ         string `json:"invitee_qq"`
	AssignTranslator  *bool  `json:"assign_translator,omitempty"`
	AssignProofreader *bool  `json:"assign_proofreader,omitempty"`
	AssignTypesetter  *bool  `json:"assign_typesetter,omitempty"`
	AssignRedrawer    *bool  `json:"assign_redrawer,omitempty"`
	AssignReviewer    *bool  `json:"assign_reviewer,omitempty"`
	AssignUploader    *bool  `json:"assign_uploader,omitempty"`
}

type CreateInvitationReply struct {
	InvCode string `json:"invitation_code"`
}
