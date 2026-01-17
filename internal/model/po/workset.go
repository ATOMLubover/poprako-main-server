package po

import (
	"time"
)

const (
	WORKSET_TABLE = "workset_tbl"
)

// Used when creating a new workset.
type NewWorkset struct {
	ID          string  `gorm:"column:id;primaryKey"`
	Name        string  `gorm:"column:name"`
	Index       int64   `gorm:"column:index"`
	ComicCount  int64   `gorm:"column:comic_count"`
	Description *string `gorm:"column:description"`
	CreatorID   string  `gorm:"column:creator_id"`
}

// Used when retrieving detailed workset info.
type DetailedWorkset struct {
	ID              string    `gorm:"column:id;primaryKey"`
	Index           int64     `gorm:"column:index"`
	Name            string    `gorm:"column:name"`
	ComicCount      int64     `gorm:"column:comic_count"`
	Description     *string   `gorm:"column:description"`
	CreatorID       string    `gorm:"column:creator_id"`
	CreatorNickname string    `gorm:"column:creator_nickname"` // 补全了映射
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

// Used when updating workset info.
// Any fields with default zero values (nil) will not be updated.
type PatchWorkset struct {
	ID          string  `gorm:"column:id;primaryKey"`
	Name        *string `gorm:"column:name"`
	Index       *int64  `gorm:"column:index"`
	ComicCount  *int64  `gorm:"column:comic_count"`
	Description *string `gorm:"column:description"`
	CreatorID   *string `gorm:"column:creator_id"`
}

func (*NewWorkset) TableName() string { return WORKSET_TABLE }

func (*DetailedWorkset) TableName() string { return WORKSET_TABLE }

func (*PatchWorkset) TableName() string { return WORKSET_TABLE }