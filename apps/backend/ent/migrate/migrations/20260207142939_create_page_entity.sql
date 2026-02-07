-- Create "pages" table
CREATE TABLE "pages" ("id" uuid NOT NULL, "title" character varying NOT NULL, "slug" character varying NOT NULL, "content" jsonb NULL, "type" character varying NOT NULL DEFAULT 'project', "is_published" boolean NOT NULL DEFAULT false, "project_id" uuid NULL, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "user_id" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "pages_users_pages" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE SET NULL);
-- Create index "page_slug" to table: "pages"
CREATE UNIQUE INDEX "page_slug" ON "pages" ("slug");
-- Create index "page_project_id" to table: "pages"
CREATE INDEX "page_project_id" ON "pages" ("project_id");
-- Create index "page_user_id" to table: "pages"
CREATE INDEX "page_user_id" ON "pages" ("user_id");
