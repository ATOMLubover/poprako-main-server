package model

type ComicBrief struct {
	ID string `json:"id"`

	WorksetID    string `json:"workset_id"`
	WorksetIndex int    `json:"workset_index"`
	Index        int64  `json:"index"`

	Author string `json:"author"`
	Title  string `json:"title"`

	PageCount int64 `json:"page_count"`

	TranslatingStartedAt    *int64 `json:"translating_started_at"`
	TranslatingCompletedAt  *int64 `json:"translating_completed_at"`
	ProofreadingStartedAt   *int64 `json:"proofreading_started_at"`
	ProofreadingCompletedAt *int64 `json:"proofreading_completed_at"`
	TypesettingStartedAt    *int64 `json:"typesetting_started_at"`
	TypesettingCompletedAt  *int64 `json:"typesetting_completed_at"`
	ReviewingCompletedAt    *int64 `json:"reviewing_completed_at"`
	UploadingCompletedAt    *int64 `json:"uploading_completed_at"`
}

type ComicInfo struct {
	ID string `json:"id"`

	WorksetID    string `json:"workset_id"`
	WorksetIndex int    `json:"workset_index"`
	Index        int64  `json:"index"`

	CreatorID       string `json:"creator_id"`
	CreatorNickname string `json:"creator_nickname"`

	Author      string  `json:"author"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Comment     *string `json:"comment,omitempty"`

	PageCount int64 `json:"page_count"`

	TranslatingStartedAt    *int64 `json:"translating_started_at"`
	TranslatingCompletedAt  *int64 `json:"translating_completed_at"`
	ProofreadingStartedAt   *int64 `json:"proofreading_started_at"`
	ProofreadingCompletedAt *int64 `json:"proofreading_completed_at"`
	TypesettingStartedAt    *int64 `json:"typesetting_started_at"`
	TypesettingCompletedAt  *int64 `json:"typesetting_completed_at"`
	ReviewingCompletedAt    *int64 `json:"reviewing_completed_at"`
	UploadingCompletedAt    *int64 `json:"uploading_completed_at"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type RetrieveComicOpt struct {
	// Two below both fuzzy.
	Author *string `url:"aut,omitempty"`
	Title  *string `url:"tit,omitempty"`

	WorksetIndex *string `url:"widx,omitempty"`
	Index        *string `url:"idx,omitempty"`

	// Every group below is only allowed to be one of three states:
	// nil (not care), true, false.
	// What's more, only one of the three states can be true at the same time in each group.
	TranslatingNotStarted *bool `url:"tsl_pending,omitempty"`
	TranslatingInProgress *bool `url:"tsl_wip,omitempty"`
	TranslatingCompleted  *bool `url:"tsl_fin,omitempty"`

	ProofreadingNotStarted *bool `url:"pr_pending,omitempty"`
	ProofreadingInProgress *bool `url:"pr_wip,omitempty"`
	ProofreadingCompleted  *bool `url:"pr_fin,omitempty"`

	TypesettingNotStarted *bool `url:"tst_pending,omitempty"`
	TypesettingInProgress *bool `url:"tst_wip,omitempty"`
	TypesettingCompleted  *bool `url:"tst_fin,omitempty"`

	ReviewingNotStarted *bool `url:"rv_pending,omitempty"`
	ReviewingCompleted  *bool `url:"rv_fin,omitempty"`

	UploadingNotStarted *bool `url:"ul_pending,omitempty"`
	UploadingCompleted  *bool `url:"ul_fin,omitempty"`

	// Accurate.
	AssignedUserID *string `url:"auid,omitempty"`

	Offset int `url:"offset"`
	Limit  int `url:"limit"`
}

type CreateComicArgs struct {
	WorksetID   string  `json:"workset_id"`
	Author      string  `json:"author"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Comment     *string `json:"comment,omitempty"`

	// Pre-assignments when creating the comic.
	// So the ComicID is not known yet, only user IDs and roles can be specified.
	PreAsgns []PreAsgnArgs `json:"pre_asgns,omitempty"`
}

type CreateComicReply struct {
	ID string `json:"id"`
}

type UpdateComicArgs struct {
	ID          string  `json:"id"`
	Author      *string `json:"author,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Comment     *string `json:"comment,omitempty"`

	// Workflow status toggles - true sets timestamp to now, false/nil ignores
	TranslatingStarted    *bool `json:"translating_started,omitempty"`
	TranslatingCompleted  *bool `json:"translating_completed,omitempty"`
	ProofreadingStarted   *bool `json:"proofreading_started,omitempty"`
	ProofreadingCompleted *bool `json:"proofreading_completed,omitempty"`
	TypesettingStarted    *bool `json:"typesetting_started,omitempty"`
	TypesettingCompleted  *bool `json:"typesetting_completed,omitempty"`
	ReviewingCompleted    *bool `json:"reviewing_completed,omitempty"`
	UploadingCompleted    *bool `json:"uploading_completed,omitempty"`
}

type ExportComicReply struct {
	ExportURI string `json:"export_uri"`
}
