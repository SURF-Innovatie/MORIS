import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  useGetOrganisationsIdPolicies,
  usePostOrganisationsIdPolicies,
  useDeletePoliciesId,
  usePutPoliciesId,
  useGetEventTypes,
  getGetOrganisationsIdPoliciesQueryKey,
  useGetOrganisationNodesIdOrganisationRoles,
  useGetOrganisationNodesIdRoles,
} from "@api/moris";
import {
  EventPolicyRequest,
  EventPolicyResponse,
  EventTypeInfo,
  OrganisationRoleResponse,
  ProjectRoleResponse,
} from "@api/model";
import { MultiUserSelect } from "@/components/user/MultiUserSelect";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import {
  Loader2,
  Plus,
  Trash2,
  Bell,
  ShieldCheck,
  ArrowUpRight,
  Users,
  Edit,
} from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

const DYNAMIC_RECIPIENTS = [
  { value: "project_members", label: "All Project Members" },
  { value: "project_owner", label: "Project Owner" },
  { value: "org_admins", label: "Organisation Admins" },
];

interface EventPoliciesTabProps {
  nodeId: string;
}

export function EventPoliciesTab({ nodeId }: EventPoliciesTabProps) {
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editPolicy, setEditPolicy] = useState<EventPolicyResponse | null>(
    null
  );
  const queryClient = useQueryClient();
  const { toast } = useToast();

  // Fetch event types from backend
  const { data: eventTypes = [] } = useGetEventTypes();

  // Fetch org roles for this node
  const { data: orgRoles = [] } =
    useGetOrganisationNodesIdOrganisationRoles(nodeId);

  // Fetch project roles available at this node (includes inherited)
  const { data: projectRoles = [] } = useGetOrganisationNodesIdRoles(nodeId);

  const { data: policies, isLoading } = useGetOrganisationsIdPolicies(nodeId, {
    inherited: true,
  });

  const deleteMutation = useDeletePoliciesId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetOrganisationsIdPoliciesQueryKey(nodeId),
        });
        toast({ title: "Policy deleted" });
      },
      onError: () => {
        toast({ title: "Failed to delete policy", variant: "destructive" });
      },
    },
  });

  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const ownPolicies =
    policies?.filter((p: EventPolicyResponse) => !p.inherited) || [];
  const inheritedPolicies =
    policies?.filter((p: EventPolicyResponse) => p.inherited) || [];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-xl font-semibold">Event Policies</h2>
          <p className="text-sm text-muted-foreground">
            Configure automatic notifications and approval workflows
          </p>
        </div>
        <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              New Policy
            </Button>
          </DialogTrigger>
          <PolicyFormDialog
            nodeId={nodeId}
            eventTypes={eventTypes}
            orgRoles={orgRoles}
            projectRoles={projectRoles}
            onClose={() => setCreateDialogOpen(false)}
          />
        </Dialog>
      </div>

      {/* Own Policies */}
      <div className="space-y-4">
        <h3 className="text-lg font-medium">Organisation Policies</h3>
        {ownPolicies.length === 0 ? (
          <Card>
            <CardContent className="py-8 text-center text-muted-foreground">
              No policies configured. Create one to get started.
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {ownPolicies.map((policy: EventPolicyResponse) => (
              <PolicyCard
                key={policy.id}
                policy={policy}
                eventTypes={eventTypes}
                projectRoles={projectRoles}
                orgRoles={orgRoles}
                onEdit={() => setEditPolicy(policy)}
                onDelete={() => deleteMutation.mutate({ id: policy.id! })}
                isDeleting={deleteMutation.isPending}
              />
            ))}
          </div>
        )}
      </div>

      {/* Inherited Policies */}
      {inheritedPolicies.length > 0 && (
        <div className="space-y-4">
          <h3 className="text-lg font-medium flex items-center gap-2">
            Inherited Policies
            <Badge variant="outline" className="font-normal">
              Read-only
            </Badge>
          </h3>
          <div className="grid gap-4">
            {inheritedPolicies.map((policy: EventPolicyResponse) => (
              <PolicyCard
                key={policy.id}
                policy={policy}
                eventTypes={eventTypes}
                projectRoles={projectRoles}
                orgRoles={orgRoles}
                inherited
              />
            ))}
          </div>
        </div>
      )}

      {/* Edit Dialog */}
      <Dialog
        open={!!editPolicy}
        onOpenChange={(open) => !open && setEditPolicy(null)}
      >
        {editPolicy && (
          <PolicyFormDialog
            nodeId={nodeId}
            eventTypes={eventTypes}
            orgRoles={orgRoles}
            projectRoles={projectRoles}
            policy={editPolicy}
            onClose={() => setEditPolicy(null)}
          />
        )}
      </Dialog>
    </div>
  );
}

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

function PolicyCard({
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

interface PolicyFormDialogProps {
  nodeId: string;
  eventTypes: EventTypeInfo[];
  orgRoles: OrganisationRoleResponse[];
  projectRoles: ProjectRoleResponse[];
  policy?: EventPolicyResponse;
  onClose: () => void;
}

function PolicyFormDialog({
  nodeId,
  eventTypes,
  orgRoles,
  projectRoles,
  policy,
  onClose,
}: PolicyFormDialogProps) {
  const isEditing = !!policy;

  const [name, setName] = useState(policy?.name || "");
  const [description, setDescription] = useState(policy?.description || "");
  const [selectedEventTypes, setSelectedEventTypes] = useState<string[]>(
    policy?.event_types || []
  );
  const [actionType, setActionType] = useState<"notify" | "request_approval">(
    (policy?.action_type as "notify" | "request_approval") || "notify"
  );
  const [dynamicRecipients, setDynamicRecipients] = useState<string[]>(
    policy?.recipient_dynamic || []
  );
  const [specificUsers, setSpecificUsers] = useState<string[]>(
    policy?.recipient_user_ids || []
  );
  const [selectedOrgRoles, setSelectedOrgRoles] = useState<string[]>(
    policy?.recipient_org_role_ids || []
  );
  const [selectedProjectRoles, setSelectedProjectRoles] = useState<string[]>(
    policy?.recipient_project_role_ids || []
  );
  const [enabled, setEnabled] = useState(policy?.enabled ?? true);

  const queryClient = useQueryClient();
  const { toast } = useToast();

  const createMutation = usePostOrganisationsIdPolicies({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetOrganisationsIdPoliciesQueryKey(nodeId),
        });
        toast({ title: "Policy created" });
        onClose();
      },
      onError: () => {
        toast({ title: "Failed to create policy", variant: "destructive" });
      },
    },
  });

  const updateMutation = usePutPoliciesId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetOrganisationsIdPoliciesQueryKey(nodeId),
        });
        toast({ title: "Policy updated" });
        onClose();
      },
      onError: () => {
        toast({ title: "Failed to update policy", variant: "destructive" });
      },
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const data: EventPolicyRequest = {
      name,
      description: description || undefined,
      event_types: selectedEventTypes,
      action_type: actionType,
      recipient_dynamic:
        dynamicRecipients.length > 0 ? dynamicRecipients : undefined,
      recipient_user_ids: specificUsers.length > 0 ? specificUsers : undefined,
      recipient_org_role_ids:
        selectedOrgRoles.length > 0 ? selectedOrgRoles : undefined,
      recipient_project_role_ids:
        selectedProjectRoles.length > 0 ? selectedProjectRoles : undefined,
      enabled,
    };

    if (isEditing && policy?.id) {
      updateMutation.mutate({ id: policy.id, data });
    } else {
      createMutation.mutate({ id: nodeId, data });
    }
  };

  const toggleEventType = (type: string) => {
    setSelectedEventTypes((prev) =>
      prev.includes(type) ? prev.filter((t) => t !== type) : [...prev, type]
    );
  };

  const toggleDynamicRecipient = (type: string) => {
    setDynamicRecipients((prev) =>
      prev.includes(type) ? prev.filter((t) => t !== type) : [...prev, type]
    );
  };

  const toggleOrgRole = (roleId: string) => {
    setSelectedOrgRoles((prev) =>
      prev.includes(roleId)
        ? prev.filter((r) => r !== roleId)
        : [...prev, roleId]
    );
  };

  const toggleProjectRole = (roleId: string) => {
    setSelectedProjectRoles((prev) =>
      prev.includes(roleId)
        ? prev.filter((r) => r !== roleId)
        : [...prev, roleId]
    );
  };

  const isPending = createMutation.isPending || updateMutation.isPending;

  return (
    <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>
          {isEditing ? "Edit Event Policy" : "Create Event Policy"}
        </DialogTitle>
      </DialogHeader>
      <form
        onSubmit={handleSubmit}
        className="grid grid-cols-1 lg:grid-cols-2 gap-8 h-full"
      >
        <div className="space-y-6 flex flex-col h-full">
          <div className="space-y-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g., Notify on project creation"
              required
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description (optional)</Label>
            <Textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="What does this policy do?"
              rows={3}
            />
          </div>

          <div className="space-y-2 flex-1 flex flex-col">
            <Label>Trigger on Events</Label>
            <div className="flex flex-wrap gap-2 flex-1 min-h-[400px] max-h-[60vh] overflow-y-auto p-2 border rounded-md content-start">
              {eventTypes.map((type) => (
                <Badge
                  key={type.type}
                  variant={
                    selectedEventTypes.includes(type.type!)
                      ? "default"
                      : "outline"
                  }
                  className="cursor-pointer hover:bg-primary/90 h-fit"
                  onClick={() => toggleEventType(type.type!)}
                >
                  {type.friendlyName || type.type}
                </Badge>
              ))}
            </div>
          </div>
        </div>

        <div className="space-y-6">
          <div className="space-y-2">
            <Label>Action Type</Label>
            <Select
              value={actionType}
              onValueChange={(val) =>
                setActionType(val as "notify" | "request_approval")
              }
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="notify">
                  <div className="flex items-center gap-2">
                    <Bell className="h-4 w-4" />
                    Send Notification
                  </div>
                </SelectItem>
                <SelectItem value="request_approval">
                  <div className="flex items-center gap-2">
                    <ShieldCheck className="h-4 w-4" />
                    Request Approval
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="border rounded-lg p-6 space-y-6 ">
            <h3 className="font-semibold text-lg">Recipients</h3>

            <div className="space-y-4">
              <div className="space-y-2">
                <Label className="text-sm font-medium text-muted-foreground">
                  Dynamic Groups
                </Label>
                <div className="flex flex-wrap gap-2">
                  {DYNAMIC_RECIPIENTS.map((type) => (
                    <Badge
                      key={type.value}
                      variant={
                        dynamicRecipients.includes(type.value)
                          ? "default"
                          : "outline"
                      }
                      className="cursor-pointer hover:bg-primary/90"
                      onClick={() => toggleDynamicRecipient(type.value)}
                    >
                      {type.label}
                    </Badge>
                  ))}
                </div>
              </div>

              <div className="space-y-2">
                <Label className="text-sm font-medium text-muted-foreground">
                  Specific Users
                </Label>
                <MultiUserSelect
                  value={specificUsers}
                  onChange={setSpecificUsers}
                  placeholder="Search users..."
                />
              </div>

              {projectRoles.length > 0 && (
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-muted-foreground">
                    Project Roles
                  </Label>
                  <div className="flex flex-wrap gap-2">
                    {projectRoles.map((role) => (
                      <Badge
                        key={role.id}
                        variant={
                          selectedProjectRoles.includes(role.id!)
                            ? "default"
                            : "outline"
                        }
                        className="cursor-pointer hover:bg-primary/90"
                        onClick={() => toggleProjectRole(role.id!)}
                      >
                        {role.name}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}

              {orgRoles.length > 0 && (
                <div className="space-y-2">
                  <Label className="text-sm font-medium text-muted-foreground">
                    Organisation Roles
                  </Label>
                  <div className="flex flex-wrap gap-2">
                    {orgRoles.map((role) => (
                      <Badge
                        key={role.id}
                        variant={
                          selectedOrgRoles.includes(role.id!)
                            ? "default"
                            : "outline"
                        }
                        className="cursor-pointer hover:bg-primary/90"
                        onClick={() => toggleOrgRole(role.id!)}
                      >
                        {role.displayName}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>

          <div className="flex items-center justify-between border-t pt-4">
            <div className="space-y-0.5">
              <Label htmlFor="enabled" className="text-base">
                Policy Enabled
              </Label>
              <p className="text-sm text-muted-foreground">
                Turn off to temporarily disable this policy
              </p>
            </div>
            <Switch
              id="enabled"
              checked={enabled}
              onCheckedChange={setEnabled}
            />
          </div>
        </div>

        <DialogFooter className="col-span-1 lg:col-span-2 mt-6">
          <Button type="button" variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={!name || selectedEventTypes.length === 0 || isPending}
          >
            {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {isEditing ? "Save Changes" : "Create Policy"}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  );
}
