package repo

import (
	"errors"
	"fmt"

	"saas-template-go/internal/model/po"
)

// Note: model structs have been moved to internal/model/po package.

// UserRepo defines the interface for user repository operations.
type UserRepo interface {
	GetUserBasicByID(ex Executor, userID string) (*po.UserBasic, error)
	GetUserBasicByEmail(ex Executor, email string) (*po.UserBasic, error)
	GetUserBasicsByIDs(ex Executor, userIDs []string) ([]po.UserBasic, error)
	GetUserBasicsByEmails(ex Executor, emails []string) ([]po.UserBasic, error)
	GetPwdHashByEmail(ex Executor, email string) (string, error)
	CreateNewUser(ex Executor, newUser *po.NewUser) error
	UpdateUserByID(ex Executor, updateUser *po.UpdateUser) error
}

// Default implementation of UserRepo.
type userRepo struct {
	ex Executor
}

func NewUserRepo(ex Executor) UserRepo {
	return &userRepo{ex: ex}
}

// Executor returns the Executor associated with the repository.
func (ur *userRepo) Executor() Executor {
	return ur.ex
}

// withTransaction returns the effective executor: tx if non-nil, otherwise repo's executor.
func (ur *userRepo) withTransaction(tx Executor) Executor {
	if tx != nil {
		return tx
	}

	return ur.ex
}

// Create a new user.
// The generated ID is returned in newUser.ID.
func (ur *userRepo) CreateNewUser(ex Executor, newUser *po.NewUser) error {
	ex = ur.withTransaction(ex)

	newID, err := genUUID()
	if err != nil {
		return err
	}

	newUser.ID = newID

	return ex.Create(newUser).Error
}

// Get user basic by ID.
// A nil UserBasic pointer is returned if no user is found.
func (ur *userRepo) GetUserBasicByID(ex Executor, userID string) (*po.UserBasic, error) {
	ex = ur.withTransaction(ex)

	ub := &po.UserBasic{}

	if err := ex.
		Where("id = ?", userID).
		First(ub).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get user by ID: %w", err)
	}

	return ub, nil
}

// Get user basic by email.
// A nil UserBasic pointer is returned if no user is found.
func (ur *userRepo) GetUserBasicByEmail(ex Executor, email string) (*po.UserBasic, error) {
	ex = ur.withTransaction(ex)

	ub := &po.UserBasic{}

	if err := ex.
		Where("email = ?", email).
		First(ub).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get user by email: %w", err)
	}

	return ub, nil
}

// Get a list of user basics by their IDs.
// A zero-length slice is returned if no users are found.
func (ur *userRepo) GetUserBasicsByIDs(ex Executor, userIDs []string) ([]po.UserBasic, error) {
	ex = ur.withTransaction(ex)

	var ubLst []po.UserBasic

	if err := ex.
		Where("id IN ?", userIDs).
		Find(&ubLst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get user basics by IDs: %w", err)
	}

	return ubLst, nil
}

// Get a list of user basics by their emails.
// A zero-length slice is returned if no users are found.
func (ur *userRepo) GetUserBasicsByEmails(ex Executor, emails []string) ([]po.UserBasic, error) {
	ex = ur.withTransaction(ex)

	var ubLst []po.UserBasic

	if err := ex.
		Where("email IN ?", emails).
		Find(&ubLst).
		Error; err != nil {
		return nil, fmt.Errorf("Failed to get user basics by emails: %w", err)
	}

	return ubLst, nil
}

// Independent functions.

// Get password hash by email. This is typically used during login.
// An empty string is returned if no user is found.
func (ur *userRepo) GetPwdHashByEmail(ex Executor, email string) (string, error) {
	ex = ur.withTransaction(ex)

	var passwordHash string

	if err := ex.
		Table(po.USER_TABLE).
		Select("password_hash").
		Where("email = ?", email).
		Scan(&passwordHash).
		Error; err != nil {
		return "", fmt.Errorf("Failed to get password hash by email: %w", err)
	}

	return passwordHash, nil
}

// Update a user's info by ID.
func (ur *userRepo) UpdateUserByID(ex Executor, updateUser *po.UpdateUser) error {
	if updateUser.ID == "" {
		return errors.New("user ID is required for update")
	}

	ex = ur.withTransaction(ex)

	return ex.Save(updateUser).Error
}
