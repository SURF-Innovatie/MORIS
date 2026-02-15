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
import { usePostOrganisationNodesIdCustomFields } from "@api/moris";

export const createFieldSchema = z.object({
  name: z.string().min(1, "Name is required"),
  type: z.enum(["TEXT", "NUMBER", "BOOLEAN", "DATE"]),
  category: z.enum(["PROJECT", "PERSON"]),
  description: z.string().optional(),
  required: z.boolean().default(false),
  validation_regex: z.string().optional(),
  example_value: z.string().optional(),
});

interface AddCustomFieldDialogProps {
  nodeId: string;
  onSuccess: () => void;
}

export const AddCustomFieldDialog = ({
  nodeId,
  onSuccess,
}: AddCustomFieldDialogProps) => {
  const [open, setOpen] = useState(false);
  const { toast } = useToast();
  const { mutateAsync: createField, isPending } =
    usePostOrganisationNodesIdCustomFields();

  const form = useForm<z.infer<typeof createFieldSchema>>({
    resolver: standardSchemaResolver(createFieldSchema),
    defaultValues: {
      name: "",
      type: "TEXT",
      category: "PROJECT",
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
        data: values,
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
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
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
              name="category"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Category</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select category" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="PROJECT">Project Field</SelectItem>
                      <SelectItem value="PERSON">Person Field</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormDescription>
                    Project fields appear on projects. Person fields appear on
                    member profiles.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description (Optional)</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Describe the purpose of this field..."
                      {...field}
                    />
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
                    <FormLabel>Required field</FormLabel>
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
