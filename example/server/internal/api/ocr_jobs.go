package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
	"ai.ro/syncra/internal/ocr"
)

var syncPathFunc = syncPath

const defaultOCRJobListSize = 20
const maxOCRJobListSize = 100
const maxOCRJobDeleteIDs = 100
const maxOCRJobPagesWithoutSchema = 1000
const maxOCRJobPagesWithSchema = 150

type ocrJobListCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        uuid.UUID `json:"id"`
	Sort      string    `json:"sort"`
}

func writeOCRJobFile(dir string, id uuid.UUID, mimeType string, data []byte) (storedPath string, err error) {
	if id == uuid.Nil {
		return "", errors.New("OCR job id is required")
	}
	ext, ok := ocrJobFileExtension(mimeType)
	if !ok {
		return "", errors.New("unsupported file type")
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}
	dir = filepath.Clean(dir)
	if err := ensureDirDurable(dir); err != nil {
		return "", err
	}
	path := filepath.Join(dir, id.String()+ext)
	claimPath := ocrJobClaimPath(dir, id)
	write := ocrJobFileWrite{
		dir:       dir,
		finalPath: path,
		claimPath: claimPath,
	}
	success := false
	defer func() {
		if !success {
			err = errors.Join(err, write.cleanup())
		}
	}()

	claim, err := createOCRJobClaim(dir, id, claimPath)
	if err != nil {
		return "", err
	}
	write.claimCreated = true
	if err := syncAndClose(claim); err != nil {
		return "", err
	}
	if err := syncPathFunc(dir); err != nil {
		return "", err
	}

	tmp, err := os.CreateTemp(dir, id.String()+"-*.tmp")
	if err != nil {
		return "", err
	}
	write.tmpPath = tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return "", err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	if err := write.linkFinalFile(); err != nil {
		return "", err
	}
	if err := syncPathFunc(dir); err != nil {
		return "", err
	}
	if err := os.Remove(write.tmpPath); err != nil {
		return "", err
	}
	write.tmpPath = ""
	if err := syncPathFunc(dir); err != nil {
		return "", err
	}
	success = true
	return path, nil
}

func removeOCRJobFile(path string) error {
	if path == "" {
		return nil
	}
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	id, err := uuid.Parse(strings.TrimSuffix(base, ext))
	if err != nil || id == uuid.Nil {
		return fmt.Errorf("invalid OCR job file path %q", path)
	}

	var errs []error
	needsDirSync := false
	for _, removePath := range []string{path, ocrJobClaimPath(dir, id)} {
		if err := os.Remove(removePath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			errs = append(errs, fmt.Errorf("remove OCR job file entry %q: %w", removePath, err))
			continue
		}
		needsDirSync = true
	}
	if needsDirSync {
		if err := syncPathFunc(dir); err != nil {
			errs = append(errs, fmt.Errorf("sync OCR job directory after removal: %w", err))
		}
	}
	return errors.Join(errs...)
}

func (h *Handler) CreateOCRJob(c *gin.Context) {
	upload, ok := h.readOCRJobUpload(c)
	if !ok {
		return
	}

	userID, err := parseRequiredUserID(c.PostForm("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	if !h.validateOCRJobUser(c, userID) {
		return
	}

	schema, ok := h.resolveSchema(c)
	if !ok {
		return
	}
	h.queueOCRJobUpload(c, userID, upload, schema)
}

func (h *Handler) CreatePublicOCRJob(c *gin.Context) {
	if rejectPublicQueryUserID(c) {
		return
	}

	upload, ok := h.readOCRJobUpload(c)
	if !ok {
		return
	}
	if rejectPublicFormUserID(c) {
		return
	}

	userID, ok := publicAPIUserID(c)
	if !ok {
		writeError(c, http.StatusInternalServerError, "authenticated user not found")
		return
	}

	schema, ok := h.resolveSchemaScoped(c, &userID)
	if !ok {
		return
	}
	h.queueOCRJobUpload(c, userID, upload, schema)
}

func (h *Handler) readOCRJobUpload(c *gin.Context) (uploadData, bool) {
	limit := h.maxUploadBytes()
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit+multipartOverheadAllowanceBytes)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			writeError(c, http.StatusBadRequest, "request body too large")
			return uploadData{}, false
		}
		writeError(c, http.StatusBadRequest, "file is required")
		return uploadData{}, false
	}
	upload, err := h.readUpload(fileHeader)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return uploadData{}, false
	}
	if utf8.RuneCountInString(upload.Filename) > maxOriginalFilenameCharacters {
		writeError(c, http.StatusBadRequest, "filename must be at most 255 characters")
		return uploadData{}, false
	}
	loggerFromGin(c).Debug("ocr.upload_accepted",
		"mime_type", upload.MimeType,
		"file_size", upload.Size,
		"max_upload_bytes", limit,
	)
	return upload, true
}

func (h *Handler) validateOCRJobUser(c *gin.Context, userID string) bool {
	var count int64
	if err := h.DB.WithContext(c.Request.Context()).Model(&auth.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to validate user")
		return false
	}
	if count == 0 {
		writeError(c, http.StatusBadRequest, "invalid user_id")
		return false
	}
	return true
}

func rejectPublicQueryUserID(c *gin.Context) bool {
	if _, ok := c.GetQuery("user_id"); ok {
		writeError(c, http.StatusBadRequest, "user_id is not allowed")
		return true
	}
	return false
}

func rejectPublicFormUserID(c *gin.Context) bool {
	if _, ok := c.GetPostForm("user_id"); ok {
		writeError(c, http.StatusBadRequest, "user_id is not allowed")
		return true
	}
	return false
}

func (h *Handler) queueOCRJobUpload(c *gin.Context, userID string, upload uploadData, schema resolvedSchema) {
	logger := loggerFromGin(c).With("domain", "ocr", "user_id", userID)
	pageCount, err := countUploadPages(upload.MimeType, upload.Bytes)
	if err != nil {
		logger.Warn("ocr.job_page_count_failed", "mime_type", upload.MimeType, "file_size", upload.Size, "error", safeLogError(err))
		writeError(c, http.StatusBadRequest, "failed to count document pages")
		return
	}
	logger.Debug("ocr.job_page_counted", "mime_type", upload.MimeType, "file_size", upload.Size, "page_count", pageCount)
	if err := validateOCRJobPageLimit(pageCount, schema); err != nil {
		logger.Warn("ocr.job_page_limit_exceeded",
			"mime_type", upload.MimeType,
			"file_size", upload.Size,
			"page_count", pageCount,
			"has_schema", len(schema.Schema) > 0,
		)
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	documentHash, err := computeDocumentHash(upload.Bytes, schema.Schema, schema.Strict)
	if err != nil {
		logger.Warn("ocr.job_document_hash_failed", "error", safeLogError(err))
		writeError(c, http.StatusBadRequest, "invalid schema")
		return
	}
	logger.Debug("ocr.job_document_hash_computed", resolvedSchemaLogAttrs(schema)...)

	fileDir, err := h.ocrJobFileDir()
	if err != nil {
		logger.Error("ocr.job_storage_resolve_failed", "error", safeLogError(err))
		writeError(c, http.StatusInternalServerError, "failed to resolve OCR job file storage")
		return
	}

	jobID := uuid.New()
	userIDPtr := userID
	var job ocr.OCRJob
	var filePath string
	var availableCredits int
	err = h.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := lockOCRJobCreditGate(c.Request.Context(), tx, userID); err != nil {
			return fmt.Errorf("lock OCR job credit gate: %w", err)
		}
		available, err := availableCreditsForNewJob(c.Request.Context(), tx, userID, time.Now().UTC())
		if err != nil {
			return fmt.Errorf("check OCR job credits: %w", err)
		}
		availableCredits = available
		logger.Debug("ocr.job_credit_gate_checked", "required_credits", pageCount, "available_credits", available)
		if available < pageCount {
			return insufficientOCRCreditsError{required: pageCount, available: available}
		}

		filePath, err = writeOCRJobFile(fileDir, jobID, upload.MimeType, upload.Bytes)
		if err != nil {
			return fmt.Errorf("write OCR job file: %w", err)
		}

		job = ocr.OCRJob{
			ID:               jobID,
			UserID:           &userIDPtr,
			OriginalFilename: upload.Filename,
			MimeType:         upload.MimeType,
			FileSize:         upload.Size,
			PageCount:        pageCount,
			DocumentHash:     documentHash,
			FilePath:         filePath,
			SchemaID:         schema.SchemaID,
			Status:           ocr.OCRJobStatusQueued,
		}
		if schema.Inline {
			job.InlineSchemaJSON = datatypes.JSON(schema.Schema)
		}
		if err := tx.Create(&job).Error; err != nil {
			return fmt.Errorf("create OCR job: %w", err)
		}
		return nil
	})
	if err != nil {
		_ = removeOCRJobFile(filePath)
		var insufficient insufficientOCRCreditsError
		if errors.As(err, &insufficient) {
			logger.Warn("ocr.job_queue_rejected_insufficient_credits",
				"required_credits", insufficient.required,
				"available_credits", insufficient.available,
			)
			c.JSON(http.StatusPaymentRequired, PaymentRequiredResponse{
				Error:            "insufficient credits",
				RequiredCredits:  insufficient.required,
				AvailableCredits: insufficient.available,
			})
			return
		}
		logger.Error("ocr.job_queue_failed",
			"job_id", jobID.String(),
			"mime_type", upload.MimeType,
			"file_size", upload.Size,
			"page_count", pageCount,
			"available_credits", availableCredits,
			"error", safeLogError(err),
		)
		writeError(c, http.StatusInternalServerError, "failed to save OCR job")
		return
	}

	notifier := h.OCRJobNotifier
	if notifier == nil {
		notifier = ocr.NotifyOCRJobQueued
	}
	if err := notifier(c.Request.Context(), h.DB, job.ID); err != nil {
		// Notification is an acceleration path only; the executor sweeper will recover.
		logger.Warn("ocr.job_notify_failed", "job_id", job.ID.String(), "error", safeLogError(err))
	} else {
		logger.Debug("ocr.job_notified", "job_id", job.ID.String())
	}
	attrs := append([]any{
		"job_id", job.ID.String(),
		"mime_type", upload.MimeType,
		"file_size", upload.Size,
		"page_count", pageCount,
	}, resolvedSchemaLogAttrs(schema)...)
	logger.Info("ocr.job_queued", attrs...)
	c.JSON(http.StatusAccepted, ocrJobResponse(job))
}

func validateOCRJobPageLimit(pageCount int, schema resolvedSchema) error {
	if len(schema.Schema) > 0 {
		if pageCount > maxOCRJobPagesWithSchema {
			return fmt.Errorf("document must have at most %d pages with a schema", maxOCRJobPagesWithSchema)
		}
		return nil
	}
	if pageCount > maxOCRJobPagesWithoutSchema {
		return fmt.Errorf("document must have at most %d pages without a schema", maxOCRJobPagesWithoutSchema)
	}
	return nil
}

type insufficientOCRCreditsError struct {
	required  int
	available int
}

func (err insufficientOCRCreditsError) Error() string {
	return "insufficient credits"
}

func lockOCRJobCreditGate(ctx context.Context, tx *gorm.DB, userID string) error {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte("syncra:ocr-job-credit-gate:"))
	_, _ = hash.Write([]byte(userID))
	return tx.WithContext(ctx).Exec("SELECT pg_advisory_xact_lock(?)", int64(hash.Sum64())).Error
}

func availableCreditsForNewJob(ctx context.Context, db *gorm.DB, userID string, now time.Time) (int, error) {
	balance, err := billing.AvailableCredits(ctx, db, userID, now)
	if err != nil {
		return 0, err
	}
	var activePages int
	if err := db.WithContext(ctx).
		Model(&ocr.OCRJob{}).
		Select("COALESCE(SUM(page_count), 0)").
		Where("user_id = ? AND status IN ?", userID, []ocr.OCRJobStatus{ocr.OCRJobStatusQueued, ocr.OCRJobStatusProcessing}).
		Scan(&activePages).Error; err != nil {
		return 0, err
	}
	effective := balance.Available - activePages
	if effective < 0 {
		return 0, nil
	}
	return effective, nil
}

func (h *Handler) ListOCRJobs(c *gin.Context) {
	userID, err := parseOptionalUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	status, hasStatus, err := parseOCRJobStatusQuery(c)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	createdFrom, err := parseOCRJobTimeQuery(c, "created_from")
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid created_from")
		return
	}
	createdTo, err := parseOCRJobTimeQuery(c, "created_to")
	if err != nil {
		writeError(c, http.StatusBadRequest, "invalid created_to")
		return
	}
	if createdFrom != nil && createdTo != nil && createdFrom.After(*createdTo) {
		writeError(c, http.StatusBadRequest, "created_from must be before or equal to created_to")
		return
	}

	sortDirection, err := parseOCRJobListSort(c.Query("sort"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	size, err := parseOCRJobListSize(c.Query("size"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}
	cursor, err := parseOCRJobListCursor(c.Query("cursor"), sortDirection)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	query := scopeByUserID(h.DB.WithContext(c.Request.Context()).Preload("Schema"), userID)
	if hasStatus {
		query = query.Where("status = ?", status)
	}
	if createdFrom != nil {
		query = query.Where("created_at >= ?", *createdFrom)
	}
	if createdTo != nil {
		query = query.Where("created_at <= ?", *createdTo)
	}
	if cursor != nil {
		operator := "<"
		if sortDirection == "asc" {
			operator = ">"
		}
		query = query.Where("(created_at, id) "+operator+" (?, ?)", cursor.CreatedAt, cursor.ID)
	}

	order := "created_at desc, id desc"
	if sortDirection == "asc" {
		order = "created_at asc, id asc"
	}
	var jobs []ocr.OCRJob
	if err := query.Order(order).Limit(size + 1).Find(&jobs).Error; err != nil {
		writeError(c, http.StatusInternalServerError, "failed to list OCR jobs")
		return
	}

	var nextCursor *string
	if len(jobs) > size {
		jobs = jobs[:size]
		if len(jobs) > 0 {
			encoded, err := encodeOCRJobListCursor(jobs[len(jobs)-1], sortDirection)
			if err != nil {
				writeError(c, http.StatusInternalServerError, "failed to encode next cursor")
				return
			}
			nextCursor = &encoded
		}
	}

	out := make([]OCRJobListItemResponse, 0, len(jobs))
	for _, job := range jobs {
		out = append(out, ocrJobListItemResponse(job))
	}
	loggerFromGin(c).Debug("ocr.jobs_listed", "result_count", len(out), "has_next_cursor", nextCursor != nil)
	c.JSON(http.StatusOK, OCRJobListResponse{Jobs: out, NextCursor: nextCursor})
}

func parseOCRJobStatusQuery(c *gin.Context) (ocr.OCRJobStatus, bool, error) {
	raw, ok := c.GetQuery("status")
	if !ok {
		return "", false, nil
	}
	status := ocr.OCRJobStatus(strings.TrimSpace(raw))
	switch status {
	case ocr.OCRJobStatusQueued, ocr.OCRJobStatusProcessing, ocr.OCRJobStatusCompleted, ocr.OCRJobStatusFailed:
		return status, true, nil
	default:
		return "", false, errors.New("invalid status")
	}
}

func parseOCRJobTimeQuery(c *gin.Context, name string) (*time.Time, error) {
	raw, ok := c.GetQuery(name)
	if !ok {
		return nil, nil
	}
	value, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(raw))
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func parseOCRJobListSort(raw string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "desc":
		return "desc", nil
	case "asc":
		return "asc", nil
	default:
		return "", errors.New("sort must be asc or desc")
	}
}

func parseOCRJobListSize(raw string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return defaultOCRJobListSize, nil
	}
	size, err := strconv.Atoi(raw)
	if err != nil || size < 1 || size > maxOCRJobListSize {
		return 0, fmt.Errorf("size must be between 1 and %d", maxOCRJobListSize)
	}
	return size, nil
}

func parseOCRJobListCursor(raw string, sortDirection string) (*ocrJobListCursor, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}
	var cursor ocrJobListCursor
	if err := json.Unmarshal(decoded, &cursor); err != nil {
		return nil, errors.New("invalid cursor")
	}
	if cursor.ID == uuid.Nil || cursor.CreatedAt.IsZero() || (cursor.Sort != "asc" && cursor.Sort != "desc") {
		return nil, errors.New("invalid cursor")
	}
	if cursor.Sort != sortDirection {
		return nil, errors.New("cursor sort does not match sort")
	}
	return &cursor, nil
}

func encodeOCRJobListCursor(job ocr.OCRJob, sortDirection string) (string, error) {
	raw, err := json.Marshal(ocrJobListCursor{
		CreatedAt: job.CreatedAt.UTC(),
		ID:        job.ID,
		Sort:      sortDirection,
	})
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func (h *Handler) GetOCRJob(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil || id == uuid.Nil {
		writeError(c, http.StatusBadRequest, "invalid OCR job id")
		return
	}

	rawUserID, hasUserID := c.GetQuery("user_id")
	userID, err := parseOptionalUserID(rawUserID)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	h.getOCRJob(c, id, userID, hasUserID)
}

func (h *Handler) GetPublicOCRJob(c *gin.Context) {
	if rejectPublicQueryUserID(c) {
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil || id == uuid.Nil {
		writeError(c, http.StatusBadRequest, "invalid OCR job id")
		return
	}
	userID, ok := publicAPIUserID(c)
	if !ok {
		writeError(c, http.StatusInternalServerError, "authenticated user not found")
		return
	}
	h.getPublicOCRJob(c, id, userID)
}

func (h *Handler) getOCRJob(c *gin.Context, id uuid.UUID, userID *string, scoped bool) {
	job, ok := h.loadOCRJob(c, id, userID, scoped)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, ocrJobResponse(job))
}

func (h *Handler) getPublicOCRJob(c *gin.Context, id uuid.UUID, userID string) {
	job, ok := h.loadOCRJob(c, id, &userID, true)
	if !ok {
		return
	}
	doc, ok := h.loadPublicOCRJobDocument(c, job, userID)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, publicOCRJobResponse(job, doc))
}

func (h *Handler) loadOCRJob(c *gin.Context, id uuid.UUID, userID *string, scoped bool) (ocr.OCRJob, bool) {
	var job ocr.OCRJob
	query := h.DB.WithContext(c.Request.Context()).Preload("Schema").Where("id = ?", id)
	if scoped {
		query = scopeByUserID(query, userID)
	}
	if err := query.First(&job).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeError(c, http.StatusNotFound, "OCR job not found")
			return ocr.OCRJob{}, false
		}
		writeError(c, http.StatusInternalServerError, "failed to load OCR job")
		return ocr.OCRJob{}, false
	}
	return job, true
}

func (h *Handler) loadPublicOCRJobDocument(c *gin.Context, job ocr.OCRJob, userID string) (*ocr.OCRDocument, bool) {
	if job.DocumentID == nil {
		return nil, true
	}
	var doc ocr.OCRDocument
	query := scopeByUserID(h.DB.WithContext(c.Request.Context()).Where("id = ?", *job.DocumentID), &userID)
	result := query.Limit(1).Find(&doc)
	if result.Error != nil {
		writeError(c, http.StatusInternalServerError, "failed to load OCR document")
		return nil, false
	}
	if result.RowsAffected == 0 {
		return nil, true
	}
	return &doc, true
}

func (h *Handler) DeleteOCRJobs(c *gin.Context) {
	userID, err := parseOptionalUserID(c.Query("user_id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	var req DeleteOCRJobsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid OCR job delete request")
		return
	}

	ids, err := parseOCRJobIDs(req.IDs)
	if err != nil {
		writeError(c, http.StatusBadRequest, err.Error())
		return
	}

	deletedIDs, err := h.deleteOCRJobs(c.Request.Context(), ids, userID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to delete OCR jobs")
		return
	}

	c.JSON(http.StatusOK, DeleteOCRJobsResponse{
		DeletedIDs:   deletedIDs,
		DeletedCount: len(deletedIDs),
	})
	loggerFromGin(c).Info("ocr.jobs_deleted", "requested_count", len(ids), "deleted_count", len(deletedIDs))
}

func parseOCRJobIDs(rawIDs []string) ([]uuid.UUID, error) {
	if len(rawIDs) == 0 {
		return nil, errors.New("ids is required")
	}
	if len(rawIDs) > maxOCRJobDeleteIDs {
		return nil, fmt.Errorf("ids must include at most %d OCR jobs", maxOCRJobDeleteIDs)
	}

	ids := make([]uuid.UUID, 0, len(rawIDs))
	seen := make(map[uuid.UUID]struct{}, len(rawIDs))
	for _, rawID := range rawIDs {
		id, err := uuid.Parse(strings.TrimSpace(rawID))
		if err != nil || id == uuid.Nil {
			return nil, errors.New("invalid OCR job id")
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}

func (h *Handler) deleteOCRJobs(ctx context.Context, ids []uuid.UUID, userID *string) ([]uuid.UUID, error) {
	if len(ids) == 0 {
		return []uuid.UUID{}, nil
	}

	query := scopeByUserID(h.DB.WithContext(ctx).Model(&ocr.OCRJob{}).Where("id IN ?", ids), userID)
	var matchedIDs []uuid.UUID
	if err := query.Pluck("id", &matchedIDs).Error; err != nil {
		return nil, err
	}
	if len(matchedIDs) == 0 {
		return []uuid.UUID{}, nil
	}

	if err := scopeByUserID(h.DB.WithContext(ctx).Where("id IN ?", matchedIDs), userID).
		Delete(&ocr.OCRJob{}).Error; err != nil {
		return nil, err
	}

	matched := make(map[uuid.UUID]struct{}, len(matchedIDs))
	for _, id := range matchedIDs {
		matched[id] = struct{}{}
	}

	deletedIDs := make([]uuid.UUID, 0, len(matchedIDs))
	for _, id := range ids {
		if _, ok := matched[id]; ok {
			deletedIDs = append(deletedIDs, id)
		}
	}
	return deletedIDs, nil
}

type ocrJobFileWrite struct {
	dir          string
	tmpPath      string
	finalPath    string
	claimPath    string
	finalLinked  bool
	claimCreated bool
}

func (w *ocrJobFileWrite) linkFinalFile() error {
	if err := os.Link(w.tmpPath, w.finalPath); err != nil {
		return err
	}
	w.finalLinked = true
	return nil
}

func (w *ocrJobFileWrite) cleanup() error {
	var errs []error
	needsDirSync := false
	if w.tmpPath != "" {
		needsDirSync = true
		if err := os.Remove(w.tmpPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			errs = append(errs, fmt.Errorf("remove temp OCR job file: %w", err))
		}
	}
	if w.finalLinked {
		needsDirSync = true
		if err := os.Remove(w.finalPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			errs = append(errs, fmt.Errorf("remove final OCR job file: %w", err))
		}
	}
	if w.claimCreated {
		needsDirSync = true
		if err := os.Remove(w.claimPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			errs = append(errs, fmt.Errorf("remove OCR job claim file: %w", err))
		}
	}
	if needsDirSync {
		if err := syncPathFunc(w.dir); err != nil {
			errs = append(errs, fmt.Errorf("sync OCR job directory after cleanup: %w", err))
		}
	}
	return errors.Join(errs...)
}

func ensureDirDurable(path string) error {
	clean, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	clean = filepath.Clean(clean)
	volume := filepath.VolumeName(clean)
	rest := strings.TrimPrefix(clean, volume)
	separator := string(os.PathSeparator)
	current := volume + separator
	parts := strings.Split(strings.Trim(rest, separator), separator)

	for _, part := range parts {
		if part == "" {
			continue
		}
		next := filepath.Join(current, part)
		info, err := os.Stat(next)
		if err == nil {
			if !info.IsDir() {
				return fmt.Errorf("directory path component %q is not a directory", next)
			}
			current = next
			continue
		}
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("inspect directory path component %q: %w", next, err)
		}
		if err := os.Mkdir(next, 0o750); err != nil {
			if errors.Is(err, os.ErrExist) {
				info, statErr := os.Stat(next)
				if statErr != nil {
					return fmt.Errorf("inspect directory path component %q: %w", next, statErr)
				}
				if !info.IsDir() {
					return fmt.Errorf("directory path component %q is not a directory", next)
				}
				current = next
				continue
			}
			return fmt.Errorf("create directory %q: %w", next, err)
		}
		if err := syncPathFunc(current); err != nil {
			return fmt.Errorf("sync parent directory %q: %w", current, err)
		}
		current = next
	}
	return nil
}

func syncPath(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func syncAndClose(file *os.File) error {
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func createOCRJobClaim(dir string, id uuid.UUID, claimPath string) (*os.File, error) {
	if existingPath, ok, err := existingOCRJobFinalFilePath(dir, id); err != nil {
		return nil, err
	} else if ok {
		return nil, fmt.Errorf("OCR job file already exists: %s", existingPath)
	}

	claim, err := os.OpenFile(claimPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o640)
	if err == nil {
		return claim, nil
	}
	if errors.Is(err, os.ErrExist) {
		return nil, errors.New("OCR job file already exists")
	}
	return nil, err
}

func ocrJobClaimPath(dir string, id uuid.UUID) string {
	return filepath.Join(dir, id.String()+".claim")
}

func existingOCRJobFinalFilePath(dir string, id uuid.UUID) (string, bool, error) {
	for _, ext := range supportedOCRJobFileExtensions() {
		path := filepath.Join(dir, id.String()+ext)
		if _, err := os.Stat(path); err == nil {
			return path, true, nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", false, err
		}
	}
	return "", false, nil
}

func supportedOCRJobFileExtensions() []string {
	return []string{".pdf", ".png", ".jpg"}
}

func ocrJobFileExtension(mimeType string) (string, bool) {
	switch mimeType {
	case "application/pdf":
		return ".pdf", true
	case "image/png":
		return ".png", true
	case "image/jpeg":
		return ".jpg", true
	default:
		return "", false
	}
}
