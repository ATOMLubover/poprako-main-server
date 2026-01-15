package repo

import (
	"errors"
	"fmt"

	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
)

// UserRepo defines the interface for user repository operations.
type UserRepo interface {
	Repo

	GetUserByID(ex Executor, userID string) (*po.BasicUser, error)
	GetUserByQQ(ex Executor, qq string) (*po.BasicUser, error)
	RetrieveUsers(ex Executor, opt model.RetrieveUserOpt) ([]po.BasicUser, error)

	GetSecretUserByQQ(ex Executor, qq string) (*po.SecretUser, error)

	CreateUser(ex Executor, newUser *po.NewUser) error

	UpdateUserByID(ex Executor, updateUser *po.PatchUser) error
}

// Default implementation of UserRepo.
type userRepo struct {
	ex Executor
}

func NewUserRepo(ex Executor) UserRepo {
	return &userRepo{ex: ex}
}

// Executor returns the Executor associated with the repository.
func (ur *userRepo) Exec() Executor {
	return ur.ex
}

// withTrx returns the effective executor: tx if non-nil, otherwise repo's executor.
func (ur *userRepo) withTrx(tx Executor) Executor {
	if tx != nil {
		return tx
	}

	return ur.ex
}

// Create a new user.
// The generated ID is returned in newUser.ID.
func (ur *userRepo) CreateUser(ex Executor, newUser *po.NewUser) error {
	ex = ur.withTrx(ex)

	return ex.Create(newUser).Error
}

// Get user basic by ID.
// A nil BasicUser pointer is returned if no user is found.
func (ur *userRepo) GetUserByID(ex Executor, userID string) (*po.BasicUser, error) {
	ex = ur.withTrx(ex)

	ub := &po.BasicUser{}

	if err := ex.
		Where("id = ?", userID).
		First(ub).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get user by ID: %w", err)
	}

	return ub, nil
}

// Get user basic by QQ.
// A nil BasicUser pointer is returned if no user is found.
func (ur *userRepo) GetUserByQQ(ex Executor, qq string) (*po.BasicUser, error) {
	ex = ur.withTrx(ex)

	ub := &po.BasicUser{}

	if err := ex.
		Where("qq = ?", qq).
		First(ub).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get user by qq: %w", err)
	}

	return ub, nil
}

func (ur *userRepo) GetSecretUserByQQ(ex Executor, qq string) (*po.SecretUser, error) {
	ex = ur.withTrx(ex)

	var user po.SecretUser

	if err := ex.
		Where("qq = ?", qq).
		First(&user).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get secret user by qq: %w", err)
	}

	return &user, nil
}

// Update a user's info by ID.
func (ur *userRepo) UpdateUserByID(ex Executor, patchUser *po.PatchUser) error {
	if patchUser.ID == "" {
		return errors.New("user ID is required for update")
	}

	ex = ur.withTrx(ex)

	updates := map[string]any{}

	if patchUser.QQ != nil {
		updates["qq"] = *patchUser.QQ
	}
	if patchUser.Nickname != nil {
		updates["nickname"] = *patchUser.Nickname
	}
	if patchUser.IsAdmin != nil {
		updates["is_admin"] = *patchUser.IsAdmin
	}

	if patchUser.AssignedTranslatorAt != nil {
		if patchUser.AssignedTranslatorAt.IsZero() {
			updates["assigned_translator_at"] = nil
		} else {
			updates["assigned_translator_at"] = *patchUser.AssignedTranslatorAt
		}
	}

	if patchUser.AssignedProofreaderAt != nil {
		if patchUser.AssignedProofreaderAt.IsZero() {
			updates["assigned_proofreader_at"] = nil
		} else {
			updates["assigned_proofreader_at"] = *patchUser.AssignedProofreaderAt
		}
	}

	if patchUser.AssignedTypesetterAt != nil {
		if patchUser.AssignedTypesetterAt.IsZero() {
			updates["assigned_typesetter_at"] = nil
		} else {
			updates["assigned_typesetter_at"] = *patchUser.AssignedTypesetterAt
		}
	}

	if patchUser.AssignedRedrawerAt != nil {
		if patchUser.AssignedRedrawerAt.IsZero() {
			updates["assigned_redrawer_at"] = nil
		} else {
			updates["assigned_redrawer_at"] = *patchUser.AssignedRedrawerAt
		}
	}

	if patchUser.AssignedReviewerAt != nil {
		if patchUser.AssignedReviewerAt.IsZero() {
			updates["assigned_reviewer_at"] = nil
		} else {
			updates["assigned_reviewer_at"] = *patchUser.AssignedReviewerAt
		}
	}

	if patchUser.AssignedUploaderAt != nil {
		if patchUser.AssignedUploaderAt.IsZero() {
			updates["assigned_uploader_at"] = nil
		} else {
			updates["assigned_uploader_at"] = *patchUser.AssignedUploaderAt
		}
	}

	if len(updates) == 0 {
		return nil
	}

	return ex.Model(&po.PatchUser{}).
		Where("id = ?", patchUser.ID).
		Updates(updates).
		Error
}

// RetrieveUsers returns a slice of BasicUser with filtering and pagination.
// A zero-length slice is returned if no users are found.
func (ur *userRepo) RetrieveUsers(ex Executor, opt model.RetrieveUserOpt) ([]po.BasicUser, error) {
	ex = ur.withTrx(ex)

	var users []po.BasicUser

	query := ex

	if opt.Nickname != nil {
		query = query.Where("nickname LIKE ?", "%"+*opt.Nickname+"%")
	}

	if opt.IsAdmin != nil {
		query = query.Where("is_admin = ?", *opt.IsAdmin)
	}

	if opt.IsTranslator != nil {
		if *opt.IsTranslator {
			query = query.Where("assigned_translator_at IS NOT NULL")
		} else {
			query = query.Where("assigned_translator_at IS NULL")
		}
	}

	if opt.IsProofreader != nil {
		if *opt.IsProofreader {
			query = query.Where("assigned_proofreader_at IS NOT NULL")
		} else {
			query = query.Where("assigned_proofreader_at IS NULL")
		}
	}

	if opt.IsTypesetter != nil {
		if *opt.IsTypesetter {
			query = query.Where("assigned_typesetter_at IS NOT NULL")
		} else {
			query = query.Where("assigned_typesetter_at IS NULL")
		}
	}

	if opt.IsRedrawer != nil {
		if *opt.IsRedrawer {
			query = query.Where("assigned_redrawer_at IS NOT NULL")
		} else {
			query = query.Where("assigned_redrawer_at IS NULL")
		}
	}

	if opt.IsReviewer != nil {
		if *opt.IsReviewer {
			query = query.Where("assigned_reviewer_at IS NOT NULL")
		} else {
			query = query.Where("assigned_reviewer_at IS NULL")
		}
	}

	if opt.IsUploader != nil {
		if *opt.IsUploader {
			query = query.Where("assigned_uploader_at IS NOT NULL")
		} else {
			query = query.Where("assigned_uploader_at IS NULL")
		}
	}

	if opt.Offset > 0 {
		query = query.Offset(opt.Offset)
	}

	if opt.Limit > 0 {
		query = query.Limit(opt.Limit)
	}

	if err := query.
		Find(&users).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to retrieve users: %w", err)
	}

	return users, nil
}
