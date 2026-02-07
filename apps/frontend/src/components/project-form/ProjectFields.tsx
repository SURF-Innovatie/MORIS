import { UseFormReturn } from "react-hook-form";

import { Badge } from "@/components/ui/badge";
import {
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { ProjectFormValues } from "@/lib/schemas/project";
import { DatePicker } from "../ui/datepicker";

interface ProjectFieldsProps {
  form: UseFormReturn<ProjectFormValues>;
  disabledFields?: {
    title?: boolean;
    description?: boolean;
    startDate?: boolean;
    endDate?: boolean;
    slug?: boolean;
  };
  pendingFields?: {
    title?: boolean;
    description?: boolean;
    startDate?: boolean;
    endDate?: boolean;
  };
}

export function ProjectFields({
  form,
  disabledFields,
  pendingFields,
}: ProjectFieldsProps) {
  const startDate = form.watch("startDate");

  const isPending = (key: keyof NonNullable<typeof pendingFields>) =>
    pendingFields?.[key];

  return (
    <div className="space-y-6">
      <FormField
        control={form.control}
        name="title"
        render={({ field }) => (
          <FormItem className="max-w-2xl">
            <FormLabel className="flex items-center gap-2">
              Title
              {isPending("title") && (
                <Badge variant="secondary" className="h-5 text-[10px]">
                  Pending Approval
                </Badge>
              )}
            </FormLabel>
            <FormControl>
              <Input
                placeholder="Project title"
                {...field}
                disabled={disabledFields?.title || isPending("title")}
                onChange={(e) => {
                  field.onChange(e);
                  // Auto-generate slug if slug is pristine?
                  // Implementation detail: we need access to setValue/getValues
                  // For now let's just let user type or implementing syncing logic inside component
                  const currentSlug = form.getValues("slug");
                  if (
                    !currentSlug ||
                    currentSlug === slugify(e.target.value.slice(0, -1))
                  ) {
                    form.setValue("slug", slugify(e.target.value), {
                      shouldValidate: true,
                    });
                  }
                }}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name="slug"
        render={({ field }) => (
          <SlugField
            field={field}
            form={form}
            disabled={disabledFields?.slug}
          />
        )}
      />

      <FormField
        control={form.control}
        name="description"
        render={({ field }) => (
          <FormItem className="max-w-2xl">
            <FormLabel className="flex items-center gap-2">
              Description
              {isPending("description") && (
                <Badge variant="secondary" className="h-5 text-[10px]">
                  Pending Approval
                </Badge>
              )}
            </FormLabel>
            <FormControl>
              <Textarea
                placeholder="Describe the project..."
                className="min-h-[120px] resize-none"
                {...field}
                disabled={
                  disabledFields?.description || isPending("description")
                }
              />
            </FormControl>
            <FormDescription>
              A brief summary of what this project is about.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <div className="grid gap-4 sm:grid-cols-2 max-w-2xl">
        <FormField
          control={form.control}
          name="startDate"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel className="flex items-center gap-2">
                Start Date
                {isPending("startDate") && (
                  <Badge variant="secondary" className="h-5 text-[10px]">
                    Pending
                  </Badge>
                )}
              </FormLabel>
              <DatePicker
                initialDate={field.value}
                onDateChange={(date) => {
                  field.onChange(date);
                }}
                disabled={(date) => date < new Date("1900-01-01")}
              />
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="endDate"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel className="flex items-center gap-2">
                End Date
                {isPending("endDate") && (
                  <Badge variant="secondary" className="h-5 text-[10px]">
                    Pending
                  </Badge>
                )}
              </FormLabel>
              <DatePicker
                initialDate={field.value}
                onDateChange={(date) => {
                  field.onChange(date);
                }}
                disabled={(date) =>
                  date < new Date("1900-01-01") ||
                  (startDate ? date < startDate : false)
                }
              />
              <FormMessage />
            </FormItem>
          )}
        />
      </div>
    </div>
  );
}

import { slugify } from "@/lib/utils";
import { useGetProjectsSlugCheck } from "@/api/generated-orval/moris";
import { useEffect } from "react";
import { Check, X, Loader2 } from "lucide-react";

function SlugField({
  field,
  form,
  disabled,
}: {
  field: any;
  form: UseFormReturn<ProjectFormValues>;
  disabled?: boolean;
}) {
  const slug = field.value;
  const { data, isLoading } = useGetProjectsSlugCheck(
    { slug },
    { query: { enabled: !!slug && slug.length > 2 && !disabled, retry: 0 } },
  );

  // The generated hook returns the response directly as data
  const isAvailable = data?.available;

  // Custom validation effect
  useEffect(() => {
    if (disabled) return;
    if (slug && data && !data.available) {
      form.setError("slug", {
        type: "manual",
        message: "Slug is already taken",
      });
    } else if (slug && data && data.available) {
      form.clearErrors("slug");
    }
  }, [data, slug, form, disabled]);

  return (
    <FormItem className="max-w-2xl">
      <FormLabel>Slug</FormLabel>
      <FormControl>
        <div className="relative">
          <Input
            placeholder="project-slug"
            {...field}
            disabled={disabled}
            onChange={(e) => field.onChange(slugify(e.target.value))}
          />
          <div className="absolute right-3 top-2.5">
            {isLoading ? (
              <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
            ) : slug && data ? (
              isAvailable ? (
                <Check className="h-4 w-4 text-green-500" />
              ) : (
                <X className="h-4 w-4 text-red-500" />
              )
            ) : null}
          </div>
        </div>
      </FormControl>
      <FormDescription>
        Unique identifier for the project URL. Auto-generated from title but can
        be customized.
      </FormDescription>
      <FormMessage />
    </FormItem>
  );
}
