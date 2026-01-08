import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useGetProjectsIdRoles } from "@api/moris";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/hooks/use-toast";
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
  onSuccess?: () => void;
}

export function EditRoleDialog({
  open,
  onOpenChange,
  member,
  projectId,
  onSuccess,
}: EditRoleDialogProps) {
  const [role, setRole] = useState(member.role_id || "");
  const [isSaving, setIsSaving] = useState(false);
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const { data: roles, isLoading: isLoadingRoles } = useGetProjectsIdRoles(projectId);

  const handleSave = async () => {
    if (!member.id) return;
    setIsSaving(true);
    try {
      // 1. Unassign current role if exists
      if (member.role_id) {
        await createProjectRoleUnassignedEvent(projectId, {
          person_id: member.id,
          project_role_id: member.role_id,
        });
      }

      // 2. Assign new role
      if (role) {
        await createProjectRoleAssignedEvent(projectId, {
          person_id: member.id,
          project_role_id: role,
        });
      }

      queryClient.invalidateQueries({ queryKey: ["/projects", projectId] });
      toast({ title: "Role updated" });
      onOpenChange(false);
      onSuccess?.();
    } catch (error) {
      console.error(error);
      toast({
        title: "Failed to update role",
        variant: "destructive",
      });
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Role</DialogTitle>
        </DialogHeader>
        <div className="py-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Role</label>
            {isLoadingRoles ? (
              <div className="flex bg-muted h-10 w-full items-center px-3 rounded-md">
                <Loader2 className="h-4 w-4 animate-spin text-muted-foreground mr-2" />
                <span className="text-sm text-muted-foreground">
                  Loading roles...
                </span>
              </div>
            ) : (
              <Select value={role} onValueChange={setRole}>
                <SelectTrigger>
                  <SelectValue placeholder="Select a role" />
                </SelectTrigger>
                <SelectContent>
                  {roles?.map((r) => (
                    <SelectItem key={r.id} value={r.id || ""}>
                      {r.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={isSaving || isLoadingRoles}>
            {isSaving ? "Saving..." : "Save Changes"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
