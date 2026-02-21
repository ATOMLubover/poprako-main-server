package repo

import (
	"fmt"

	"poprako-main-server/internal/model/po"
)

// ComicUnitRepo defines repository operations for comic units.
type ComicUnitRepo interface {
	Repo

	GetUnitsByPageID(ex Exct, pageID string) ([]po.BasicComicUnit, error)

	GetUnitCountsByPageID(ex Exct, pageID string) (po.UnitCounts, error)

	GetUnitCountsByPageIDs(ex Exct, pageIDs []string) (map[string]po.UnitCounts, error)

	CreateUnits(ex Exct, newUnits []po.NewComicUnit) error

	UpdateUnitsByIDs(ex Exct, patchUnits []po.PatchComicUnit) error

	DeleteUnitByIDs(ex Exct, unitIDs []string) error
}

type comicUnitRepo struct {
	ex Exct
}

func NewComicUnitRepo(ex Exct) ComicUnitRepo {
	return &comicUnitRepo{ex: ex}
}

func (cur *comicUnitRepo) Exct() Exct { return cur.ex }

func (cur *comicUnitRepo) withTrx(tx Exct) Exct {
	if tx != nil {
		return tx
	}

	return cur.ex
}

func (cur *comicUnitRepo) CreateUnits(ex Exct, newUnits []po.NewComicUnit) error {
	if len(newUnits) == 0 {
		return nil
	}

	ex = cur.withTrx(ex)

	return ex.Create(newUnits).Error
}

func (cur *comicUnitRepo) GetUnitsByPageID(ex Exct, pageID string) ([]po.BasicComicUnit, error) {
	ex = cur.withTrx(ex)

	var lst []po.BasicComicUnit

	if err := ex.
		Where("page_id = ?", pageID).
		Find(&lst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get units by page ID: %w", err)
	}

	return lst, nil
}

func (cur *comicUnitRepo) UpdateUnitsByIDs(ex Exct, patchUnits []po.PatchComicUnit) error {
	if len(patchUnits) == 0 {
		return nil
	}

	ex = cur.withTrx(ex)

	// Update each unit individually to handle optional fields correctly
	for _, patchUnit := range patchUnits {
		if patchUnit.ID == "" {
			continue
		}

		updates := map[string]any{}

		if patchUnit.Index != nil {
			updates["index"] = *patchUnit.Index
		}
		if patchUnit.XCoordinate != nil {
			updates["x_coordinate"] = *patchUnit.XCoordinate
		}
		if patchUnit.YCoordinate != nil {
			updates["y_coordinate"] = *patchUnit.YCoordinate
		}
		if patchUnit.IsInBox != nil {
			updates["is_in_box"] = *patchUnit.IsInBox
		}
		if patchUnit.TranslatedText != nil {
			updates["translated_text"] = *patchUnit.TranslatedText
		}
		if patchUnit.TranslatorID != nil {
			updates["translator_id"] = *patchUnit.TranslatorID
		}
		if patchUnit.TranslatorComment != nil {
			updates["translator_comment"] = *patchUnit.TranslatorComment
		}
		if patchUnit.ProvedText != nil {
			updates["proved_text"] = *patchUnit.ProvedText
		}
		if patchUnit.Proved != nil {
			updates["proved"] = *patchUnit.Proved
		}
		if patchUnit.ProofreaderID != nil {
			updates["proofreader_id"] = *patchUnit.ProofreaderID
		}
		if patchUnit.ProofreaderComment != nil {
			updates["proofreader_comment"] = *patchUnit.ProofreaderComment
		}
		if patchUnit.CreatorID != nil {
			updates["creator_id"] = *patchUnit.CreatorID
		}

		if len(updates) == 0 {
			continue
		}

		if err := ex.Model(&po.PatchComicUnit{}).
			Where("id = ?", patchUnit.ID).
			Updates(updates).
			Error; err != nil {
			return err
		}
	}

	return nil
}

func (cur *comicUnitRepo) DeleteUnitByIDs(ex Exct, unitIDs []string) error {
	if len(unitIDs) == 0 {
		return nil
	}

	ex = cur.withTrx(ex)

	return ex.Where("id IN ?", unitIDs).Delete(&po.BasicComicUnit{}).Error
}

func (cur *comicUnitRepo) GetUnitCountsByPageID(ex Exct, pageID string) (po.UnitCounts, error) {
	ex = cur.withTrx(ex)

	var result struct {
		Inbox      int64
		Outbox     int64
		Translated int64
		Proved     int64
	}

	err := ex.Table("comic_unit_tbl").
		Select(`
			SUM(CASE WHEN is_in_box = true THEN 1 ELSE 0 END) AS inbox,
			SUM(CASE WHEN is_in_box = false THEN 1 ELSE 0 END) AS outbox,
			SUM(CASE WHEN translated_text IS NOT NULL AND translated_text != '' THEN 1 ELSE 0 END) AS translated,
			SUM(CASE WHEN proved = true THEN 1 ELSE 0 END) AS proved
		`).
		Where("page_id = ?", pageID).
		Scan(&result).
		Error
	if err != nil {
		return po.UnitCounts{}, fmt.Errorf("Failed to get unit counts by page ID: %w", err)
	}

	return po.UnitCounts{
		Inbox:      result.Inbox,
		Outbox:     result.Outbox,
		Translated: result.Translated,
		Proved:     result.Proved,
	}, nil
}

func (cur *comicUnitRepo) GetUnitCountsByPageIDs(ex Exct, pageIDs []string) (map[string]po.UnitCounts, error) {
	if len(pageIDs) == 0 {
		return make(map[string]po.UnitCounts), nil
	}

	ex = cur.withTrx(ex)

	var results []struct {
		PageID     string
		Inbox      int64
		Outbox     int64
		Translated int64
		Proved     int64
	}

	err := ex.Table("comic_unit_tbl").
		Select(`
			page_id,
			SUM(CASE WHEN is_in_box = true THEN 1 ELSE 0 END) AS inbox,
			SUM(CASE WHEN is_in_box = false THEN 1 ELSE 0 END) AS outbox,
			SUM(CASE WHEN translated_text IS NOT NULL AND translated_text != '' THEN 1 ELSE 0 END) AS translated,
			SUM(CASE WHEN proved = true THEN 1 ELSE 0 END) AS proved
		`).
		Where("page_id IN ?", pageIDs).
		Group("page_id").
		Scan(&results).
		Error
	if err != nil {
		return nil, fmt.Errorf("Failed to get unit counts by page IDs: %w", err)
	}

	countsMap := make(map[string]po.UnitCounts, len(results))
	for _, r := range results {
		countsMap[r.PageID] = po.UnitCounts{
			Inbox:      r.Inbox,
			Outbox:     r.Outbox,
			Translated: r.Translated,
			Proved:     r.Proved,
		}
	}

	return countsMap, nil
}
