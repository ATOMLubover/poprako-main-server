package po

const (
	COMIC_PAGE_TABLE = "comic_page_tbl"
)

type NewComicPage struct {
	ID        string `gorm:"id;primaryKey"`
	ComicID   string `gorm:"comic_id"`
	Index     int64  `gorm:"index"`
	OssKey    string `gorm:"oss_key"`
	SizeBytes int64  `gorm:"size_bytes"`
	Uploaded  *bool  `gorm:"uploaded"`
}

type BasicComicPage struct {
	ID        string `gorm:"id;primaryKey"`
	ComicID   string `gorm:"comic_id"`
	Index     int64  `gorm:"index"`
	OssKey    string `gorm:"oss_key"`
	SizeBytes int64  `gorm:"size_bytes"`
	Uploaded  bool   `gorm:"uploaded"`
	CreatedAt int64  `gorm:"created_at"`
	UpdatedAt int64  `gorm:"updated_at"`
}

type PatchComicPage struct {
	ID        string  `gorm:"id;primaryKey"`
	ComicID   *string `gorm:"comic_id"`
	Index     *int64  `gorm:"index"`
	OssKey    *string `gorm:"oss_key"`
	SizeBytes *int64  `gorm:"size_bytes"`
	Uploaded  *bool   `gorm:"uploaded"`
}

func (*NewComicPage) TableName() string { return COMIC_PAGE_TABLE }

func (*BasicComicPage) TableName() string { return COMIC_PAGE_TABLE }

func (*PatchComicPage) TableName() string { return COMIC_PAGE_TABLE }
