package repo

import (
	"errors"
	"fmt"

	"poprako-main-server/internal/model/po"
)

// ComicUnitRepo defines repository operations for comic units.
type ComicUnitRepo interface {
	Repo

	GetUnitByID(ex Executor, unitID string) (*po.BasicComicUnit, error)
	GetUnitsByPageID(ex Executor, pageID string) ([]po.BasicComicUnit, error)
	GetUnitsByIDs(ex Executor, unitIDs []string) ([]po.BasicComicUnit, error)

	CreateUnit(ex Executor, newUnit *po.NewComicUnit) error

	UpdateUnitByID(ex Executor, patchUnit *po.PatchComicUnit) error
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

func (cur *comicUnitRepo) CreateUnit(ex Executor, newUnit *po.NewComicUnit) error {
	ex = cur.withTrx(ex)

	return ex.Create(newUnit).Error
}

func (cur *comicUnitRepo) GetUnitByID(ex Executor, unitID string) (*po.BasicComicUnit, error) {
	ex = cur.withTrx(ex)

	u := &po.BasicComicUnit{}

	if err := ex.
		Where("id = ?", unitID).
		First(u).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get unit by ID: %w", err)
	}

	return u, nil
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

func (cur *comicUnitRepo) GetUnitsByIDs(ex Executor, unitIDs []string) ([]po.BasicComicUnit, error) {
	ex = cur.withTrx(ex)

	var lst []po.BasicComicUnit

	if err := ex.
		Where("id IN ?", unitIDs).
		Find(&lst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get units by IDs: %w", err)
	}

	return lst, nil
}

func (cur *comicUnitRepo) UpdateUnitByID(ex Executor, patchUnit *po.PatchComicUnit) error {
	if patchUnit.ID == "" {
		return errors.New("unit ID is required for update")
	}

	ex = cur.withTrx(ex)

	return ex.Save(patchUnit).Error
}
