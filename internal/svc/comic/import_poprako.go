package comic

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"

	"github.com/google/uuid"
)

type poprakoImportProject struct {
	Author string              `json:"author"`
	Title  string              `json:"title"`
	Pages  []poprakoImportPage `json:"pages"`
}

type poprakoImportPage struct {
	ImageFilename string              `json:"image_filename"`
	Units         []poprakoImportUnit `json:"units"`
}

type poprakoImportUnit struct {
	ID             string  `json:"id"`
	X              float64 `json:"x"`
	Y              float64 `json:"y"`
	IndexInPage    uint32  `json:"index_in_page"`
	IsInbox        bool    `json:"is_inbox"`
	TranslatedText *string `json:"translated_text,omitempty"`
	ProovedText    *string `json:"prooved_text,omitempty"`
	IsProoved      bool    `json:"is_prooved"`
	Comment        *string `json:"comment,omitempty"`
	IsLocal        bool    `json:"is_local"`
}

type jsonImportPage struct {
	units []jsonImportUnit
}

type jsonImportUnit struct {
	index              int
	x                  float64
	y                  float64
	isInBox            bool
	translatedText     *string
	provedText         *string
	isProved           bool
	translatorComment  *string
	proofreaderComment *string
}

// ImportPoprakoComic imports a .poprako.json file into the database.
// It follows the same page-level transactional behavior as LabelPlus import.
func ImportPoprakoComic(
	file io.Reader,
	comicID string,
	pageRepo repo.ComicPageRepo,
	unitRepo repo.ComicUnitRepo,
	opts ImportOptions,
) error {
	project, err := decodeAndValidatePoprakoJSON(file)
	if err != nil {
		return fmt.Errorf("failed to parse poprako json: %w", err)
	}

	parsedPages, err := normalizePoprakoProject(project)
	if err != nil {
		return fmt.Errorf("failed to normalize poprako json: %w", err)
	}

	tx := pageRepo.Exct().Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	dbPages, err := pageRepo.GetPagesByComicID(tx, comicID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get pages from database: %w", err)
	}

	sort.Slice(dbPages, func(i, j int) bool {
		return dbPages[i].Index < dbPages[j].Index
	})

	if len(parsedPages) != len(dbPages) {
		tx.Rollback()
		return fmt.Errorf("page count mismatch: file has %d pages, database has %d pages", len(parsedPages), len(dbPages))
	}

	for i, parsedPage := range parsedPages {
		dbPage := dbPages[i]

		existingUnits, err := unitRepo.GetUnitsByPageID(tx, dbPage.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get existing units for page %s: %w", dbPage.ID, err)
		}

		if !opts.IsProofreader {
			hasProvedUnit := false
			for _, u := range existingUnits {
				if u.Proved {
					hasProvedUnit = true
					break
				}
			}
			if hasProvedUnit {
				continue
			}
		}

		if len(existingUnits) > 0 {
			unitIDs := make([]string, len(existingUnits))
			for j, unit := range existingUnits {
				unitIDs[j] = unit.ID
			}

			if err := unitRepo.DeleteUnitByIDs(tx, unitIDs); err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to delete existing units for page %s: %w", dbPage.ID, err)
			}
		}

		newUnits := make([]po.NewComicUnit, len(parsedPage.units))
		for j, parsedUnit := range parsedPage.units {
			unitID, err := uuid.NewV7()
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to generate UUID for unit: %w", err)
			}

			translatedText := parsedUnit.translatedText
			provedText := parsedUnit.provedText
			translatorComment := parsedUnit.translatorComment
			proofreaderComment := parsedUnit.proofreaderComment

			var translatorID *string
			var proofreaderID *string
			proved := parsedUnit.isProved

			if opts.IsProofreader {
				proofreaderID = &opts.UserID
			} else {
				translatorID = &opts.UserID
			}

			newUnits[j] = po.NewComicUnit{
				ID:                 unitID.String(),
				PageID:             dbPage.ID,
				Index:              int64(parsedUnit.index),
				XCoordinate:        parsedUnit.x,
				YCoordinate:        parsedUnit.y,
				IsInBox:            parsedUnit.isInBox,
				TranslatedText:     translatedText,
				TranslatorComment:  translatorComment,
				ProvedText:         provedText,
				Proved:             proved,
				ProofreaderComment: proofreaderComment,
				TranslatorID:       translatorID,
				ProofreaderID:      proofreaderID,
				CreatorID:          &opts.UserID,
			}
		}

		if err := unitRepo.CreateUnits(tx, newUnits); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create units for page %s: %w", dbPage.ID, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func decodeAndValidatePoprakoJSON(file io.Reader) (*poprakoImportProject, error) {
	var raw map[string]json.RawMessage
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("invalid json: %w", err)
	}

	if _, ok := raw["author"]; !ok {
		return nil, fmt.Errorf("missing required field: author")
	}
	if _, ok := raw["title"]; !ok {
		return nil, fmt.Errorf("missing required field: title")
	}
	pagesRaw, ok := raw["pages"]
	if !ok {
		return nil, fmt.Errorf("missing required field: pages")
	}
	if string(pagesRaw) == "null" {
		return nil, fmt.Errorf("pages must not be null")
	}

	var project poprakoImportProject
	if err := json.Unmarshal(pagesRaw, &project.Pages); err != nil {
		return nil, fmt.Errorf("invalid pages field: %w", err)
	}

	if err := json.Unmarshal(raw["author"], &project.Author); err != nil {
		return nil, fmt.Errorf("invalid author field: %w", err)
	}
	if err := json.Unmarshal(raw["title"], &project.Title); err != nil {
		return nil, fmt.Errorf("invalid title field: %w", err)
	}

	if strings.TrimSpace(project.Author) == "" {
		return nil, fmt.Errorf("author must not be empty")
	}
	if strings.TrimSpace(project.Title) == "" {
		return nil, fmt.Errorf("title must not be empty")
	}

	return &project, nil
}

func normalizePoprakoProject(project *poprakoImportProject) ([]jsonImportPage, error) {
	result := make([]jsonImportPage, len(project.Pages))

	for i, page := range project.Pages {
		if strings.TrimSpace(page.ImageFilename) == "" {
			return nil, fmt.Errorf("page %d: image_filename must not be empty", i+1)
		}

		seenIndex := make(map[uint32]struct{}, len(page.Units))
		seenIDLocal := make(map[string]struct{}, len(page.Units))

		normalizedUnits := make([]jsonImportUnit, len(page.Units))
		for j, unit := range page.Units {
			if strings.TrimSpace(unit.ID) == "" {
				return nil, fmt.Errorf("page %d unit %d: id must not be empty", i+1, j+1)
			}
			if unit.IndexInPage < 1 {
				return nil, fmt.Errorf("page %d unit %d: index_in_page must be >= 1", i+1, j+1)
			}
			if !isFinite(unit.X) || !isFinite(unit.Y) {
				return nil, fmt.Errorf("page %d unit %d: x/y must be finite numbers", i+1, j+1)
			}

			if _, exists := seenIndex[unit.IndexInPage]; exists {
				return nil, fmt.Errorf("page %d: duplicated index_in_page %d", i+1, unit.IndexInPage)
			}
			seenIndex[unit.IndexInPage] = struct{}{}

			idLocalKey := fmt.Sprintf("%s|%t", unit.ID, unit.IsLocal)
			if _, exists := seenIDLocal[idLocalKey]; exists {
				return nil, fmt.Errorf("page %d: duplicated (id,is_local) pair for id=%s", i+1, unit.ID)
			}
			seenIDLocal[idLocalKey] = struct{}{}

			translatorComment, proofreaderComment := parseComments(derefString(normalizeOptionalText(unit.Comment)))

			normalizedUnits[j] = jsonImportUnit{
				index:              int(unit.IndexInPage),
				x:                  unit.X,
				y:                  unit.Y,
				isInBox:            unit.IsInbox,
				translatedText:     normalizeOptionalText(unit.TranslatedText),
				provedText:         normalizeOptionalText(unit.ProovedText),
				isProved:           unit.IsProoved,
				translatorComment:  translatorComment,
				proofreaderComment: proofreaderComment,
			}
		}

		sort.Slice(normalizedUnits, func(a, b int) bool {
			return normalizedUnits[a].index < normalizedUnits[b].index
		})

		result[i] = jsonImportPage{units: normalizedUnits}
	}

	return result, nil
}

func isFinite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}
