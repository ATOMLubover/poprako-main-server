package repo

import (
	"errors"
	"fmt"

	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"

	"gorm.io/gorm"
)

// ComicRepo defines repository operations for comics.
type ComicRepo interface {
	Repo

	GetComicByID(ex Executor, comicID string) (*po.BasicComic, error)
	GetComicsByWorksetID(ex Executor, worksetID string, offset, limit int) ([]po.BriefComic, error)
	RetrieveComics(ex Executor, opt model.RetrieveComicOpt) ([]po.BriefComic, error)

	CreateComic(newComic *po.NewComic) error

	UpdateComicByID(ex Executor, patchComic *po.PatchComic) error

	DeleteComicByID(ex Executor, comicID string) error
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

func (cr *comicRepo) CreateComic(newComic *po.NewComic) error {
	if err := cr.Exec().Transaction(func(ex Executor) error {
		// Query workset index
		var workset po.DetailedWorkset
		if err := ex.Model(&po.DetailedWorkset{}).
			Select("index").
			Where("id = ?", newComic.WorksetID).
			First(&workset).
			Error; err != nil {
			return fmt.Errorf("Failed to get workset index: %w", err)
		}
		newComic.WorksetIndex = int(workset.Index)

		// Count total comic count in the workset.
		// A optimistic lock based on unqiue index is expected.
		var cnt int64

		if err := ex.Model(&po.BasicComic{}).
			Where("workset_id = ?", newComic.WorksetID).
			Count(&cnt).
			Error; err != nil {
			return fmt.Errorf("Failed to count comics in workset: %w", err)
		}

		newComic.Index = cnt + 1

		// Create the comic
		if err := ex.Create(newComic).Error; err != nil {
			return err
		}

		// Update workset comic_count
		if err := ex.Model(&po.DetailedWorkset{}).
			Where("id = ?", newComic.WorksetID).
			UpdateColumn("comic_count", gorm.Expr("comic_count + ?", 1)).
			Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (cr *comicRepo) GetComicByID(ex Executor, comicID string) (*po.BasicComic, error) {
	ex = cr.withTrx(ex)

	c := &po.BasicComic{}

	if err := ex.
		Model(&po.BasicComic{}).
		Select(`comic_tbl.*, 
			user_tbl.nickname AS creator_nickname`).
		Joins("LEFT JOIN user_tbl ON comic_tbl.creator_id = user_tbl.id").
		Where("comic_tbl.id = ?", comicID).
		First(c).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get comic by ID: %w", err)
	}

	return c, nil
}

func (cr *comicRepo) GetComicsByWorksetID(ex Executor, worksetID string, offset, limit int) ([]po.BriefComic, error) {
	ex = cr.withTrx(ex)

	var lst []po.BriefComic

	q := ex.Model(&po.BriefComic{}).
		Where("comic_tbl.workset_id = ?", worksetID)

	if offset > 0 {
		q = q.Offset(offset)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}

	if err := q.
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

	updates := map[string]any{}

	if patchComic.Author != nil {
		updates["author"] = *patchComic.Author
	}
	if patchComic.Title != nil {
		updates["title"] = *patchComic.Title
	}
	if patchComic.Comment != nil {
		updates["comment"] = *patchComic.Comment
	}
	if patchComic.Description != nil {
		updates["description"] = *patchComic.Description
	}

	if patchComic.TranslatingStartedAt != nil {
		if patchComic.TranslatingStartedAt.IsZero() {
			updates["translating_started_at"] = nil
		} else {
			updates["translating_started_at"] = *patchComic.TranslatingStartedAt
		}
	}
	if patchComic.TranslatingCompletedAt != nil {
		if patchComic.TranslatingCompletedAt.IsZero() {
			updates["translating_completed_at"] = nil
		} else {
			updates["translating_completed_at"] = *patchComic.TranslatingCompletedAt
		}
	}

	if patchComic.ProofreadingStartedAt != nil {
		if patchComic.ProofreadingStartedAt.IsZero() {
			updates["proofreading_started_at"] = nil
		} else {
			updates["proofreading_started_at"] = *patchComic.ProofreadingStartedAt
		}
	}
	if patchComic.ProofreadingCompletedAt != nil {
		if patchComic.ProofreadingCompletedAt.IsZero() {
			updates["proofreading_completed_at"] = nil
		} else {
			updates["proofreading_completed_at"] = *patchComic.ProofreadingCompletedAt
		}
	}

	if patchComic.TypesettingStartedAt != nil {
		if patchComic.TypesettingStartedAt.IsZero() {
			updates["typesetting_started_at"] = nil
		} else {
			updates["typesetting_started_at"] = *patchComic.TypesettingStartedAt
		}
	}
	if patchComic.TypesettingCompletedAt != nil {
		if patchComic.TypesettingCompletedAt.IsZero() {
			updates["typesetting_completed_at"] = nil
		} else {
			updates["typesetting_completed_at"] = *patchComic.TypesettingCompletedAt
		}
	}

	if patchComic.ReviewingCompletedAt != nil {
		if patchComic.ReviewingCompletedAt.IsZero() {
			updates["reviewing_completed_at"] = nil
		} else {
			updates["reviewing_completed_at"] = *patchComic.ReviewingCompletedAt
		}
	}

	if patchComic.UploadingCompletedAt != nil {
		if patchComic.UploadingCompletedAt.IsZero() {
			updates["uploading_completed_at"] = nil
		} else {
			updates["uploading_completed_at"] = *patchComic.UploadingCompletedAt
		}
	}

	if len(updates) == 0 {
		return nil
	}

	return ex.Model(&po.PatchComic{}).
		Where("id = ?", patchComic.ID).
		Updates(updates).
		Error
}

// RetrieveComics returns a slice of BriefComic with filtering and pagination.
// A zero-length slice is returned if no comics are found.
func (cr *comicRepo) RetrieveComics(ex Executor, opt model.RetrieveComicOpt) ([]po.BriefComic, error) {
	ex = cr.withTrx(ex)

	var lst []po.BriefComic

	query := ex.Model(&po.BriefComic{})

	if opt.Author != nil {
		query = query.Where("comic_tbl.author LIKE ?", "%"+*opt.Author+"%")
	}

	if opt.Title != nil {
		query = query.Where("comic_tbl.title LIKE ?", "%"+*opt.Title+"%")
	}

	if opt.WorksetIndex != nil {
		query = query.Where("comic_tbl.workset_index = ?", *opt.WorksetIndex)
	}

	if opt.Index != nil {
		query = query.Where("comic_tbl.index = ?", *opt.Index)
	}

	if opt.TranslatingNotStarted != nil && *opt.TranslatingNotStarted {
		query = query.Where("comic_tbl.translating_started_at IS NULL")
	}
	if opt.TranslatingInProgress != nil && *opt.TranslatingInProgress {
		query = query.Where("comic_tbl.translating_started_at IS NOT NULL AND comic_tbl.translating_completed_at IS NULL")
	}
	if opt.TranslatingCompleted != nil && *opt.TranslatingCompleted {
		query = query.Where("comic_tbl.translating_completed_at IS NOT NULL")
	}

	if opt.ProofreadingNotStarted != nil && *opt.ProofreadingNotStarted {
		query = query.Where("comic_tbl.proofreading_started_at IS NULL")
	}
	if opt.ProofreadingInProgress != nil && *opt.ProofreadingInProgress {
		query = query.Where("comic_tbl.proofreading_started_at IS NOT NULL AND comic_tbl.proofreading_completed_at IS NULL")
	}
	if opt.ProofreadingCompleted != nil && *opt.ProofreadingCompleted {
		query = query.Where("comic_tbl.proofreading_completed_at IS NOT NULL")
	}

	if opt.TypesettingNotStarted != nil && *opt.TypesettingNotStarted {
		query = query.Where("comic_tbl.typesetting_started_at IS NULL")
	}
	if opt.TypesettingInProgress != nil && *opt.TypesettingInProgress {
		query = query.Where("comic_tbl.typesetting_started_at IS NOT NULL AND comic_tbl.typesetting_completed_at IS NULL")
	}
	if opt.TypesettingCompleted != nil && *opt.TypesettingCompleted {
		query = query.Where("comic_tbl.typesetting_completed_at IS NOT NULL")
	}

	if opt.ReviewingNotStarted != nil && *opt.ReviewingNotStarted {
		query = query.Where("comic_tbl.reviewing_completed_at IS NULL")
	}
	if opt.ReviewingCompleted != nil && *opt.ReviewingCompleted {
		query = query.Where("comic_tbl.reviewing_completed_at IS NOT NULL")
	}

	if opt.UploadingNotStarted != nil && *opt.UploadingNotStarted {
		query = query.Where("comic_tbl.uploading_completed_at IS NULL")
	}
	if opt.UploadingCompleted != nil && *opt.UploadingCompleted {
		query = query.Where("comic_tbl.uploading_completed_at IS NOT NULL")
	}

	if opt.AssignedUserID != nil {
		query = query.Where(`EXISTS (
			SELECT 1 FROM comic_assignment_tbl 
			WHERE comic_assignment_tbl.comic_id = comic_tbl.id 
			AND comic_assignment_tbl.user_id = ?
		)`, *opt.AssignedUserID)
	}

	if opt.Offset > 0 {
		query = query.Offset(opt.Offset)
	}

	if opt.Limit > 0 {
		query = query.Limit(opt.Limit)
	}

	if err := query.
		Order("comic_tbl.updated_at DESC").
		Find(&lst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to retrieve comics: %w", err)
	}

	return lst, nil
}

func (cr *comicRepo) DeleteComicByID(ex Executor, comicID string) error {
	return cr.Exec().Transaction(func(tx Executor) error {
		// Get comic first to get workset_id
		comic := &po.BasicComic{}
		if err := tx.Where("id = ?", comicID).First(comic).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return REC_NOT_FOUND
			}
			return fmt.Errorf("Failed to get comic for deletion: %w", err)
		}

		// Delete the comic
		if err := tx.Where("id = ?", comicID).Delete(&po.BasicComic{}).Error; err != nil {
			return fmt.Errorf("Failed to delete comic: %w", err)
		}

		// Update workset comic_count
		if err := tx.Model(&po.DetailedWorkset{}).
			Where("id = ?", comic.WorksetID).
			UpdateColumn("comic_count", gorm.Expr("comic_count - ?", 1)).
			Error; err != nil {
			return fmt.Errorf("Failed to update workset comic_count: %w", err)
		}

		return nil
	})
}
