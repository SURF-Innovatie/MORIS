import { useState } from "react";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
} from "@/components/ui/dialog";
import { Checkbox } from "@/components/ui/checkbox";
import { toast } from "sonner";
import {
  useGetEventTypes,
  usePatchOrganisationNodesIdRolesRoleId,
} from "@api/moris";
import { ProjectRoleResponse } from "@/api/generated-orval/model";

interface EditRolePermissionsDialogProps {
  role: ProjectRoleResponse;
  nodeId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}

export const EditRolePermissionsDialog = ({
  role,
  nodeId,
  open,
  onOpenChange,
  onSuccess,
}: EditRolePermissionsDialogProps) => {
  const { data: eventTypes, isLoading: isLoadingEventTypes } =
    useGetEventTypes();
  const { mutateAsync: updateRole, isPending } =
    usePatchOrganisationNodesIdRolesRoleId();

  const [selectedTypes, setSelectedTypes] = useState<Set<string>>(
    new Set(role.allowedEventTypes || []),
  );

  const handleToggle = (eventType: string) => {
    const newSelected = new Set(selectedTypes);
    if (newSelected.has(eventType)) {
      newSelected.delete(eventType);
    } else {
      newSelected.add(eventType);
    }
    setSelectedTypes(newSelected);
  };

  const handleSelectAll = () => {
    if (eventTypes) {
      setSelectedTypes(new Set(eventTypes.map((e) => e.type!)));
    }
  };

  const handleSelectNone = () => {
    setSelectedTypes(new Set());
  };

  const handleSave = async () => {
    try {
      await updateRole({
        id: nodeId,
        roleId: role.id!,
        data: {
          allowedEventTypes: Array.from(selectedTypes),
        },
      });
      toast.success("Role permissions updated");
      onSuccess();
    } catch (error) {
      console.error("Failed to update role permissions", error);
      toast.error("Error updating permissions");
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg max-h-[80vh] overflow-hidden flex flex-col">
        <DialogHeader>
          <DialogTitle>Edit Permissions: {role.name}</DialogTitle>
          <DialogDescription>
            Select which event types this role is allowed to use. Users with
            this role can only create events that are checked below.
          </DialogDescription>
        </DialogHeader>

        {isLoadingEventTypes ? (
          <div className="flex justify-center p-8">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : (
          <>
            <div className="flex gap-2 mb-4">
              <Button variant="outline" size="sm" onClick={handleSelectAll}>
                Select All
              </Button>
              <Button variant="outline" size="sm" onClick={handleSelectNone}>
                Select None
              </Button>
            </div>
            <div className="flex-1 overflow-y-auto space-y-3 pr-2">
              {eventTypes?.map((eventType) => (
                <label
                  key={eventType.type}
                  className="flex items-center gap-3 p-3 rounded-lg border hover:bg-muted/50 cursor-pointer transition-colors"
                >
                  <Checkbox
                    checked={selectedTypes.has(eventType.type!)}
                    onCheckedChange={() => handleToggle(eventType.type!)}
                  />
                  <div className="flex-1">
                    <div className="font-medium text-sm">
                      {eventType.friendlyName}
                    </div>
                    <div className="text-xs text-muted-foreground font-mono">
                      {eventType.type}
                    </div>
                  </div>
                </label>
              ))}
            </div>
          </>
        )}

        <DialogFooter className="mt-4">
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button
            onClick={handleSave}
            disabled={isPending || isLoadingEventTypes}
          >
            {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Save Permissions
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
