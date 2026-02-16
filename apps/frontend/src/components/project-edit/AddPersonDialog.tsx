import { useState } from "react";
import { useForm } from "react-hook-form";
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema";
import { z } from "zod";
import { Loader2, Plus, UserPlus, Users } from "lucide-react";

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
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { Checkbox } from "@/components/ui/checkbox";
import { useToast } from "@/hooks/use-toast";
import { usePostPeople, useGetProjectsIdRoles } from "@api/moris";
import { createProjectRoleAssignedEvent } from "@/api/events";
import { OrcidSearchSelect } from "@/components/profile/OrcidSearchSelect";
import { OrcidPerson } from "@api/model";
import { MultiUserSelect } from "@/components/user/MultiUserSelect";
import { cn } from "@/lib/utils";

const addPersonSchema = z.object({
  name: z.string().min(1, "Name is required"),
  email: z.string().email("Invalid email address"),
  roles: z.array(z.string()).min(1, "At least one role is required"),
  orcid: z.string().optional(),
});

type Mode = "create" | "existing";

interface AddPersonDialogProps {
  projectId: string;
  onPersonAdded: () => void;
}

export function AddPersonDialog({
  projectId,
  onPersonAdded,
}: AddPersonDialogProps) {
  const [open, setOpen] = useState(false);
  const [mode, setMode] = useState<Mode>("create");
  const { toast } = useToast();

  // State for existing users mode
  const [selectedPersonIds, setSelectedPersonIds] = useState<string[]>([]);
  const [selectedRoles, setSelectedRoles] = useState<string[]>([]);

  const { mutateAsync: createPerson, isPending: isCreatingPerson } =
    usePostPeople();
  const { data: roles, isLoading: isLoadingRoles } =
    useGetProjectsIdRoles(projectId);

  const form = useForm<z.infer<typeof addPersonSchema>>({
    resolver: standardSchemaResolver(addPersonSchema),
    defaultValues: {
      name: "",
      email: "",
      roles: [],
      orcid: "",
    },
  });

  const [isAssigningRole, setIsAssigningRole] = useState(false);

  const isPending = isCreatingPerson || isAssigningRole;

  const resetDialog = () => {
    form.reset();
    setSelectedPersonIds([]);
    setSelectedRoles([]);
    setMode("create");
  };

  async function onSubmitCreate(values: z.infer<typeof addPersonSchema>) {
    try {
      // 1. Create the person
      const nameParts = values.name.split(" ");
      const givenName = nameParts[0];
      const familyName = nameParts.slice(1).join(" ") || "Unknown";

      const personData: any = {
        name: values.name,
        email: values.email,
        givenName: givenName,
        familyName: familyName,
        orcid: values.orcid,
      };

      const person = await createPerson({
        data: personData,
      });

      // 2. Assign all selected roles to the person
      if (person && person.id) {
        setIsAssigningRole(true);

        for (const roleId of values.roles) {
          await createProjectRoleAssignedEvent(projectId, {
            person_id: person.id,
            project_role_id: roleId,
          });
        }

        setIsAssigningRole(false);

        const roleCount = values.roles.length;
        toast({
          title: "Member added",
          description: `${values.name} has been successfully added to the project with ${roleCount} ${roleCount === 1 ? "role" : "roles"}.`,
        });

        setOpen(false);
        resetDialog();
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

  async function onSubmitExisting() {
    if (selectedPersonIds.length === 0) {
      toast({
        variant: "destructive",
        title: "No users selected",
        description: "Please select at least one user to add.",
      });
      return;
    }

    if (selectedRoles.length === 0) {
      toast({
        variant: "destructive",
        title: "No roles selected",
        description: "Please select at least one role for the users.",
      });
      return;
    }

    try {
      setIsAssigningRole(true);

      // Assign all selected roles to each selected person
      for (const personId of selectedPersonIds) {
        for (const roleId of selectedRoles) {
          await createProjectRoleAssignedEvent(projectId, {
            person_id: personId,
            project_role_id: roleId,
          });
        }
      }

      setIsAssigningRole(false);

      const userCount = selectedPersonIds.length;
      const roleCount = selectedRoles.length;
      toast({
        title: userCount === 1 ? "Member added" : "Members added",
        description: `${userCount} ${userCount === 1 ? "user" : "users"} added with ${roleCount} ${roleCount === 1 ? "role" : "roles"}.`,
      });

      setOpen(false);
      resetDialog();
      onPersonAdded();
    } catch (error) {
      console.error(error);
      setIsAssigningRole(false);
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to add members. Please try again.",
      });
    }
  }

  const handleOrcidSelect = (item: OrcidPerson) => {
    if (item.orcid) {
      form.setValue("orcid", item.orcid);
    }

    let name = "";
    if (item.credit_name) {
      name = item.credit_name;
    } else if (item.first_name || item.last_name) {
      name = [item.first_name, item.last_name].filter(Boolean).join(" ");
    }

    if (name) {
      form.setValue("name", name);
    }
  };

  const handleOpenChange = (isOpen: boolean) => {
    setOpen(isOpen);
    if (!isOpen) {
      resetDialog();
    }
  };

  const toggleRole = (roleId: string, checked: boolean) => {
    if (checked) {
      setSelectedRoles((prev) => [...prev, roleId]);
    } else {
      setSelectedRoles((prev) => prev.filter((id) => id !== roleId));
    }
  };

  const toggleFormRole = (roleId: string, checked: boolean) => {
    const currentRoles = form.getValues("roles");
    if (checked) {
      form.setValue("roles", [...currentRoles, roleId], {
        shouldValidate: true,
      });
    } else {
      form.setValue(
        "roles",
        currentRoles.filter((id) => id !== roleId),
        { shouldValidate: true },
      );
    }
  };

  const getRoleNames = (roleIds: string[]) => {
    return roleIds
      .map((id) => roles?.find((r) => r.id === id)?.name)
      .filter(Boolean);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button size="sm">
          <Plus className="mr-2 h-4 w-4" />
          Add Member
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[480px]">
        <DialogHeader>
          <DialogTitle>Add Team Member</DialogTitle>
          <DialogDescription>
            Create a new person or add existing users to this project.
          </DialogDescription>
        </DialogHeader>

        {/* Mode Toggle */}
        <div className="flex gap-2 p-1 bg-muted rounded-lg">
          <button
            type="button"
            onClick={() => setMode("create")}
            className={cn(
              "flex-1 flex items-center justify-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-colors",
              mode === "create"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground",
            )}
          >
            <UserPlus className="h-4 w-4" />
            Create New
          </button>
          <button
            type="button"
            onClick={() => setMode("existing")}
            className={cn(
              "flex-1 flex items-center justify-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-colors",
              mode === "existing"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground",
            )}
          >
            <Users className="h-4 w-4" />
            Add Existing
          </button>
        </div>

        {mode === "create" ? (
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmitCreate)}
              className="space-y-4"
            >
              <FormField
                control={form.control}
                name="orcid"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Search ORCID (Autofill)</FormLabel>
                    <FormControl>
                      <OrcidSearchSelect
                        value={field.value}
                        onSelect={handleOrcidSelect}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

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
                name="roles"
                render={() => (
                  <FormItem>
                    <FormLabel>Roles</FormLabel>
                    <div className="border rounded-md p-3 space-y-2">
                      {isLoadingRoles ? (
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <Loader2 className="h-4 w-4 animate-spin" />
                          Loading roles...
                        </div>
                      ) : (
                        roles?.map((role) => (
                          <div
                            key={role.id}
                            className="flex items-center space-x-2"
                          >
                            <Checkbox
                              id={`create-role-${role.id}`}
                              checked={form
                                .watch("roles")
                                .includes(role.id || "")}
                              onCheckedChange={(checked) =>
                                toggleFormRole(role.id || "", checked === true)
                              }
                            />
                            <label
                              htmlFor={`create-role-${role.id}`}
                              className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
                            >
                              {role.name}
                            </label>
                          </div>
                        ))
                      )}
                    </div>
                    {form.watch("roles").length > 0 && (
                      <div className="flex flex-wrap gap-1 mt-2">
                        {getRoleNames(form.watch("roles")).map((name) => (
                          <Badge
                            key={name}
                            variant="secondary"
                            className="text-xs"
                          >
                            {name}
                          </Badge>
                        ))}
                      </div>
                    )}
                    <FormMessage />
                  </FormItem>
                )}
              />
              <DialogFooter>
                <Button type="submit" disabled={isPending}>
                  {isPending && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  Add Member
                </Button>
              </DialogFooter>
            </form>
          </Form>
        ) : (
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>Select Users</Label>
              <MultiUserSelect
                value={selectedPersonIds}
                onChange={setSelectedPersonIds}
                placeholder="Search for users..."
              />
            </div>

            <div className="space-y-2">
              <Label>Roles</Label>
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
                        id={`existing-role-${role.id}`}
                        checked={selectedRoles.includes(role.id || "")}
                        onCheckedChange={(checked) =>
                          toggleRole(role.id || "", checked === true)
                        }
                      />
                      <label
                        htmlFor={`existing-role-${role.id}`}
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
            </div>

            <DialogFooter>
              <Button
                type="button"
                onClick={onSubmitExisting}
                disabled={
                  isPending ||
                  selectedPersonIds.length === 0 ||
                  selectedRoles.length === 0
                }
              >
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Add{" "}
                {selectedPersonIds.length > 0
                  ? selectedPersonIds.length
                  : ""}{" "}
                Member
                {selectedPersonIds.length !== 1 ? "s" : ""}
              </Button>
            </DialogFooter>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
