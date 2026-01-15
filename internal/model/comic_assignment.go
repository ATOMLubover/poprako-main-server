package model

type CreateComicAsgnArgs struct {
	ComicID       string `json:"comic_id"`
	AssigneeID    string `json:"assignee_id"`
	IsTranslator  *bool  `json:"is_translator"`
	IsProofreader *bool  `json:"is_proofreader"`
	IsTypesetter  *bool  `json:"is_typesetter"`
	IsRedrawer    *bool  `json:"is_redrawer"`
	IsReviewer    *bool  `json:"is_reviewer"`
}

type PreAsgnArgs struct {
	AssigneeID    string `json:"assignee_id"`
	IsTranslator  *bool  `json:"is_translator"`
	IsProofreader *bool  `json:"is_proofreader"`
	IsTypesetter  *bool  `json:"is_typesetter"`
	IsRedrawer    *bool  `json:"is_redrawer"`
	IsReviewer    *bool  `json:"is_reviewer"`
}
