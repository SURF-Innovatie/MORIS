-- Modify "persons" table
ALTER TABLE "persons" ADD COLUMN "zenodo_access_token" character varying NULL, ADD COLUMN "zenodo_refresh_token" character varying NULL;
-- Modify "products" table
ALTER TABLE "products" ADD COLUMN "zenodo_deposition_id" bigint NULL;
