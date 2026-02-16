-- Create "affiliated_organisations" table
CREATE TABLE "affiliated_organisations" ("id" uuid NOT NULL, "name" character varying NOT NULL, "kvk_number" character varying NULL, "ror_id" character varying NULL, "vat_number" character varying NULL, "city" character varying NULL, "country" character varying NULL, PRIMARY KEY ("id"));
