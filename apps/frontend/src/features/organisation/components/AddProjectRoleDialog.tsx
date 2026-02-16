import { useState } from "react";
import { useForm } from "react-hook-form";
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema";
import { z } from "zod";
import { Loader2, Plus } from "lucide-react";
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
import { toast } from "sonner";
import { usePostOrganisationNodesIdRoles } from "@api/moris";

const createRoleSchema = z.object({
  key: z
    .string()
    .min(1, "Key is required")
    .regex(
      /^[a-z0-9_-]+$/,
      "Key must be lowercase alphanumeric, dashes, or underscores",
    ),
  name: z.string().min(1, "Name is required"),
});

interface AddProjectRoleDialogProps {
  nodeId: string;
  onSuccess: () => void;
}

export const AddProjectRoleDialog = ({
  nodeId,
  onSuccess,
}: AddProjectRoleDialogProps) => {
  const [open, setOpen] = useState(false);
  const { mutateAsync: createRole, isPending } =
    usePostOrganisationNodesIdRoles();

  const form = useForm<z.infer<typeof createRoleSchema>>({
    resolver: standardSchemaResolver(createRoleSchema),
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
      toast.success("Role created", {
        description:
          "Project role has been successfully created. Configure permissions to enable event access.",
      });
      setOpen(false);
      form.reset();
      onSuccess();
    } catch (error) {
      console.error(error);
      toast.error("Error", {
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
