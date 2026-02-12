-- Create index "products_doi_key" to table: "products"
CREATE UNIQUE INDEX "products_doi_key" ON "products" ("doi");
-- Create "product_authors" table
CREATE TABLE "product_authors" ("product_id" uuid NOT NULL, "person_id" uuid NOT NULL, PRIMARY KEY ("product_id", "person_id"), CONSTRAINT "product_authors_product_id" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE, CONSTRAINT "product_authors_person_id" FOREIGN KEY ("person_id") REFERENCES "persons" ("id") ON DELETE CASCADE);
