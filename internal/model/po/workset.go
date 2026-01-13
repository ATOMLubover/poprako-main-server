package po

const (
	WORKSET_TABLE = "workset_tbl"
)

type NewWorkset struct {
	ID          string  `gorm:"id;primaryKey"`
	Index       int64   `gorm:"index"`
	ComicCount  int64   `gorm:"comic_count"`
	Description *string `gorm:"description"`
	CreatorID   *string `gorm:"creator_id"`
}

type BasicWorkset struct {
	ID          string  `gorm:"id;primaryKey"`
	Index       int64   `gorm:"index"`
	ComicCount  int64   `gorm:"comic_count"`
	Description *string `gorm:"description"`
	CreatorID   *string `gorm:"creator_id"`
	CreatedAt   int64   `gorm:"created_at"`
	UpdatedAt   int64   `gorm:"updated_at"`
}

type PatchWorkset struct {
	ID          string  `gorm:"id;primaryKey"`
	Index       *int64  `gorm:"index"`
	ComicCount  *int64  `gorm:"comic_count"`
	Description *string `gorm:"description"`
	CreatorID   *string `gorm:"creator_id"`
}

func (*NewWorkset) TableName() string { return WORKSET_TABLE }

func (*BasicWorkset) TableName() string { return WORKSET_TABLE }

func (*PatchWorkset) TableName() string { return WORKSET_TABLE }
