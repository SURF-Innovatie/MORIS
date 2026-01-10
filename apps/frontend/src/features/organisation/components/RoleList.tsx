import { useState } from "react";
import {
  useGetOrganisationNodesIdOrganisationRoles,
  useDeleteOrganisationRolesId,
  getGetOrganisationNodesIdOrganisationRolesQueryKey,
  useGetOrganisationPermissions,
} from "@api/moris";
import { OrganisationRoleResponse } from "@api/model";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Plus, MoreVertical, Pencil, Trash2 } from "lucide-react";
import { RoleDialog } from "./RoleDialog";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";
import { useToast } from "@/hooks/use-toast";
import { useQueryClient } from "@tanstack/react-query";

interface RoleListProps {
  nodeId: string;
}

export function RoleList({ nodeId }: RoleListProps) {
  const { data: roles, isLoading } =
    useGetOrganisationNodesIdOrganisationRoles(nodeId);
  const { data: permissionsData } = useGetOrganisationPermissions();
  const permissions = permissionsData?.permissions || [];

  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<
    OrganisationRoleResponse | undefined
  >(undefined);

  const [roleToDelete, setRoleToDelete] =
    useState<OrganisationRoleResponse | null>(null);

  const { toast } = useToast();
  const queryClient = useQueryClient();
  const { mutateAsync: deleteRole, isPending: isDeleting } =
    useDeleteOrganisationRolesId();

  const handleCreate = () => {
    setEditingRole(undefined);
    setIsDialogOpen(true);
  };

  const handleEdit = (role: OrganisationRoleResponse) => {
    setEditingRole(role);
    setIsDialogOpen(true);
  };

  const handleDelete = async () => {
    if (!roleToDelete?.id) return;
    try {
      await deleteRole({ id: roleToDelete.id });
      toast({ title: "Role deleted successfully" });
      queryClient.invalidateQueries({
        queryKey: getGetOrganisationNodesIdOrganisationRolesQueryKey(nodeId),
      });
      setRoleToDelete(null);
    } catch (error: any) {
      console.error("Failed to delete role", error);
      const msg =
        error?.response?.data?.message ||
        error?.message ||
        "Failed to delete role";
      toast({
        title: "Cannot delete role",
        description: msg,
        variant: "destructive",
      });
    }
  };

  if (isLoading) return <div>Loading roles...</div>;

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h2 className="text-lg font-semibold">Organisation Roles</h2>
        <Button onClick={handleCreate}>
          <Plus className="mr-2 h-4 w-4" /> Create Role
        </Button>
      </div>

      <div className="border rounded-lg overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Display Name</TableHead>
              <TableHead>Key</TableHead>
              <TableHead>Permissions</TableHead>
              <TableHead className="w-[100px]">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {roles?.map((role) => (
              <TableRow key={role.id}>
                <TableCell className="font-medium">
                  {role.displayName}
                </TableCell>
                <TableCell className="font-mono text-xs text-muted-foreground">
                  {role.key}
                </TableCell>
                <TableCell>
                  <div className="flex flex-wrap gap-1">
                    {role.permissions?.map((p) => {
                      const label =
                        permissions.find((perm) => perm.key === p)?.label || p;
                      return (
                        <span
                          key={p}
                          className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-secondary text-secondary-foreground"
                        >
                          {label}
                        </span>
                      );
                    })}
                    {(!role.permissions || role.permissions.length === 0) && (
                      <span className="text-xs text-muted-foreground italic">
                        No permissions
                      </span>
                    )}
                  </div>
                </TableCell>
                <TableCell>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" className="h-8 w-8 p-0">
                        <MoreVertical className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem onClick={() => handleEdit(role)}>
                        <Pencil className="mr-2 h-4 w-4" /> Edit
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => setRoleToDelete(role)}
                        className="text-destructive focus:text-destructive"
                      >
                        <Trash2 className="mr-2 h-4 w-4" /> Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
            {roles?.length === 0 && (
              <TableRow>
                <TableCell
                  colSpan={4}
                  className="text-center py-8 text-muted-foreground"
                >
                  No roles defined for this organisation. Create one to get
                  started.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <RoleDialog
        open={isDialogOpen}
        onOpenChange={setIsDialogOpen}
        nodeId={nodeId}
        role={editingRole}
      />

      <ConfirmationModal
        isOpen={!!roleToDelete}
        onClose={() => setRoleToDelete(null)}
        onConfirm={handleDelete}
        title="Delete Role"
        description={`Are you sure you want to delete the role "${roleToDelete?.displayName}"? This action cannot be undone.`}
        confirmLabel="Delete"
        variant="destructive"
        isLoading={isDeleting}
      />
    </div>
  );
}
