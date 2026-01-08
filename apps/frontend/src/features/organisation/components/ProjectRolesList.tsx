import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2, Plus, Trash2 } from "lucide-react";
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
import { useToast } from "@/hooks/use-toast";
import {
    useGetOrganisationNodesIdRoles,
    usePostOrganisationNodesIdRoles,
    useDeleteOrganisationNodesIdRolesRoleId,
} from "@api/moris";

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";

// Schema for creating a role
const createRoleSchema = z.object({
    key: z.string().min(1, "Key is required").regex(/^[a-z0-9_-]+$/, "Key must be lowercase alphanumeric, dashes, or underscores"),
    name: z.string().min(1, "Name is required"),
});

interface ProjectRolesListProps {
    nodeId: string;
}

export const ProjectRolesList = ({ nodeId }: ProjectRolesListProps) => {   
    const { data: roles, isLoading, refetch } = useGetOrganisationNodesIdRoles(nodeId);
    const { mutateAsync: deleteRole } = useDeleteOrganisationNodesIdRolesRoleId();
    const { toast } = useToast();

    const handleDelete = async (roleId: string) => {
        if (!confirm("Are you sure you want to delete this role? This cannot be undone.")) {
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
                            <TableHead className="text-right">Actions</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {isLoading ? (
                            <TableRow>
                                <TableCell colSpan={3} className="h-24 text-center">
                                    <div className="flex items-center justify-center">
                                        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                                    </div>
                                </TableCell>
                            </TableRow>
                        ) : roles?.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={3} className="h-24 text-center text-muted-foreground">
                                    No custom roles defined.
                                </TableCell>
                            </TableRow>
                        ) : (
                            roles?.map((role) => (
                                <TableRow key={role.id}>
                                    <TableCell className="font-medium">{role.name}</TableCell>
                                    <TableCell className="font-mono text-xs text-muted-foreground">{role.key}</TableCell>
                                    <TableCell className="text-right">
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
        </div>
    );
};

const AddProjectRoleDialog = ({ nodeId, onSuccess }: { nodeId: string; onSuccess: () => void }) => {
    const [open, setOpen] = useState(false);
    const { toast } = useToast();
    const { mutateAsync: createRole, isPending } = usePostOrganisationNodesIdRoles();

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
                data: values
            });
            toast({
                title: "Role created",
                description: "Project role has been successfully created.",
            });
            setOpen(false);
            form.reset();
            onSuccess();
        } catch (error) {
            console.error(error);
            toast({
                variant: "destructive",
                title: "Error",
                description: "Failed to create role. Ensure key is unique within the organisation.",
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
                                    <p className="text-xs text-muted-foreground">Used internally. Lowercase, smooth.</p>
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
