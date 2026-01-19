package po

import "time"

const (
	INVITATION_TABLE = "invitation_tbl"
)

type BasicInvitation struct {
	ID string `gorm:"column:id;primaryKey"`

	InvitorID string `gorm:"column:invitor_id"`
	InviteeQQ string `gorm:"column:invitee_qq"`
	InvCode   string `gorm:"column:invitation_code"`

	AssignTranslator  bool `gorm:"column:assign_translator"`
	AssignProofreader bool `gorm:"column:assign_proofreader"`
	AssignTypesetter  bool `gorm:"column:assign_typesetter"`
	AssignRedrawer    bool `gorm:"column:assign_redrawer"`
	AssignReviewer    bool `gorm:"column:assign_reviewer"`
	AssignUploader    bool `gorm:"column:assign_uploader"`

	Pending bool `gorm:"column:pending"`

	CreatedAt time.Time `gorm:"column:created_at"`
}

type NewInvitation struct {
	ID string `gorm:"column:id;primaryKey"`

	InvitorID string `gorm:"column:invitor_id"`
	InviteeQQ string `gorm:"column:invitee_qq"`
	InvCode   string `gorm:"column:invitation_code"`

	AssignTranslator  bool `gorm:"column:assign_translator"`
	AssignProofreader bool `gorm:"column:assign_proofreader"`
	AssignTypesetter  bool `gorm:"column:assign_typesetter"`
	AssignRedrawer    bool `gorm:"column:assign_redrawer"`
	AssignReviewer    bool `gorm:"column:assign_reviewer"`
	AssignUploader    bool `gorm:"column:assign_uploader"`
}

func (*BasicInvitation) TableName() string { return INVITATION_TABLE }

func (*NewInvitation) TableName() string { return INVITATION_TABLE }
