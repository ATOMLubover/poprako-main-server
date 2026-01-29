package model

type ComicAsgnInfo struct {
	ID                    string `json:"id"`
	ComicID               string `json:"comic_id"`
	UserID                string `json:"user_id"`
	UserNickname          string `json:"user_nickname"`
	AssignedTranslatorAt  *int64 `json:"assigned_translator_at,omitempty"`
	AssignedProofreaderAt *int64 `json:"assigned_proofreader_at,omitempty"`
	AssignedTypesetterAt  *int64 `json:"assigned_typesetter_at,omitempty"`
	AssignedRedrawerAt    *int64 `json:"assigned_redrawer_at,omitempty"`
	AssignedReviewerAt    *int64 `json:"assigned_reviewer_at,omitempty"`
	CreatedAt             int64  `json:"created_at"`
	UpdatedAt             int64  `json:"updated_at"`
}

type CreateComicAsgnArgs struct {
	ComicID       string `json:"comic_id"`
	AssigneeID    string `json:"assignee_id"`
	IsTranslator  *bool  `json:"is_translator"`
	IsProofreader *bool  `json:"is_proofreader"`
	IsTypesetter  *bool  `json:"is_typesetter"`
	IsRedrawer    *bool  `json:"is_redrawer"`
	IsReviewer    *bool  `json:"is_reviewer"`
}

type UpdateComicAsgnArgs struct {
	ID            string `json:"id"`
	IsTranslator  *bool  `json:"is_translator,omitempty"`
	IsProofreader *bool  `json:"is_proofreader,omitempty"`
	IsTypesetter  *bool  `json:"is_typesetter,omitempty"`
	IsRedrawer    *bool  `json:"is_redrawer,omitempty"`
	IsReviewer    *bool  `json:"is_reviewer,omitempty"`
}

type PreAsgnArgs struct {
	AssigneeID    string `json:"assignee_id"`
	IsTranslator  *bool  `json:"is_translator"`
	IsProofreader *bool  `json:"is_proofreader"`
	IsTypesetter  *bool  `json:"is_typesetter"`
	IsRedrawer    *bool  `json:"is_redrawer"`
	IsReviewer    *bool  `json:"is_reviewer"`
}
