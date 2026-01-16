import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2, Plus, Trash2, Settings } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
  DialogDescription,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { useToast } from "@/hooks/use-toast";
import {
  useGetOrganisationNodesIdRoles,
  usePostOrganisationNodesIdRoles,
  useDeleteOrganisationNodesIdRolesRoleId,
  usePatchOrganisationNodesIdRolesRoleId,
  useGetEventTypes,
} from "@api/moris";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { ProjectRoleResponse } from "@/api/generated-orval/model";

// Schema for creating a role
const createRoleSchema = z.object({
  key: z
    .string()
    .min(1, "Key is required")
    .regex(
      /^[a-z0-9_-]+$/,
      "Key must be lowercase alphanumeric, dashes, or underscores"
    ),
  name: z.string().min(1, "Name is required"),
});

interface ProjectRolesListProps {
  nodeId: string;
}

export const ProjectRolesList = ({ nodeId }: ProjectRolesListProps) => {
  const {
    data: roles,
    isLoading,
    refetch,
  } = useGetOrganisationNodesIdRoles(nodeId);
  const { mutateAsync: deleteRole } = useDeleteOrganisationNodesIdRolesRoleId();
  const { toast } = useToast();

  const [selectedRole, setSelectedRole] = useState<ProjectRoleResponse | null>(
    null
  );

  const handleDelete = async (roleId: string) => {
    if (
      !confirm(
        "Are you sure you want to delete this role? This cannot be undone."
      )
    ) {
      return;
    }
    try {
      await deleteRole({ id: nodeId, roleId });
      toast({ title: "Role deleted" });
      refetch();
    } catch (error) {
      console.error("Failed to delete role", error);
      toast({ title: "Error deleting role", variant: "destructive" });
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
                      onClick={() => handleDelete(role.id as string)}
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
    </div>
  );
};

const EditRolePermissionsDialog = ({
  role,
  nodeId,
  open,
  onOpenChange,
  onSuccess,
}: {
  role: ProjectRoleResponse;
  nodeId: string;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
}) => {
  const { toast } = useToast();
  const { data: eventTypes, isLoading: isLoadingEventTypes } =
    useGetEventTypes();
  const { mutateAsync: updateRole, isPending } =
    usePatchOrganisationNodesIdRolesRoleId();

  const [selectedTypes, setSelectedTypes] = useState<Set<string>>(
    new Set(role.allowedEventTypes || [])
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
      toast({ title: "Role permissions updated" });
      onSuccess();
    } catch (error) {
      console.error("Failed to update role permissions", error);
      toast({
        title: "Error updating permissions",
        variant: "destructive",
      });
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

const AddProjectRoleDialog = ({
  nodeId,
  onSuccess,
}: {
  nodeId: string;
  onSuccess: () => void;
}) => {
  const [open, setOpen] = useState(false);
  const { toast } = useToast();
  const { mutateAsync: createRole, isPending } =
    usePostOrganisationNodesIdRoles();

  const form = useForm<z.infer<typeof createRoleSchema>>({
    resolver: zodResolver(createRoleSchema),
    defaultValues: {
      key: "",
      name: "",
    },
  });

  async function onSubmit(values: z.infer<typeof createRoleSchema>) {
    try {
      await createRole({
        id: nodeId,
        data: values,
      });
      toast({
        title: "Role created",
        description:
          "Project role has been successfully created. Configure permissions to enable event access.",
      });
      setOpen(false);
      form.reset();
      onSuccess();
    } catch (error) {
      console.error(error);
      toast({
        variant: "destructive",
        title: "Error",
        description:
          "Failed to create role. Ensure key is unique within the organisation.",
      });
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm">
          <Plus className="mr-2 h-4 w-4" />
          Create Role
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Project Role</DialogTitle>
          <DialogDescription>
            New roles start with no permissions. After creating, use the
            settings button to configure allowed event types.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name</FormLabel>
                  <FormControl>
                    <Input placeholder="e.g. Data Steward" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="key"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Key</FormLabel>
                  <FormControl>
                    <Input placeholder="e.g. data_steward" {...field} />
                  </FormControl>
                  <p className="text-xs text-muted-foreground">
                    Used internally. Lowercase, smooth.
                  </p>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button type="submit" disabled={isPending}>
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Create
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
};
