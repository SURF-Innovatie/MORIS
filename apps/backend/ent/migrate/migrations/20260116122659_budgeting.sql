-- Create "api_keys" table
CREATE TABLE "api_keys" ("id" uuid NOT NULL, "name" character varying NOT NULL, "key_hash" character varying NOT NULL, "key_prefix" character varying NOT NULL, "created_at" timestamptz NOT NULL, "last_used_at" timestamptz NULL, "expires_at" timestamptz NULL, "is_active" boolean NOT NULL DEFAULT true, "user_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "api_keys_users_api_keys" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION);
-- Create index "apikey_user_id" to table: "api_keys"
CREATE INDEX "apikey_user_id" ON "api_keys" ("user_id");
-- Create index "apikey_key_prefix" to table: "api_keys"
CREATE INDEX "apikey_key_prefix" ON "api_keys" ("key_prefix");
-- Create index "apikey_key_hash" to table: "api_keys"
CREATE UNIQUE INDEX "apikey_key_hash" ON "api_keys" ("key_hash");
-- Create "budgets" table
CREATE TABLE "budgets" ("id" uuid NOT NULL, "project_id" uuid NOT NULL, "title" character varying NOT NULL, "description" text NULL, "status" character varying NOT NULL DEFAULT 'draft', "total_amount" double precision NOT NULL DEFAULT 0, "currency" character varying NOT NULL DEFAULT 'EUR', "version" bigint NOT NULL DEFAULT 1, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, PRIMARY KEY ("id"));
-- Create index "budget_project_id" to table: "budgets"
CREATE UNIQUE INDEX "budget_project_id" ON "budgets" ("project_id");
-- Create "budget_line_items" table
CREATE TABLE "budget_line_items" ("id" uuid NOT NULL, "category" character varying NOT NULL, "description" character varying NOT NULL, "budgeted_amount" double precision NOT NULL, "year" bigint NOT NULL, "funding_source" character varying NOT NULL, "budget_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "budget_line_items_budgets_line_items" FOREIGN KEY ("budget_id") REFERENCES "budgets" ("id") ON DELETE NO ACTION);
-- Create index "budgetlineitem_budget_id" to table: "budget_line_items"
CREATE INDEX "budgetlineitem_budget_id" ON "budget_line_items" ("budget_id");
-- Create index "budgetlineitem_category" to table: "budget_line_items"
CREATE INDEX "budgetlineitem_category" ON "budget_line_items" ("category");
-- Create index "budgetlineitem_year" to table: "budget_line_items"
CREATE INDEX "budgetlineitem_year" ON "budget_line_items" ("year");
-- Create index "budgetlineitem_funding_source" to table: "budget_line_items"
CREATE INDEX "budgetlineitem_funding_source" ON "budget_line_items" ("funding_source");
-- Create "budget_actuals" table
CREATE TABLE "budget_actuals" ("id" uuid NOT NULL, "amount" double precision NOT NULL, "description" character varying NULL, "recorded_date" timestamptz NOT NULL, "source" character varying NOT NULL DEFAULT 'manual', "external_ref" character varying NULL, "line_item_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "budget_actuals_budget_line_items_actuals" FOREIGN KEY ("line_item_id") REFERENCES "budget_line_items" ("id") ON DELETE NO ACTION);
-- Create index "budgetactual_line_item_id" to table: "budget_actuals"
CREATE INDEX "budgetactual_line_item_id" ON "budget_actuals" ("line_item_id");
-- Create index "budgetactual_recorded_date" to table: "budget_actuals"
CREATE INDEX "budgetactual_recorded_date" ON "budget_actuals" ("recorded_date");
-- Create index "budgetactual_source" to table: "budget_actuals"
CREATE INDEX "budgetactual_source" ON "budget_actuals" ("source");
