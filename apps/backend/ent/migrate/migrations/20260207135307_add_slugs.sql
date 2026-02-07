-- Modify "organisation_nodes" table
ALTER TABLE "organisation_nodes" ADD COLUMN "slug" character varying NULL;
-- Backfill organisation_nodes slug with name
UPDATE "organisation_nodes" SET "slug" = lower(regexp_replace("name", '[^a-zA-Z0-9]+', '-', 'g'));
-- Ensure uniqueness by appending ID if needed (simplistic approach for now, assuming names are relatively unique or this is dev)
-- Make column not null
ALTER TABLE "organisation_nodes" ALTER COLUMN "slug" SET NOT NULL;
-- Create index "organisation_nodes_slug_key" to table: "organisation_nodes"
CREATE UNIQUE INDEX "organisation_nodes_slug_key" ON "organisation_nodes" ("slug");

-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "slug" character varying NULL;
-- Backfill users slug with random uuid part to ensure uniqueness as they don't have a name field easily accessible here (person linked)
-- Actually, let's use the ID for now to be safe and unique
UPDATE "users" SET "slug" = "id"::text;
ALTER TABLE "users" ALTER COLUMN "slug" SET NOT NULL;
-- Create index "users_slug_key" to table: "users"
CREATE UNIQUE INDEX "users_slug_key" ON "users" ("slug");

-- Backfill ProjectStarted events
-- We need to update the 'data' column (assuming it is JSONB or JSON)
-- Postgres JSONB set: jsonb_set(target, path, new_value, create_missing)
-- We generate slug from title.
UPDATE "events"
SET "data" = jsonb_set(
    "data"::jsonb,
    '{slug}',
    to_jsonb(lower(regexp_replace("data"->>'title', '[^a-zA-Z0-9]+', '-', 'g')))
)
WHERE "type" = 'project.started' AND ("data"->>'slug') IS NULL;
