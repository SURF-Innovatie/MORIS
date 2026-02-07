import { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema";
import { z } from "zod";
import { Loader2 } from "lucide-react";
import { v4 as uuidv4 } from "uuid";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { useToast } from "@/hooks/use-toast";
import {
  createProjectStartedEvent,
  createCustomFieldValueSetEvent,
} from "@/api/events";

import { projectFormSchema } from "@/lib/schemas/project";
import { ProjectFields } from "@/components/project-form/ProjectFields";
import { CustomFieldsRenderer } from "@/components/project-form/CustomFieldsRenderer";
import { OrganisationSearchSelect } from "@/components/organisation/OrganisationSearchSelect";
import { EMPTY_UUID } from "@/lib/constants";
import { slugify } from "@/lib/utils";

export default function CreateProjectRoute() {
  const navigate = useNavigate();
  const { toast } = useToast();
  const [isCreating, setIsCreating] = useState(false);

  const form = useForm<z.infer<typeof projectFormSchema>>({
    resolver: standardSchemaResolver(projectFormSchema),
    defaultValues: {
      title: "",
      description: "",
      // TODO: This should be dynamic or selected from a list
      organisationID: EMPTY_UUID,
    },
  });

  async function onSubmit(values: z.infer<typeof projectFormSchema>) {
    setIsCreating(true);
    try {
      const newProjectId = uuidv4();
      await createProjectStartedEvent(newProjectId, {
        title: values.title,
        slug: slugify(values.title),
        description: values.description,
        start_date: values.startDate.toISOString(),
        end_date: values.endDate.toISOString(),
        members_ids: [],
        owning_org_node_id: values.organisationID,
      });

      // Handle Custom Fields separately
      if (values.customFields) {
        const promises = Object.entries(values.customFields).map(
          ([defId, value]) => {
            let valStr = String(value);
            if (value instanceof Date) valStr = value.toISOString();

            return createCustomFieldValueSetEvent(newProjectId, {
              definition_id: defId,
              value: valStr, // Sending string for now.
            });
          },
        );
        await Promise.all(promises);
      }

      toast({
        title: "Project created",
        description: "The new project has been successfully created.",
      });
      navigate("/dashboard");
    } catch (error: any) {
      toast({
        variant: "destructive",
        title: "Error",
        description: error.message || "Something went wrong. Please try again.",
      });
    } finally {
      setIsCreating(false);
    }
  }

  return (
    <div className="mx-auto max-w-2xl py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Create Project</h1>
        <p className="text-muted-foreground">
          Fill in the details to start a new project.
        </p>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <ProjectFields form={form} />

          {/* Organisation ID is still specific here as it might be editable during creation but not update */}
          <FormField
            control={form.control}
            name="organisationID"
            render={({ field }) => (
              <FormItem className="max-w-2xl">
                <FormLabel>Organisation</FormLabel>
                <FormControl>
                  <OrganisationSearchSelect
                    value={field.value}
                    onSelect={(organisationId) =>
                      field.onChange(organisationId)
                    }
                    disabled={field.disabled}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <CustomFieldsRenderer
            form={form}
            organisationId={form.watch("organisationID")}
          />

          <div className="flex justify-end gap-4">
            <Button
              type="button"
              variant="ghost"
              onClick={() => navigate("/dashboard")}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isCreating}>
              {isCreating ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Saving...
                </>
              ) : (
                "Create Project"
              )}
            </Button>
          </div>
        </form>
      </Form>
    </div>
  );
}
