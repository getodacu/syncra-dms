-- Modify "ocr_jobs" table
ALTER TABLE "ocr_jobs" ADD COLUMN "deleted_at" timestamptz;
-- Create index "idx_ocr_jobs_deleted_at" to table: "ocr_jobs"
CREATE INDEX "idx_ocr_jobs_deleted_at" ON "ocr_jobs" ("deleted_at");
