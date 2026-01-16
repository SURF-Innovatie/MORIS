import { UseFormReturn } from "react-hook-form";
import { useGetOrganisationNodesIdCustomFields } from "@api/moris";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { ProjectFormValues } from "@/lib/schemas/project";
import { format, parseISO } from "date-fns";
import { CalendarIcon } from "lucide-react";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { cn } from "@/lib/utils";

interface CustomFieldsRendererProps {
  form: UseFormReturn<ProjectFormValues>;
  organisationId: string;
}

export const CustomFieldsRenderer = ({
  form,
  organisationId,
}: CustomFieldsRendererProps) => {
  const { data: fields, isLoading } = useGetOrganisationNodesIdCustomFields(
    organisationId,
    undefined,
    {
      query: {
        enabled:
          !!organisationId &&
          organisationId !== "00000000-0000-0000-0000-000000000000",
      }, // Skip if empty UUID
    }
  );

  if (isLoading || !fields || fields.length === 0) return null;

  return (
    <div className="space-y-4 border-t pt-4">
      <h3 className="text-lg font-medium">Additional Information</h3>
      <div className="grid gap-4 sm:grid-cols-2">
        {fields.map((field) => (
          <FormField
            key={field.id}
            control={form.control}
            name={`customFields.${field.id}`}
            rules={{
              required: field.required ? "This field is required" : false,
              pattern: field.validation_regex
                ? {
                    value: new RegExp(field.validation_regex),
                    message: "Invalid format",
                  }
                : undefined,
            }}
            render={({ field: formField }) => (
              <FormItem className="flex flex-col">
                <FormLabel>
                  {field.name}
                  {field.required && (
                    <span className="text-red-500 ml-1">*</span>
                  )}
                </FormLabel>
                <div className="flex-1">
                  <CustomFieldInput
                    type={field.type as any}
                    field={formField}
                    placeholder={field.example_value || ""}
                  />
                </div>
                {field.description && (
                  <FormDescription>{field.description}</FormDescription>
                )}
                <FormMessage />
              </FormItem>
            )}
          />
        ))}
      </div>
    </div>
  );
};

const CustomFieldInput = ({
  type,
  field,
  placeholder,
}: {
  type: "TEXT" | "NUMBER" | "BOOLEAN" | "DATE";
  field: any;
  placeholder: string;
}) => {
  if (type === "BOOLEAN") {
    return (
      <div className="flex items-center space-x-2 h-10">
        <Checkbox
          checked={field.value === true || field.value === "true"}
          onCheckedChange={(checked: boolean) => field.onChange(checked)}
        />
      </div>
    );
  }

  if (type === "DATE") {
    return (
      <Popover>
        <PopoverTrigger asChild>
          <FormControl>
            <Button
              variant={"outline"}
              className={cn(
                "w-full pl-3 text-left font-normal",
                !field.value && "text-muted-foreground"
              )}
            >
              {field.value ? (
                format(
                  typeof field.value === "string"
                    ? parseISO(field.value)
                    : field.value,
                  "PPP"
                )
              ) : (
                <span>Pick a date</span>
              )}
              <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
            </Button>
          </FormControl>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0" align="start">
          <Calendar
            mode="single"
            selected={
              field.value
                ? typeof field.value === "string"
                  ? parseISO(field.value)
                  : field.value
                : undefined
            }
            onSelect={(date) => field.onChange(date?.toISOString())}
            initialFocus
          />
        </PopoverContent>
      </Popover>
    );
  }

  return (
    <FormControl>
      <Input
        {...field}
        type={type === "NUMBER" ? "number" : "text"}
        placeholder={placeholder}
        value={field.value || ""} // Ensure controlled input
      />
    </FormControl>
  );
};
