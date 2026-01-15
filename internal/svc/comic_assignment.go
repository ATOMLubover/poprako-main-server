package svc

import (
	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"

	"go.uber.org/zap"
)

// ComicAsgnSvc defines service operations for comic assignments.
type ComicAsgnSvc interface {
	GetAsgnByID(assignmentID string) (SvcRslt[*po.BasicComicAsgn], SvcErr)
	GetAsgnsByComicID(comicID string, offset, limit int) (SvcRslt[[]po.BasicComicAsgn], SvcErr)
	GetAsgnsByUserID(userID string, offset, limit int) (SvcRslt[[]po.BasicComicAsgn], SvcErr)

	CreateAsgn(args model.CreateComicAsgnArgs) (SvcRslt[string], SvcErr)

	UpdateAsgnByID(patchAssign *po.PatchComicAsgn) SvcErr
}

type comicAsgnSvc struct {
	repo repo.ComicAsgnRepo
}

// NewComicAsgnSvc creates a new ComicAsgnSvc. r must not be nil.
func NewComicAsgnSvc(r repo.ComicAsgnRepo) ComicAsgnSvc {
	if r == nil {
		panic("ComicAsgnRepo cannot be nil")
	}

	return &comicAsgnSvc{repo: r}
}

// GetAsgnByID retrieves a comic assignment by ID.
func (cas *comicAsgnSvc) GetAsgnByID(assignmentID string) (SvcRslt[*po.BasicComicAsgn], SvcErr) {
	asgn, err := cas.repo.GetAsgnByID(nil, assignmentID)
	if err != nil {
		zap.L().Error("Failed to get assignment by ID", zap.String("assignmentID", assignmentID), zap.Error(err))
		return SvcRslt[*po.BasicComicAsgn]{}, DB_FAILURE
	}

	return accept(200, asgn), NO_ERROR
}

// GetAsgnsByComicID retrieves comic assignments by comic ID with pagination.
func (cas *comicAsgnSvc) GetAsgnsByComicID(comicID string, offset, limit int) (SvcRslt[[]po.BasicComicAsgn], SvcErr) {
	asgnList, err := cas.repo.GetAsgnsByComicID(nil, comicID, offset, limit)
	if err != nil {
		zap.L().Error("Failed to get assignments by comic ID", zap.String("comicID", comicID), zap.Error(err))
		return SvcRslt[[]po.BasicComicAsgn]{}, DB_FAILURE
	}

	return accept(200, asgnList), NO_ERROR
}

// GetAsgnsByUserID retrieves comic assignments by user ID with pagination.
func (cas *comicAsgnSvc) GetAsgnsByUserID(userID string, offset, limit int) (SvcRslt[[]po.BasicComicAsgn], SvcErr) {
	asgnList, err := cas.repo.GetAsgnsByUserID(nil, userID, offset, limit)
	if err != nil {
		zap.L().Error("Failed to get assignments by user ID", zap.String("userID", userID), zap.Error(err))
		return SvcRslt[[]po.BasicComicAsgn]{}, DB_FAILURE
	}

	return accept(200, asgnList), NO_ERROR
}

// CreateAsgn creates a new comic assignment.
func (cas *comicAsgnSvc) CreateAsgn(args model.CreateComicAsgnArgs) (SvcRslt[string], SvcErr) {
	// Generate ID for the assignment
	id, err := genUUID()
	if err != nil {
		zap.L().Error("Failed to generate UUID for assignment", zap.Error(err))
		return SvcRslt[string]{}, ID_GEN_FAILURE
	}

	newAssign := &po.NewComicAsgn{
		ID:      id,
		ComicID: args.ComicID,
		UserID:  args.AssigneeID,
	}

	if err := cas.repo.CreateAsgn(nil, newAssign); err != nil {
		zap.L().Error("Failed to create assignment", zap.Error(err))
		return SvcRslt[string]{}, DB_FAILURE
	}

	return accept(201, id), NO_ERROR
}

// UpdateAsgnByID updates a comic assignment by ID.
func (cas *comicAsgnSvc) UpdateAsgnByID(patchAssign *po.PatchComicAsgn) SvcErr {
	if err := cas.repo.UpdateAsgnByID(nil, patchAssign); err != nil {
		zap.L().Error("Failed to update assignment", zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}
