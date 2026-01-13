package repo

import (
	"errors"
	"fmt"

	"poprako-main-server/internal/model/po"
)

// ComicRepo defines repository operations for comics.
type ComicRepo interface {
	Repo

	GetComicByID(ex Executor, comicID string) (*po.BasicComic, error)
	GetComicsByWorksetID(ex Executor, worksetID string) ([]po.BasicComic, error)
	GetComicsByIDs(ex Executor, comicIDs []string) ([]po.BasicComic, error)

	CreateComic(ex Executor, newComic *po.NewComic) error
	UpdateComicByID(ex Executor, patchComic *po.PatchComic) error
}

type comicRepo struct {
	ex Executor
}

func NewComicRepo(ex Executor) ComicRepo {
	return &comicRepo{ex: ex}
}

func (cr *comicRepo) Exec() Executor { return cr.ex }

func (cr *comicRepo) withTrx(tx Executor) Executor {
	if tx != nil {
		return tx
	}

	return cr.ex
}

func (cr *comicRepo) CreateComic(ex Executor, newComic *po.NewComic) error {
	ex = cr.withTrx(ex)

	return ex.Create(newComic).Error
}

func (cr *comicRepo) GetComicByID(ex Executor, comicID string) (*po.BasicComic, error) {
	ex = cr.withTrx(ex)

	c := &po.BasicComic{}

	if err := ex.
		Where("id = ?", comicID).
		First(c).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get comic by ID: %w", err)
	}

	return c, nil
}

func (cr *comicRepo) GetComicsByWorksetID(ex Executor, worksetID string) ([]po.BasicComic, error) {
	ex = cr.withTrx(ex)

	var lst []po.BasicComic

	if err := ex.
		Where("workset_id = ?", worksetID).
		Find(&lst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get comics by workset ID: %w", err)
	}

	return lst, nil
}

func (cr *comicRepo) GetComicsByIDs(ex Executor, comicIDs []string) ([]po.BasicComic, error) {
	ex = cr.withTrx(ex)

	var lst []po.BasicComic

	if err := ex.
		Where("id IN ?", comicIDs).
		Find(&lst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get comics by IDs: %w", err)
	}

	return lst, nil
}

func (cr *comicRepo) UpdateComicByID(ex Executor, patchComic *po.PatchComic) error {
	if patchComic.ID == "" {
		return errors.New("comic ID is required for update")
	}

	ex = cr.withTrx(ex)

	return ex.Save(patchComic).Error
}
