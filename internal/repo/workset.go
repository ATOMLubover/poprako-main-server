package repo

import (
	"errors"
	"fmt"

	"poprako-main-server/internal/model/po"
)

// WorksetRepo defines repository operations for worksets.
type WorksetRepo interface {
	Repo

	GetWorksetByID(ex Executor, worksetID string) (*po.BasicWorkset, error)
	GetWorksetsByIDs(ex Executor, worksetIDs []string) ([]po.BasicWorkset, error)

	CreateWorkset(ex Executor, newWorkset *po.NewWorkset) error

	UpdateWorksetByID(ex Executor, patchWorkset *po.PatchWorkset) error
}

type worksetRepo struct {
	ex Executor
}

func NewWorksetRepo(ex Executor) WorksetRepo {
	return &worksetRepo{ex: ex}
}

func (wr *worksetRepo) Exec() Executor { return wr.ex }

func (wr *worksetRepo) withTrx(tx Executor) Executor {
	if tx != nil {
		return tx
	}

	return wr.ex
}

func (wr *worksetRepo) CreateWorkset(ex Executor, newWorkset *po.NewWorkset) error {
	ex = wr.withTrx(ex)

	return ex.Create(newWorkset).Error
}

func (wr *worksetRepo) GetWorksetByID(ex Executor, worksetID string) (*po.BasicWorkset, error) {
	ex = wr.withTrx(ex)

	w := &po.BasicWorkset{}

	if err := ex.
		Where("id = ?", worksetID).
		First(w).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get workset by ID: %w", err)
	}

	return w, nil
}

func (wr *worksetRepo) GetWorksetsByIDs(ex Executor, worksetIDs []string) ([]po.BasicWorkset, error) {
	ex = wr.withTrx(ex)

	var lst []po.BasicWorkset

	if err := ex.
		Where("id IN ?", worksetIDs).
		Find(&lst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get worksets by IDs: %w", err)
	}

	return lst, nil
}

func (wr *worksetRepo) UpdateWorksetByID(ex Executor, patchWorkset *po.PatchWorkset) error {
	if patchWorkset.ID == "" {
		return errors.New("workset ID is required for update")
	}

	ex = wr.withTrx(ex)

	return ex.Save(patchWorkset).Error
}
