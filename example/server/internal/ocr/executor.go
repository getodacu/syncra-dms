package ocr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ai.ro/syncra/internal/billing"
	"ai.ro/syncra/internal/logging"
	"ai.ro/syncra/internal/webhooks"
)

const (
	defaultStaleProcessingTimeout = 30 * time.Minute
	defaultWebhookDispatchSlots   = 8
	maxOCRJobErrorMessageBytes    = 500
)

type ExecutorConfig struct {
	DB                     *gorm.DB
	DSN                    string
	Processor              Processor
	WebhookDispatcher      WebhookDispatcher
	WebhookDispatchSlots   int
	NotifyChannel          string
	WorkerCount            int
	QueueBuffer            int
	PollInterval           time.Duration
	StaleProcessingTimeout time.Duration
	Logger                 *slog.Logger
}

type WebhookDispatcher interface {
	DispatchJobEvent(context.Context, webhooks.JobEventInput) error
}

type Executor struct {
	db                     *gorm.DB
	dsn                    string
	processor              Processor
	webhookDispatcher      WebhookDispatcher
	webhookDispatchSlots   chan struct{}
	notifyChannel          string
	workerCount            int
	pollInterval           time.Duration
	staleProcessingTimeout time.Duration
	jobs                   chan uuid.UUID
	logger                 *slog.Logger
	listenerReady          func()
}

type jobSchema struct {
	Schema   json.RawMessage
	Strict   bool
	Inline   bool
	SchemaID *uuid.UUID
}

func NewExecutor(cfg ExecutorConfig) *Executor {
	workerCount := cfg.WorkerCount
	if workerCount <= 0 {
		workerCount = 2
	}
	queueBuffer := cfg.QueueBuffer
	if queueBuffer <= 0 {
		queueBuffer = workerCount * 4
	}
	pollInterval := cfg.PollInterval
	if pollInterval <= 0 {
		pollInterval = 10 * time.Second
	}
	staleProcessingTimeout := cfg.StaleProcessingTimeout
	if staleProcessingTimeout <= 0 {
		staleProcessingTimeout = defaultStaleProcessingTimeout
	}
	webhookDispatchSlotCount := cfg.WebhookDispatchSlots
	if webhookDispatchSlotCount <= 0 {
		webhookDispatchSlotCount = defaultWebhookDispatchSlots
	}
	logger := cfg.Logger
	if logger == nil {
		logger = logging.Nop()
	}
	logger = logger.With("component", "ocr_executor")
	notifyChannel := strings.TrimSpace(cfg.NotifyChannel)
	if notifyChannel == "" {
		notifyChannel = OCRJobsNotifyChannel
	}
	return &Executor{
		db:                     cfg.DB,
		dsn:                    cfg.DSN,
		processor:              cfg.Processor,
		webhookDispatcher:      cfg.WebhookDispatcher,
		webhookDispatchSlots:   make(chan struct{}, webhookDispatchSlotCount),
		notifyChannel:          notifyChannel,
		workerCount:            workerCount,
		pollInterval:           pollInterval,
		staleProcessingTimeout: staleProcessingTimeout,
		jobs:                   make(chan uuid.UUID, queueBuffer),
		logger:                 logger,
	}
}

func (e *Executor) Run(ctx context.Context) error {
	if e.db == nil {
		return errors.New("OCR executor database is not configured")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	e.logger.Info("ocr.executor_start",
		"workers", e.workerCount,
		"queue_buffer", cap(e.jobs),
		"poll_interval_ms", e.pollInterval.Milliseconds(),
		"stale_processing_timeout_ms", e.staleProcessingTimeout.Milliseconds(),
		"notify_channel", e.notifyChannel,
		"listener_enabled", strings.TrimSpace(e.dsn) != "",
	)
	var wg sync.WaitGroup
	for i := 0; i < e.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			e.runWorker(ctx, workerID)
		}(i + 1)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.runSweeper(ctx)
	}()
	if strings.TrimSpace(e.dsn) != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.runListener(ctx)
		}()
	}
	<-ctx.Done()
	wg.Wait()
	e.logger.Info("ocr.executor_stop")
	return nil
}

func (e *Executor) enqueue(id uuid.UUID) {
	if id == uuid.Nil {
		e.logger.Warn("ocr.job_enqueue_dropped", "reason", "nil_id")
		return
	}
	select {
	case e.jobs <- id:
		e.jobLogger(id).Debug("ocr.job_enqueued", "queue_depth", len(e.jobs), "queue_capacity", cap(e.jobs))
	default:
		e.jobLogger(id).Warn("ocr.job_enqueue_dropped", "reason", "queue_full", "queue_capacity", cap(e.jobs))
	}
}

func (e *Executor) runWorker(ctx context.Context, workerID int) {
	e.logger.Debug("ocr.worker_start", "worker_id", workerID)
	defer e.logger.Debug("ocr.worker_stop", "worker_id", workerID)
	for {
		select {
		case <-ctx.Done():
			return
		case id := <-e.jobs:
			e.processJobFromWorker(ctx, id, workerID)
		}
	}
}

func (e *Executor) processJobFromWorker(ctx context.Context, id uuid.UUID, workerID int) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err := fmt.Errorf("OCR job panic: %v", recovered)
			e.jobLogger(id).Error("ocr.job_panic", "worker_id", workerID, "error", err)
			if _, failErr := e.failJob(context.Background(), id, err); failErr != nil {
				e.jobLogger(id).Error("ocr.job_panic_fail_update_failed", "worker_id", workerID, "error", failErr)
			}
		}
	}()
	e.ProcessJob(ctx, id)
}

func (e *Executor) runSweeper(ctx context.Context) {
	e.logger.Debug("ocr.sweeper_start")
	defer e.logger.Debug("ocr.sweeper_stop")
	ticker := time.NewTicker(e.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			e.sweepQueuedJobs(ctx)
		}
	}
}

func (e *Executor) sweepQueuedJobs(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}
	if e.db == nil {
		e.logger.Error("ocr.sweep_failed", "error", "database is not configured")
		return
	}
	limit := cap(e.jobs)
	if limit <= 0 {
		limit = e.workerCount
	}
	if limit <= 0 {
		limit = 1
	}
	e.logger.Debug("ocr.sweep_started", "limit", limit)
	if err := e.requeueStaleProcessingJobs(ctx, limit); err != nil {
		if ctx.Err() == nil {
			e.logger.Error("ocr.stale_jobs_requeue_failed", "error", err)
		}
		return
	}
	if ctx.Err() != nil {
		return
	}
	var ids []uuid.UUID
	err := e.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Model(&OCRJob{}).
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ?", OCRJobStatusQueued).
			Order("created_at asc").
			Limit(limit).
			Pluck("id", &ids).
			Error
	})
	if err != nil {
		if ctx.Err() == nil {
			e.logger.Error("ocr.sweep_failed", "error", err)
		}
		return
	}
	e.logger.Debug("ocr.sweep_completed", "queued_jobs_found", len(ids), "limit", limit)
	for _, id := range ids {
		e.enqueue(id)
	}
}

func (e *Executor) requeueStaleProcessingJobs(ctx context.Context, limit int) error {
	if e.staleProcessingTimeout <= 0 {
		return nil
	}
	cutoff := time.Now().Add(-e.staleProcessingTimeout)
	var ids []uuid.UUID
	err := e.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&OCRJob{}).
			Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("status = ? AND updated_at < ?", OCRJobStatusProcessing, cutoff).
			Order("updated_at asc").
			Limit(limit).
			Pluck("id", &ids).
			Error; err != nil {
			return err
		}
		if len(ids) == 0 {
			return nil
		}
		result := tx.Model(&OCRJob{}).
			Where("id IN ? AND status = ?", ids, OCRJobStatusProcessing).
			Updates(map[string]any{
				"status":        OCRJobStatusQueued,
				"error_message": "",
				"document_id":   nil,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != int64(len(ids)) {
			return fmt.Errorf("requeue stale OCR jobs affected %d rows, want %d", result.RowsAffected, len(ids))
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, id := range ids {
		e.jobLogger(id).Info("ocr.job_requeued_stale")
	}
	return nil
}

func (e *Executor) runListener(ctx context.Context) {
	e.logger.Debug("ocr.listener_start", "notify_channel", e.notifyChannel)
	defer e.logger.Debug("ocr.listener_stop", "notify_channel", e.notifyChannel)
	backoff := 100 * time.Millisecond
	for ctx.Err() == nil {
		err := e.listenOnce(ctx)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			e.logger.Error("ocr.listener_error", "error", err, "backoff_ms", backoff.Milliseconds())
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}
		if backoff < 5*time.Second {
			backoff *= 2
			if backoff > 5*time.Second {
				backoff = 5 * time.Second
			}
		}
	}
}

func (e *Executor) listenOnce(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, e.dsn)
	if err != nil {
		return err
	}
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			e.logger.Warn("ocr.listener_close_failed", "error", err)
		}
	}()
	if _, err := conn.Exec(ctx, "LISTEN "+pgx.Identifier{e.notifyChannel}.Sanitize()); err != nil {
		return err
	}
	e.logger.Info("ocr.listener_ready", "notify_channel", e.notifyChannel)
	if e.listenerReady != nil {
		e.listenerReady()
	}
	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			return err
		}
		id, err := uuid.Parse(notification.Payload)
		if err != nil || id == uuid.Nil {
			e.logger.Warn("ocr.notification_dropped", "reason", "invalid_uuid")
			continue
		}
		e.jobLogger(id).Debug("ocr.notification_received", "notify_channel", notification.Channel)
		e.enqueue(id)
	}
}

func (e *Executor) ProcessJob(ctx context.Context, id uuid.UUID) {
	if ctx == nil {
		ctx = context.Background()
	}
	logger := e.jobLogger(id)
	start := time.Now()
	logger.Debug("ocr.job_process_started")
	job, claimed, err := e.claimJob(ctx, id)
	if err != nil {
		logger.Error("ocr.job_claim_failed", "error", ocrLogError(err))
		return
	}
	if !claimed {
		logger.Debug("ocr.job_claim_skipped")
		return
	}
	logger = e.jobLogger(job.ID).With(
		"status", string(job.Status),
		"mime_type", job.MimeType,
		"file_size", job.FileSize,
		"page_count", job.PageCount,
	)
	if job.UserID != nil {
		logger = logger.With("user_id", *job.UserID)
	}
	logger.Info("ocr.job_claimed")
	e.dispatchJobStarted(context.Background(), job)
	doc, err := e.executeClaimedJob(ctx, job)
	if err != nil {
		if isOCRJobLifecycleInterrupt(err) || ctx.Err() != nil {
			if requeueErr := e.requeueJob(context.Background(), job.ID); requeueErr != nil {
				logger.Error("ocr.job_requeue_failed", "processing_error", ocrLogError(err), "error", ocrLogError(requeueErr))
				return
			}
			logger.Info("ocr.job_requeued", "reason", lifecycleInterruptReason(err, ctx), "duration_ms", time.Since(start).Milliseconds())
			return
		}
		errorMessage, failErr := e.failJob(context.Background(), job.ID, err)
		if failErr != nil {
			logger.Error("ocr.job_fail_update_failed", "processing_error", ocrLogError(err), "error", ocrLogError(failErr))
			return
		}
		job.Status = OCRJobStatusFailed
		job.ErrorMessage = errorMessage
		e.dispatchJobFailed(context.Background(), job, errorMessage)
		logger.Error("ocr.job_failed", "error", ocrLogError(err), "duration_ms", time.Since(start).Milliseconds())
		return
	}
	job.Status = OCRJobStatusCompleted
	job.ErrorMessage = ""
	if doc != nil {
		job.DocumentID = &doc.ID
	}
	e.dispatchJobSucceeded(context.Background(), job)
	logger.Info("ocr.job_completed", "duration_ms", time.Since(start).Milliseconds())
}

func (e *Executor) claimJob(ctx context.Context, id uuid.UUID) (*OCRJob, bool, error) {
	if e.db == nil {
		return nil, false, errors.New("OCR executor database is not configured")
	}
	var job OCRJob
	err := e.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
			Where("id = ? AND status = ?", id, OCRJobStatusQueued).
			Limit(1).
			Find(&job)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}
		result = tx.Model(&OCRJob{}).
			Where("id = ? AND status = ?", job.ID, OCRJobStatusQueued).
			Updates(map[string]any{
				"status":        OCRJobStatusProcessing,
				"error_message": "",
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return fmt.Errorf("claim OCR job affected %d rows", result.RowsAffected)
		}
		job.Status = OCRJobStatusProcessing
		job.ErrorMessage = ""
		return nil
	})
	if err != nil {
		return nil, false, err
	}
	if job.ID == uuid.Nil {
		return nil, false, nil
	}
	return &job, true, nil
}

func (e *Executor) executeClaimedJob(ctx context.Context, job *OCRJob) (*OCRDocument, error) {
	if job == nil {
		return nil, errors.New("OCR job is required")
	}
	if e.processor == nil {
		return nil, errors.New("OCR processor is not configured")
	}
	logger := e.jobLogger(job.ID)
	data, err := os.ReadFile(job.FilePath)
	if err != nil {
		return nil, fmt.Errorf("read OCR job file: %w", err)
	}
	logger.Debug("ocr.job_file_read", "storage", "ocr_job_file", "file_bytes", len(data))
	schema, err := e.resolveJobSchema(ctx, job)
	if err != nil {
		return nil, err
	}
	logger.Debug("ocr.job_schema_resolved", jobSchemaLogAttrs(schema)...)
	processorStart := time.Now()
	response, rawResponse, err := e.processor(ctx, ProcessInput{
		Filename: job.OriginalFilename,
		MimeType: job.MimeType,
		DataURL:  DataURL(job.MimeType, data),
		Schema:   schema.Schema,
		Strict:   schema.Strict,
	})
	if err != nil {
		logger.Error("ocr.upstream_request_failed", "duration_ms", time.Since(processorStart).Milliseconds(), "error", ocrLogError(err))
		return nil, err
	}
	logger.Debug("ocr.upstream_request_completed", "duration_ms", time.Since(processorStart).Milliseconds(), "response_bytes", len(rawResponse))
	if response == nil {
		return nil, errors.New("OCR processor returned nil response")
	}
	if len(rawResponse) == 0 {
		return nil, errors.New("OCR processor returned empty raw response")
	}
	pageCount, err := CountRawResponsePages(rawResponse)
	if err != nil {
		return nil, fmt.Errorf("count OCR response pages: %w", err)
	}
	logger.Debug("ocr.raw_response_pages_counted", "page_count", pageCount)
	annotation, err := ParseAnnotationJSON(response.DocumentAnnotation, len(schema.Schema) > 0)
	if err != nil {
		return nil, fmt.Errorf("parse OCR annotation JSON: %w", err)
	}
	logger.Debug("ocr.annotation_parsed", "has_annotation", len(annotation) > 0)
	jobID := job.ID
	doc := OCRDocument{
		UserID:           job.UserID,
		JobID:            &jobID,
		OriginalFilename: job.OriginalFilename,
		MimeType:         job.MimeType,
		FileSize:         job.FileSize,
		PageCount:        pageCount,
		DocumentHash:     job.DocumentHash,
		SchemaID:         schema.SchemaID,
		Markdown:         JoinMarkdown(response.Pages),
		AnnotationJSON:   datatypes.JSON(annotation),
		RawResponseJSON:  datatypes.JSON(rawResponse),
	}
	if schema.Inline {
		doc.InlineSchemaJSON = datatypes.JSON(schema.Schema)
	}
	if err := e.completeJob(ctx, job.ID, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func isOCRJobLifecycleInterrupt(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func (e *Executor) resolveJobSchema(ctx context.Context, job *OCRJob) (jobSchema, error) {
	if job == nil {
		return jobSchema{}, errors.New("OCR job is required")
	}
	if job.SchemaID != nil {
		var schema ExtractionSchema
		if err := e.db.WithContext(ctx).First(&schema, "id = ?", *job.SchemaID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return jobSchema{}, errors.New("OCR job schema not found")
			}
			return jobSchema{}, fmt.Errorf("load OCR job schema: %w", err)
		}
		return jobSchema{
			Schema:   json.RawMessage(schema.SchemaJSON),
			Strict:   schema.Strict,
			SchemaID: &schema.ID,
		}, nil
	}
	if len(job.InlineSchemaJSON) > 0 {
		return jobSchema{
			Schema: json.RawMessage(job.InlineSchemaJSON),
			Strict: true,
			Inline: true,
		}, nil
	}
	return jobSchema{}, nil
}

func (e *Executor) completeJob(ctx context.Context, id uuid.UUID, doc *OCRDocument) error {
	if e.db == nil {
		return errors.New("OCR executor database is not configured")
	}
	if doc == nil {
		return errors.New("OCR document is required")
	}
	logger := e.jobLogger(id)
	if doc.UserID != nil {
		logger = logger.With("user_id", *doc.UserID)
	}
	err := e.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(doc).Error; err != nil {
			return fmt.Errorf("create OCR document: %w", err)
		}
		logger.Debug("ocr.document_created",
			"document_id", doc.ID.String(),
			"page_count", doc.PageCount,
			"file_size", doc.FileSize,
			"mime_type", doc.MimeType,
		)
		if err := LinkDocumentToMatchingCollections(ctx, tx, *doc); err != nil {
			return fmt.Errorf("link OCR document collections: %w", err)
		}
		if doc.UserID != nil && doc.PageCount > 0 {
			if err := billing.DebitCreditsForJobTx(ctx, tx, billing.DebitCreditsInput{
				UserID:         *doc.UserID,
				RelatedJobID:   id,
				Credits:        doc.PageCount,
				IdempotencyKey: "ocr_job_completed:" + id.String(),
				Now:            time.Now().UTC(),
			}); err != nil {
				return fmt.Errorf("debit OCR job credits: %w", err)
			}
			logger.Debug("ocr.job_credits_debited", "credits", doc.PageCount)
		}
		result := tx.Model(&OCRJob{}).
			Where("id = ? AND status = ?", id, OCRJobStatusProcessing).
			Updates(map[string]any{
				"status":        OCRJobStatusCompleted,
				"document_id":   doc.ID,
				"error_message": "",
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return fmt.Errorf("complete OCR job affected %d rows", result.RowsAffected)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if e.logger.Enabled(ctx, slog.LevelDebug) {
		var linkCount int64
		if countErr := e.db.WithContext(ctx).Model(&CollectionDocument{}).Where("document_id = ?", doc.ID).Count(&linkCount).Error; countErr != nil {
			logger.Warn("ocr.document_collection_link_count_failed", "document_id", doc.ID.String(), "error", countErr)
		} else {
			logger.Debug("ocr.document_collections_linked", "document_id", doc.ID.String(), "collection_link_count", linkCount)
		}
	}
	return nil
}

func (e *Executor) failJob(ctx context.Context, id uuid.UUID, cause error) (string, error) {
	if e.db == nil {
		return "", errors.New("OCR executor database is not configured")
	}
	message := boundedOCRJobErrorMessage(cause)
	result := e.db.WithContext(ctx).
		Model(&OCRJob{}).
		Where("id = ? AND status = ?", id, OCRJobStatusProcessing).
		Updates(map[string]any{
			"status":        OCRJobStatusFailed,
			"error_message": message,
		})
	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected != 1 {
		return "", fmt.Errorf("fail OCR job affected %d rows", result.RowsAffected)
	}
	return message, nil
}

func (e *Executor) requeueJob(ctx context.Context, id uuid.UUID) error {
	if e.db == nil {
		return errors.New("OCR executor database is not configured")
	}
	result := e.db.WithContext(ctx).
		Model(&OCRJob{}).
		Where("id = ? AND status = ?", id, OCRJobStatusProcessing).
		Updates(map[string]any{
			"status":        OCRJobStatusQueued,
			"error_message": "",
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return fmt.Errorf("requeue OCR job affected %d rows", result.RowsAffected)
	}
	return nil
}

func (e *Executor) dispatchJobStarted(ctx context.Context, job *OCRJob) {
	e.dispatchJobWebhook(ctx, webhooks.EventJobStarted, job, OCRJobStatusProcessing, nil)
}

func (e *Executor) dispatchJobSucceeded(ctx context.Context, job *OCRJob) {
	if job == nil {
		return
	}
	e.dispatchJobWebhook(ctx, webhooks.EventJobSucceeded, job, OCRJobStatusCompleted, nil)
}

func (e *Executor) dispatchJobFailed(ctx context.Context, job *OCRJob, errorMessage string) {
	e.dispatchJobWebhook(ctx, webhooks.EventJobFailed, job, OCRJobStatusFailed, &errorMessage)
}

func (e *Executor) dispatchJobWebhook(ctx context.Context, event webhooks.Event, job *OCRJob, status OCRJobStatus, errorMessage *string) {
	if e.webhookDispatcher == nil || job == nil || job.UserID == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	userID := *job.UserID
	jobPayload := webhooks.JobPayload{ID: job.ID.String()}
	if event != webhooks.EventJobSucceeded {
		jobPayload.Status = string(status)
		jobPayload.OriginalFilename = job.OriginalFilename
	}
	input := webhooks.JobEventInput{
		Event:  event,
		UserID: &userID,
		Job:    jobPayload,
	}
	if errorMessage != nil {
		message := *errorMessage
		input.Job.ErrorMessage = &message
	}
	if e.webhookDispatchSlots == nil {
		e.dispatchJobWebhookNow(ctx, job.ID, input)
		return
	}
	select {
	case e.webhookDispatchSlots <- struct{}{}:
	default:
		e.jobLogger(job.ID).Warn("ocr.webhook_dispatch_dropped",
			"event", string(event),
			"reason", "dispatcher_busy",
		)
		return
	}
	go func() {
		defer func() { <-e.webhookDispatchSlots }()
		e.dispatchJobWebhookNow(ctx, job.ID, input)
	}()
}

func (e *Executor) dispatchJobWebhookNow(ctx context.Context, jobID uuid.UUID, input webhooks.JobEventInput) {
	if err := e.webhookDispatcher.DispatchJobEvent(ctx, input); err != nil {
		e.jobLogger(jobID).Warn("ocr.webhook_dispatch_failed",
			"event", string(input.Event),
			"error", sanitizedWebhookDispatchError(err),
		)
	}
}

func (e *Executor) jobLogger(id uuid.UUID) *slog.Logger {
	logger := e.logger
	if logger == nil {
		logger = logging.Nop()
	}
	if id != uuid.Nil {
		logger = logger.With("job_id", id.String())
	}
	return logger
}

func jobSchemaLogAttrs(schema jobSchema) []any {
	source := "none"
	switch {
	case schema.SchemaID != nil:
		source = "saved"
	case schema.Inline:
		source = "inline"
	}
	attrs := []any{
		"schema_source", source,
		"has_schema", len(schema.Schema) > 0,
		"strict", schema.Strict,
	}
	if schema.SchemaID != nil {
		attrs = append(attrs, "schema_id", schema.SchemaID.String())
	}
	return attrs
}

func lifecycleInterruptReason(err error, ctx context.Context) string {
	switch {
	case ctx != nil && ctx.Err() != nil:
		return "context_done"
	case errors.Is(err, context.Canceled):
		return "context_canceled"
	case errors.Is(err, context.DeadlineExceeded):
		return "context_deadline_exceeded"
	default:
		return "lifecycle_interrupt"
	}
}

func ocrLogError(err error) string {
	if err == nil {
		return ""
	}
	message := err.Error()
	if strings.Contains(message, "read OCR job file") {
		return "read OCR job file"
	}
	return message
}

func boundedOCRJobErrorMessage(cause error) string {
	message := strings.TrimSpace(fmt.Sprint(cause))
	if message == "" || message == "<nil>" {
		message = "OCR job failed"
	}
	if len(message) <= maxOCRJobErrorMessageBytes {
		return message
	}
	return message[:maxOCRJobErrorMessageBytes-3] + "..."
}

func sanitizedWebhookDispatchError(err error) string {
	if err == nil {
		return ""
	}
	switch {
	case errors.Is(err, context.Canceled):
		return context.Canceled.Error()
	case errors.Is(err, context.DeadlineExceeded):
		return context.DeadlineExceeded.Error()
	default:
		return "webhook dispatch failed"
	}
}
