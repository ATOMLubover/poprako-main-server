package repo

import (
	"errors"
	"fmt"

	"poprako-main-server/internal/model/po"
)

// ComicAssignmentRepo defines repository operations for comic assignments.
type ComicAssignmentRepo interface {
	Repo

	GetAssignmentByID(ex Executor, assignmentID string) (*po.BasicComicAssignment, error)
	GetAssignmentsByComicID(ex Executor, comicID string) ([]po.BasicComicAssignment, error)
	GetAssignmentsByUserID(ex Executor, userID string) ([]po.BasicComicAssignment, error)

	CreateAssignment(ex Executor, newAssign *po.NewComicAssignment) error

	UpdateAssignmentByID(ex Executor, patchAssign *po.PatchComicAssignment) error
}

type comicAssignmentRepo struct {
	ex Executor
}

func NewComicAssignmentRepo(ex Executor) ComicAssignmentRepo {
	return &comicAssignmentRepo{ex: ex}
}

func (car *comicAssignmentRepo) Exec() Executor { return car.ex }

func (car *comicAssignmentRepo) withTrx(tx Executor) Executor {
	if tx != nil {
		return tx
	}

	return car.ex
}

func (car *comicAssignmentRepo) CreateAssignment(ex Executor, newAssign *po.NewComicAssignment) error {
	ex = car.withTrx(ex)

	return ex.Create(newAssign).Error
}

func (car *comicAssignmentRepo) GetAssignmentByID(ex Executor, assignmentID string) (*po.BasicComicAssignment, error) {
	ex = car.withTrx(ex)

	a := &po.BasicComicAssignment{}

	if err := ex.
		Where("id = ?", assignmentID).
		First(a).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get assignment by ID: %w", err)
	}

	return a, nil
}

func (car *comicAssignmentRepo) GetAssignmentsByComicID(ex Executor, comicID string) ([]po.BasicComicAssignment, error) {
	ex = car.withTrx(ex)

	var lst []po.BasicComicAssignment

	if err := ex.
		Where("comic_id = ?", comicID).
		Find(&lst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get assignments by comic ID: %w", err)
	}

	return lst, nil
}

func (car *comicAssignmentRepo) GetAssignmentsByUserID(ex Executor, userID string) ([]po.BasicComicAssignment, error) {
	ex = car.withTrx(ex)

	var lst []po.BasicComicAssignment

	if err := ex.
		Where("user_id = ?", userID).
		Find(&lst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get assignments by user ID: %w", err)
	}

	return lst, nil
}

func (car *comicAssignmentRepo) UpdateAssignmentByID(ex Executor, patchAssign *po.PatchComicAssignment) error {
	if patchAssign.ID == "" {
		return errors.New("assignment ID is required for update")
	}

	ex = car.withTrx(ex)

	updates := map[string]interface{}{}

	if patchAssign.ComicID != nil {
		updates["comic_id"] = *patchAssign.ComicID
	}
	if patchAssign.UserID != nil {
		updates["user_id"] = *patchAssign.UserID
	}

	if patchAssign.AssignedTranslatorAt != nil {
		if *patchAssign.AssignedTranslatorAt == 0 {
			updates["assigned_translator_at"] = nil
		} else {
			updates["assigned_translator_at"] = *patchAssign.AssignedTranslatorAt
		}
	}

	if patchAssign.AssignedProofreaderAt != nil {
		if *patchAssign.AssignedProofreaderAt == 0 {
			updates["assigned_proofreader_at"] = nil
		} else {
			updates["assigned_proofreader_at"] = *patchAssign.AssignedProofreaderAt
		}
	}

	if patchAssign.AssignedTypesetterAt != nil {
		if *patchAssign.AssignedTypesetterAt == 0 {
			updates["assigned_typesetter_at"] = nil
		} else {
			updates["assigned_typesetter_at"] = *patchAssign.AssignedTypesetterAt
		}
	}

	if patchAssign.AssignedRedrawerAt != nil {
		if *patchAssign.AssignedRedrawerAt == 0 {
			updates["assigned_redrawer_at"] = nil
		} else {
			updates["assigned_redrawer_at"] = *patchAssign.AssignedRedrawerAt
		}
	}

	if patchAssign.AssignedReviewerAt != nil {
		if *patchAssign.AssignedReviewerAt == 0 {
			updates["assigned_reviewer_at"] = nil
		} else {
			updates["assigned_reviewer_at"] = *patchAssign.AssignedReviewerAt
		}
	}

	if len(updates) == 0 {
		return nil
	}

	return ex.Model(&po.PatchComicAssignment{}).
		Where("id = ?", patchAssign.ID).
		Updates(updates).
		Error
}
