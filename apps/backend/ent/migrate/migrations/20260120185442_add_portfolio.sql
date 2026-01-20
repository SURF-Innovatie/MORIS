-- Create "portfolios" table
CREATE TABLE "portfolios" ("id" uuid NOT NULL, "headline" character varying NULL, "summary" character varying NULL, "website" character varying NULL, "show_email" boolean NOT NULL DEFAULT true, "show_orcid" boolean NOT NULL DEFAULT true, "pinned_project_ids" jsonb NULL, "pinned_product_ids" jsonb NULL, "person_id" uuid NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "portfolios_persons_portfolio" FOREIGN KEY ("person_id") REFERENCES "persons" ("id") ON DELETE NO ACTION);
-- Create index "portfolios_person_id_key" to table: "portfolios"
CREATE UNIQUE INDEX "portfolios_person_id_key" ON "portfolios" ("person_id");
