import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Loader2,
  Trash2,
  Bell,
  ShieldCheck,
  ArrowUpRight,
  Users,
  Edit,
} from "lucide-react";
import {
  EventPolicyResponse,
  EventTypeInfo,
  ProjectRoleResponse,
  OrganisationRoleResponse,
} from "@api/model";

interface PolicyCardProps {
  policy: EventPolicyResponse;
  eventTypes: EventTypeInfo[];
  projectRoles: ProjectRoleResponse[];
  orgRoles?: OrganisationRoleResponse[];
  onEdit?: () => void;
  onDelete?: () => void;
  isDeleting?: boolean;
  inherited?: boolean;
}

export function PolicyCard({
  policy,
  eventTypes,
  projectRoles,
  orgRoles = [],
  onEdit,
  onDelete,
  isDeleting,
  inherited,
}: PolicyCardProps) {
  const getEventTypeLabel = (type: string) => {
    const found = eventTypes.find((e) => e.type === type);
    return found?.friendlyName || type;
  };

  return (
    <Card className={inherited ? "opacity-75" : ""}>
      <CardHeader className="pb-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            {policy.action_type === "notify" ? (
              <Bell className="h-5 w-5 text-blue-500" />
            ) : (
              <ShieldCheck className="h-5 w-5 text-amber-500" />
            )}
            <div>
              <CardTitle className="text-base flex items-center gap-2">
                {policy.name}
                {!policy.enabled && <Badge variant="secondary">Disabled</Badge>}
              </CardTitle>
              {inherited && policy.source_org_node_name && (
                <CardDescription className="flex items-center gap-1">
                  <ArrowUpRight className="h-3 w-3" />
                  Inherited from {policy.source_org_node_name}
                </CardDescription>
              )}
            </div>
          </div>
          {!inherited && (
            <div className="flex gap-1">
              {onEdit && (
                <Button variant="ghost" size="icon" onClick={onEdit}>
                  <Edit className="h-4 w-4" />
                </Button>
              )}
              {onDelete && (
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={onDelete}
                  disabled={isDeleting}
                >
                  {isDeleting ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <Trash2 className="h-4 w-4 text-destructive" />
                  )}
                </Button>
              )}
            </div>
          )}
        </div>
      </CardHeader>
      <CardContent className="pt-0">
        <div className="space-y-2 text-sm">
          <div className="flex flex-wrap gap-1">
            {policy.event_types?.map((type) => (
              <Badge key={type} variant="outline">
                {getEventTypeLabel(type)}
              </Badge>
            ))}
          </div>
          {policy.description && (
            <p className="text-muted-foreground">{policy.description}</p>
          )}
          <div className="flex gap-2 text-muted-foreground">
            <span>Action:</span>
            <span className="font-medium">
              {policy.action_type === "notify"
                ? "Send Notification"
                : "Request Approval"}
            </span>
          </div>
          {(policy.recipient_user_ids?.length ?? 0) > 0 && (
            <div className="flex items-center gap-2 text-muted-foreground">
              <Users className="h-3 w-3" />
              <span>{policy.recipient_user_ids?.length} specific user(s)</span>
            </div>
          )}

          {(policy.recipient_project_role_ids?.length ?? 0) > 0 && (
            <div className="flex flex-wrap gap-1 mt-1">
              {policy.recipient_project_role_ids?.map((roleId) => {
                const role = projectRoles.find((r) => r.id === roleId);
                return role ? (
                  <Badge
                    key={roleId}
                    variant="secondary"
                    className="text-xs bg-blue-100 text-blue-800 hover:bg-blue-100/80 dark:bg-blue-900/50 dark:text-blue-300"
                  >
                    Project Role: {role.name}
                  </Badge>
                ) : null;
              })}
            </div>
          )}

          {(policy.recipient_org_role_ids?.length ?? 0) > 0 && (
            <div className="flex flex-wrap gap-1 mt-1">
              {policy.recipient_org_role_ids?.map((roleId) => {
                const role = orgRoles?.find((r) => r.id === roleId);
                return role ? (
                  <Badge
                    key={roleId}
                    variant="secondary"
                    className="text-xs bg-purple-100 text-purple-800 hover:bg-purple-100/80 dark:bg-purple-900/50 dark:text-purple-300"
                  >
                    Org Role: {role.displayName}
                  </Badge>
                ) : null;
              })}
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
