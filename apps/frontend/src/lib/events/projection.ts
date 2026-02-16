import { ProjectEventType } from "@/api/events";

/**
 * Applies pending events to a project, projecting what the project state will be
 * after the events are approved and applied. This is critical for event-sourced
 * state management in MORIS.
 *
 * @param project - The current project state
 * @param events - Array of pending events to apply
 * @returns A new project object with pending events applied
 */
export function applyPendingEvents(project: any, events: any[]): any {
  if (!events || events.length === 0) return project;

  const p = JSON.parse(JSON.stringify(project));

  for (const e of events) {
    if (e.status !== "pending") continue;

    switch (e.type) {
      case ProjectEventType.TitleChanged:
        if (e.data?.title) p.title = e.data.title;
        break;
      case ProjectEventType.DescriptionChanged:
        if (e.data?.description) p.description = e.data.description;
        break;
      case ProjectEventType.StartDateChanged:
        if (e.data?.start_date) p.start_date = e.data.start_date;
        break;
      case ProjectEventType.EndDateChanged:
        if (e.data?.end_date) p.end_date = e.data.end_date;
        break;
      case ProjectEventType.OwningOrgNodeChanged:
        if (e.data?.owning_org_node_id) {
          p.owning_org_node = {
            ...(p.owning_org_node || {}),
            id: e.data.owning_org_node_id,
          };
        }
        break;
      case ProjectEventType.CustomFieldValueSet:
        if (e.data?.definition_id) {
          p.custom_fields = p.custom_fields || {};
          p.custom_fields[e.data.definition_id] = e.data.value;
        }
        break;
      case ProjectEventType.ProductAdded:
        if (e.product && e.product.id) {
          p.products = p.products || [];
          p.products.push({ ...e.product, pending: true });
        }
        break;
      case ProjectEventType.ProductRemoved:
        if (e.data?.product_id) {
          p.products = (p.products || []).filter(
            (prod: any) => prod.id !== e.data.product_id,
          );
        }
        break;
      case ProjectEventType.ProjectRoleAssigned:
        if (e.person && e.projectRole) {
          p.members = p.members || [];
          const newMember = {
            id: `pending-${e.person.id}-${e.projectRole.id}`,
            user_id: e.person.id,
            name: `${e.person.givenName} ${e.person.familyName}`.trim(),
            email: e.person.email,
            avatarUrl: e.person.avatarUrl,
            role: e.projectRole.slug,
            role_id: e.projectRole.id,
            role_name: e.projectRole.name,
            pending: true,
          };
          p.members.push(newMember);
        }
        break;
      case ProjectEventType.ProjectRoleUnassigned:
        if (e.data?.person_id && e.data?.project_role_id) {
          p.members = (p.members || []).filter(
            (m: any) =>
              !(
                m.user_id === e.data.person_id &&
                m.role_id === e.data.project_role_id
              ),
          );
        }
        break;
    }
  }
  return p;
}
