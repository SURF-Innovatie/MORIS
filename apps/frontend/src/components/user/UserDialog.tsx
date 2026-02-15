import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema";
import { z } from "zod";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { useToast } from "@/hooks/use-toast";
import { UserResponse } from "@api/model";
import {
  usePostPeople,
  usePostUsers,
  usePutPeopleId,
  usePutUsersId,
  getGetAdminUsersListQueryKey,
} from "@api/moris";
import { useQueryClient } from "@tanstack/react-query";

const userSchema = z.object({
  name: z.string().min(1, "Name is required"),
  email: z.string().email("Invalid email address"),
  password: z.string().optional(),
  is_sys_admin: z.boolean().default(false),
  is_active: z.boolean().default(true),
});

type UserFormValues = z.infer<typeof userSchema>;

interface UserDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  user?: UserResponse | null; // If null, we are creating a new user
}

export function UserDialog({ open, onOpenChange, user }: UserDialogProps) {
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const isEditMode = !!user;

  const { mutateAsync: createPerson, isPending: isCreatingPerson } =
    usePostPeople();
  const { mutateAsync: createUser, isPending: isCreatingUser } = usePostUsers();
  const { mutateAsync: updatePerson, isPending: isUpdatingPerson } =
    usePutPeopleId();
  const { mutateAsync: updateUser, isPending: isUpdatingUser } =
    usePutUsersId();

  const form = useForm<UserFormValues>({
    resolver: standardSchemaResolver(userSchema),
    defaultValues: {
      name: "",
      email: "",
      password: "",
      is_sys_admin: false,
      is_active: true,
    },
  });

  useEffect(() => {
    if (open) {
      if (user) {
        form.reset({
          name: user.name,
          email: user.email,
          password: "", // Always clear password on edit
          is_sys_admin: user.is_sys_admin,
          is_active: user.is_active,
        });
      } else {
        form.reset({
          name: "",
          email: "",
          password: "",
          is_sys_admin: false,
          is_active: true,
        });
      }
    }
  }, [open, user, form]);

  const onSubmit = async (data: UserFormValues) => {
    try {
      if (isEditMode && user) {
        // Update existing
        // 1. Update Person details
        await updatePerson({
          id: user.person_id!,
          data: {
            name: data.name,
            email: data.email, // Email change might require verification in real app, but allowed here
          },
        });

        // 2. Update User details (only send password if provided)
        await updateUser({
          id: user.id!,
          data: {
            person_id: user.person_id!,
            is_sys_admin: data.is_sys_admin,
            ...(data.password ? { password: data.password } : {}),
          },
        });

        toast({ title: "Success", description: "User updated successfully" });
      } else {
        // Create new
        // 1. Create Person
        const person = await createPerson({
          data: {
            name: data.name,
            email: data.email,
          },
        });

        if (!person?.id) {
          throw new Error("Failed to create person");
        }

        // 2. Create User (password is optional for OAuth-only users)
        await createUser({
          data: {
            person_id: person.id,
            is_sys_admin: data.is_sys_admin,
            ...(data.password ? { password: data.password } : {}),
          },
        });

        toast({ title: "Success", description: "User created successfully" });
      }

      queryClient.invalidateQueries({
        queryKey: getGetAdminUsersListQueryKey(),
      });
      onOpenChange(false);
    } catch (error: any) {
      console.error("User operation failed", error);
      toast({
        variant: "destructive",
        title: "Error",
        description: error.response?.data?.message || "Failed to save user",
      });
    }
  };

  const isSubmitting =
    isCreatingPerson || isCreatingUser || isUpdatingPerson || isUpdatingUser;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{isEditMode ? "Edit User" : "Create User"}</DialogTitle>
          <DialogDescription>
            {isEditMode
              ? "Update user details, roles, or reset password."
              : "Create a new user account. This will also create a new Person profile."}
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
                    <Input placeholder="john@example.com" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    {isEditMode
                      ? "New Password (Optional)"
                      : "Password (Optional)"}
                  </FormLabel>
                  <FormControl>
                    <Input
                      type="password"
                      placeholder={
                        isEditMode
                          ? "Leave blank to keep current"
                          : "Leave blank for OAuth-only user"
                      }
                      {...field}
                    />
                  </FormControl>
                  <FormDescription>
                    {isEditMode
                      ? "Only enter a value if you want to change the password."
                      : "Leave blank if the user will only login via SURFconext or other OAuth providers."}
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="is_sys_admin"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                  <div className="space-y-0.5">
                    <FormLabel className="text-base">
                      System Administrator
                    </FormLabel>
                    <FormDescription>
                      Grant full administrative access.
                    </FormDescription>
                  </div>
                  <FormControl>
                    <Checkbox
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting ? "Saving..." : "Save User"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
