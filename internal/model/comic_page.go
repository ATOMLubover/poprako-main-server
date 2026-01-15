package model

type ComicPageInfo struct {
	ID       string `json:"id"`
	ComicID  string `json:"comic_id"`
	Index    int64  `json:"index"`
	OSSURL   string `json:"oss_url"`
	Uploaded bool   `json:"uploaded"`
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

type PatchComicPageArgs struct {
	ID string `json:"id"`
	// After upload, client should report the page as uploaded
	Uploaded *bool `json:"uploaded,omitempty"`
}
