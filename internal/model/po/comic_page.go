package po

import "time"

const (
	COMIC_PAGE_TABLE = "comic_page_tbl"
)

type NewComicPage struct {
	ID       string `gorm:"id;primaryKey"`
	ComicID  string `gorm:"comic_id"`
	Index    int64  `gorm:"index"`
	OSSKey   string `gorm:"oss_key"`
	Uploaded *bool  `gorm:"uploaded"`
}

type BasicComicPage struct {
	ID        string    `gorm:"id;primaryKey"`
	ComicID   string    `gorm:"comic_id"`
	Index     int64     `gorm:"index"`
	OSSKey    string    `gorm:"oss_key"`
	Uploaded  bool      `gorm:"uploaded"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

type PatchComicPage struct {
	ID       string  `gorm:"id;primaryKey"`
	ComicID  *string `gorm:"comic_id"`
	Index    *int64  `gorm:"index"`
	OSSKey   *string `gorm:"oss_key"`
	Uploaded *bool   `gorm:"uploaded"`
}

func (*NewComicPage) TableName() string { return COMIC_PAGE_TABLE }

func (*BasicComicPage) TableName() string { return COMIC_PAGE_TABLE }

func (*PatchComicPage) TableName() string { return COMIC_PAGE_TABLE }
