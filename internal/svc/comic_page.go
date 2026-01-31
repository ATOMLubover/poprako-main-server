package svc

import (
	"fmt"

	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/oss"
	"poprako-main-server/internal/repo"

	"go.uber.org/zap"
)

// FIXME: track periodically whether some pages have been uploaded completely,
// but not yet marked as uploaded in DB.

type ComicPageSvc interface {
	GetPageByID(pageID string) (SvcRslt[model.ComicPageInfo], SvcErr)
	GetPagesByComicID(comicID string) (SvcRslt[[]model.ComicPageInfo], SvcErr)

	CreatePages(
		opID string,
		args []model.CreateComicPageArgs,
	) (
		SvcRslt[[]model.CreateComicPageReply],
		SvcErr,
	)

	UpdatePageByID(opID string, args *model.PatchComicPageArgs) SvcErr

	DeletePageByID(pageID string) SvcErr
}

type comicPageSvc struct {
	pageRepo  repo.ComicPageRepo
	comicRepo repo.ComicRepo
	unitRepo  repo.ComicUnitRepo
	ossClient oss.OSSClient
}

func NewComicPageSvc(
	pageRepo repo.ComicPageRepo,
	comicRepo repo.ComicRepo,
	unitRepo repo.ComicUnitRepo,
	ossClient oss.OSSClient,
) ComicPageSvc {
	return &comicPageSvc{
		pageRepo:  pageRepo,
		comicRepo: comicRepo,
		unitRepo:  unitRepo,
		ossClient: ossClient,
	}
}

func (cps *comicPageSvc) GetPageByID(pageID string) (SvcRslt[model.ComicPageInfo], SvcErr) {
	page, err := cps.pageRepo.GetPageByID(nil, pageID)
	if err != nil {
		zap.L().Error("Failed to get page by ID", zap.String("pageID", pageID), zap.Error(err))
		return SvcRslt[model.ComicPageInfo]{}, DB_FAILURE
	}

	// Get unit counts for the page
	counts, err := cps.unitRepo.GetUnitCountsByPageID(nil, pageID)
	if err != nil {
		zap.L().Error("Failed to get unit counts for page", zap.String("pageID", pageID), zap.Error(err))
		return SvcRslt[model.ComicPageInfo]{}, DB_FAILURE
	}

	// Only generate OSS URL if page is uploaded
	var ossURL string
	if page.Uploaded {
		// Dynamically generate OSS key: comic/{comic_id}/page_{index}
		ossKey := fmt.Sprintf("comic/%s/page_%d", page.ComicID, page.Index)

		// Get presigned URL for download
		var err error
		ossURL, err = cps.ossClient.PresignGet(ossKey)
		if err != nil {
			zap.L().Error("Failed to generate presigned URL for page", zap.String("pageID", pageID), zap.String("ossKey", ossKey), zap.Error(err))
			return SvcRslt[model.ComicPageInfo]{}, DB_FAILURE
		}
	}

	pageInfo := model.ComicPageInfo{
		ID:                  page.ID,
		ComicID:             page.ComicID,
		Index:               page.Index,
		OSSURL:              ossURL,
		Uploaded:            page.Uploaded,
		InboxUnitCount:      counts.Inbox,
		OutboxUnitCount:     counts.Outbox,
		TranslatedUnitCount: counts.Translated,
		ProvedUnitCount:     counts.Proved,
	}

	return accept(200, pageInfo), NO_ERROR
}

func (cps *comicPageSvc) GetPagesByComicID(comicID string) (SvcRslt[[]model.ComicPageInfo], SvcErr) {
	pages, err := cps.pageRepo.GetPagesByComicID(nil, comicID)
	if err != nil {
		zap.L().Error("Failed to get pages by comic ID", zap.String("comicID", comicID), zap.Error(err))
		return SvcRslt[[]model.ComicPageInfo]{}, DB_FAILURE
	}

	// Collect page IDs for batch unit counts query
	pageIDs := make([]string, len(pages))
	for i, page := range pages {
		pageIDs[i] = page.ID
	}

	// Get unit counts for all pages in one query
	countsMap, err := cps.unitRepo.GetUnitCountsByPageIDs(nil, pageIDs)
	if err != nil {
		zap.L().Error("Failed to get unit counts for pages", zap.String("comicID", comicID), zap.Error(err))
		return SvcRslt[[]model.ComicPageInfo]{}, DB_FAILURE
	}

	pageInfos := make([]model.ComicPageInfo, len(pages))
	for i, page := range pages {
		// Only generate OSS URL if page is uploaded
		var ossURL string
		if page.Uploaded {
			// Dynamically generate OSS key for each page
			ossKey := fmt.Sprintf("comic/%s/page_%d", page.ComicID, page.Index)

			// Get presigned URL for each page
			var err error
			ossURL, err = cps.ossClient.PresignGet(ossKey)
			if err != nil {
				zap.L().Error("Failed to generate presigned URL for page", zap.String("pageID", page.ID), zap.String("ossKey", ossKey), zap.Error(err))
				return SvcRslt[[]model.ComicPageInfo]{}, DB_FAILURE
			}
		}

		// Get counts from map, default to zero if not found
		counts := countsMap[page.ID]

		pageInfos[i] = model.ComicPageInfo{
			ID:                  page.ID,
			ComicID:             page.ComicID,
			Index:               page.Index,
			OSSURL:              ossURL,
			Uploaded:            page.Uploaded,
			InboxUnitCount:      counts.Inbox,
			OutboxUnitCount:     counts.Outbox,
			TranslatedUnitCount: counts.Translated,
			ProvedUnitCount:     counts.Proved,
		}
	}

	return accept(200, pageInfos), NO_ERROR
}

func (cps *comicPageSvc) CreatePages(
	opID string,
	args []model.CreateComicPageArgs,
) (
	SvcRslt[[]model.CreateComicPageReply],
	SvcErr,
) {
	if len(args) == 0 {
		return SvcRslt[[]model.CreateComicPageReply]{}, INVALID_PAGE_DATA
	}

	// Verify all pages belong to the same comic
	comicID := args[0].ComicID
	for _, arg := range args {
		if arg.ComicID != comicID {
			return SvcRslt[[]model.CreateComicPageReply]{}, INVALID_PAGE_DATA
		}
	}

	// Verify comic exists
	_, err := cps.comicRepo.GetComicByID(nil, comicID)
	if err != nil {
		zap.L().Error("Failed to verify comic exists", zap.String("comicID", comicID), zap.Error(err))
		return SvcRslt[[]model.CreateComicPageReply]{}, DB_FAILURE
	}

	// Create new pages with generated IDs and OSS keys
	newPages := make([]po.NewComicPage, len(args))
	replies := make([]model.CreateComicPageReply, len(args))

	uploadedFalse := false

	for i, arg := range args {
		pageID, err := genUUID()
		if err != nil {
			zap.L().Error("Failed to generate UUID for page", zap.Error(err))
			return SvcRslt[[]model.CreateComicPageReply]{}, ID_GEN_FAILURE
		}

		// Dynamically generate OSS key: comic/{comic_id}/page_{index}.{ext}
		ossKey := fmt.Sprintf("comic/%s/page_%d.%s", arg.ComicID, arg.Index, arg.ImageExt)

		newPages[i] = po.NewComicPage{
			ID:       pageID,
			ComicID:  arg.ComicID,
			Index:    arg.Index,
			Uploaded: &uploadedFalse,
		}

		// Generate presigned upload URL
		uploadURL, err := cps.ossClient.PresignPut(ossKey)
		if err != nil {
			zap.L().Error("Failed to generate presigned upload URL", zap.String("ossKey", ossKey), zap.Error(err))
			return SvcRslt[[]model.CreateComicPageReply]{}, DB_FAILURE
		}

		replies[i] = model.CreateComicPageReply{
			ID:     pageID,
			OSSURL: uploadURL,
		}
	}

	// Save to database
	if err := cps.pageRepo.CreatePages(newPages); err != nil {
		zap.L().Error("Failed to create pages", zap.String("comicID", comicID), zap.Error(err))
		return SvcRslt[[]model.CreateComicPageReply]{}, DB_FAILURE
	}

	return accept(201, replies), NO_ERROR
}

func (cps *comicPageSvc) UpdatePageByID(opID string, args *model.PatchComicPageArgs) SvcErr {
	if args.ID == "" {
		return INVALID_PAGE_DATA
	}

	// Verify page exists
	_, err := cps.pageRepo.GetPageByID(nil, args.ID)
	if err != nil {
		zap.L().Error("Failed to get page for update", zap.String("pageID", args.ID), zap.Error(err))
		return DB_FAILURE
	}

	// Build patch object
	patchPage := &po.PatchComicPage{
		ID:       args.ID,
		Uploaded: args.Uploaded,
	}

	if err := cps.pageRepo.UpdatePageByID(nil, patchPage); err != nil {
		zap.L().Error("Failed to update page", zap.String("pageID", args.ID), zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}

func (cps *comicPageSvc) DeletePageByID(pageID string) SvcErr {
	if err := cps.pageRepo.DeletePageByID(nil, pageID); err != nil {
		if err == repo.REC_NOT_FOUND {
			zap.L().Warn("Page not found for deletion", zap.String("pageID", pageID))
			return NOT_FOUND
		}
		zap.L().Error("Failed to delete page", zap.String("pageID", pageID), zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}
