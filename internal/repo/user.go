package repo

import (
	"errors"
	"fmt"

	"poprako-main-server/internal/model/po"
)

// UserRepo defines the interface for user repository operations.
type UserRepo interface {
	Repo

	GetUserByID(ex Executor, userID string) (*po.BasicUser, error)
	GetUserByQQ(ex Executor, qq string) (*po.BasicUser, error)
	// GetUsersByIDs(ex Executor, userIDs []string) ([]po.BasicUser, error)
	// GetUsersByQQs(ex Executor, qqs []string) ([]po.BasicUser, error)

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

// // Get a list of user basics by their IDs.
// // A zero-length slice is returned if no users are found.
// func (ur *userRepo) GetUsersByIDs(ex Executor, userIDs []string) ([]po.BasicUser, error) {
// 	ex = ur.withTrx(ex)

// 	var ubLst []po.BasicUser

// 	if err := ex.
// 		Where("id IN ?", userIDs).
// 		Find(&ubLst).
// 		Error; err != nil {
// 		return nil, fmt.Errorf("Failed to get user basics by IDs: %w", err)
// 	}

// 	return ubLst, nil
// }

// // Get a list of user basics by their qqs.
// // A zero-length slice is returned if no users are found.
// func (ur *userRepo) GetUsersByQQs(ex Executor, qqs []string) ([]po.BasicUser, error) {
// 	ex = ur.withTrx(ex)

// 	var ubLst []po.BasicUser

// 	if err := ex.
// 		Where("qq = ANY(?)", qqs).
// 		Find(&ubLst).
// 		Error; err != nil {
// 		return nil, fmt.Errorf("Failed to get user basics by QQs: %w", err)
// 	}

// 	return ubLst, nil
// }

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

	updates := map[string]interface{}{}

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
		if *patchUser.AssignedTranslatorAt == 0 {
			updates["assigned_translator_at"] = nil
		} else {
			updates["assigned_translator_at"] = *patchUser.AssignedTranslatorAt
		}
	}

	if patchUser.AssignedProofreaderAt != nil {
		if *patchUser.AssignedProofreaderAt == 0 {
			updates["assigned_proofreader_at"] = nil
		} else {
			updates["assigned_proofreader_at"] = *patchUser.AssignedProofreaderAt
		}
	}

	if patchUser.AssignedTypesetterAt != nil {
		if *patchUser.AssignedTypesetterAt == 0 {
			updates["assigned_typesetter_at"] = nil
		} else {
			updates["assigned_typesetter_at"] = *patchUser.AssignedTypesetterAt
		}
	}

	if patchUser.AssignedRedrawerAt != nil {
		if *patchUser.AssignedRedrawerAt == 0 {
			updates["assigned_redrawer_at"] = nil
		} else {
			updates["assigned_redrawer_at"] = *patchUser.AssignedRedrawerAt
		}
	}

	if patchUser.AssignedReviewerAt != nil {
		if *patchUser.AssignedReviewerAt == 0 {
			updates["assigned_reviewer_at"] = nil
		} else {
			updates["assigned_reviewer_at"] = *patchUser.AssignedReviewerAt
		}
	}

	if patchUser.AssignedUploaderAt != nil {
		if *patchUser.AssignedUploaderAt == 0 {
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
