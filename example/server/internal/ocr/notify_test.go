package ocr

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gorm.io/gorm"

	"ai.ro/syncra/internal/testsupport"
)

func TestNotifyOCRJobQueuedPublishesJobID(t *testing.T) {
	db := testsupport.OpenPostgresDB(t)
	dsn := testsupport.PostgresTestDSN(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatalf("connect listener: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close(context.Background()) })
	if _, err := conn.Exec(ctx, "LISTEN "+OCRJobsNotifyChannel); err != nil {
		t.Fatalf("listen: %v", err)
	}

	jobID := uuid.New()
	if err := NotifyOCRJobQueued(ctx, db, jobID); err != nil {
		t.Fatalf("NotifyOCRJobQueued() error = %v", err)
	}

	notification, err := conn.WaitForNotification(ctx)
	if err != nil {
		t.Fatalf("wait notification: %v", err)
	}
	if notification.Channel != OCRJobsNotifyChannel {
		t.Fatalf("channel = %q, want %q", notification.Channel, OCRJobsNotifyChannel)
	}
	if notification.Payload != jobID.String() {
		t.Fatalf("payload = %q, want %q", notification.Payload, jobID)
	}
}

func TestNotifyOCRJobQueuedRejectsNilID(t *testing.T) {
	err := NotifyOCRJobQueued(context.Background(), &gorm.DB{}, uuid.Nil)
	if err == nil {
		t.Fatal("NotifyOCRJobQueued() error = nil, want error")
	}
}
