-- Modify "organisation_nodes" table
ALTER TABLE "organisation_nodes" DROP COLUMN "slug";
-- Modify "portfolios" table
ALTER TABLE "portfolios" ADD COLUMN "recent_project_ids" jsonb NULL;
-- Modify "users" table
ALTER TABLE "users" DROP COLUMN "slug";
