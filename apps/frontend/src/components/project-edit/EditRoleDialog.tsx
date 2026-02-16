import { useState, useEffect, useMemo } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useGetProjectsIdRoles } from "@api/moris";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Checkbox } from "@/components/ui/checkbox";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { Loader2 } from "lucide-react";
import {
  createProjectRoleAssignedEvent,
  createProjectRoleUnassignedEvent,
} from "@/api/events";
import { ProjectMemberResponse } from "@/api/generated-orval/model";

interface EditRoleDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  member: ProjectMemberResponse;
  projectId: string;
  /** All members for this person (to see all their current roles) */
  allMembersForPerson?: ProjectMemberResponse[];
  onSuccess?: () => void;
}

export function EditRoleDialog({
  open,
  onOpenChange,
  member,
  projectId,
  allMembersForPerson,
  onSuccess,
}: EditRoleDialogProps) {
  // Compute current roles from all member entries for this person
  const currentRoleIds = useMemo(() => {
    if (allMembersForPerson && allMembersForPerson.length > 0) {
      return allMembersForPerson
        .map((m) => m.role_id)
        .filter((id): id is string => !!id);
    }
    // Fallback to just the current member's role
    return member.role_id ? [member.role_id] : [];
  }, [allMembersForPerson, member.role_id]);

  const [selectedRoles, setSelectedRoles] = useState<string[]>(currentRoleIds);
  const [isSaving, setIsSaving] = useState(false);
  const queryClient = useQueryClient();

  const { data: roles, isLoading: isLoadingRoles } =
    useGetProjectsIdRoles(projectId);

  // Reset selected roles when dialog opens or member changes
  useEffect(() => {
    if (open) {
      setSelectedRoles(currentRoleIds);
    }
  }, [open, currentRoleIds]);

  const toggleRole = (roleId: string, checked: boolean) => {
    if (checked) {
      setSelectedRoles((prev) => [...prev, roleId]);
    } else {
      setSelectedRoles((prev) => prev.filter((id) => id !== roleId));
    }
  };

  const handleSave = async () => {
    if (!member.id) return;
    setIsSaving(true);

    try {
      // Determine which roles to add and remove
      const rolesToAdd = selectedRoles.filter(
        (r) => !currentRoleIds.includes(r),
      );
      const rolesToRemove = currentRoleIds.filter(
        (r) => !selectedRoles.includes(r),
      );

      // Unassign removed roles
      for (const roleId of rolesToRemove) {
        await createProjectRoleUnassignedEvent(projectId, {
          person_id: member.id,
          project_role_id: roleId,
        });
      }

      // Assign new roles
      for (const roleId of rolesToAdd) {
        await createProjectRoleAssignedEvent(projectId, {
          person_id: member.id,
          project_role_id: roleId,
        });
      }

      queryClient.invalidateQueries({ queryKey: ["/projects", projectId] });

      const totalChanges = rolesToAdd.length + rolesToRemove.length;
      if (totalChanges > 0) {
        toast.success("Roles updated", {
          description: `${rolesToAdd.length} role(s) added, ${rolesToRemove.length} role(s) removed.`,
        });
      } else {
        toast.success("No changes made");
      }

      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error(error);
      toast.error("Failed to update roles");
    } finally {
      setIsSaving(false);
    }
  };

  const getRoleNames = (roleIds: string[]) => {
    return roleIds
      .map((id) => roles?.find((r) => r.id === id)?.name)
      .filter(Boolean);
  };

  const hasChanges = useMemo(() => {
    if (selectedRoles.length !== currentRoleIds.length) return true;
    return selectedRoles.some((r) => !currentRoleIds.includes(r));
  }, [selectedRoles, currentRoleIds]);

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Roles for {member.name || "Member"}</DialogTitle>
          <DialogDescription>
            Select one or more roles for this team member.
          </DialogDescription>
        </DialogHeader>
        <div className="py-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Roles</label>
            <div className="border rounded-md p-3 space-y-2">
              {isLoadingRoles ? (
                <div className="flex items-center gap-2 text-sm text-muted-foreground">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Loading roles...
                </div>
              ) : (
                roles?.map((role) => (
                  <div key={role.id} className="flex items-center space-x-2">
                    <Checkbox
                      id={`edit-role-${role.id}`}
                      checked={selectedRoles.includes(role.id || "")}
                      onCheckedChange={(checked) =>
                        toggleRole(role.id || "", checked === true)
                      }
                    />
                    <label
                      htmlFor={`edit-role-${role.id}`}
                      className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                    >
                      {role.name}
                    </label>
                  </div>
                ))
              )}
            </div>
            {selectedRoles.length > 0 && (
              <div className="flex flex-wrap gap-1 mt-2">
                {getRoleNames(selectedRoles).map((name) => (
                  <Badge key={name} variant="secondary" className="text-xs">
                    {name}
                  </Badge>
                ))}
              </div>
            )}
            {selectedRoles.length === 0 && (
              <p className="text-sm text-destructive mt-1">
                At least one role is required.
              </p>
            )}
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={handleSave}
            disabled={
              isSaving ||
              isLoadingRoles ||
              selectedRoles.length === 0 ||
              !hasChanges
            }
          >
            {isSaving ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Saving...
              </>
            ) : (
              "Save Changes"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
