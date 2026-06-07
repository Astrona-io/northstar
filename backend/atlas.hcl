data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "cmd/atlas-provider/main.go",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "sqlite://dev?mode=memory"
  migration {
    dir = "file://internal/database/migrations"
  }
  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
