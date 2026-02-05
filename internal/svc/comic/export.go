package comic

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"
)

// ExportLabelplusComic exports a comic to LabelPlus format file.
// Returns the absolute file path on success.
func ExportLabelplusComic(
	comicID string,
	exportDir string,
	comicRepo repo.ComicRepo,
	comicPageRepo repo.ComicPageRepo,
	comicUnitRepo repo.ComicUnitRepo,
) (string, error) {
	// 1. Get comic basic info
	comic, err := comicRepo.GetComicByID(nil, comicID)
	if err != nil {
		return "", fmt.Errorf("failed to get comic: %w", err)
	}

	// 2. Fetch all pages and sort by index
	pages, err := comicPageRepo.GetPagesByComicID(nil, comicID)
	if err != nil {
		return "", fmt.Errorf("failed to get pages: %w", err)
	}

	sort.Slice(pages, func(i, j int) bool {
		return pages[i].Index < pages[j].Index
	})

	// 3. Create export file
	now := time.Now()
	safeAuthor := truncateRunes(sanitizeFilename(comic.Author), 20)
	safeTitle := truncateRunes(sanitizeFilename(comic.Title), 60)

	fileName := fmt.Sprintf("【%s】%s-%s.labelplus.txt", safeAuthor, safeTitle, now.Format("20060102150405"))
	
	// Ensure single filename component stays within filesystem limits (use 255 bytes as safe limit)
	const maxComponentBytes = 255
	
	for len([]byte(fileName)) > maxComponentBytes {
		if len([]rune(safeTitle)) > 1 {
			safeTitle = shortenByOneRune(safeTitle)
		} else if len([]rune(safeAuthor)) > 1 {
			safeAuthor = shortenByOneRune(safeAuthor)
		} else {
			break
		}

		fileName = fmt.Sprintf("【%s】%s-%s.labelplus.txt", safeAuthor, safeTitle, now.Format("20060102150405"))
	}

	filePath := filepath.Join(exportDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	defer file.Close()

	// 4. Write header
	if err := writeLabelPlusHeader(file); err != nil {
		return "", fmt.Errorf("failed to write header: %w", err)
	}

	// 5. Stream write pages and their units
	for _, page := range pages {
		if err := writeLabelPlusPage(file, page, comicUnitRepo); err != nil {
			return "", fmt.Errorf("failed to write page %s: %w", page.ID, err)
		}
	}

	return filePath, nil
}

// writeLabelPlusHeader writes the fixed LabelPlus format header.
func writeLabelPlusHeader(w io.Writer) error {
	_, err := fmt.Fprintf(w, "1,0\n-\n框内\n框外\n-\nExported by PopRaKo Web\n\n")

	return err
}

// writeLabelPlusPage writes a page header and all its units.
func writeLabelPlusPage(w io.Writer, page po.BasicComicPage, unitRepo repo.ComicUnitRepo) error {
	// Write page header: \n\n>>>>>>>>[filename]<<<<<<<<\n
	imgExt := filepath.Ext(page.OSSKey)
	imgName := fmt.Sprintf("page_%d%s", page.Index, imgExt)
	
	if _, err := fmt.Fprintf(w, "\n\n>>>>>>>>[%s]<<<<<<<<\n", imgName); err != nil {
		return err
	}

	// Fetch and sort units by index
	units, err := unitRepo.GetUnitsByPageID(nil, page.ID)
	if err != nil {
		return fmt.Errorf("failed to get units: %w", err)
	}

	sort.Slice(units, func(i, j int) bool {
		return units[i].Index < units[j].Index
	})

	// Write each unit
	for i, unit := range units {
		if err := writeLabelPlusUnit(w, unit, i+1); err != nil {
			return fmt.Errorf("failed to write unit %d: %w", i+1, err)
		}
	}

	return nil
}

// writeLabelPlusUnit writes a single unit in LabelPlus format.
func writeLabelPlusUnit(w io.Writer, unit po.BasicComicUnit, n int) error {
	// Determine group ID: 1=inbox, 2=outbox
	g := 2
	if unit.IsInBox {
		g = 1
	}

	// Write unit header: ----------------[N]----------------[X,Y,G]
	if _, err := fmt.Fprintf(w, "----------------[%d]----------------[%.4f,%.4f,%d]\n",
		n, unit.XCoordinate, unit.YCoordinate, g); err != nil {
		return err
	}

	// Write main text (priority: ProvedText > TranslatedText)
	mainText := selectMainText(unit.ProvedText, unit.TranslatedText)
	if mainText != "" {
		if _, err := fmt.Fprintf(w, "%s\n", mainText); err != nil {
			return err
		}
	}

	// Write comment if exists
	comment := formatComment(unit.TranslatorComment, unit.ProofreaderComment)
	if comment != "" {
		if _, err := fmt.Fprintf(w, "\n#[翻校注释]：%s\n", comment); err != nil {
			return err
		}
	}

	// Write unit ending newline
	if _, err := fmt.Fprintf(w, "\n"); err != nil {
		return err
	}

	return nil
}

// selectMainText returns the main text with priority: proved > translated.
func selectMainText(provedText, translatedText *string) string {
	if provedText != nil && *provedText != "" {
		return *provedText
	}

	if translatedText != nil && *translatedText != "" {
		return *translatedText
	}

	return ""
}

// formatComment merges translator and proofreader comments.
func formatComment(translatorComment, proofreaderComment *string) string {
	var parts []string

	if translatorComment != nil && *translatorComment != "" {
		parts = append(parts, "【翻译】"+*translatorComment)
	}
	if proofreaderComment != nil && *proofreaderComment != "" {
		parts = append(parts, "【校对】"+*proofreaderComment)
	}

	return strings.Join(parts, "\n")
}

// sanitizeFilename removes invalid characters for filenames.
func sanitizeFilename(s string) string {
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}

	result := s
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}

	return result
}

// truncateRunes truncates a string to at most `limit` runes and appends an ellipsis if truncated.
func truncateRunes(s string, limit int) string {
	r := []rune(s)
	if len(r) <= limit {
		return s
	}

	return string(r[:limit]) + "…"
}

// shortenByOneRune removes the last rune from the string (used for iterative byte-limit enforcement).
func shortenByOneRune(s string) string {
	r := []rune(s)
	if len(r) <= 1 {
		return ""
	}
	return string(r[:len(r)-1])
}
