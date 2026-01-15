package model

type WorksetInfo struct {
	ID              string  `json:"id"`
	Index           int64   `json:"index"`
	Name            string  `json:"name"`
	ComicCount      int64   `json:"comic_count"`
	Description     *string `json:"description"`
	CreatorID       string  `json:"creator_id"`
	CreatorNickname string  `json:"creator_nickname"`
	CreatedAt       int64   `json:"created_at"`
	UpdatedAt       int64   `json:"updated_at"`
}

type CreateWorksetArgs struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type CreateWorksetReply struct {
	ID string `json:"id"`
}

type UpdateWorksetArgs struct {
	ID          string  `json:"id"`
	Description *string `json:"description,omitempty"`
}
