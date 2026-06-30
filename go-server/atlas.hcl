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
  }

  diff {
    skip {
      rename_constraint = true
    }
  }
}
