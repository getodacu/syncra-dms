package ocr

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const OCRJobsNotifyChannel = "ocr_jobs"

type JobNotifier func(context.Context, *gorm.DB, uuid.UUID) error

func NotifyOCRJobQueued(ctx context.Context, db *gorm.DB, id uuid.UUID) error {
	return notifyOCRJobQueued(ctx, db, OCRJobsNotifyChannel, id)
}

func notifyOCRJobQueued(ctx context.Context, db *gorm.DB, channel string, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("OCR job id is required")
	}
	channel = strings.TrimSpace(channel)
	if channel == "" {
		return errors.New("OCR job notify channel is required")
	}
	return db.WithContext(ctx).
		Exec("SELECT pg_notify(?, ?)", channel, id.String()).
		Error
}
