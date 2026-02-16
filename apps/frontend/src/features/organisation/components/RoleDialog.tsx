import { useEffect } from "react";
import { useForm } from "react-hook-form";
import {
  usePostOrganisationNodesIdOrganisationRoles,
  usePutOrganisationRolesId,
  getGetOrganisationNodesIdOrganisationRolesQueryKey,
  useGetOrganisationPermissions,
} from "@api/moris";
import { OrganisationRoleResponse } from "@api/model";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { toast } from "sonner";
import { useQueryClient } from "@tanstack/react-query";
import { Loader2 } from "lucide-react";

interface RoleDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  nodeId: string;
  role?: OrganisationRoleResponse; // If provided, edit mode
}

interface RoleFormValues {
  key: string;
  displayName: string;
  permissions: string[];
}

export function RoleDialog({
  open,
  onOpenChange,
  nodeId,
  role,
}: RoleDialogProps) {
  const isEdit = !!role;
  const queryClient = useQueryClient();

  const { data: permissionsData } = useGetOrganisationPermissions();
  const permissions = permissionsData?.permissions || [];

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<RoleFormValues>({
    defaultValues: {
      key: "",
      displayName: "",
      permissions: [],
    },
  });

  const selectedPermissions = watch("permissions");

  useEffect(() => {
    if (open) {
      if (role) {
        reset({
          key: role.key,
          displayName: role.displayName,
          permissions: role.permissions || [],
        });
      } else {
        reset({
          key: "",
          displayName: "",
          permissions: [],
        });
      }
    }
  }, [open, role, reset]);

  const { mutateAsync: createRole, isPending: isCreating } =
    usePostOrganisationNodesIdOrganisationRoles();
  const { mutateAsync: updateRole, isPending: isUpdating } =
    usePutOrganisationRolesId();

  const onSubmit = async (data: RoleFormValues) => {
    try {
      if (isEdit && role) {
        await updateRole({
          id: role.id!,
          data: {
            displayName: data.displayName,
            permissions: data.permissions,
          },
        });
        toast.success("Role updated successfully");
      } else {
        await createRole({
          id: nodeId,
          data: {
            key: data.key,
            displayName: data.displayName,
            permissions: data.permissions,
          },
        });
        toast.success("Role created successfully");
      }

      queryClient.invalidateQueries({
        queryKey: getGetOrganisationNodesIdOrganisationRolesQueryKey(nodeId),
      });
      onOpenChange(false);
    } catch (error: any) {
      console.error("Failed to save role", error);
      const msg =
        error?.response?.data?.message ||
        error?.message ||
        "Failed to save role";
      toast.error("Error", {
        description: msg,
      });
    }
  };

  const handlePermissionChange = (permKey: string, checked: boolean) => {
    const current = selectedPermissions || [];
    if (checked) {
      setValue("permissions", [...current, permKey]);
    } else {
      setValue(
        "permissions",
        current.filter((p) => p !== permKey),
      );
    }
  };

  const isPending = isCreating || isUpdating;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-xl">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Edit Role" : "Create Role"}</DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="key">Key</Label>
              <Input
                id="key"
                {...register("key", {
                  required: "Key is required",
                  pattern: /^[a-z_]+$/,
                })}
                placeholder="e.g. project_manager"
                disabled={isEdit || isPending}
              />
              {errors.key && (
                <p className="text-sm text-destructive">{errors.key.message}</p>
              )}
              <p className="text-xs text-muted-foreground">
                Unique identifier (lowercase, underscores).
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="displayName">Display Name</Label>
              <Input
                id="displayName"
                {...register("displayName", {
                  required: "Display name is required",
                })}
                placeholder="e.g. Project Manager"
                disabled={isPending}
              />
              {errors.displayName && (
                <p className="text-sm text-destructive">
                  {errors.displayName.message}
                </p>
              )}
            </div>
          </div>

          <div className="space-y-2">
            <Label>Permissions</Label>
            <div className="border rounded-md p-4 space-y-3 max-h-60 overflow-y-auto">
              {permissions.map((perm) => {
                const isChecked = (selectedPermissions || []).includes(
                  perm.key || "",
                );
                return (
                  <div key={perm.key} className="flex items-start space-x-2">
                    <Checkbox
                      id={`perm-${perm.key}`}
                      checked={isChecked}
                      onCheckedChange={(c) =>
                        handlePermissionChange(perm.key || "", c as boolean)
                      }
                      disabled={isPending}
                    />
                    <div className="grid gap-1.5 leading-none">
                      <label
                        htmlFor={`perm-${perm.key}`}
                        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                      >
                        {perm.label}
                      </label>
                      <p className="text-xs text-muted-foreground">
                        {perm.description}
                      </p>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isPending}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {isEdit ? "Update Role" : "Create Role"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
