package repo

import (
	"errors"
	"fmt"

	"poprako-main-server/internal/model/po"
)

// ComicAsgnRepo defines repository operations for comic assignments.
type ComicAsgnRepo interface {
	Repo

	GetAsgnByID(ex Exct, assignmentID string) (*po.BasicComicAsgn, error)
	GetAsgnsByComicID(ex Exct, comicID string, offset, limit int) ([]po.BasicComicAsgn, error)
	GetAsgnsByUserID(ex Exct, userID string, offset, limit int) ([]po.BasicComicAsgn, error)
	GetAsgnsByUserAndComicID(ex Exct, userID, comicID string) (*po.BasicComicAsgn, error)

	CreateAsgn(ex Exct, newAssign *po.NewComicAsgn) error

	UpdateAsgnByID(ex Exct, patchAssign *po.PatchComicAsgn) error

	DeleteAsgnByID(ex Exct, assignmentID string) error
}

type comicAsgnRepo struct {
	ex Exct
}

func NewComicAsgnRepo(ex Exct) ComicAsgnRepo {
	return &comicAsgnRepo{ex: ex}
}

func (car *comicAsgnRepo) Exct() Exct { return car.ex }

func (car *comicAsgnRepo) withTrx(tx Exct) Exct {
	if tx != nil {
		return tx
	}

	return car.ex
}

func (car *comicAsgnRepo) CreateAsgn(ex Exct, newAssign *po.NewComicAsgn) error {
	ex = car.withTrx(ex)

	return ex.Create(newAssign).Error
}

func (car *comicAsgnRepo) GetAsgnByID(ex Exct, assignmentID string) (*po.BasicComicAsgn, error) {
	ex = car.withTrx(ex)

	a := &po.BasicComicAsgn{}

	if err := ex.
		Table(po.COMIC_ASSIGNMENT_TABLE).
		Select(po.COMIC_ASSIGNMENT_TABLE+".*, "+po.USER_TABLE+".nickname AS user_nickname").
		Joins("LEFT JOIN "+po.USER_TABLE+" ON "+po.COMIC_ASSIGNMENT_TABLE+".user_id = "+po.USER_TABLE+".id").
		Where(po.COMIC_ASSIGNMENT_TABLE+".id = ?", assignmentID).
		First(a).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get assignment by ID: %w", err)
	}

	return a, nil
}

func (car *comicAsgnRepo) GetAsgnsByComicID(ex Exct, comicID string, offset, limit int) ([]po.BasicComicAsgn, error) {
	ex = car.withTrx(ex)

	var lst []po.BasicComicAsgn

	query := ex.
		Table(po.COMIC_ASSIGNMENT_TABLE).
		Select(po.COMIC_ASSIGNMENT_TABLE+".*, "+po.USER_TABLE+".nickname AS user_nickname").
		Joins("LEFT JOIN "+po.USER_TABLE+" ON "+po.COMIC_ASSIGNMENT_TABLE+".user_id = "+po.USER_TABLE+".id").
		Where(po.COMIC_ASSIGNMENT_TABLE+".comic_id = ?", comicID)

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.
		Find(&lst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get assignments by comic ID: %w", err)
	}

	return lst, nil
}

func (car *comicAsgnRepo) GetAsgnsByUserID(ex Exct, userID string, offset, limit int) ([]po.BasicComicAsgn, error) {
	ex = car.withTrx(ex)

	var lst []po.BasicComicAsgn

	q := ex.
		Table(po.COMIC_ASSIGNMENT_TABLE).
		Select(po.COMIC_ASSIGNMENT_TABLE+".*, "+po.USER_TABLE+".nickname AS user_nickname").
		Joins("LEFT JOIN "+po.USER_TABLE+" ON "+po.COMIC_ASSIGNMENT_TABLE+".user_id = "+po.USER_TABLE+".id").
		Where(po.COMIC_ASSIGNMENT_TABLE+".user_id = ?", userID)

	if offset > 0 {
		q = q.Offset(offset)
	}

	if limit > 0 {
		q = q.Limit(limit)
	}

	if err := q.
		Find(&lst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get assignments by user ID: %w", err)
	}

	return lst, nil
}

func (car *comicAsgnRepo) GetAsgnsByUserAndComicID(ex Exct, userID, comicID string) (*po.BasicComicAsgn, error) {
	ex = car.withTrx(ex)

	a := &po.BasicComicAsgn{}

	if err := ex.
		Table(po.COMIC_ASSIGNMENT_TABLE).
		Select(po.COMIC_ASSIGNMENT_TABLE+".*, "+po.USER_TABLE+".nickname AS user_nickname").
		Joins("LEFT JOIN "+po.USER_TABLE+" ON "+po.COMIC_ASSIGNMENT_TABLE+".user_id = "+po.USER_TABLE+".id").
		Where(po.COMIC_ASSIGNMENT_TABLE+".user_id = ? AND "+po.COMIC_ASSIGNMENT_TABLE+".comic_id = ?", userID, comicID).
		First(a).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get assignment by user ID and comic ID: %w", err)
	}

	return a, nil
}

func (car *comicAsgnRepo) UpdateAsgnByID(ex Exct, patchAssign *po.PatchComicAsgn) error {
	if patchAssign.ID == "" {
		return errors.New("assignment ID is required for update")
	}

	ex = car.withTrx(ex)

	updates := map[string]any{}

	if patchAssign.ComicID != nil {
		updates["comic_id"] = *patchAssign.ComicID
	}
	if patchAssign.UserID != nil {
		updates["user_id"] = *patchAssign.UserID
	}

	if patchAssign.AssignedTranslatorAt != nil {
		if patchAssign.AssignedTranslatorAt.IsZero() {
			updates["assigned_translator_at"] = nil
		} else {
			updates["assigned_translator_at"] = *patchAssign.AssignedTranslatorAt
		}
	}

	if patchAssign.AssignedProofreaderAt != nil {
		if patchAssign.AssignedProofreaderAt.IsZero() {
			updates["assigned_proofreader_at"] = nil
		} else {
			updates["assigned_proofreader_at"] = *patchAssign.AssignedProofreaderAt
		}
	}

	if patchAssign.AssignedTypesetterAt != nil {
		if patchAssign.AssignedTypesetterAt.IsZero() {
			updates["assigned_typesetter_at"] = nil
		} else {
			updates["assigned_typesetter_at"] = *patchAssign.AssignedTypesetterAt
		}
	}

	if patchAssign.AssignedRedrawerAt != nil {
		if patchAssign.AssignedRedrawerAt.IsZero() {
			updates["assigned_redrawer_at"] = nil
		} else {
			updates["assigned_redrawer_at"] = *patchAssign.AssignedRedrawerAt
		}
	}

	if patchAssign.AssignedReviewerAt != nil {
		if patchAssign.AssignedReviewerAt.IsZero() {
			updates["assigned_reviewer_at"] = nil
		} else {
			updates["assigned_reviewer_at"] = *patchAssign.AssignedReviewerAt
		}
	}

	if len(updates) == 0 {
		return nil
	}

	return ex.Model(&po.PatchComicAsgn{}).
		Where("id = ?", patchAssign.ID).
		Updates(updates).
		Error
}

func (car *comicAsgnRepo) DeleteAsgnByID(ex Exct, assignmentID string) error {
	ex = car.withTrx(ex)

	result := ex.Where("id = ?", assignmentID).Delete(&po.BasicComicAsgn{})
	if result.Error != nil {
		return fmt.Errorf("Failed to delete assignment: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return REC_NOT_FOUND
	}

	return nil
}
