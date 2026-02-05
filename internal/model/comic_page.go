package model

type ComicPageInfo struct {
	ID       string `json:"id"`
	ComicID  string `json:"comic_id"`
	Index    int64  `json:"index"`
	OSSURL   string `json:"oss_url"`
	Uploaded bool   `json:"uploaded"`

	InboxUnitCount      int64 `json:"inbox_unit_count"`
	OutboxUnitCount     int64 `json:"outbox_unit_count"`
	TranslatedUnitCount int64 `json:"translated_unit_count"`
	ProvedUnitCount     int64 `json:"proved_unit_count"`
}

type CreateComicPageArgs struct {
	// OSSKey/URL is composed of "comic/{comic_id}/page_{index}.ext"
	ComicID  string `json:"comic_id"`
	Index    int64  `json:"index"`
	ImageExt string `json:"image_ext"`
}

type CreateComicPageReply struct {
	ID     string `json:"id"`
	OSSURL string `json:"oss_url"`
}

type RecreateComicPageArgs struct {
	ID       string `json:"id"`
	ImageExt string `json:"image_ext"`
}

type PatchComicPageArgs struct {
	ID string `json:"id"`
	// After upload, client should report the page as uploaded
	ImageExt *string `json:"image_ext,omitempty"`
	Uploaded *bool   `json:"uploaded,omitempty"`
}
