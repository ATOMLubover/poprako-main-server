package comic

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"

	"github.com/google/uuid"
)

// parsedPage represents a page parsed from LabelPlus format.
type parsedPage struct {
	units []parsedUnit
}

// parsedUnit represents a unit parsed from LabelPlus format.
type parsedUnit struct {
	index              int // 1-based index in page
	x                  float64
	y                  float64
	isInBox            bool
	text               string
	translatorComment  *string
	proofreaderComment *string
}

// ImportOptions defines options for importing a comic.
type ImportOptions struct {
	IsProofreader bool
	UserID        string
}

var (
	// Regex patterns for parsing LabelPlus format
	pageHeaderRegex = regexp.MustCompile(`^>>>>>>>>\[.+\]<<<<<<<<$`)
	unitHeaderRegex = regexp.MustCompile(`^----------------\[(\d+)\]----------------\[(-?[\d.]+),(-?[\d.]+),([12])\]$`)
	commentRegex    = regexp.MustCompile(`^#\[翻校注释\]：(.*)$`)
)

// ImportLabelplusComic imports a LabelPlus format file into the database.
// The number of pages in the file must match the number of pages in the database.
// Pages are matched by order: first parsed page -> first DB page by index, etc.
func ImportLabelplusComic(
	file io.Reader,
	comicID string,
	pageRepo repo.ComicPageRepo,
	unitRepo repo.ComicUnitRepo,
	opts ImportOptions,
) error {
	// 1. Parse the LabelPlus file
	parsedPages, err := parseLabelPlusFile(file)
	if err != nil {
		return fmt.Errorf("failed to parse LabelPlus file: %w", err)
	}

	// 2. Start transaction
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

	// 3. Get all pages from database and sort by index
	dbPages, err := pageRepo.GetPagesByComicID(tx, comicID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get pages from database: %w", err)
	}

	sort.Slice(dbPages, func(i, j int) bool {
		return dbPages[i].Index < dbPages[j].Index
	})

	// 4. Validate page count match
	if len(parsedPages) != len(dbPages) {
		tx.Rollback()
		return fmt.Errorf("page count mismatch: file has %d pages, database has %d pages", len(parsedPages), len(dbPages))
	}

	// 5. Process each page: delete old units and insert new ones
	for i, parsedPage := range parsedPages {
		dbPage := dbPages[i]

		// Get existing units to delete
		existingUnits, err := unitRepo.GetUnitsByPageID(tx, dbPage.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to get existing units for page %s: %w", dbPage.ID, err)
		}

		// Check if translator can overwrite this page
		if !opts.IsProofreader {
			hasProvedUnit := false
			for _, u := range existingUnits {
				if u.Proved {
					hasProvedUnit = true
					break
				}
			}
			if hasProvedUnit {
				// Skip this page - translator cannot overwrite proofread content
				continue
			}
		}

		// Delete all existing units
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

		// Create new units
		newUnits := make([]po.NewComicUnit, len(parsedPage.units))
		for j, parsedUnit := range parsedPage.units {
			unitID, err := uuid.NewV7()
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to generate UUID for unit: %w", err)
			}

			// Map fields based on importer role
			var translatedText *string
			var provedText *string
			var translatorID *string
			var proofreaderID *string
			var proved bool

			if opts.IsProofreader {
				// Proofreader: write to proved layer
				provedText = stringPtrOrNil(parsedUnit.text)
				proved = true
				proofreaderID = &opts.UserID
			} else {
				// Translator: write to translation layer
				translatedText = stringPtrOrNil(parsedUnit.text)
				proved = false
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
				TranslatorComment:  parsedUnit.translatorComment,
				ProvedText:         provedText,
				Proved:             proved,
				ProofreaderComment: parsedUnit.proofreaderComment,
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

	// 6. Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// parseLabelPlusFile parses a LabelPlus format file and returns the parsed pages.
func parseLabelPlusFile(r io.Reader) ([]parsedPage, error) {
	scanner := bufio.NewScanner(r)

	// Validate header
	if err := validateHeader(scanner); err != nil {
		return nil, err
	}

	// Parse pages and units
	var pages []parsedPage
	var currentPage *parsedPage
	var currentUnit *parsedUnit
	var mainTextBuffer []string
	var commentBuffer []string

	for scanner.Scan() {
		line := scanner.Text()

		// Check for page header
		if pageHeaderRegex.MatchString(line) {
			// Save previous unit if exists
			if currentUnit != nil {
				applyParsedData(currentUnit, mainTextBuffer, commentBuffer)
				currentPage.units = append(currentPage.units, *currentUnit)
				currentUnit = nil
				mainTextBuffer = nil
				commentBuffer = nil
			}

			// Save previous page if exists
			if currentPage != nil {
				pages = append(pages, *currentPage)
			}

			// Start new page
			currentPage = &parsedPage{}
			continue
		}

		// Check for unit header
		if matches := unitHeaderRegex.FindStringSubmatch(line); matches != nil {
			// Save previous unit if exists
			if currentUnit != nil {
				applyParsedData(currentUnit, mainTextBuffer, commentBuffer)
				currentPage.units = append(currentPage.units, *currentUnit)
				mainTextBuffer = nil
				commentBuffer = nil
			}

			// Parse unit header
			index, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, fmt.Errorf("failed to parse unit index: %w", err)
			}

			x, err := strconv.ParseFloat(matches[2], 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse X coordinate: %w", err)
			}

			y, err := strconv.ParseFloat(matches[3], 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse Y coordinate: %w", err)
			}

			groupID := matches[4]

			currentUnit = &parsedUnit{
				index:   index,
				x:       x,
				y:       y,
				isInBox: groupID == "1",
			}
			continue
		}

		// Check for comment line (marks transition to comment section)
		if matches := commentRegex.FindStringSubmatch(line); matches != nil {
			commentText := matches[1]
			commentBuffer = append(commentBuffer, commentText)
			continue
		}

		// Regular text line or empty line
		if currentUnit != nil {
			if line != "" {
				// If we have a comment buffer, new text goes to comments
				// Otherwise, new text goes to main text
				if len(commentBuffer) > 0 {
					commentBuffer = append(commentBuffer, line)
				} else {
					mainTextBuffer = append(mainTextBuffer, line)
				}
			}
		}
	}

	// Save last unit and page
	if currentUnit != nil {
		applyParsedData(currentUnit, mainTextBuffer, commentBuffer)
		currentPage.units = append(currentPage.units, *currentUnit)
	}
	if currentPage != nil {
		pages = append(pages, *currentPage)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return pages, nil
}

// validateHeader validates the LabelPlus file header.
func validateHeader(scanner *bufio.Scanner) error {
	expectedLines := []string{
		"1,0",
		"-",
		"框内",
		"框外",
		"-",
	}

	for i, expected := range expectedLines {
		if !scanner.Scan() {
			return fmt.Errorf("unexpected end of file at header line %d", i+1)
		}
		if scanner.Text() != expected {
			return fmt.Errorf("invalid header at line %d: expected '%s', got '%s'", i+1, expected, scanner.Text())
		}
	}

	// Skip user comment line and empty line
	if !scanner.Scan() {
		return fmt.Errorf("unexpected end of file after header")
	}
	// Skip comment content (any text)

	if !scanner.Scan() {
		return fmt.Errorf("unexpected end of file after header comment")
	}
	// Should be empty line
	if scanner.Text() != "" {
		return fmt.Errorf("expected empty line after header comment, got '%s'", scanner.Text())
	}

	return nil
}

// applyParsedData applies the parsed main text and comment data to the unit.
func applyParsedData(unit *parsedUnit, mainTextBuffer []string, commentBuffer []string) {
	// Apply main text
	if len(mainTextBuffer) > 0 {
		unit.text = strings.Join(mainTextBuffer, "\n")
	}

	// Apply comments
	if len(commentBuffer) > 0 {
		commentText := strings.Join(commentBuffer, "\n")
		unit.translatorComment, unit.proofreaderComment = parseComments(commentText)
	}
}

// parseComments splits comment text into translator and proofreader comments.
func parseComments(commentText string) (*string, *string) {
	var translatorComment, proofreaderComment *string

	lines := strings.Split(commentText, "\n")
	var translatorLines []string
	var proofreaderLines []string
	var currentTarget *[]string

	for _, line := range lines {
		if strings.HasPrefix(line, "【翻译】") {
			content := strings.TrimPrefix(line, "【翻译】")
			translatorLines = append(translatorLines, content)
			currentTarget = &translatorLines
		} else if strings.HasPrefix(line, "【校对】") {
			content := strings.TrimPrefix(line, "【校对】")
			proofreaderLines = append(proofreaderLines, content)
			currentTarget = &proofreaderLines
		} else if currentTarget != nil {
			*currentTarget = append(*currentTarget, line)
		} else {
			// No prefix, treat as translator comment by default
			translatorLines = append(translatorLines, line)
			currentTarget = &translatorLines
		}
	}

	if len(translatorLines) > 0 {
		text := strings.Join(translatorLines, "\n")
		translatorComment = &text
	}
	if len(proofreaderLines) > 0 {
		text := strings.Join(proofreaderLines, "\n")
		proofreaderComment = &text
	}

	return translatorComment, proofreaderComment
}

// stringPtrOrNil returns a pointer to the string if non-empty, otherwise nil.
func stringPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
