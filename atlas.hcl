data "external_schema" "gorm" {
   program = [
      "go",
      "run",
      "-mod=mod",
      "./cmd/atlas",
   ]
}

env "local" {
   src = data.external_schema.gorm.url

   url = "postgres://postgres:postgres@localhost:5432/jiandco_db?sslmode=disable"

   dev = "docker://postgres/16/dev?search_path=public"

   migration {
      dir = "file://migrations"
   }
}