# Atlas configuration for versioned migrations
# https://atlasgo.io/atlas-schema/projects

variable "db_user" {
  type    = string
  default = getenv("DB_USER")
}

variable "db_password" {
  type    = string
  default = getenv("DB_PASSWORD")
}

variable "db_host" {
  type    = string
  default = getenv("DB_HOST")
}

variable "db_port" {
  type    = string
  default = getenv("DB_PORT")
}

variable "db_name" {
  type    = string
  default = getenv("DB_NAME")
}

env "local" {
  # Source schema from ent
  src = "ent://ent/schema"
  
  # Migration directory
  migration {
    dir = "file://ent/migrate/migrations"
  }
  
  # Dev database for computing schema diffs
  # Uses a separate 'atlas_dev' database on the same server
  # Create with: CREATE DATABASE atlas_dev;
  dev = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/atlas_dev?sslmode=disable"
  
  # Local development database URL (for applying migrations)
  url = "postgres://${var.db_user}:${var.db_password}@${var.db_host}:${var.db_port}/${var.db_name}?sslmode=disable"
}
