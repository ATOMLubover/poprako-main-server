package model

type ComicUnitInfo struct {
	ID string `json:"id"`

	PageID string `json:"page_id"`
	Index  int64  `json:"index"`

	XCoordinate float64 `json:"x_coordinate"`
	YCoordinate float64 `json:"y_coordinate"`

	IsInBox bool `json:"is_in_box"`

	TranslatedText    *string `json:"translated_text,omitempty"`
	TranslatorID      *string `json:"translator_id,omitempty"`
	TranslatorComment *string `json:"translator_comment,omitempty"`

	ProvedText         *string `json:"proved_text,omitempty"`
	Proved             bool    `json:"proved"`
	ProofreaderID      *string `json:"proofreader_id,omitempty"`
	ProofreaderComment *string `json:"proofreader_comment,omitempty"`

	// TODO
	CreatorID *string `json:"creator_id,omitempty"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type NewComicUnitArgs struct {
	PageID      string  `json:"page_id"`
	Index       int64   `json:"index"`
	XCoordinate float64 `json:"x_coordinate"`
	YCoordinate float64 `json:"y_coordinate"`
	IsInBox     bool    `json:"is_in_box"`

	TranslatedText    *string `json:"translated_text,omitempty"`
	TranslatorComment *string `json:"translator_comment,omitempty"`

	ProvedText         *string `json:"proved_text,omitempty"`
	Proved             bool    `json:"proved"`
	ProofreaderComment *string `json:"proofreader_comment,omitempty"`
}

type PatchComicUnitArgs struct {
	ID string `json:"id"`

	Index       *int64   `json:"index,omitempty"`
	XCoordinate *float64 `json:"x_coordinate,omitempty"`
	YCoordinate *float64 `json:"y_coordinate,omitempty"`
	IsInBox     *bool    `json:"is_in_box,omitempty"`

	TranslatedText    *string `json:"translated_text,omitempty"`
	TranslatorComment *string `json:"translator_comment,omitempty"`

	ProvedText         *string `json:"proved_text,omitempty"`
	Proved             *bool   `json:"proved,omitempty"`
	ProofreaderComment *string `json:"proofreader_comment,omitempty"`
}
