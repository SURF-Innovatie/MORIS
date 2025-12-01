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
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/use-toast";
import { usePostPeople, usePostProjectsIdPeoplePersonId } from "@api/moris";

const addPersonSchema = z.object({
    name: z.string().min(1, "Name is required"),
    email: z.string().email("Invalid email address"),
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

    const { mutateAsync: addPersonToProject, isPending: isAddingToProject } =
        usePostProjectsIdPeoplePersonId();
    const { mutateAsync: createPerson, isPending: isCreatingPerson } =
        usePostPeople();

    const form = useForm<z.infer<typeof addPersonSchema>>({
        resolver: zodResolver(addPersonSchema),
        defaultValues: {
            name: "",
            email: "",
        },
    });

    const isPending = isCreatingPerson || isAddingToProject;

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
                    user_id: "00000000-0000-0000-0000-000000000000", // Placeholder, backend handles this or it should be optional
                },
            });

            // 2. Add person to project
            if (person && person.id) {
                await addPersonToProject({
                    id: projectId,
                    personId: person.id,
                });

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
