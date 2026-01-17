package po

import (
	"time"
)

const (
	COMIC_PAGE_TABLE = "comic_page_tbl"
)

// Used when creating a new comic page.
type NewComicPage struct {
	ID      string `gorm:"column:id;primaryKey"`
	ComicID string `gorm:"column:comic_id"`
	Index   int64  `gorm:"column:index"`
	OSSKey  string `gorm:"column:oss_key"`
	Uploaded *bool `gorm:"column:uploaded"`
}

// Used when retrieving basic comic page info.
type BasicComicPage struct {
	ID        string    `gorm:"column:id;primaryKey"`
	ComicID   string    `gorm:"column:comic_id"`
	Index     int64     `gorm:"column:index"`
	OSSKey    string    `gorm:"column:oss_key"`
	Uploaded  bool      `gorm:"column:uploaded"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// Used when updating comic page info.
type PatchComicPage struct {
	ID       string  `gorm:"column:id;primaryKey"`
	ComicID  *string `gorm:"column:comic_id"`
	Index    *int64  `gorm:"column:index"`
	OSSKey   *string `gorm:"column:oss_key"`
	Uploaded *bool   `gorm:"column:uploaded"`
}

func (*NewComicPage) TableName() string { return COMIC_PAGE_TABLE }

func (*BasicComicPage) TableName() string { return COMIC_PAGE_TABLE }

func (*PatchComicPage) TableName() string { return COMIC_PAGE_TABLE }