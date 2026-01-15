package repo

import (
	"errors"
	"fmt"

	"poprako-main-server/internal/model/po"
)

// WorksetRepo defines repository operations for worksets.
type WorksetRepo interface {
	Repo

	GetWorksetByID(ex Executor, worksetID string) (*po.DetailedWorkset, error)
	RetrieveWorksets(ex Executor, limit, offset int) ([]po.DetailedWorkset, error)

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
	// Count total workset count.
	// A optimistic lock based on unqiue index is expected.
	var cnt int64

	if err := wr.Exec().Model(&po.DetailedWorkset{}).
		Count(&cnt).
		Error; err != nil {
		return fmt.Errorf("Failed to count worksets: %w", err)
	}

	newWorkset.Index = cnt + 1

	ex = wr.withTrx(ex)

	return ex.Create(newWorkset).Error
}

func (wr *worksetRepo) GetWorksetByID(ex Executor, worksetID string) (*po.DetailedWorkset, error) {
	ex = wr.withTrx(ex)

	w := &po.DetailedWorkset{}

	if err := ex.
		Model(&po.DetailedWorkset{}).
		Select("workset_tbl.*, workset_tbl.name AS name, user_tbl.nickname AS creator_nickname").
		Joins("LEFT JOIN user_tbl ON workset_tbl.creator_id = user_tbl.id").
		Where("workset_tbl.id = ?", worksetID).
		First(w).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get workset by ID: %w", err)
	}

	return w, nil
}

func (wr *worksetRepo) UpdateWorksetByID(ex Executor, patchWorkset *po.PatchWorkset) error {
	if patchWorkset.ID == "" {
		return errors.New("workset ID is required for update")
	}

	ex = wr.withTrx(ex)

	updates := map[string]any{}

	if patchWorkset.Name != nil {
		updates["name"] = *patchWorkset.Name
	}
	if patchWorkset.Index != nil {
		updates["index"] = *patchWorkset.Index
	}
	if patchWorkset.ComicCount != nil {
		updates["comic_count"] = *patchWorkset.ComicCount
	}
	if patchWorkset.Description != nil {
		updates["description"] = *patchWorkset.Description
	}
	if patchWorkset.CreatorID != nil {
		updates["creator_id"] = *patchWorkset.CreatorID
	}

	if len(updates) == 0 {
		return nil
	}

	return ex.Model(&po.PatchWorkset{}).
		Where("id = ?", patchWorkset.ID).
		Updates(updates).
		Error
}

// RetrieveWorksets returns a list of DetailedWorkset with pagination (limit, offset).
// A zero-length slice is returned if no worksets are found.
func (wr *worksetRepo) RetrieveWorksets(ex Executor, limit, offset int) ([]po.DetailedWorkset, error) {
	ex = wr.withTrx(ex)

	var lst []po.DetailedWorkset

	q := ex.Model(&po.DetailedWorkset{}).
		Select("workset_tbl.*, workset_tbl.name AS name, user_tbl.nickname AS creator_nickname").
		Joins("LEFT JOIN user_tbl ON workset_tbl.creator_id = user_tbl.id").
		Order("workset_tbl.updated_at DESC")

	if limit > 0 {
		q = q.Limit(limit)
	}

	if offset > 0 {
		q = q.Offset(offset)
	}

	if err := q.
		Find(&lst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to retrieve worksets: %w", err)
	}

	return lst, nil
}
