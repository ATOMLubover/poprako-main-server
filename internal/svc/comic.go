package svc

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/oss"
	"poprako-main-server/internal/repo"
	comicPkg "poprako-main-server/internal/svc/comic"

	"go.uber.org/zap"
)

type ComicSvc interface {
	GetComicInfoByID(comicID string) (SvcRslt[model.ComicInfo], SvcErr)
	GetComicBriefsByWorksetID(worksetID string, offset, limit int) (SvcRslt[[]model.ComicBrief], SvcErr)
	RetrieveComics(opt model.RetrieveComicOpt) (SvcRslt[[]model.ComicBrief], SvcErr)

	ExportComic(comicID string) (SvcRslt[model.ExportComicReply], SvcErr)
	ExportBaseURI() string

	ImportComic(opID string, comicID string, fileName string, reader io.Reader) SvcErr

	CreateComic(opID string, args model.CreateComicArgs) (SvcRslt[model.CreateComicReply], SvcErr)

	UpdateComicByID(args model.UpdateComicArgs) SvcErr

	DeleteComicByID(comicID string) SvcErr
}

type comicSvc struct {
	repo          repo.ComicRepo
	userRepo      repo.UserRepo
	comicAsgnRepo repo.ComicAsgnRepo
	comicPageRepo repo.ComicPageRepo
	comicUnitRepo repo.ComicUnitRepo
	exportDir     string
	ossClient     oss.OSSClient
}

func NewComicSvc(
	r repo.ComicRepo,
	ur repo.UserRepo,
	car repo.ComicAsgnRepo,
	cpr repo.ComicPageRepo,
	cur repo.ComicUnitRepo,
	exportDir string,
	ossClient oss.OSSClient,
) ComicSvc {
	if r == nil {
		panic("ComicRepo cannot be nil")
	}
	if ur == nil {
		panic("UserRepo cannot be nil")
	}
	if car == nil {
		panic("ComicAsgnRepo cannot be nil")
	}
	if cpr == nil {
		panic("ComicPageRepo cannot be nil")
	}
	if cur == nil {
		panic("ComicUnitRepo cannot be nil")
	}
	if exportDir == "" {
		panic("exportDir cannot be empty")
	}
	if ossClient == nil {
		panic("ossClient cannot be nil")
	}

	return &comicSvc{
		repo:          r,
		userRepo:      ur,
		comicAsgnRepo: car,
		comicPageRepo: cpr,
		comicUnitRepo: cur,
		exportDir:     exportDir,
		ossClient:     ossClient,
	}
}

// GetComicInfoByID retrieves detailed comic info by ID.
func (cs *comicSvc) GetComicInfoByID(comicID string) (SvcRslt[model.ComicInfo], SvcErr) {
	basic, err := cs.repo.GetComicByID(nil, comicID)
	if err != nil {
		zap.L().Error("Failed to get comic by ID", zap.String("comicID", comicID), zap.Error(err))
		return SvcRslt[model.ComicInfo]{}, DB_FAILURE
	}

	info := model.ComicInfo{
		ID:              basic.ID,
		WorksetID:       basic.WorksetID,
		WorksetIndex:    basic.WorksetIndex,
		Index:           basic.Index,
		CreatorID:       basic.CreatorID,
		CreatorNickname: basic.CreatorNickname,
		Author:          basic.Author,
		Title:           basic.Title,
		Description:     basic.Description,
		Comment:         basic.Comment,
		PageCount:       basic.PageCount,
		CreatedAt:       basic.CreatedAt.Unix(),
		UpdatedAt:       basic.UpdatedAt.Unix(),
	}

	// Handle optional timestamp fields
	info.TranslatingStartedAt = timePtrToInt64Ptr(basic.TranslatingStartedAt)
	info.TranslatingCompletedAt = timePtrToInt64Ptr(basic.TranslatingCompletedAt)
	info.ProofreadingStartedAt = timePtrToInt64Ptr(basic.ProofreadingStartedAt)
	info.ProofreadingCompletedAt = timePtrToInt64Ptr(basic.ProofreadingCompletedAt)
	info.TypesettingStartedAt = timePtrToInt64Ptr(basic.TypesettingStartedAt)
	info.TypesettingCompletedAt = timePtrToInt64Ptr(basic.TypesettingCompletedAt)
	info.ReviewingCompletedAt = timePtrToInt64Ptr(basic.ReviewingCompletedAt)
	info.UploadingCompletedAt = timePtrToInt64Ptr(basic.UploadingCompletedAt)

	return accept(200, info), NO_ERROR
}

// GetComicBriefsByWorksetID retrieves brief comic info by workset ID with pagination.
func (cs *comicSvc) GetComicBriefsByWorksetID(worksetID string, offset, limit int) (SvcRslt[[]model.ComicBrief], SvcErr) {
	briefs, err := cs.repo.GetComicsByWorksetID(nil, worksetID, offset, limit)
	if err != nil {
		zap.L().Error("Failed to get comics by workset ID", zap.String("worksetID", worksetID), zap.Error(err))
		return SvcRslt[[]model.ComicBrief]{}, DB_FAILURE
	}

	lst := make([]model.ComicBrief, 0, len(briefs))

	for _, cb := range briefs {
		brief := model.ComicBrief{
			ID:           cb.ID,
			WorksetID:    cb.WorksetID,
			WorksetIndex: cb.WorksetIndex,
			Index:        cb.Index,
			Author:       cb.Author,
			Title:        cb.Title,
			PageCount:    cb.PageCount,
		}

		// Handle optional timestamp fields
		brief.TranslatingStartedAt = timePtrToInt64Ptr(cb.TranslatingStartedAt)
		brief.TranslatingCompletedAt = timePtrToInt64Ptr(cb.TranslatingCompletedAt)
		brief.ProofreadingStartedAt = timePtrToInt64Ptr(cb.ProofreadingStartedAt)
		brief.ProofreadingCompletedAt = timePtrToInt64Ptr(cb.ProofreadingCompletedAt)
		brief.TypesettingStartedAt = timePtrToInt64Ptr(cb.TypesettingStartedAt)
		brief.TypesettingCompletedAt = timePtrToInt64Ptr(cb.TypesettingCompletedAt)
		brief.ReviewingCompletedAt = timePtrToInt64Ptr(cb.ReviewingCompletedAt)
		brief.UploadingCompletedAt = timePtrToInt64Ptr(cb.UploadingCompletedAt)
		lst = append(lst, brief)
	}

	return accept(200, lst), NO_ERROR
}

// RetrieveComics retrieves comics with filtering and pagination.
func (cs *comicSvc) RetrieveComics(opt model.RetrieveComicOpt) (SvcRslt[[]model.ComicBrief], SvcErr) {
	briefs, err := cs.repo.RetrieveComics(nil, opt)
	if err != nil {
		zap.L().Error("Failed to retrieve comics", zap.Error(err))
		return SvcRslt[[]model.ComicBrief]{}, DB_FAILURE
	}

	lst := make([]model.ComicBrief, 0, len(briefs))
	for _, cb := range briefs {
		brief := model.ComicBrief{
			ID:           cb.ID,
			WorksetID:    cb.WorksetID,
			WorksetIndex: cb.WorksetIndex,
			Index:        cb.Index,
			Author:       cb.Author,
			Title:        cb.Title,
			PageCount:    cb.PageCount,
		}

		// Handle optional timestamp fields
		brief.TranslatingStartedAt = timePtrToInt64Ptr(cb.TranslatingStartedAt)
		brief.TranslatingCompletedAt = timePtrToInt64Ptr(cb.TranslatingCompletedAt)
		brief.ProofreadingStartedAt = timePtrToInt64Ptr(cb.ProofreadingStartedAt)
		brief.ProofreadingCompletedAt = timePtrToInt64Ptr(cb.ProofreadingCompletedAt)
		brief.TypesettingStartedAt = timePtrToInt64Ptr(cb.TypesettingStartedAt)
		brief.TypesettingCompletedAt = timePtrToInt64Ptr(cb.TypesettingCompletedAt)
		brief.ReviewingCompletedAt = timePtrToInt64Ptr(cb.ReviewingCompletedAt)
		brief.UploadingCompletedAt = timePtrToInt64Ptr(cb.UploadingCompletedAt)

		lst = append(lst, brief)
	}

	return accept(200, lst), NO_ERROR
}

// ExportComic exports a comic to LabelPlus format.
func (cs *comicSvc) ExportComic(comicID string) (SvcRslt[model.ExportComicReply], SvcErr) {
	// Ensure export directory exists
	if err := os.MkdirAll(cs.exportDir, 0o755); err != nil {
		zap.L().Error("Failed to create export directory", zap.String("dir", cs.exportDir), zap.Error(err))
		return SvcRslt[model.ExportComicReply]{}, DB_FAILURE
	}

	// Clean old exports if exceeding limit
	if err := cs.cleanOldExports(); err != nil {
		zap.L().Warn("Failed to clean old exports", zap.Error(err))
		// Continue anyway - this is not critical
	}

	// Export the comic using the comic package
	filePath, err := comicPkg.ExportLabelplusComic(
		comicID,
		cs.exportDir,
		cs.repo,
		cs.comicPageRepo,
		cs.comicUnitRepo,
	)
	if err != nil {
		if err == repo.REC_NOT_FOUND {
			return SvcRslt[model.ExportComicReply]{}, NOT_FOUND
		}
		zap.L().Error("Failed to export comic", zap.String("comicID", comicID), zap.Error(err))
		return SvcRslt[model.ExportComicReply]{}, DB_FAILURE
	}

	// Return relative URI
	fileName := filepath.Base(filePath)
	exportURI := cs.ExportBaseURI() + url.PathEscape(fileName)

	return accept(200, model.ExportComicReply{ExportURI: exportURI}), NO_ERROR
}

func (*comicSvc) ExportBaseURI() string {
	return "/comics/export/"
}

func (cs *comicSvc) ImportComic(
	opID string,
	comicID string,
	fileName string,
	reader io.Reader,
) SvcErr {
	// Check operation permission: opID must be assigned to the comic
	asgn, err := cs.comicAsgnRepo.GetAsgnsByUserAndComicID(nil, opID, comicID)
	if err != nil {
		zap.L().Error("Failed to get comic assignment for import", zap.String("userID", opID), zap.String("comicID", comicID), zap.Error(err))
		return PERMISSION_DENIED
	}
	if asgn == nil {
		zap.L().Warn("User not assigned to comic for import", zap.String("userID", opID), zap.String("comicID", comicID))
		return PERMISSION_DENIED
	}

	// Only allow import by users assigned as reviewer, translator or proofreader
	if asgn.AssignedReviewerAt == nil && asgn.AssignedTranslatorAt == nil && asgn.AssignedProofreaderAt == nil {
		zap.L().Warn("User does not have required role for importing comic", zap.String("userID", opID), zap.String("comicID", comicID))
		return PERMISSION_DENIED
	}

	// Determine import role: Proofreader has priority over Translator
	isProofreader := asgn.AssignedProofreaderAt != nil

	importOpts := comicPkg.ImportOptions{
		IsProofreader: isProofreader,
		UserID:        opID,
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {
	case ".txt":
		// Import LabelPlus format
		if err := comicPkg.ImportLabelplusComic(reader, comicID, cs.comicPageRepo, cs.comicUnitRepo, importOpts); err != nil {
			zap.L().Error("Failed to import LabelPlus comic",
				zap.String("comicID", comicID),
				zap.String("fileName", fileName),
				zap.Error(err))
			return INVALID_PROJ_DATA
		}
		return NO_ERROR

	default:
		zap.L().Warn("Unsupported project file extension",
			zap.String("comicID", comicID),
			zap.String("fileName", fileName),
			zap.String("extension", ext))
		return INVALID_PROJ_EXT
	}
}

// cleanOldExports removes oldest export files if count exceeds 30.
func (cs *comicSvc) cleanOldExports() error {
	const maxExports = 30

	// List all files in export directory
	files, err := os.ReadDir(cs.exportDir)
	if err != nil {
		return fmt.Errorf("failed to read export directory: %w", err)
	}

	// Filter out directories and collect file info
	type fileInfo struct {
		name    string
		modTime time.Time
	}

	var fileList []fileInfo
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		info, err := f.Info()
		if err != nil {
			continue
		}

		fileList = append(fileList, fileInfo{
			name:    f.Name(),
			modTime: info.ModTime(),
		})
	}

	// If under limit, nothing to do
	if len(fileList) <= maxExports {
		return nil
	}

	// Sort by modification time (oldest first)
	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].modTime.Before(fileList[j].modTime)
	})

	// Delete oldest files to bring count down to maxExports
	toDelete := len(fileList) - maxExports
	for i := 0; i < toDelete; i++ {
		filePath := filepath.Join(cs.exportDir, fileList[i].name)
		if err := os.Remove(filePath); err != nil {
			zap.L().Warn("Failed to delete old export file",
				zap.String("file", filePath),
				zap.Error(err))
		}
	}

	return nil
}

// CreateComic creates a new comic.
func (cs *comicSvc) CreateComic(opID string, args model.CreateComicArgs) (SvcRslt[model.CreateComicReply], SvcErr) {
	// Check if creator is admin
	creator, err := cs.userRepo.GetUserByID(nil, opID)
	if err != nil {
		zap.L().Error("Failed to get creator info for comic creation", zap.String("userID", opID), zap.Error(err))
		return SvcRslt[model.CreateComicReply]{}, DB_FAILURE
	}

	if !creator.IsAdmin {
		zap.L().Warn("Non-admin user attempted to create comic", zap.String("userID", opID))
		return SvcRslt[model.CreateComicReply]{}, PERMISSION_DENIED
	}

	// Validate pre-assignments
	if len(args.PreAsgns) > 0 {
		if svcErr := cs.validatePreAssignments(args.PreAsgns); svcErr != NO_ERROR {
			return SvcRslt[model.CreateComicReply]{}, svcErr
		}
	}

	// Generate UUID for the new comic
	newID, err := genUUID()
	if err != nil {
		zap.L().Error("Failed to generate UUID for new comic", zap.Error(err))
		return SvcRslt[model.CreateComicReply]{}, ID_GEN_FAILURE
	}

	// Create comic and assignments in a transaction
	if err := cs.repo.Exec().Transaction(func(tx repo.Executor) error {
		// Create the comic
		newComic := &po.NewComic{
			ID:          newID,
			WorksetID:   args.WorksetID,
			CreatorID:   opID,
			Author:      args.Author,
			Title:       args.Title,
			Description: args.Description,
			Comment:     args.Comment,
		}

		if err := cs.repo.CreateComic(newComic); err != nil {
			return fmt.Errorf("failed to create comic: %w", err)
		}

		// Create pre-assignments
		if err := cs.createPreAssignments(tx, newID, args.PreAsgns); err != nil {
			return err
		}

		return nil
	}); err != nil {
		zap.L().Error("Failed to create comic with assignments", zap.String("worksetID", args.WorksetID), zap.Error(err))
		return SvcRslt[model.CreateComicReply]{}, DB_FAILURE
	}

	return accept(201, model.CreateComicReply{ID: newID}), NO_ERROR
}

// UpdateComicByID updates comic info by ID.
func (cs *comicSvc) UpdateComicByID(args model.UpdateComicArgs) SvcErr {
	now := time.Now()

	patch := &po.PatchComic{
		ID:          args.ID,
		Author:      args.Author,
		Title:       args.Title,
		Description: args.Description,
		Comment:     args.Comment,
	}

	// Handle workflow timestamp toggles - only set when true
	if args.TranslatingStarted != nil && *args.TranslatingStarted {
		patch.TranslatingStartedAt = &now
	}
	if args.TranslatingCompleted != nil && *args.TranslatingCompleted {
		patch.TranslatingCompletedAt = &now
	}
	if args.ProofreadingStarted != nil && *args.ProofreadingStarted {
		patch.ProofreadingStartedAt = &now
	}
	if args.ProofreadingCompleted != nil && *args.ProofreadingCompleted {
		patch.ProofreadingCompletedAt = &now
	}
	if args.TypesettingStarted != nil && *args.TypesettingStarted {
		patch.TypesettingStartedAt = &now
	}
	if args.TypesettingCompleted != nil && *args.TypesettingCompleted {
		patch.TypesettingCompletedAt = &now
	}
	if args.ReviewingCompleted != nil && *args.ReviewingCompleted {
		patch.ReviewingCompletedAt = &now
	}
	if args.UploadingCompleted != nil && *args.UploadingCompleted {
		patch.UploadingCompletedAt = &now
	}

	if err := cs.repo.UpdateComicByID(nil, patch); err != nil {
		zap.L().Error("Failed to update comic", zap.String("comicID", args.ID), zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}

// validatePreAssignments validates that all pre-assigned users have the required qualifications.
func (cs *comicSvc) validatePreAssignments(preAsgns []model.PreAsgnArgs) SvcErr {
	for _, preAsgn := range preAsgns {
		user, err := cs.userRepo.GetUserByID(nil, preAsgn.AssigneeID)
		if err != nil {
			zap.L().Error("Failed to get user info for pre-assignment validation",
				zap.String("userID", preAsgn.AssigneeID), zap.Error(err))
			return DB_FAILURE
		}

		// Check if user has required qualifications
		if preAsgn.IsTranslator != nil && *preAsgn.IsTranslator {
			if user.AssignedTranslatorAt == nil {
				zap.L().Warn("User does not have translator qualification",
					zap.String("userID", preAsgn.AssigneeID))
				return PERMISSION_DENIED
			}
		}

		if preAsgn.IsProofreader != nil && *preAsgn.IsProofreader {
			if user.AssignedProofreaderAt == nil {
				zap.L().Warn("User does not have proofreader qualification",
					zap.String("userID", preAsgn.AssigneeID))
				return PERMISSION_DENIED
			}
		}

		if preAsgn.IsTypesetter != nil && *preAsgn.IsTypesetter {
			if user.AssignedTypesetterAt == nil {
				zap.L().Warn("User does not have typesetter qualification",
					zap.String("userID", preAsgn.AssigneeID))
				return PERMISSION_DENIED
			}
		}

		if preAsgn.IsRedrawer != nil && *preAsgn.IsRedrawer {
			if user.AssignedRedrawerAt == nil {
				zap.L().Warn("User does not have redrawer qualification",
					zap.String("userID", preAsgn.AssigneeID))
				return PERMISSION_DENIED
			}
		}

		if preAsgn.IsReviewer != nil && *preAsgn.IsReviewer {
			if user.AssignedReviewerAt == nil {
				zap.L().Warn("User does not have reviewer qualification",
					zap.String("userID", preAsgn.AssigneeID))
				return PERMISSION_DENIED
			}
		}
	}

	return NO_ERROR
}

// createPreAssignments creates comic assignments for pre-assigned users.
// For each pre-assignment, it creates a new assignment record and sets the role timestamps.
func (cs *comicSvc) createPreAssignments(tx repo.Executor, comicID string, preAsgns []model.PreAsgnArgs) error {
	now := time.Now()

	for _, preAsgn := range preAsgns {
		// Generate assignment ID
		asgnID, err := genUUID()
		if err != nil {
			return fmt.Errorf("failed to generate assignment ID: %w", err)
		}

		// Create the assignment record
		newAsgn := &po.NewComicAsgn{
			ID:      asgnID,
			ComicID: comicID,
			UserID:  preAsgn.AssigneeID,
		}

		if err := cs.comicAsgnRepo.CreateAsgn(tx, newAsgn); err != nil {
			return fmt.Errorf("failed to create assignment: %w", err)
		}

		// Set role timestamps based on pre-assignment roles
		patchAsgn := &po.PatchComicAsgn{
			ID: asgnID,
		}

		if preAsgn.IsTranslator != nil && *preAsgn.IsTranslator {
			patchAsgn.AssignedTranslatorAt = &now
		}
		if preAsgn.IsProofreader != nil && *preAsgn.IsProofreader {
			patchAsgn.AssignedProofreaderAt = &now
		}
		if preAsgn.IsTypesetter != nil && *preAsgn.IsTypesetter {
			patchAsgn.AssignedTypesetterAt = &now
		}
		if preAsgn.IsRedrawer != nil && *preAsgn.IsRedrawer {
			patchAsgn.AssignedRedrawerAt = &now
		}
		if preAsgn.IsReviewer != nil && *preAsgn.IsReviewer {
			patchAsgn.AssignedReviewerAt = &now
		}

		// Update the assignment with role timestamps
		if err := cs.comicAsgnRepo.UpdateAsgnByID(tx, patchAsgn); err != nil {
			return fmt.Errorf("failed to update assignment roles: %w", err)
		}
	}

	return nil
}

func (cs *comicSvc) DeleteComicByID(comicID string) SvcErr {
	// First, get all pages for this comic
	pages, err := cs.comicPageRepo.GetPagesByComicID(nil, comicID)
	if err != nil {
		zap.L().Error("Failed to get pages for comic deletion", zap.String("comicID", comicID), zap.Error(err))
		return DB_FAILURE
	}

	// Delete all pages concurrently
	// For each page: delete OSS object first, then delete page from DB
	// OSS deletion failure must abort the entire operation
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	ctx := context.Background()

	for _, page := range pages {
		wg.Add(1)
		go func(p po.BasicComicPage) {
			defer wg.Done()

			// Delete OSS object first (if exists)
			if p.OSSKey != "" {
				if err := cs.ossClient.DeleteObject(ctx, p.OSSKey); err != nil {
					zap.L().Error("Failed to delete OSS object for page",
						zap.String("comicID", comicID),
						zap.String("pageID", p.ID),
						zap.String("ossKey", p.OSSKey),
						zap.Error(err))
					mu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					mu.Unlock()
					return
				}
			}

			// Then delete page from database
			// CASCADE will automatically delete units
			if err := cs.comicPageRepo.DeletePageByID(nil, p.ID); err != nil {
				zap.L().Error("Failed to delete page from DB",
					zap.String("comicID", comicID),
					zap.String("pageID", p.ID),
					zap.Error(err))
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}
		}(page)
	}

	// Wait for all page deletions to complete
	wg.Wait()

	// If any page deletion failed, abort
	if firstErr != nil {
		zap.L().Error("Failed to delete pages for comic", zap.String("comicID", comicID), zap.Error(firstErr))
		return DB_FAILURE
	}

	// All pages deleted successfully, now delete the comic
	if err := cs.repo.DeleteComicByID(nil, comicID); err != nil {
		if err == repo.REC_NOT_FOUND {
			zap.L().Warn("Comic not found for deletion", zap.String("comicID", comicID))
			return NOT_FOUND
		}
		zap.L().Error("Failed to delete comic", zap.String("comicID", comicID), zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}
