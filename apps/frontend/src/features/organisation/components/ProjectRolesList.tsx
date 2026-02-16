import { useState } from "react";
import { Loader2, Trash2, Settings } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";
import { useToast } from "@/hooks/use-toast";
import {
  useGetOrganisationNodesIdRoles,
  useDeleteOrganisationNodesIdRolesRoleId,
} from "@api/moris";
import { ProjectRoleResponse } from "@/api/generated-orval/model";
import { EditRolePermissionsDialog } from "./EditRolePermissionsDialog";
import { AddProjectRoleDialog } from "./AddProjectRoleDialog";

interface ProjectRolesListProps {
  nodeId: string;
}

export const ProjectRolesList = ({ nodeId }: ProjectRolesListProps) => {
  const {
    data: roles,
    isLoading,
    refetch,
  } = useGetOrganisationNodesIdRoles(nodeId);
  const { mutateAsync: deleteRole, isPending: isDeleting } =
    useDeleteOrganisationNodesIdRolesRoleId();
  const { toast } = useToast();

  const [selectedRole, setSelectedRole] = useState<ProjectRoleResponse | null>(
    null
  );
  const [roleToDelete, setRoleToDelete] = useState<string | null>(null);

  const handleDelete = async () => {
    if (!roleToDelete) return;
    try {
      await deleteRole({ id: nodeId, roleId: roleToDelete });
      toast({ title: "Role deleted" });
      refetch();
    } catch (error) {
      console.error("Failed to delete role", error);
      toast({ title: "Error deleting role", variant: "destructive" });
    } finally {
      setRoleToDelete(null);
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-medium">Project Roles</h3>
        <AddProjectRoleDialog nodeId={nodeId} onSuccess={refetch} />
      </div>

      <div className="border rounded-lg overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[200px]">Name</TableHead>
              <TableHead>Key</TableHead>
              <TableHead>Allowed Events</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={4} className="h-24 text-center">
                  <div className="flex items-center justify-center">
                    <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                  </div>
                </TableCell>
              </TableRow>
            ) : roles?.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={4}
                  className="h-24 text-center text-muted-foreground"
                >
                  No custom roles defined.
                </TableCell>
              </TableRow>
            ) : (
              roles?.map((role) => (
                <TableRow key={role.id}>
                  <TableCell className="font-medium">{role.name}</TableCell>
                  <TableCell className="font-mono text-xs text-muted-foreground">
                    {role.key}
                  </TableCell>
                  <TableCell>
                    <span className="text-sm text-muted-foreground">
                      {role.allowedEventTypes?.length || 0} events
                    </span>
                  </TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => setSelectedRole(role)}
                      title="Edit Permissions"
                    >
                      <Settings className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => setRoleToDelete(role.id as string)}
                      title="Delete Role"
                    >
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
      <p className="text-xs text-muted-foreground mt-2">
        Note: This list includes roles inherited from parent organisations.
      </p>

      {selectedRole && (
        <EditRolePermissionsDialog
          role={selectedRole}
          nodeId={nodeId}
          open={!!selectedRole}
          onOpenChange={(open) => !open && setSelectedRole(null)}
          onSuccess={() => {
            setSelectedRole(null);
            refetch();
          }}
        />
      )}

      <ConfirmationModal
        isOpen={!!roleToDelete}
        onClose={() => setRoleToDelete(null)}
        onConfirm={handleDelete}
        title="Delete Role"
        description="Are you sure you want to delete this role? This cannot be undone."
        confirmLabel="Delete"
        variant="destructive"
        isLoading={isDeleting}
      />
    </div>
  );
};
