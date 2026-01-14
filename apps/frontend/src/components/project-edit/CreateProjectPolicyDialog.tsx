import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { getGetProjectsIdPoliciesQueryKey } from "@api/moris";
import {
  EventTypeInfo,
  ProjectRoleResponse,
  OrganisationRoleResponse,
} from "@api/model";
import {
  createEventPolicyAddedEvent,
  EventPolicyAddedInput,
} from "@/api/events";
import { MultiUserSelect } from "@/components/user/MultiUserSelect";
import { Button } from "@/components/ui/button";
import {
  DialogContent,
  DialogHeader,
  DialogTitle,
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
import { Loader2, Bell, ShieldCheck } from "lucide-react";
import { useToast } from "@/hooks/use-toast";

const DYNAMIC_RECIPIENTS = [
  { value: "project_members", label: "All Project Members" },
  { value: "project_owner", label: "Project Owner" },
  { value: "org_admins", label: "Organisation Admins" },
];

interface CreateProjectPolicyDialogProps {
  projectId: string;
  eventTypes: EventTypeInfo[];
  projectRoles: ProjectRoleResponse[];
  orgRoles: OrganisationRoleResponse[];
  onClose: () => void;
}

export function CreateProjectPolicyDialog({
  projectId,
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

  const toggleProjectRole = (roleId: string) => {
    setSelectedProjectRoles((prev) =>
      prev.includes(roleId)
        ? prev.filter((r) => r !== roleId)
        : [...prev, roleId]
    );
  };

  const toggleOrgRole = (roleId: string) => {
    setSelectedOrgRoles((prev) =>
      prev.includes(roleId)
        ? prev.filter((r) => r !== roleId)
        : [...prev, roleId]
    );
  };

  return (
    <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Create Project Event Policy</DialogTitle>
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

          <div className="border rounded-lg p-6 space-y-6 bg-muted/20">
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
