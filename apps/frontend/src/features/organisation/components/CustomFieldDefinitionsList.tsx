import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2, Plus, Trash2, Info } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
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
import { Checkbox } from "@/components/ui/checkbox";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";

import { useToast } from "@/hooks/use-toast";
import {
    useGetOrganisationNodesIdCustomFields,
    usePostOrganisationNodesIdCustomFields,
    useDeleteOrganisationNodesIdCustomFieldsFieldId,
} from "@api/moris";

import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";

const createFieldSchema = z.object({
    name: z.string().min(1, "Name is required"),
    type: z.enum(["TEXT", "NUMBER", "BOOLEAN", "DATE"]),
    description: z.string().optional(),
    required: z.boolean().default(false),
    validation_regex: z.string().optional(),
    example_value: z.string().optional(),
});

interface CustomFieldDefinitionsListProps {
    nodeId: string;
}

export const CustomFieldDefinitionsList = ({ nodeId }: CustomFieldDefinitionsListProps) => {   
    const { data: fields, isLoading, refetch } = useGetOrganisationNodesIdCustomFields(nodeId);
    const { mutateAsync: deleteField } = useDeleteOrganisationNodesIdCustomFieldsFieldId();
    const { toast } = useToast();

    const handleDelete = async (fieldId: string) => {
        if (!confirm("Are you sure you want to delete this field? Data in projects might be lost or hidden.")) {
            return;
        }
        try {
            await deleteField({ id: nodeId, fieldId });
            toast({ title: "Custom field deleted" });
            refetch();
        } catch (error) {
            console.error("Failed to delete field", error);
            toast({ title: "Error deleting field", variant: "destructive" });
        }
    };

    return (
        <div className="space-y-4">
            <div className="flex justify-between items-center">
                <h3 className="text-lg font-medium">Custom Project Fields</h3>
                <AddCustomFieldDialog nodeId={nodeId} onSuccess={refetch} />
            </div>

            <div className="border rounded-lg overflow-hidden">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-[200px]">Name</TableHead>
                            <TableHead>Type</TableHead>
                            <TableHead>Required</TableHead>
                            <TableHead>Description</TableHead>
                            <TableHead className="text-right">Actions</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {isLoading ? (
                            <TableRow>
                                <TableCell colSpan={5} className="h-24 text-center">
                                    <div className="flex items-center justify-center">
                                        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                                    </div>
                                </TableCell>
                            </TableRow>
                        ) : fields?.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={5} className="h-24 text-center text-muted-foreground">
                                    No custom fields defined.
                                </TableCell>
                            </TableRow>
                        ) : (
                            fields?.map((field: any) => (
                                <TableRow key={field.id}>
                                    <TableCell className="font-medium">
                                        <div className="flex items-center gap-2">
                                            {field.name}
                                            {field.validation_regex && (
                                                <TooltipProvider>
                                                    <Tooltip>
                                                        <TooltipTrigger asChild>
                                                            <Info className="h-3 w-3 text-muted-foreground cursor-help" />
                                                        </TooltipTrigger>
                                                        <TooltipContent>
                                                            <p className="text-xs">Regex: {field.validation_regex}</p>
                                                        </TooltipContent>
                                                    </Tooltip>
                                                </TooltipProvider>
                                            )}
                                        </div>
                                    </TableCell>
                                    <TableCell className="text-xs font-mono bg-gray-50 rounded px-2 py-1 inline-block mt-2">{field.type}</TableCell>
                                    <TableCell>{field.required ? "Yes" : "No"}</TableCell>
                                    <TableCell className="text-muted-foreground max-w-[300px] truncate" title={field.description}>
                                        {field.description}
                                    </TableCell>
                                    <TableCell className="text-right">
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            onClick={() => handleDelete(field.id)}
                                            title="Delete Field"
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
                Note: These fields will be available to all projects within this organisation and its sub-organisations.
            </p>
        </div>
    );
};

const AddCustomFieldDialog = ({ nodeId, onSuccess }: { nodeId: string; onSuccess: () => void }) => {
    const [open, setOpen] = useState(false);
    const { toast } = useToast();
    const { mutateAsync: createField, isPending } = usePostOrganisationNodesIdCustomFields();

    const form = useForm<z.infer<typeof createFieldSchema>>({
        resolver: zodResolver(createFieldSchema),
        defaultValues: {
            name: "",
            type: "TEXT",
            description: "",
            required: false,
            validation_regex: "",
            example_value: "",
        },
    });

    async function onSubmit(values: z.infer<typeof createFieldSchema>) {
        try {
            await createField({
                id: nodeId,
                data: values
            });
            toast({
                title: "Field created",
                description: "Custom field has been successfully created.",
            });
            setOpen(false);
            form.reset();
            onSuccess();
        } catch (error) {
            console.error(error);
            toast({
                variant: "destructive",
                title: "Error",
                description: "Failed to create custom field.",
            });
        }
    }

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button size="sm">
                    <Plus className="mr-2 h-4 w-4" />
                    Create Field
                </Button>
            </DialogTrigger>
            <DialogContent className="max-w-lg">
                <DialogHeader>
                    <DialogTitle>Create Custom Field</DialogTitle>
                </DialogHeader>
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                        <div className="grid grid-cols-2 gap-4">
                            <FormField
                                control={form.control}
                                name="name"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Name</FormLabel>
                                        <FormControl>
                                            <Input placeholder="e.g. Cost Center" {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="type"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Type</FormLabel>
                                        <Select onValueChange={field.onChange} defaultValue={field.value}>
                                            <FormControl>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Select type" />
                                                </SelectTrigger>
                                            </FormControl>
                                            <SelectContent>
                                                <SelectItem value="TEXT">Text</SelectItem>
                                                <SelectItem value="NUMBER">Number</SelectItem>
                                                <SelectItem value="BOOLEAN">Boolean</SelectItem>
                                                <SelectItem value="DATE">Date</SelectItem>
                                            </SelectContent>
                                        </Select>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>

                        <FormField
                            control={form.control}
                            name="description"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>Description (Optional)</FormLabel>
                                    <FormControl>
                                        <Textarea placeholder="Describe the purpose of this field..." {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <div className="grid grid-cols-2 gap-4">
                            <FormField
                                control={form.control}
                                name="validation_regex"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Validation Regex (Optional)</FormLabel>
                                        <FormControl>
                                            <Input placeholder="e.g. ^[A-Z]{3}-\d{3}$" {...field} />
                                        </FormControl>
                                        <FormDescription className="text-xs">
                                            For TEXT fields only.
                                        </FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                            <FormField
                                control={form.control}
                                name="example_value"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel>Example Value (Optional)</FormLabel>
                                        <FormControl>
                                            <Input placeholder="e.g. ABC-123" {...field} />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        </div>

                        <FormField
                            control={form.control}
                            name="required"
                            render={({ field }) => (
                                <FormItem className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4">
                                    <FormControl>
                                        <Checkbox
                                            checked={field.value}
                                            onCheckedChange={field.onChange}
                                        />
                                    </FormControl>
                                    <div className="space-y-1 leading-none">
                                        <FormLabel>
                                            Required field
                                        </FormLabel>
                                        <FormDescription>
                                            Projects must provide a value for this field.
                                        </FormDescription>
                                    </div>
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
