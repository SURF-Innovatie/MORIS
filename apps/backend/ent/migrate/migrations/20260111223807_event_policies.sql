-- Create "event_policies" table
CREATE TABLE "event_policies" ("id" uuid NOT NULL, "name" character varying NOT NULL, "description" character varying NULL, "event_types" jsonb NOT NULL, "conditions" jsonb NULL, "action_type" character varying NOT NULL, "message_template" character varying NULL, "recipient_user_ids" jsonb NULL, "recipient_project_role_ids" jsonb NULL, "recipient_org_role_ids" jsonb NULL, "recipient_dynamic" jsonb NULL, "project_id" uuid NULL, "enabled" boolean NOT NULL DEFAULT true, "created_at" timestamptz NOT NULL, "updated_at" timestamptz NOT NULL, "org_node_id" uuid NULL, PRIMARY KEY ("id"), CONSTRAINT "event_policies_organisation_nodes_org_node" FOREIGN KEY ("org_node_id") REFERENCES "organisation_nodes" ("id") ON DELETE SET NULL);
-- Create index "eventpolicy_org_node_id" to table: "event_policies"
CREATE INDEX "eventpolicy_org_node_id" ON "event_policies" ("org_node_id");
-- Create index "eventpolicy_project_id" to table: "event_policies"
CREATE INDEX "eventpolicy_project_id" ON "event_policies" ("project_id");
