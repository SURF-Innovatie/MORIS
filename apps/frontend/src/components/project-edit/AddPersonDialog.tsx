import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2, Plus } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/use-toast";
import { usePostPeople, useGetProjectsRoles } from "@api/moris";
import { createProjectRoleAssignedEvent } from "@/api/events";

const addPersonSchema = z.object({
  name: z.string().min(1, "Name is required"),
  email: z.string().email("Invalid email address"),
  role: z.string().min(1, "Role is required"),
});

interface AddPersonDialogProps {
  projectId: string;
  onPersonAdded: () => void;
}

export function AddPersonDialog({
  projectId,
  onPersonAdded,
}: AddPersonDialogProps) {
  const [open, setOpen] = useState(false);
  const { toast } = useToast();

  const { mutateAsync: createPerson, isPending: isCreatingPerson } =
    usePostPeople();
  const { data: roles, isLoading: isLoadingRoles } = useGetProjectsRoles();

  const form = useForm<z.infer<typeof addPersonSchema>>({
    resolver: zodResolver(addPersonSchema),
    defaultValues: {
      name: "",
      email: "",
      role: "",
    },
  });

  // We can handle the specific role assignment 'loading' state separately if needed,
  // but for now grouping it makes sense for the button disabled state.
  const [isAssigningRole, setIsAssigningRole] = useState(false);

  const isPending = isCreatingPerson || isAssigningRole;

  async function onSubmit(values: z.infer<typeof addPersonSchema>) {
    try {
      // 1. Create the person
      const nameParts = values.name.split(" ");
      const givenName = nameParts[0];
      const familyName = nameParts.slice(1).join(" ") || "Unknown";

      const person = await createPerson({
        data: {
          name: values.name,
          email: values.email,
          givenName: givenName,
          familyName: familyName,
        },
      });

      // 2. Assign role to project (effectively adding them)
      if (person && person.id) {
        setIsAssigningRole(true);
        await createProjectRoleAssignedEvent(projectId, {
          person_id: person.id,
          project_role_id: values.role,
        });
        setIsAssigningRole(false);

        toast({
          title: "Member added",
          description: `${values.name} has been successfully added to the project.`,
        });

        setOpen(false);
        form.reset();
        onPersonAdded();
      }
    } catch (error) {
      console.error(error);
      setIsAssigningRole(false);
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to add member. Please try again.",
      });
    }
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button size="sm">
          <Plus className="mr-2 h-4 w-4" />
          Add Member
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Add Team Member</DialogTitle>
          <DialogDescription>
            Create a new person and add them to this project.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Full Name</FormLabel>
                  <FormControl>
                    <Input placeholder="John Doe" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="john.doe@example.com"
                      type="email"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="role"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Role (Optional)</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select a role" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {isLoadingRoles ? (
                        <div className="flex bg-muted h-10 w-full items-center px-3 rounded-md">
                          <Loader2 className="h-4 w-4 animate-spin text-muted-foreground mr-2" />
                          <span className="text-sm text-muted-foreground">
                            Loading roles...
                          </span>
                        </div>
                      ) : (
                        roles?.map((r) => (
                            <SelectItem key={r.id} value={r.id || ""}>
                              {r.name}
                            </SelectItem>
                        ))
                      )}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button type="submit" disabled={isPending}>
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Add Member
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
