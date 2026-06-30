package ocr

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/auth"
	"ai.ro/syncra/internal/billing"
	"ai.ro/syncra/internal/logging"
	"ai.ro/syncra/internal/testsupport"
	"ai.ro/syncra/internal/webhooks"
)

func TestExecutorProcessesQueuedJob(t *testing.T) {
	db := executorTestDB(t)
	user := createExecutorTestUser(t, db)
	grantExecutorTestCredits(t, db, user.ID, 5)
	fileData := []byte("invoice bytes")
	job := createQueuedExecutorJob(t, db, executorJobFixture{
		UserID:   &user.ID,
		FileData: fileData,
	})
	annotation := `{"invoice_id":"INV-001"}`
	rawResponse := []byte(`{"pages":[{"index":0,"markdown":"# Invoice"},{"index":1,"markdown":"Raw page two"}],"document_annotation":{"invoice_id":"INV-001"}}`)
	processor := &recordingExecutorProcessor{
		response: &MistralResponse{
			Pages: []MistralPage{
				{
					Index:    0,
					Markdown: "# Invoice",
					Images:   json.RawMessage(`[{"id":"invoice.png","image_base64":"iVBORw0KGgo="}]`),
				},
			},
			DocumentAnnotation: &annotation,
		},
		raw: rawResponse,
	}

	NewExecutor(ExecutorConfig{DB: db, Processor: processor.Process}).ProcessJob(context.Background(), job.ID)

	if got := processor.CallCount(); got != 1 {
		t.Fatalf("processor calls = %d, want 1", got)
	}
	input := processor.InputAt(t, 0)
	if input.Filename != job.OriginalFilename || input.MimeType != job.MimeType {
		t.Fatalf("processor input metadata = %#v, want filename %q mime %q", input, job.OriginalFilename, job.MimeType)
	}
	if input.DataURL != DataURL(job.MimeType, fileData) {
		t.Fatalf("processor DataURL = %q, want stored file data URL", input.DataURL)
	}

	var gotJob OCRJob
	if err := db.First(&gotJob, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if gotJob.Status != OCRJobStatusCompleted {
		t.Fatalf("job status = %q, want %q", gotJob.Status, OCRJobStatusCompleted)
	}
	if gotJob.ErrorMessage != "" {
		t.Fatalf("job error_message = %q, want empty", gotJob.ErrorMessage)
	}
	if gotJob.DocumentID == nil {
		t.Fatal("job document_id is nil")
	}

	var doc OCRDocument
	if err := db.First(&doc, "id = ?", *gotJob.DocumentID).Error; err != nil {
		t.Fatalf("load OCR document: %v", err)
	}
	if doc.UserID == nil || *doc.UserID != user.ID {
		t.Fatalf("document user_id = %#v, want %s", doc.UserID, user.ID)
	}
	if doc.JobID == nil || *doc.JobID != job.ID {
		t.Fatalf("document job_id = %#v, want %s", doc.JobID, job.ID)
	}
	if doc.OriginalFilename != job.OriginalFilename || doc.MimeType != job.MimeType || doc.FileSize != job.FileSize {
		t.Fatalf("document metadata = %#v, want job metadata", doc)
	}
	if doc.PageCount != 2 {
		t.Fatalf("document page_count = %d, want 2", doc.PageCount)
	}
	if doc.DocumentHash != job.DocumentHash {
		t.Fatalf("document hash = %q, want %q", doc.DocumentHash, job.DocumentHash)
	}
	wantMarkdown := "# Invoice\n\n![invoice.png](data:image/png;base64,iVBORw0KGgo=)"
	if doc.Markdown != wantMarkdown {
		t.Fatalf("document markdown = %q, want %q", doc.Markdown, wantMarkdown)
	}
	assertJSONEqual(t, doc.AnnotationJSON, datatypes.JSON([]byte(annotation)))
	assertJSONEqual(t, doc.RawResponseJSON, datatypes.JSON(rawResponse))
	assertExecutorCreditBalance(t, db, user.ID, 3)
	assertExecutorCreditLedgerCount(t, db, 1, "related_job_id = ? AND entry_type = ? AND credits_delta = ?", job.ID, billing.CreditLedgerEntryDebit, -2)
}

func TestExecutorDispatchesWebhookOnSuccessfulJob(t *testing.T) {
	db := executorTestDB(t)
	user := createExecutorTestUser(t, db)
	grantExecutorTestCredits(t, db, user.ID, 5)
	job := createQueuedExecutorJob(t, db, executorJobFixture{UserID: &user.ID, FileData: []byte("invoice bytes")})
	annotation := `{"invoice_id":"INV-001"}`
	rawResponse := []byte(`{"pages":[{"index":0,"markdown":"# Invoice"},{"index":1,"markdown":"# Page two"}],"document_annotation":{"invoice_id":"INV-001"}}`)
	processor := &recordingExecutorProcessor{
		response: &MistralResponse{
			Pages:              []MistralPage{{Index: 0, Markdown: "# Invoice"}},
			DocumentAnnotation: &annotation,
		},
		raw: rawResponse,
	}
	dispatcher := &fakeWebhookDispatcher{}

	NewExecutor(ExecutorConfig{
		DB:                db,
		Processor:         processor.Process,
		WebhookDispatcher: dispatcher,
	}).ProcessJob(context.Background(), job.ID)

	var gotJob OCRJob
	if err := db.First(&gotJob, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if gotJob.Status != OCRJobStatusCompleted {
		t.Fatalf("job status = %q, want %q", gotJob.Status, OCRJobStatusCompleted)
	}
	if gotJob.DocumentID == nil {
		t.Fatal("job document_id is nil")
	}

	events := waitForWebhookEvents(t, dispatcher, 2)
	assertExecutorWebhookJobEvent(t, events[0], webhooks.EventJobStarted, user.ID, job.ID, OCRJobStatusProcessing, job.OriginalFilename)
	if events[0].Job.ErrorMessage != nil {
		t.Fatalf("started error_message = %#v, want nil", events[0].Job.ErrorMessage)
	}

	assertExecutorWebhookSucceededJobIDOnly(t, events[1], user.ID, job.ID)
	if events[1].Job.ErrorMessage != nil {
		t.Fatalf("succeeded error_message = %#v, want nil", events[1].Job.ErrorMessage)
	}
}

func TestExecutorDispatchesWebhookOnProcessorFailure(t *testing.T) {
	db := executorTestDB(t)
	user := createExecutorTestUser(t, db)
	job := createQueuedExecutorJob(t, db, executorJobFixture{UserID: &user.ID, FileData: []byte("invoice bytes")})
	processor := &recordingExecutorProcessor{err: errors.New("processor unavailable")}
	dispatcher := &fakeWebhookDispatcher{}

	NewExecutor(ExecutorConfig{
		DB:                db,
		Processor:         processor.Process,
		WebhookDispatcher: dispatcher,
	}).ProcessJob(context.Background(), job.ID)

	var gotJob OCRJob
	if err := db.First(&gotJob, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if gotJob.Status != OCRJobStatusFailed {
		t.Fatalf("job status = %q, want %q", gotJob.Status, OCRJobStatusFailed)
	}

	events := waitForWebhookEvents(t, dispatcher, 2)
	assertExecutorWebhookJobEvent(t, events[0], webhooks.EventJobStarted, user.ID, job.ID, OCRJobStatusProcessing, job.OriginalFilename)
	assertExecutorWebhookJobEvent(t, events[1], webhooks.EventJobFailed, user.ID, job.ID, OCRJobStatusFailed, job.OriginalFilename)
	if events[1].Job.ErrorMessage == nil {
		t.Fatal("failed error_message is nil")
	}
	if *events[1].Job.ErrorMessage != gotJob.ErrorMessage {
		t.Fatalf("failed error_message = %q, want stored %q", *events[1].Job.ErrorMessage, gotJob.ErrorMessage)
	}
}

func TestExecutorDispatchesWebhookOnFileFailure(t *testing.T) {
	db := executorTestDB(t)
	user := createExecutorTestUser(t, db)
	job := createQueuedExecutorJob(t, db, executorJobFixture{UserID: &user.ID, FileData: []byte("invoice bytes")})
	if err := os.Remove(job.FilePath); err != nil {
		t.Fatalf("remove stored file: %v", err)
	}
	processor := &recordingExecutorProcessor{}
	dispatcher := &fakeWebhookDispatcher{}

	NewExecutor(ExecutorConfig{
		DB:                db,
		Processor:         processor.Process,
		WebhookDispatcher: dispatcher,
	}).ProcessJob(context.Background(), job.ID)

	if got := processor.CallCount(); got != 0 {
		t.Fatalf("processor calls = %d, want 0", got)
	}
	var gotJob OCRJob
	if err := db.First(&gotJob, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if gotJob.Status != OCRJobStatusFailed {
		t.Fatalf("job status = %q, want %q", gotJob.Status, OCRJobStatusFailed)
	}

	events := waitForWebhookEvents(t, dispatcher, 2)
	assertExecutorWebhookJobEvent(t, events[0], webhooks.EventJobStarted, user.ID, job.ID, OCRJobStatusProcessing, job.OriginalFilename)
	assertExecutorWebhookJobEvent(t, events[1], webhooks.EventJobFailed, user.ID, job.ID, OCRJobStatusFailed, job.OriginalFilename)
	if events[1].Job.ErrorMessage == nil || *events[1].Job.ErrorMessage != gotJob.ErrorMessage {
		t.Fatalf("failed error_message = %#v, want stored %q", events[1].Job.ErrorMessage, gotJob.ErrorMessage)
	}
}

func TestExecutorWebhookDispatcherErrorDoesNotPreventFinalState(t *testing.T) {
	t.Run("completed", func(t *testing.T) {
		db := executorTestDB(t)
		user := createExecutorTestUser(t, db)
		grantExecutorTestCredits(t, db, user.ID, 5)
		job := createQueuedExecutorJob(t, db, executorJobFixture{UserID: &user.ID, FileData: []byte("invoice bytes")})
		processor := &recordingExecutorProcessor{}
		dispatcher := &fakeWebhookDispatcher{err: errors.New("webhook dispatch failed")}

		NewExecutor(ExecutorConfig{
			DB:                db,
			Processor:         processor.Process,
			WebhookDispatcher: dispatcher,
		}).ProcessJob(context.Background(), job.ID)

		var got OCRJob
		if err := db.First(&got, "id = ?", job.ID).Error; err != nil {
			t.Fatalf("load job: %v", err)
		}
		if got.Status != OCRJobStatusCompleted || got.DocumentID == nil {
			t.Fatalf("job status/document = %q/%#v, want completed with document", got.Status, got.DocumentID)
		}
		if events := waitForWebhookEvents(t, dispatcher, 2); len(events) != 2 {
			t.Fatalf("webhook events = %#v, want dispatch attempts despite errors", events)
		}
	})

	t.Run("failed", func(t *testing.T) {
		db := executorTestDB(t)
		user := createExecutorTestUser(t, db)
		job := createQueuedExecutorJob(t, db, executorJobFixture{UserID: &user.ID, FileData: []byte("invoice bytes")})
		processor := &recordingExecutorProcessor{err: errors.New("processor unavailable")}
		dispatcher := &fakeWebhookDispatcher{err: errors.New("webhook dispatch failed")}

		NewExecutor(ExecutorConfig{
			DB:                db,
			Processor:         processor.Process,
			WebhookDispatcher: dispatcher,
		}).ProcessJob(context.Background(), job.ID)

		var got OCRJob
		if err := db.First(&got, "id = ?", job.ID).Error; err != nil {
			t.Fatalf("load job: %v", err)
		}
		if got.Status != OCRJobStatusFailed {
			t.Fatalf("job status = %q, want %q", got.Status, OCRJobStatusFailed)
		}
		if events := waitForWebhookEvents(t, dispatcher, 2); len(events) != 2 {
			t.Fatalf("webhook events = %#v, want dispatch attempts despite errors", events)
		}
	})
}

func TestExecutorWebhookDispatchDoesNotBlockJobCompletion(t *testing.T) {
	db := executorTestDB(t)
	user := createExecutorTestUser(t, db)
	grantExecutorTestCredits(t, db, user.ID, 5)
	job := createQueuedExecutorJob(t, db, executorJobFixture{UserID: &user.ID, FileData: []byte("invoice bytes")})
	processor := &recordingExecutorProcessor{}
	blockDispatch := make(chan struct{})
	dispatcher := &fakeWebhookDispatcher{blockUntil: blockDispatch}
	t.Cleanup(func() {
		closeOnce(blockDispatch)
	})

	done := make(chan struct{})
	go func() {
		defer close(done)
		NewExecutor(ExecutorConfig{
			DB:                   db,
			Processor:            processor.Process,
			WebhookDispatcher:    dispatcher,
			WebhookDispatchSlots: 2,
		}).ProcessJob(context.Background(), job.ID)
	}()

	select {
	case <-done:
	case <-time.After(750 * time.Millisecond):
		t.Fatal("ProcessJob blocked on webhook dispatch")
	}

	var got OCRJob
	if err := db.First(&got, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if got.Status != OCRJobStatusCompleted || got.DocumentID == nil {
		t.Fatalf("job status/document = %q/%#v, want completed with document", got.Status, got.DocumentID)
	}
}

func TestExecutorWebhookSkipsNilUserID(t *testing.T) {
	db := executorTestDB(t)
	job := createQueuedExecutorJob(t, db, executorJobFixture{FileData: []byte("invoice bytes")})
	dispatcher := &fakeWebhookDispatcher{}
	processor := &recordingExecutorProcessor{}

	NewExecutor(ExecutorConfig{
		DB:                db,
		Processor:         processor.Process,
		WebhookDispatcher: dispatcher,
	}).ProcessJob(context.Background(), job.ID)

	var got OCRJob
	if err := db.First(&got, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if got.Status != OCRJobStatusCompleted {
		t.Fatalf("job status = %q, want %q", got.Status, OCRJobStatusCompleted)
	}
	if events := dispatcher.Events(); len(events) != 0 {
		t.Fatalf("webhook events = %#v, want none for nil user", events)
	}
}

func TestExecutorWebhookLifecycleRequeueEvents(t *testing.T) {
	t.Run("claimed requeue emits only started", func(t *testing.T) {
		db := executorTestDB(t)
		user := createExecutorTestUser(t, db)
		job := createQueuedExecutorJob(t, db, executorJobFixture{UserID: &user.ID, FileData: []byte("invoice bytes")})
		processor := &recordingExecutorProcessor{err: context.Canceled}
		dispatcher := &fakeWebhookDispatcher{}

		NewExecutor(ExecutorConfig{
			DB:                db,
			Processor:         processor.Process,
			WebhookDispatcher: dispatcher,
		}).ProcessJob(context.Background(), job.ID)

		var got OCRJob
		if err := db.First(&got, "id = ?", job.ID).Error; err != nil {
			t.Fatalf("load job: %v", err)
		}
		if got.Status != OCRJobStatusQueued {
			t.Fatalf("job status = %q, want %q", got.Status, OCRJobStatusQueued)
		}
		events := waitForWebhookEvents(t, dispatcher, 1)
		assertExecutorWebhookJobEvent(t, events[0], webhooks.EventJobStarted, user.ID, job.ID, OCRJobStatusProcessing, job.OriginalFilename)
	})

	t.Run("not claimed emits none", func(t *testing.T) {
		db := executorTestDB(t)
		user := createExecutorTestUser(t, db)
		job := createQueuedExecutorJob(t, db, executorJobFixture{UserID: &user.ID, FileData: []byte("invoice bytes")})
		if err := db.Model(&OCRJob{}).Where("id = ?", job.ID).Update("status", OCRJobStatusProcessing).Error; err != nil {
			t.Fatalf("mark job processing: %v", err)
		}
		dispatcher := &fakeWebhookDispatcher{}

		NewExecutor(ExecutorConfig{
			DB:                db,
			Processor:         (&recordingExecutorProcessor{}).Process,
			WebhookDispatcher: dispatcher,
		}).ProcessJob(context.Background(), job.ID)

		if events := dispatcher.Events(); len(events) != 0 {
			t.Fatalf("webhook events = %#v, want none when claim is skipped", events)
		}
	})
}

func TestExecutorLogsJobLifecycleAndDebugMetadata(t *testing.T) {
	db := executorTestDB(t)
	job := createQueuedExecutorJob(t, db, executorJobFixture{FileData: []byte("invoice bytes")})
	var logs bytes.Buffer
	processor := &recordingExecutorProcessor{
		response: &MistralResponse{Pages: []MistralPage{{Index: 0, Markdown: "# Page one"}}},
		raw:      []byte(`{"pages":[{"index":0,"markdown":"# Page one"}]}`),
	}

	NewExecutor(ExecutorConfig{
		DB:        db,
		Processor: processor.Process,
		Logger:    logging.NewJSONLogger(&logs, true),
	}).ProcessJob(context.Background(), job.ID)

	out := logs.String()
	for _, want := range []string{
		`"msg":"ocr.job_claimed"`,
		`"msg":"ocr.job_file_read"`,
		`"msg":"ocr.job_schema_resolved"`,
		`"msg":"ocr.upstream_request_completed"`,
		`"msg":"ocr.job_completed"`,
		`"job_id":"` + job.ID.String() + `"`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("logs missing %s in:\n%s", want, out)
		}
	}
	for _, forbidden := range []string{"invoice bytes", "# Page one", "document_annotation"} {
		if strings.Contains(out, forbidden) {
			t.Fatalf("logs leaked %q in:\n%s", forbidden, out)
		}
	}
}

func TestExecutorLoggerFiltersDebugMetadataAtInfoLevel(t *testing.T) {
	db := executorTestDB(t)
	job := createQueuedExecutorJob(t, db, executorJobFixture{FileData: []byte("invoice bytes")})
	var logs bytes.Buffer
	processor := &recordingExecutorProcessor{
		response: &MistralResponse{Pages: []MistralPage{{Index: 0, Markdown: "# Page one"}}},
		raw:      []byte(`{"pages":[{"index":0,"markdown":"# Page one"}]}`),
	}

	NewExecutor(ExecutorConfig{
		DB:        db,
		Processor: processor.Process,
		Logger:    logging.NewJSONLogger(&logs, false),
	}).ProcessJob(context.Background(), job.ID)

	out := logs.String()
	if !strings.Contains(out, `"msg":"ocr.job_completed"`) {
		t.Fatalf("info lifecycle log missing in:\n%s", out)
	}
	for _, hidden := range []string{
		`"msg":"ocr.job_file_read"`,
		`"msg":"ocr.job_schema_resolved"`,
		`"msg":"ocr.upstream_request_completed"`,
	} {
		if strings.Contains(out, hidden) {
			t.Fatalf("debug log %s was emitted at info level:\n%s", hidden, out)
		}
	}
}

func TestExecutorMarksJobFailedOnProcessorError(t *testing.T) {
	db := executorTestDB(t)
	job := createQueuedExecutorJob(t, db, executorJobFixture{FileData: []byte("invoice bytes")})
	processor := &recordingExecutorProcessor{err: errors.New("processor unavailable")}

	NewExecutor(ExecutorConfig{DB: db, Processor: processor.Process}).ProcessJob(context.Background(), job.ID)

	assertExecutorJobFailed(t, db, job.ID, "processor unavailable")
	assertExecutorDocumentCount(t, db, 0)
}

func TestExecutorDoesNotDebitCreditsWhenProcessorFails(t *testing.T) {
	db := executorTestDB(t)
	user := createExecutorTestUser(t, db)
	grantExecutorTestCredits(t, db, user.ID, 5)
	job := createQueuedExecutorJob(t, db, executorJobFixture{UserID: &user.ID, FileData: []byte("invoice bytes")})
	processor := &recordingExecutorProcessor{err: errors.New("processor unavailable")}

	NewExecutor(ExecutorConfig{DB: db, Processor: processor.Process}).ProcessJob(context.Background(), job.ID)

	assertExecutorJobFailed(t, db, job.ID, "processor unavailable")
	assertExecutorDocumentCount(t, db, 0)
	assertExecutorCreditBalance(t, db, user.ID, 5)
	assertExecutorCreditLedgerCount(t, db, 0, "related_job_id = ? AND entry_type = ?", job.ID, billing.CreditLedgerEntryDebit)
}

func TestExecutorFailsJobWhenCompletionDebitHasInsufficientCredits(t *testing.T) {
	db := executorTestDB(t)
	user := createExecutorTestUser(t, db)
	grantExecutorTestCredits(t, db, user.ID, 1)
	job := createQueuedExecutorJob(t, db, executorJobFixture{UserID: &user.ID, FileData: []byte("invoice bytes")})
	processor := &recordingExecutorProcessor{
		response: &MistralResponse{
			Pages: []MistralPage{{Index: 0, Markdown: "# Page one"}, {Index: 1, Markdown: "# Page two"}},
		},
		raw: []byte(`{"pages":[{"index":0,"markdown":"# Page one"},{"index":1,"markdown":"# Page two"}]}`),
	}

	NewExecutor(ExecutorConfig{DB: db, Processor: processor.Process}).ProcessJob(context.Background(), job.ID)

	if got := processor.CallCount(); got != 1 {
		t.Fatalf("processor calls = %d, want 1", got)
	}
	assertExecutorJobFailed(t, db, job.ID, "debit OCR job credits")
	assertExecutorDocumentCount(t, db, 0)
	assertExecutorCreditBalance(t, db, user.ID, 1)
	assertExecutorCreditLedgerCount(t, db, 0, "related_job_id = ? AND entry_type = ?", job.ID, billing.CreditLedgerEntryDebit)
}

func TestExecutorRequeuesJobOnProcessorContextCancellation(t *testing.T) {
	db := executorTestDB(t)
	job := createQueuedExecutorJob(t, db, executorJobFixture{FileData: []byte("invoice bytes")})
	processor := &recordingExecutorProcessor{err: context.Canceled}

	NewExecutor(ExecutorConfig{DB: db, Processor: processor.Process}).ProcessJob(context.Background(), job.ID)

	if got := processor.CallCount(); got != 1 {
		t.Fatalf("processor calls = %d, want 1", got)
	}
	var got OCRJob
	if err := db.First(&got, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if got.Status != OCRJobStatusQueued {
		t.Fatalf("job status = %q, want %q", got.Status, OCRJobStatusQueued)
	}
	if got.ErrorMessage != "" {
		t.Fatalf("job error_message = %q, want empty", got.ErrorMessage)
	}
	if got.DocumentID != nil {
		t.Fatalf("job document_id = %s, want nil", *got.DocumentID)
	}
	assertExecutorDocumentCount(t, db, 0)
}

func TestExecutorRequeuesJobWhenContextCanceledWithOpaqueProcessorError(t *testing.T) {
	db := executorTestDB(t)
	job := createQueuedExecutorJob(t, db, executorJobFixture{FileData: []byte("invoice bytes")})
	ctx, cancel := context.WithCancel(context.Background())
	processor := &recordingExecutorProcessor{
		beforeReturn: cancel,
		err:          errors.New("opaque worker shutdown"),
	}

	NewExecutor(ExecutorConfig{DB: db, Processor: processor.Process}).ProcessJob(ctx, job.ID)

	if got := processor.CallCount(); got != 1 {
		t.Fatalf("processor calls = %d, want 1", got)
	}
	var got OCRJob
	if err := db.First(&got, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if got.Status != OCRJobStatusQueued {
		t.Fatalf("job status = %q, want %q", got.Status, OCRJobStatusQueued)
	}
	if got.ErrorMessage != "" {
		t.Fatalf("job error_message = %q, want empty", got.ErrorMessage)
	}
	if got.DocumentID != nil {
		t.Fatalf("job document_id = %s, want nil", *got.DocumentID)
	}
	assertExecutorDocumentCount(t, db, 0)
}

func TestExecutorMarksJobFailedWhenFileIsMissing(t *testing.T) {
	db := executorTestDB(t)
	job := createQueuedExecutorJob(t, db, executorJobFixture{FileData: []byte("invoice bytes")})
	if err := os.Remove(job.FilePath); err != nil {
		t.Fatalf("remove stored file: %v", err)
	}
	processor := &recordingExecutorProcessor{
		response: &MistralResponse{Pages: []MistralPage{{Index: 0, Markdown: "should not run"}}},
		raw:      []byte(`{"pages":[]}`),
	}

	NewExecutor(ExecutorConfig{DB: db, Processor: processor.Process}).ProcessJob(context.Background(), job.ID)

	if got := processor.CallCount(); got != 0 {
		t.Fatalf("processor calls = %d, want 0", got)
	}
	assertExecutorJobFailed(t, db, job.ID, "read OCR job file")
	assertExecutorDocumentCount(t, db, 0)
}

func TestExecutorRacingProcessesClaimJobOnlyOnce(t *testing.T) {
	db := executorTestDB(t)
	installOCRJobRaceSleepTrigger(t, db)
	dbA := openExecutorRaceDB(t)
	dbB := openExecutorRaceDB(t)
	job := createQueuedExecutorJob(t, db, executorJobFixture{FileData: []byte("invoice bytes")})
	processor := &recordingExecutorProcessor{
		delay: 50 * time.Millisecond,
		response: &MistralResponse{
			Pages: []MistralPage{{Index: 0, Markdown: "# Once"}},
		},
		raw: []byte(`{"pages":[{"index":0,"markdown":"# Once"}]}`),
	}
	executorA := NewExecutor(ExecutorConfig{DB: dbA, Processor: processor.Process})
	executorB := NewExecutor(ExecutorConfig{DB: dbB, Processor: processor.Process})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		executorA.ProcessJob(context.Background(), job.ID)
	}()
	go func() {
		defer wg.Done()
		executorB.ProcessJob(context.Background(), job.ID)
	}()
	wg.Wait()

	if got := processor.CallCount(); got != 1 {
		t.Fatalf("processor calls = %d, want 1", got)
	}
	assertExecutorDocumentCount(t, db, 1)
	var gotJob OCRJob
	if err := db.First(&gotJob, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if gotJob.Status != OCRJobStatusCompleted {
		t.Fatalf("job status = %q, want %q", gotJob.Status, OCRJobStatusCompleted)
	}
	if gotJob.DocumentID == nil {
		t.Fatal("job document_id is nil")
	}
}

func TestExecutorUsesSavedSchema(t *testing.T) {
	db := executorTestDB(t)
	user := createExecutorTestUser(t, db)
	grantExecutorTestCredits(t, db, user.ID, 5)
	schema := ExtractionSchema{
		UserID:     &user.ID,
		Name:       "invoice",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object","properties":{"total":{"type":"number"}}}`)),
		Strict:     false,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	if err := db.Model(&schema).Update("strict", false).Error; err != nil {
		t.Fatalf("update schema strict: %v", err)
	}
	schema.Strict = false
	job := createQueuedExecutorJob(t, db, executorJobFixture{
		UserID:   &user.ID,
		SchemaID: &schema.ID,
		FileData: []byte("invoice bytes"),
	})
	annotation := `{"total":42}`
	processor := &recordingExecutorProcessor{
		response: &MistralResponse{
			Pages:              []MistralPage{{Index: 0, Markdown: "# Invoice"}},
			DocumentAnnotation: &annotation,
		},
		raw: []byte(`{"pages":[{"index":0,"markdown":"# Invoice"}],"document_annotation":{"total":42}}`),
	}

	NewExecutor(ExecutorConfig{DB: db, Processor: processor.Process}).ProcessJob(context.Background(), job.ID)

	if got := processor.CallCount(); got != 1 {
		t.Fatalf("processor calls = %d, want 1", got)
	}
	input := processor.InputAt(t, 0)
	assertJSONEqual(t, datatypes.JSON(input.Schema), schema.SchemaJSON)
	if input.Strict {
		t.Fatal("processor strict = true, want false from saved schema")
	}
	var gotJob OCRJob
	if err := db.First(&gotJob, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if gotJob.Status != OCRJobStatusCompleted || gotJob.DocumentID == nil {
		t.Fatalf("job status/document = %q/%#v, want completed with document", gotJob.Status, gotJob.DocumentID)
	}
	var doc OCRDocument
	if err := db.First(&doc, "id = ?", *gotJob.DocumentID).Error; err != nil {
		t.Fatalf("load OCR document: %v", err)
	}
	if doc.SchemaID == nil || *doc.SchemaID != schema.ID {
		t.Fatalf("document schema_id = %#v, want %s", doc.SchemaID, schema.ID)
	}
	if doc.JobID == nil || *doc.JobID != job.ID {
		t.Fatalf("document job_id = %#v, want %s", doc.JobID, job.ID)
	}
}

func TestExecutorLinksCompletedDocumentToCollection(t *testing.T) {
	db := executorTestDB(t)
	user := createExecutorTestUser(t, db)
	grantExecutorTestCredits(t, db, user.ID, 5)
	schema := ExtractionSchema{
		UserID: &user.ID, Name: "Invoice",
		SchemaJSON: datatypes.JSON([]byte(`{"type":"object"}`)), Strict: true,
	}
	if err := db.Create(&schema).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	collection := Collection{UserID: user.ID, Name: "Invoices"}
	if err := db.Create(&collection).Error; err != nil {
		t.Fatalf("create collection: %v", err)
	}
	if err := db.Create(&CollectionSchema{CollectionID: collection.ID, SchemaID: schema.ID}).Error; err != nil {
		t.Fatalf("create collection schema: %v", err)
	}
	annotation := `{"invoice_id":"INV-001"}`
	job := createQueuedExecutorJob(t, db, executorJobFixture{UserID: &user.ID, SchemaID: &schema.ID, FileData: []byte("invoice bytes")})
	processor := &recordingExecutorProcessor{
		response: &MistralResponse{Pages: []MistralPage{{Index: 0, Markdown: "# Invoice"}}, DocumentAnnotation: &annotation},
		raw:      []byte(`{"pages":[{"index":0,"markdown":"# Invoice"}],"document_annotation":{"invoice_id":"INV-001"}}`),
	}

	NewExecutor(ExecutorConfig{DB: db, Processor: processor.Process}).ProcessJob(context.Background(), job.ID)

	var gotJob OCRJob
	if err := db.First(&gotJob, "id = ?", job.ID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if gotJob.DocumentID == nil {
		t.Fatal("document_id is nil")
	}
	var count int64
	if err := db.Model(&CollectionDocument{}).Where("collection_id = ? AND document_id = ?", collection.ID, *gotJob.DocumentID).Count(&count).Error; err != nil {
		t.Fatalf("count collection document: %v", err)
	}
	if count != 1 {
		t.Fatalf("collection document count = %d, want 1", count)
	}
}

func TestExecutorRunProcessesNotification(t *testing.T) {
	db := executorTestDB(t)
	executorDB := openExecutorRaceDB(t)
	job := createQueuedExecutorJob(t, db, executorJobFixture{FileData: []byte("invoice bytes")})
	processor := &recordingExecutorProcessor{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ready := make(chan struct{})
	notifyChannel := "ocr_jobs_test_" + strings.ReplaceAll(uuid.NewString(), "-", "_")
	executor := NewExecutor(ExecutorConfig{
		DB:            executorDB,
		DSN:           executorTestDSN(t),
		Processor:     processor.Process,
		NotifyChannel: notifyChannel,
		WorkerCount:   1,
		QueueBuffer:   4,
		PollInterval:  time.Hour,
	})
	executor.listenerReady = func() {
		closeOnce(ready)
	}
	runErr := make(chan error, 1)
	go func() {
		runErr <- executor.Run(ctx)
	}()
	var waitRunOnce sync.Once
	waitRun := func() {
		waitRunOnce.Do(func() {
			select {
			case err := <-runErr:
				if err != nil {
					t.Fatalf("executor Run error: %v", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("executor Run did not exit after cancellation")
			}
		})
	}
	t.Cleanup(func() {
		cancel()
		waitRun()
	})

	waitForExecutorTestCondition(t, "listener ready", func() bool {
		select {
		case <-ready:
			return true
		default:
			return false
		}
	})
	if err := notifyOCRJobQueued(context.Background(), db, notifyChannel, job.ID); err != nil {
		t.Fatalf("notify OCR job queued: %v", err)
	}
	waitForExecutorJobStatus(t, db, job.ID, OCRJobStatusCompleted)
	if got := processor.CallCount(); got != 1 {
		t.Fatalf("processor calls = %d, want 1", got)
	}
	cancel()
	waitRun()
}

func TestExecutorSweeperProcessesQueuedJob(t *testing.T) {
	db := executorTestDB(t)
	executorDB := openExecutorRaceDB(t)
	job := createQueuedExecutorJob(t, db, executorJobFixture{FileData: []byte("invoice bytes")})
	processor := &recordingExecutorProcessor{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	executor := NewExecutor(ExecutorConfig{
		DB:           executorDB,
		Processor:    processor.Process,
		WorkerCount:  1,
		QueueBuffer:  4,
		PollInterval: 20 * time.Millisecond,
	})
	runErr := make(chan error, 1)
	go func() {
		runErr <- executor.Run(ctx)
	}()
	var waitRunOnce sync.Once
	waitRun := func() {
		waitRunOnce.Do(func() {
			select {
			case err := <-runErr:
				if err != nil {
					t.Fatalf("executor Run error: %v", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("executor Run did not exit after cancellation")
			}
		})
	}
	t.Cleanup(func() {
		cancel()
		waitRun()
	})

	waitForExecutorJobStatus(t, db, job.ID, OCRJobStatusCompleted)
	if got := processor.CallCount(); got != 1 {
		t.Fatalf("processor calls = %d, want 1", got)
	}
	cancel()
	waitRun()
}

func TestExecutorSweeperRequeuesStaleProcessingJob(t *testing.T) {
	db := executorTestDB(t)
	executorDB := openExecutorRaceDB(t)
	job := createQueuedExecutorJob(t, db, executorJobFixture{FileData: []byte("invoice bytes")})
	staleDoc := OCRDocument{
		UserID:           job.UserID,
		OriginalFilename: job.OriginalFilename,
		MimeType:         job.MimeType,
		FileSize:         job.FileSize,
		DocumentHash:     "stale-" + strings.Repeat("b", 58),
		Markdown:         "stale",
		RawResponseJSON:  datatypes.JSON([]byte(`{"pages":[]}`)),
	}
	if err := db.Create(&staleDoc).Error; err != nil {
		t.Fatalf("create stale OCR document: %v", err)
	}
	staleUpdatedAt := time.Now().Add(-time.Hour)
	if err := db.Model(&OCRJob{}).
		Where("id = ?", job.ID).
		Updates(map[string]any{
			"status":        OCRJobStatusProcessing,
			"error_message": "worker disappeared",
			"document_id":   staleDoc.ID,
			"updated_at":    staleUpdatedAt,
		}).Error; err != nil {
		t.Fatalf("mark job stale processing: %v", err)
	}
	processor := &recordingExecutorProcessor{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	executor := NewExecutor(ExecutorConfig{
		DB:                     executorDB,
		Processor:              processor.Process,
		WorkerCount:            1,
		QueueBuffer:            4,
		PollInterval:           20 * time.Millisecond,
		StaleProcessingTimeout: 10 * time.Millisecond,
	})
	runErr := make(chan error, 1)
	go func() {
		runErr <- executor.Run(ctx)
	}()
	var waitRunOnce sync.Once
	waitRun := func() {
		waitRunOnce.Do(func() {
			select {
			case err := <-runErr:
				if err != nil {
					t.Fatalf("executor Run error: %v", err)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("executor Run did not exit after cancellation")
			}
		})
	}
	t.Cleanup(func() {
		cancel()
		waitRun()
	})

	got := waitForExecutorJobStatus(t, db, job.ID, OCRJobStatusCompleted)
	if got.ErrorMessage != "" {
		t.Fatalf("job error_message = %q, want empty", got.ErrorMessage)
	}
	if got.DocumentID == nil || *got.DocumentID == staleDoc.ID {
		t.Fatalf("job document_id = %#v, want newly completed document", got.DocumentID)
	}
	if got.UpdatedAt.Before(staleUpdatedAt) || got.UpdatedAt.Equal(staleUpdatedAt) {
		t.Fatalf("job updated_at = %s, want after stale timestamp %s", got.UpdatedAt, staleUpdatedAt)
	}
	if got := processor.CallCount(); got != 1 {
		t.Fatalf("processor calls = %d, want 1", got)
	}
	cancel()
	waitRun()
}

func executorTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	return testsupport.OpenPostgresDB(t, &auth.User{}, &billing.BillingOrder{}, &billing.CreditBucket{}, &billing.CreditLedgerEntry{}, &ExtractionSchema{}, &OCRDocument{}, &OCRJob{}, &Collection{}, &CollectionSchema{}, &CollectionDocument{})
}

func executorTestDSN(t *testing.T) string {
	t.Helper()
	return testsupport.PostgresTestDSN(t)
}

func openExecutorRaceDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(postgres.Open(executorTestDSN(t)), &gorm.Config{})
	if err != nil {
		t.Fatalf("open race postgres handle: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get race postgres sql db: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})
	return db
}

func installOCRJobRaceSleepTrigger(t *testing.T, db *gorm.DB) {
	t.Helper()
	if err := db.Exec(`
CREATE OR REPLACE FUNCTION ocr_job_race_sleep()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
	PERFORM pg_sleep(0.15);
	RETURN NEW;
END;
$$;
`).Error; err != nil {
		t.Fatalf("create race sleep function: %v", err)
	}
	if err := db.Exec(`DROP TRIGGER IF EXISTS ocr_job_race_sleep_trigger ON ocr_jobs`).Error; err != nil {
		t.Fatalf("drop existing race sleep trigger: %v", err)
	}
	if err := db.Exec(`
CREATE TRIGGER ocr_job_race_sleep_trigger
BEFORE UPDATE ON ocr_jobs
FOR EACH ROW
WHEN (OLD.status = 'queued' AND NEW.status = 'processing')
EXECUTE FUNCTION ocr_job_race_sleep()
`).Error; err != nil {
		t.Fatalf("create race sleep trigger: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Exec(`DROP TRIGGER IF EXISTS ocr_job_race_sleep_trigger ON ocr_jobs`).Error; err != nil {
			t.Fatalf("drop race sleep trigger: %v", err)
		}
		if err := db.Exec(`DROP FUNCTION IF EXISTS ocr_job_race_sleep()`).Error; err != nil {
			t.Fatalf("drop race sleep function: %v", err)
		}
	})
}

func waitForExecutorJobStatus(t *testing.T, db *gorm.DB, id uuid.UUID, want OCRJobStatus) OCRJob {
	t.Helper()
	var got OCRJob
	var lastErr error
	waitForExecutorTestCondition(t, "job status "+string(want), func() bool {
		if err := db.First(&got, "id = ?", id).Error; err != nil {
			lastErr = err
			return false
		}
		lastErr = nil
		return got.Status == want
	}, func() string {
		if lastErr != nil {
			return lastErr.Error()
		}
		return fmt.Sprintf("last status %q document_id %#v error %q", got.Status, got.DocumentID, got.ErrorMessage)
	})
	return got
}

func waitForExecutorTestCondition(t *testing.T, name string, check func() bool, details ...func() string) {
	t.Helper()
	deadline := time.After(5 * time.Second)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		if check() {
			return
		}
		select {
		case <-deadline:
			if len(details) > 0 {
				t.Fatalf("timed out waiting for %s: %s", name, details[0]())
			}
			t.Fatalf("timed out waiting for %s", name)
		case <-ticker.C:
		}
	}
}

func closeOnce(ch chan struct{}) {
	defer func() {
		_ = recover()
	}()
	close(ch)
}

type executorJobFixture struct {
	UserID   *string
	SchemaID *uuid.UUID
	FileData []byte
}

func createQueuedExecutorJob(t *testing.T, db *gorm.DB, fixture executorJobFixture) OCRJob {
	t.Helper()
	data := fixture.FileData
	if data == nil {
		data = []byte("ocr job data")
	}
	id := uuid.New()
	path := filepath.Join(t.TempDir(), id.String()+".txt")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write stored file: %v", err)
	}
	job := OCRJob{
		ID:               id,
		UserID:           fixture.UserID,
		OriginalFilename: "invoice.txt",
		MimeType:         "text/plain",
		FileSize:         int64(len(data)),
		PageCount:        1,
		DocumentHash:     strings.Repeat("a", 64),
		FilePath:         path,
		SchemaID:         fixture.SchemaID,
		Status:           OCRJobStatusQueued,
	}
	if err := db.Create(&job).Error; err != nil {
		t.Fatalf("create OCR job: %v", err)
	}
	return job
}

func createExecutorTestUser(t *testing.T, db *gorm.DB) auth.User {
	t.Helper()
	user := auth.User{Name: "Executor Owner", Email: uuid.NewString() + "@example.com"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	return user
}

func grantExecutorTestCredits(t *testing.T, db *gorm.DB, userID string, credits int) billing.CreditBucket {
	t.Helper()
	bucket := billing.CreditBucket{
		UserID:           userID,
		SourceType:       billing.CreditSourceAdjustment,
		CreditsGranted:   credits,
		CreditsRemaining: credits,
		ValidFrom:        time.Now().UTC().Add(-time.Hour),
	}
	if err := db.Create(&bucket).Error; err != nil {
		t.Fatalf("create credit bucket: %v", err)
	}
	entry := billing.CreditLedgerEntry{
		UserID:         userID,
		BucketID:       &bucket.ID,
		EntryType:      billing.CreditLedgerEntryAdjustment,
		CreditsDelta:   credits,
		IdempotencyKey: "executor_test_credit:" + uuid.NewString(),
	}
	if err := db.Create(&entry).Error; err != nil {
		t.Fatalf("create credit ledger entry: %v", err)
	}
	return bucket
}

type recordingExecutorProcessor struct {
	mu           sync.Mutex
	inputs       []ProcessInput
	response     *MistralResponse
	raw          []byte
	err          error
	delay        time.Duration
	beforeReturn func()
}

func (p *recordingExecutorProcessor) Process(ctx context.Context, input ProcessInput) (*MistralResponse, []byte, error) {
	p.mu.Lock()
	p.inputs = append(p.inputs, input)
	p.mu.Unlock()
	if p.delay > 0 {
		select {
		case <-time.After(p.delay):
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		}
	}
	if p.err != nil {
		if p.beforeReturn != nil {
			p.beforeReturn()
		}
		return nil, nil, p.err
	}
	response := p.response
	if response == nil {
		response = &MistralResponse{Pages: []MistralPage{{Index: 0, Markdown: "# OCR"}}}
	}
	raw := p.raw
	if raw == nil {
		raw = []byte(`{"pages":[{"index":0,"markdown":"# OCR"}]}`)
	}
	if p.beforeReturn != nil {
		p.beforeReturn()
	}
	return response, raw, nil
}

func (p *recordingExecutorProcessor) CallCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.inputs)
}

func (p *recordingExecutorProcessor) InputAt(t *testing.T, index int) ProcessInput {
	t.Helper()
	p.mu.Lock()
	defer p.mu.Unlock()
	if index < 0 || index >= len(p.inputs) {
		t.Fatalf("processor input index %d out of range for %d calls", index, len(p.inputs))
	}
	return p.inputs[index]
}

type fakeWebhookDispatcher struct {
	mu         sync.Mutex
	events     []webhooks.JobEventInput
	err        error
	blockUntil <-chan struct{}
}

func (d *fakeWebhookDispatcher) DispatchJobEvent(_ context.Context, input webhooks.JobEventInput) error {
	d.mu.Lock()
	d.events = append(d.events, cloneWebhookJobEventInput(input))
	d.mu.Unlock()
	if d.blockUntil != nil {
		<-d.blockUntil
	}
	return d.err
}

func (d *fakeWebhookDispatcher) Events() []webhooks.JobEventInput {
	d.mu.Lock()
	defer d.mu.Unlock()
	events := make([]webhooks.JobEventInput, 0, len(d.events))
	for _, event := range d.events {
		events = append(events, cloneWebhookJobEventInput(event))
	}
	return events
}

func waitForWebhookEvents(t *testing.T, dispatcher *fakeWebhookDispatcher, want int) []webhooks.JobEventInput {
	t.Helper()
	var events []webhooks.JobEventInput
	waitForExecutorTestCondition(t, fmt.Sprintf("%d webhook events", want), func() bool {
		events = dispatcher.Events()
		return len(events) >= want
	}, func() string {
		return fmt.Sprintf("last event count %d", len(events))
	})
	if len(events) != want {
		t.Fatalf("webhook events = %#v, want %d", events, want)
	}
	return events
}

func cloneWebhookJobEventInput(input webhooks.JobEventInput) webhooks.JobEventInput {
	clone := input
	if input.UserID != nil {
		userID := *input.UserID
		clone.UserID = &userID
	}
	if input.Job.ErrorMessage != nil {
		message := *input.Job.ErrorMessage
		clone.Job.ErrorMessage = &message
	}
	return clone
}

func assertExecutorWebhookSucceededJobIDOnly(t *testing.T, got webhooks.JobEventInput, wantUserID string, wantJobID uuid.UUID) {
	t.Helper()
	if got.Event != webhooks.EventJobSucceeded {
		t.Fatalf("webhook event = %q, want %q", got.Event, webhooks.EventJobSucceeded)
	}
	if got.UserID == nil || *got.UserID != wantUserID {
		t.Fatalf("webhook user_id = %#v, want %s", got.UserID, wantUserID)
	}
	if got.Job.ID != wantJobID.String() {
		t.Fatalf("webhook job id = %q, want %s", got.Job.ID, wantJobID)
	}
	if got.Job.Status != "" {
		t.Fatalf("webhook job status = %q, want empty", got.Job.Status)
	}
	if got.Job.OriginalFilename != "" {
		t.Fatalf("webhook original_filename = %q, want empty", got.Job.OriginalFilename)
	}
}

func assertExecutorWebhookJobEvent(t *testing.T, got webhooks.JobEventInput, wantEvent webhooks.Event, wantUserID string, wantJobID uuid.UUID, wantStatus OCRJobStatus, wantOriginalFilename string) {
	t.Helper()
	if got.Event != wantEvent {
		t.Fatalf("webhook event = %q, want %q", got.Event, wantEvent)
	}
	if got.UserID == nil || *got.UserID != wantUserID {
		t.Fatalf("webhook user_id = %#v, want %s", got.UserID, wantUserID)
	}
	if got.Job.ID != wantJobID.String() {
		t.Fatalf("webhook job id = %q, want %s", got.Job.ID, wantJobID)
	}
	if got.Job.Status != string(wantStatus) {
		t.Fatalf("webhook job status = %q, want %q", got.Job.Status, wantStatus)
	}
	if got.Job.OriginalFilename != wantOriginalFilename {
		t.Fatalf("webhook original_filename = %q, want %q", got.Job.OriginalFilename, wantOriginalFilename)
	}
}

func assertExecutorJobFailed(t *testing.T, db *gorm.DB, id uuid.UUID, wantMessagePart string) {
	t.Helper()
	var got OCRJob
	if err := db.First(&got, "id = ?", id).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if got.Status != OCRJobStatusFailed {
		t.Fatalf("job status = %q, want %q", got.Status, OCRJobStatusFailed)
	}
	if got.ErrorMessage == "" {
		t.Fatal("job error_message is empty")
	}
	if !strings.Contains(got.ErrorMessage, wantMessagePart) {
		t.Fatalf("job error_message = %q, want it to contain %q", got.ErrorMessage, wantMessagePart)
	}
	if got.DocumentID != nil {
		t.Fatalf("job document_id = %s, want nil", *got.DocumentID)
	}
}

func assertExecutorDocumentCount(t *testing.T, db *gorm.DB, want int64) {
	t.Helper()
	var count int64
	if err := db.Model(&OCRDocument{}).Count(&count).Error; err != nil {
		t.Fatalf("count OCR documents: %v", err)
	}
	if count != want {
		t.Fatalf("OCR document count = %d, want %d", count, want)
	}
}

func assertExecutorCreditBalance(t *testing.T, db *gorm.DB, userID string, want int) {
	t.Helper()
	balance, err := billing.AvailableCredits(context.Background(), db, userID, time.Now().UTC())
	if err != nil {
		t.Fatalf("available credits: %v", err)
	}
	if balance.Available != want {
		t.Fatalf("available credits = %d, want %d", balance.Available, want)
	}
}

func assertExecutorCreditLedgerCount(t *testing.T, db *gorm.DB, want int64, query any, args ...any) {
	t.Helper()
	var count int64
	if err := db.Model(&billing.CreditLedgerEntry{}).Where(query, args...).Count(&count).Error; err != nil {
		t.Fatalf("count credit ledger entries: %v", err)
	}
	if count != want {
		t.Fatalf("credit ledger entry count = %d, want %d", count, want)
	}
}
