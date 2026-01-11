import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  useGetProjectsIdPolicies,
  useGetEventTypes,
  getGetProjectsIdPoliciesQueryKey,
  useGetProjectsIdRoles,
  useGetOrganisationNodesIdOrganisationRoles,
} from "@api/moris";
import {
  EventPolicyResponse,
  EventTypeInfo,
  ProjectRoleResponse,
  OrganisationRoleResponse,
} from "@api/model";
import {
  createEventPolicyAddedEvent,
  createEventPolicyRemovedEvent,
  EventPolicyAddedInput,
  ProjectEventType,
} from "@/api/events";
import { useAccess } from "@/context/AccessContext";
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
import { Checkbox } from "@/components/ui/checkbox";
import { Badge } from "@/components/ui/badge";
import {
  Loader2,
  Plus,
  Trash2,
  Bell,
  ShieldCheck,
  ArrowUpRight,
  Lock,
  Users,
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

interface ProjectEventPoliciesTabProps {
  projectId: string;
  orgNodeId: string;
}

interface CreateProjectPolicyDialogProps {
  projectId: string;
  orgNodeId: string;
  eventTypes: EventTypeInfo[];
  projectRoles: ProjectRoleResponse[];
  orgRoles: OrganisationRoleResponse[];
  onClose: () => void;
}

export function ProjectEventPoliciesTab({
  projectId,
  orgNodeId,
}: ProjectEventPoliciesTabProps) {
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState<string | null>(null);
  const queryClient = useQueryClient();
  const { toast } = useToast();
  const { hasAccess } = useAccess();

  const canAddPolicy = hasAccess(ProjectEventType.EventPolicyAdded);
  const canRemovePolicy = hasAccess(ProjectEventType.EventPolicyRemoved);

  // Fetch event types from backend
  const { data: eventTypes = [] } = useGetEventTypes();

  // Fetch project roles
  const { data: projectRoles = [] } = useGetProjectsIdRoles(projectId);

  // Fetch org roles
  const { data: orgRoles = [] } =
    useGetOrganisationNodesIdOrganisationRoles(orgNodeId);

  const { data: policies, isLoading } = useGetProjectsIdPolicies(projectId, {
    org_node_id: orgNodeId,
    inherited: true,
  });

  const handleDeletePolicy = async (policy: EventPolicyResponse) => {
    if (!policy.id) return;
    setIsDeleting(policy.id);
    try {
      await createEventPolicyRemovedEvent(projectId, {
        policy_id: policy.id,
        name: policy.name || "",
      });
      queryClient.invalidateQueries({
        queryKey: getGetProjectsIdPoliciesQueryKey(projectId),
      });
      toast({ title: "Policy removal requested" });
    } catch (error) {
      toast({ title: "Failed to remove policy", variant: "destructive" });
    } finally {
      setIsDeleting(null);
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const projectPolicies =
    policies?.filter((p: EventPolicyResponse) => p.project_id === projectId) ||
    [];
  const inheritedPolicies =
    policies?.filter((p: EventPolicyResponse) => p.project_id !== projectId) ||
    [];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-xl font-semibold">Event Policies</h2>
          <p className="text-sm text-muted-foreground">
            Configure automatic notifications and approval workflows for this
            project
          </p>
        </div>
        {canAddPolicy ? (
          <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="h-4 w-4 mr-2" />
                New Policy
              </Button>
            </DialogTrigger>
            <CreateProjectPolicyDialog
              projectId={projectId}
              orgNodeId={orgNodeId}
              eventTypes={eventTypes}
              projectRoles={projectRoles}
              orgRoles={orgRoles}
              onClose={() => setCreateDialogOpen(false)}
            />
          </Dialog>
        ) : (
          <Button disabled>
            <Lock className="h-4 w-4 mr-2" />
            No Permission
          </Button>
        )}
      </div>

      {/* Project Policies */}
      <div className="space-y-4">
        <h3 className="text-lg font-medium">Project Policies</h3>
        {projectPolicies.length === 0 ? (
          <Card>
            <CardContent className="py-8 text-center text-muted-foreground">
              No project-specific policies configured.
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {projectPolicies.map((policy: EventPolicyResponse) => (
              <PolicyCard
                key={policy.id}
                policy={policy}
                eventTypes={eventTypes}
                onDelete={
                  canRemovePolicy ? () => handleDeletePolicy(policy) : undefined
                }
                isDeleting={isDeleting === policy.id}
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
              From Organisation
            </Badge>
          </h3>
          <div className="grid gap-4">
            {inheritedPolicies.map((policy: EventPolicyResponse) => (
              <PolicyCard
                key={policy.id}
                policy={policy}
                eventTypes={eventTypes}
                inherited
              />
            ))}
          </div>
        </div>
      )}
    </div>
  );
}

interface PolicyCardProps {
  policy: EventPolicyResponse;
  eventTypes: EventTypeInfo[];
  onDelete?: () => void;
  isDeleting?: boolean;
  inherited?: boolean;
}

function PolicyCard({
  policy,
  eventTypes,
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
          {!inherited && onDelete && (
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
        </div>
      </CardContent>
    </Card>
  );
}

function CreateProjectPolicyDialog({
  projectId,
  orgNodeId,
  eventTypes,
  projectRoles,
  orgRoles,
  onClose,
}: CreateProjectPolicyDialogProps) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [selectedEventTypes, setSelectedEventTypes] = useState<string[]>([]);
  const [actionType, setActionType] = useState<"notify" | "request_approval">(
    "notify"
  );
  const [dynamicRecipients, setDynamicRecipients] = useState<string[]>([]);
  const [specificUsers, setSpecificUsers] = useState<string[]>([]);
  const [selectedProjectRoles, setSelectedProjectRoles] = useState<string[]>(
    []
  );
  const [selectedOrgRoles, setSelectedOrgRoles] = useState<string[]>([]);
  const [enabled, setEnabled] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const queryClient = useQueryClient();
  const { toast } = useToast();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      const input: EventPolicyAddedInput = {
        name,
        description: description || "",
        event_types: selectedEventTypes,
        action_type: actionType,
        recipient_user_ids: specificUsers,
        recipient_project_role_ids: selectedProjectRoles,
        recipient_org_role_ids: selectedOrgRoles,
        recipient_dynamic: dynamicRecipients,
        enabled,
      };
      await createEventPolicyAddedEvent(projectId, input);
      queryClient.invalidateQueries({
        queryKey: getGetProjectsIdPoliciesQueryKey(projectId),
      });
      toast({ title: "Policy creation requested" });
      onClose();
    } catch (error) {
      toast({ title: "Failed to create policy", variant: "destructive" });
    } finally {
      setIsSubmitting(false);
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

  return (
    <DialogContent className="max-w-lg max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Create Project Event Policy</DialogTitle>
      </DialogHeader>
      <form onSubmit={handleSubmit} className="space-y-4">
        <div className="space-y-2">
          <Label htmlFor="name">Name</Label>
          <Input
            id="name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="e.g., Notify on role changes"
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
            rows={2}
          />
        </div>

        <div className="space-y-2">
          <Label>Trigger on Events</Label>
          <div className="flex flex-wrap gap-2 max-h-32 overflow-y-auto p-1">
            {eventTypes.map((type) => (
              <Badge
                key={type.type}
                variant={
                  selectedEventTypes.includes(type.type!)
                    ? "default"
                    : "outline"
                }
                className="cursor-pointer"
                onClick={() => toggleEventType(type.type!)}
              >
                {type.friendlyName || type.type}
              </Badge>
            ))}
          </div>
        </div>

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

        <div className="space-y-2">
          <Label>Dynamic Recipients</Label>
          <div className="flex flex-wrap gap-2">
            {DYNAMIC_RECIPIENTS.map((type) => (
              <Badge
                key={type.value}
                variant={
                  dynamicRecipients.includes(type.value) ? "default" : "outline"
                }
                className="cursor-pointer"
                onClick={() => toggleDynamicRecipient(type.value)}
              >
                {type.label}
              </Badge>
            ))}
          </div>
        </div>

        {projectRoles.length > 0 && (
          <div className="space-y-2">
            <Label>Project Roles (optional)</Label>
            <div className="flex flex-wrap gap-2">
              {projectRoles.map((role) => (
                <Badge
                  key={role.id}
                  variant={
                    selectedProjectRoles.includes(role.id!)
                      ? "default"
                      : "outline"
                  }
                  className="cursor-pointer"
                  onClick={() => {
                    setSelectedProjectRoles((prev) =>
                      prev.includes(role.id!)
                        ? prev.filter((id) => id !== role.id)
                        : [...prev, role.id!]
                    );
                  }}
                >
                  {role.name}
                </Badge>
              ))}
            </div>
          </div>
        )}

        {orgRoles.length > 0 && (
          <div className="space-y-2">
            <Label>Organisation Roles (optional)</Label>
            <div className="flex flex-wrap gap-2">
              {orgRoles.map((role) => (
                <Badge
                  key={role.id}
                  variant={
                    selectedOrgRoles.includes(role.id!) ? "default" : "outline"
                  }
                  className="cursor-pointer"
                  onClick={() => {
                    setSelectedOrgRoles((prev) =>
                      prev.includes(role.id!)
                        ? prev.filter((id) => id !== role.id)
                        : [...prev, role.id!]
                    );
                  }}
                >
                  {role.displayName}
                </Badge>
              ))}
            </div>
          </div>
        )}

        <div className="space-y-2">
          <Label>Specific Users (optional)</Label>
          <MultiUserSelect
            value={specificUsers}
            onChange={setSpecificUsers}
            placeholder="Search and add specific users..."
          />
        </div>

        <div className="flex items-center justify-between">
          <Label htmlFor="enabled">Enabled</Label>
          <Checkbox
            id="enabled"
            checked={enabled}
            onCheckedChange={(checked) => setEnabled(checked === true)}
          />
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={!name || selectedEventTypes.length === 0 || isSubmitting}
          >
            {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Create Policy
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  );
}
