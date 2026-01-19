package po

import (
	"time"
)

const (
	COMIC_UNIT_TABLE = "comic_unit_tbl"
)

// Used when creating a new comic unit.
type NewComicUnit struct {
	ID          string  `gorm:"column:id;primaryKey"`
	PageID      string  `gorm:"column:page_id"`
	Index       int64   `gorm:"column:index"`
	XCoordinate float64 `gorm:"column:x_coordinate"`
	YCoordinate float64 `gorm:"column:y_coordinate"`
	IsInBox     bool    `gorm:"column:is_in_box"`

	TranslatedText    *string `gorm:"column:translated_text"`
	TranslatorID      *string `gorm:"column:translator_id"`
	TranslatorComment *string `gorm:"column:translator_comment"`

	ProvedText         *string `gorm:"column:proved_text"`
	Proved             bool    `gorm:"column:proved"`
	ProofreaderID      *string `gorm:"column:proofreader_id"`
	ProofreaderComment *string `gorm:"column:proofreader_comment"`

	CreatorID *string `gorm:"column:creator_id"`
}

// Used when retrieving basic comic unit info.
type BasicComicUnit struct {
	ID          string  `gorm:"column:id;primaryKey"`
	PageID      string  `gorm:"column:page_id"`
	Index       int64   `gorm:"column:index"`
	XCoordinate float64 `gorm:"column:x_coordinate"`
	YCoordinate float64 `gorm:"column:y_coordinate"`
	IsInBox     bool    `gorm:"column:is_in_box"`

	TranslatedText    *string `gorm:"column:translated_text"`
	TranslatorID      *string `gorm:"column:translator_id"`
	TranslatorComment *string `gorm:"column:translator_comment"`

	ProvedText         *string `gorm:"column:proved_text"`
	Proved             bool    `gorm:"column:proved"`
	ProofreaderID      *string `gorm:"column:proofreader_id"`
	ProofreaderComment *string `gorm:"column:proofreader_comment"`

	CreatorID *string `gorm:"column:creator_id"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// Used when updating comic unit info.
type PatchComicUnit struct {
	ID string `gorm:"column:id;primaryKey"`

	Index       *int64   `gorm:"column:index"`
	XCoordinate *float64 `gorm:"column:x_coordinate"`
	YCoordinate *float64 `gorm:"column:y_coordinate"`
	IsInBox     *bool    `gorm:"column:is_in_box"`

	TranslatedText    *string `gorm:"column:translated_text"`
	TranslatorID      *string `gorm:"column:translator_id"`
	TranslatorComment *string `gorm:"column:translator_comment"`

	ProvedText         *string `gorm:"column:proved_text"`
	Proved             *bool   `gorm:"column:proved"`
	ProofreaderID      *string `gorm:"column:proofreader_id"`
	ProofreaderComment *string `gorm:"column:proofreader_comment"`

	CreatorID *string `gorm:"column:creator_id"`
}

func (*NewComicUnit) TableName() string { return COMIC_UNIT_TABLE }

func (*BasicComicUnit) TableName() string { return COMIC_UNIT_TABLE }

func (*PatchComicUnit) TableName() string { return COMIC_UNIT_TABLE }
