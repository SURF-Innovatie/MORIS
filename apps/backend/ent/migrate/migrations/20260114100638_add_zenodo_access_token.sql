-- Modify "persons" table
ALTER TABLE "persons" DROP COLUMN "zenodo_access_token", DROP COLUMN "zenodo_refresh_token";
-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "zenodo_access_token" character varying NULL, ADD COLUMN "zenodo_refresh_token" character varying NULL;
