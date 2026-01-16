-- Create "error_logs" table
CREATE TABLE "error_logs"
(
    "id"            uuid              NOT NULL,
    "user_id"       character varying NULL,
    "http_method"   character varying NOT NULL,
    "route"         character varying NOT NULL,
    "status_code"   bigint            NOT NULL,
    "error_message" text              NOT NULL,
    "stack_trace"   text NULL,
    "timestamp"     timestamptz       NOT NULL,
    PRIMARY KEY ("id")
);
-- Create "organisation_nodes" table
CREATE TABLE "organisation_nodes"
(
    "id"          uuid              NOT NULL,
    "name"        character varying NOT NULL,
    "description" character varying NULL,
    "avatar_url"  character varying NULL,
    "ror_id"      character varying NULL,
    "parent_id"   uuid NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "organisation_nodes_organisation_nodes_children" FOREIGN KEY ("parent_id") REFERENCES "organisation_nodes" ("id") ON DELETE SET NULL
);
-- Create index "organisationnode_parent_id" to table: "organisation_nodes"
CREATE INDEX "organisationnode_parent_id" ON "organisation_nodes" ("parent_id");
-- Create "custom_field_definitions" table
CREATE TABLE "custom_field_definitions"
(
    "id"                   uuid              NOT NULL,
    "name"                 character varying NOT NULL,
    "type"                 character varying NOT NULL,
    "category"             character varying NOT NULL DEFAULT 'PROJECT',
    "description"          character varying NULL,
    "required"             boolean           NOT NULL DEFAULT false,
    "validation_regex"     character varying NULL,
    "example_value"        character varying NULL,
    "organisation_node_id" uuid              NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "custom_field_definitions_organ_727e91b3b49f6dc828c8df4b773e34ec" FOREIGN KEY ("organisation_node_id") REFERENCES "organisation_nodes" ("id") ON DELETE NO ACTION
);
-- Create "persons" table
CREATE TABLE "persons"
(
    "id"                uuid              NOT NULL,
    "user_id"           uuid NULL,
    "orcid_id"          character varying NULL,
    "name"              character varying NOT NULL,
    "given_name"        character varying NULL,
    "family_name"       character varying NULL,
    "email"             character varying NOT NULL,
    "avatar_url"        character varying NULL,
    "description"       character varying NULL,
    "org_custom_fields" jsonb NULL,
    PRIMARY KEY ("id")
);
-- Create index "persons_user_id_key" to table: "persons"
CREATE UNIQUE INDEX "persons_user_id_key" ON "persons" ("user_id");
-- Create index "persons_orcid_id_key" to table: "persons"
CREATE UNIQUE INDEX "persons_orcid_id_key" ON "persons" ("orcid_id");
-- Create index "persons_email_key" to table: "persons"
CREATE UNIQUE INDEX "persons_email_key" ON "persons" ("email");
-- Create "organisation_roles" table
CREATE TABLE "organisation_roles"
(
    "id"                   uuid              NOT NULL,
    "key"                  character varying NOT NULL,
    "display_name"         character varying NOT NULL,
    "description"          character varying NULL,
    "permissions"          jsonb NULL,
    "organisation_node_id" uuid              NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "organisation_roles_organisation_nodes_organisation_roles" FOREIGN KEY ("organisation_node_id") REFERENCES "organisation_nodes" ("id") ON DELETE NO ACTION
);
-- Create index "organisationrole_key_organisation_node_id" to table: "organisation_roles"
CREATE UNIQUE INDEX "organisationrole_key_organisation_node_id" ON "organisation_roles" ("key", "organisation_node_id");
-- Create "role_scopes" table
CREATE TABLE "role_scopes"
(
    "id"           uuid NOT NULL,
    "role_id"      uuid NOT NULL,
    "root_node_id" uuid NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "role_scopes_organisation_roles_scopes" FOREIGN KEY ("role_id") REFERENCES "organisation_roles" ("id") ON DELETE NO ACTION,
    CONSTRAINT "role_scopes_organisation_nodes_root_node" FOREIGN KEY ("root_node_id") REFERENCES "organisation_nodes" ("id") ON DELETE NO ACTION
);
-- Create index "rolescope_role_id_root_node_id" to table: "role_scopes"
CREATE UNIQUE INDEX "rolescope_role_id_root_node_id" ON "role_scopes" ("role_id", "root_node_id");
-- Create "memberships" table
CREATE TABLE "memberships"
(
    "id"            uuid NOT NULL,
    "person_id"     uuid NOT NULL,
    "role_scope_id" uuid NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "memberships_persons_person" FOREIGN KEY ("person_id") REFERENCES "persons" ("id") ON DELETE NO ACTION,
    CONSTRAINT "memberships_role_scopes_role_scope" FOREIGN KEY ("role_scope_id") REFERENCES "role_scopes" ("id") ON DELETE NO ACTION
);
-- Create index "membership_person_id_role_scope_id" to table: "memberships"
CREATE UNIQUE INDEX "membership_person_id_role_scope_id" ON "memberships" ("person_id", "role_scope_id");
-- Create "events" table
CREATE TABLE "events"
(
    "id"          uuid              NOT NULL,
    "project_id"  uuid              NOT NULL,
    "version"     bigint            NOT NULL,
    "type"        character varying NOT NULL,
    "status"      character varying NOT NULL DEFAULT 'pending',
    "created_by"  uuid NULL,
    "occurred_at" timestamptz       NOT NULL,
    "data"        jsonb             NOT NULL,
    PRIMARY KEY ("id")
);
-- Create "users" table
CREATE TABLE "users"
(
    "id"           uuid              NOT NULL,
    "person_id"    uuid              NOT NULL,
    "password"     character varying NOT NULL,
    "is_sys_admin" boolean           NOT NULL DEFAULT false,
    "is_active"    boolean           NOT NULL DEFAULT true,
    PRIMARY KEY ("id")
);
-- Create index "users_person_id_key" to table: "users"
CREATE UNIQUE INDEX "users_person_id_key" ON "users" ("person_id");
-- Create "notifications" table
CREATE TABLE "notifications"
(
    "id"       uuid              NOT NULL,
    "message"  character varying NOT NULL,
    "type"     character varying NOT NULL DEFAULT 'info',
    "read"     boolean           NOT NULL DEFAULT false,
    "sent_at"  timestamptz       NOT NULL,
    "event_id" uuid NULL,
    "user_id"  uuid              NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "notifications_events_notifications" FOREIGN KEY ("event_id") REFERENCES "events" ("id") ON DELETE SET NULL,
    CONSTRAINT "notifications_users_notifications" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION
);
-- Create "organisation_node_closures" table
CREATE TABLE "organisation_node_closures"
(
    "id"            uuid   NOT NULL,
    "depth"         bigint NOT NULL,
    "ancestor_id"   uuid   NOT NULL,
    "descendant_id" uuid   NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "organisation_node_closures_organisation_nodes_ancestor" FOREIGN KEY ("ancestor_id") REFERENCES "organisation_nodes" ("id") ON DELETE NO ACTION,
    CONSTRAINT "organisation_node_closures_organisation_nodes_descendant" FOREIGN KEY ("descendant_id") REFERENCES "organisation_nodes" ("id") ON DELETE NO ACTION
);
-- Create index "organisationnodeclosure_ancestor_id_descendant_id" to table: "organisation_node_closures"
CREATE UNIQUE INDEX "organisationnodeclosure_ancestor_id_descendant_id" ON "organisation_node_closures" ("ancestor_id", "descendant_id");
-- Create "products" table
CREATE TABLE "products"
(
    "id"       uuid              NOT NULL,
    "name"     character varying NOT NULL,
    "language" character varying NULL,
    "type"     bigint NULL,
    "doi"      character varying NULL,
    PRIMARY KEY ("id")
);
-- Create "person_products" table
CREATE TABLE "person_products"
(
    "person_id"  uuid NOT NULL,
    "product_id" uuid NOT NULL,
    PRIMARY KEY ("person_id", "product_id"),
    CONSTRAINT "person_products_person_id" FOREIGN KEY ("person_id") REFERENCES "persons" ("id") ON DELETE CASCADE,
    CONSTRAINT "person_products_product_id" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON DELETE CASCADE
);
-- Create "project_roles" table
CREATE TABLE "project_roles"
(
    "id"                   uuid              NOT NULL,
    "key"                  character varying NOT NULL,
    "name"                 character varying NOT NULL,
    "archived_at"          timestamptz NULL,
    "allowed_event_types"  jsonb NULL,
    "organisation_node_id" uuid              NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "project_roles_organisation_nodes_project_roles" FOREIGN KEY ("organisation_node_id") REFERENCES "organisation_nodes" ("id") ON DELETE NO ACTION
);
-- Create index "projectrole_key_organisation_node_id" to table: "project_roles"
CREATE UNIQUE INDEX "projectrole_key_organisation_node_id" ON "project_roles" ("key", "organisation_node_id");
