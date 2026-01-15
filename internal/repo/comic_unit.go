package repo

import (
	"fmt"

	"poprako-main-server/internal/model/po"
)

// ComicUnitRepo defines repository operations for comic units.
type ComicUnitRepo interface {
	Repo

	GetUnitsByPageID(ex Executor, pageID string) ([]po.BasicComicUnit, error)

	CreateUnits(ex Executor, newUnits []po.NewComicUnit) error

	UpdateUnitsByIDs(ex Executor, patchUnits []po.PatchComicUnit) error

	DeleteUnitByIDs(ex Executor, unitIDs []string) error
}

type comicUnitRepo struct {
	ex Executor
}

func NewComicUnitRepo(ex Executor) ComicUnitRepo {
	return &comicUnitRepo{ex: ex}
}

func (cur *comicUnitRepo) Exec() Executor { return cur.ex }

func (cur *comicUnitRepo) withTrx(tx Executor) Executor {
	if tx != nil {
		return tx
	}

	return cur.ex
}

func (cur *comicUnitRepo) CreateUnits(ex Executor, newUnits []po.NewComicUnit) error {
	if len(newUnits) == 0 {
		return nil
	}

	ex = cur.withTrx(ex)

	return ex.Create(newUnits).Error
}

func (cur *comicUnitRepo) GetUnitsByPageID(ex Executor, pageID string) ([]po.BasicComicUnit, error) {
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

func (cur *comicUnitRepo) UpdateUnitsByIDs(ex Executor, patchUnits []po.PatchComicUnit) error {
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

func (cur *comicUnitRepo) DeleteUnitByIDs(ex Executor, unitIDs []string) error {
	if len(unitIDs) == 0 {
		return nil
	}

	ex = cur.withTrx(ex)

	return ex.Where("id IN ?", unitIDs).Delete(&po.BasicComicUnit{}).Error
}
