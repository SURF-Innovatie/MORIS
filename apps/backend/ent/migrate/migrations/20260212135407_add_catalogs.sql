-- Create "catalogs" table
CREATE TABLE "catalogs" ("id" uuid NOT NULL, "name" character varying NOT NULL, "description" character varying NULL, "rich_description" text NULL, "project_ids" jsonb NULL, "title" character varying NOT NULL, "logo_url" character varying NULL, "primary_color" character varying NULL, "secondary_color" character varying NULL, "accent_color" character varying NULL, "favicon" character varying NULL, "font_family" character varying NULL, PRIMARY KEY ("id"));
-- Create index "catalogs_name_key" to table: "catalogs"
CREATE UNIQUE INDEX "catalogs_name_key" ON "catalogs" ("name");
