data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./cmd/atlas-loader",
  ]
}

env "local" {
  src = data.external_schema.gorm.url
  dev = getenv("ATLAS_DEV_DATABASE_URL")

  migration {
    dir = "file://migrations"
    exclude = [
      "*[type=extension]",
      "*.fk_ocr_documents_job",
      "*.chk_ocr_jobs_status",
      "*.idx_ocr_documents_original_filename_trgm",
      "*.prevent_credit_ledger_entry_update",
      "*.trg_prevent_credit_ledger_entry_update",
    ]
  }

  diff {
    skip {
      rename_constraint = true
    }
  }
}
