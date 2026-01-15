package po

import "time"

const (
	COMIC_UNIT_TABLE = "comic_unit_tbl"
)

type NewComicUnit struct {
	ID          string  `gorm:"id;primaryKey"`
	PageID      string  `gorm:"page_id"`
	Index       int64   `gorm:"index"`
	XCoordinate float64 `gorm:"x_coordinate"`
	YCoordinate float64 `gorm:"y_coordinate"`
	IsInBox     bool    `gorm:"is_in_box"`

	TranslatedText    *string `gorm:"translated_text"`
	TranslatorID      *string `gorm:"translator_id"`
	TranslatorComment *string `gorm:"translator_comment"`

	ProvedText         *string `gorm:"proved_text"`
	Proved             bool    `gorm:"proved"`
	ProofreaderID      *string `gorm:"proofreader_id"`
	ProofreaderComment *string `gorm:"proofreader_comment"`

	CreatorID *string `gorm:"creator_id"`
}

type BasicComicUnit struct {
	ID          string  `gorm:"id;primaryKey"`
	PageID      string  `gorm:"page_id"`
	Index       int64   `gorm:"index"`
	XCoordinate float64 `gorm:"x_coordinate"`
	YCoordinate float64 `gorm:"y_coordinate"`
	IsInBox     bool    `gorm:"is_in_box"`

	TranslatedText    *string `gorm:"translated_text"`
	TranslatorID      *string `gorm:"translator_id"`
	TranslatorComment *string `gorm:"translator_comment"`

	ProvedText         *string `gorm:"proved_text"`
	Proved             bool    `gorm:"proved"`
	ProofreaderID      *string `gorm:"proofreader_id"`
	ProofreaderComment *string `gorm:"proofreader_comment"`

	CreatorID *string `gorm:"creator_id"`

	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

type PatchComicUnit struct {
	ID string `gorm:"id;primaryKey"`

	Index       *int64   `gorm:"index"`
	XCoordinate *float64 `gorm:"x_coordinate"`
	YCoordinate *float64 `gorm:"y_coordinate"`
	IsInBox     *bool    `gorm:"is_in_box"`

	TranslatedText    *string `gorm:"translated_text"`
	TranslatorID      *string `gorm:"translator_id"`
	TranslatorComment *string `gorm:"translator_comment"`

	ProvedText         *string `gorm:"proved_text"`
	Proved             *bool   `gorm:"proved"`
	ProofreaderID      *string `gorm:"proofreader_id"`
	ProofreaderComment *string `gorm:"proofreader_comment"`

	CreatorID *string `gorm:"creator_id"`
}

func (*NewComicUnit) TableName() string { return COMIC_UNIT_TABLE }

func (*BasicComicUnit) TableName() string { return COMIC_UNIT_TABLE }

func (*PatchComicUnit) TableName() string { return COMIC_UNIT_TABLE }
