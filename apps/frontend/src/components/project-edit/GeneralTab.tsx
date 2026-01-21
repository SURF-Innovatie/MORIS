import { UseFormReturn } from "react-hook-form";
import { Loader2, Save } from "lucide-react";

import { OrganisationSearchSelect } from "@/components/organisation/OrganisationSearchSelect";

import { Button } from "@/components/ui/button";
import { Form } from "@/components/ui/form";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ProjectFormValues } from "@/lib/schemas/project";
import { ProjectResponse } from "@api/model";
import {
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Badge } from "@/components/ui/badge";
import { ProjectFields } from "@/components/project-form/ProjectFields";
import { CustomFieldsRenderer } from "@/components/project-form/CustomFieldsRenderer";
import { useAccess } from "@/context/AccessContext";
import { ProjectEventType, ProjectEvent } from "@/api/events";

interface GeneralTabProps {
  form: UseFormReturn<ProjectFormValues>;
  onSubmit: (values: ProjectFormValues) => Promise<void>;
  isUpdating: boolean;
  project?: ProjectResponse;
  pendingEvents?: ProjectEvent[];
}

export function GeneralTab({
  form,
  onSubmit,
  isUpdating,
  project,
  pendingEvents,
}: GeneralTabProps) {
  const { hasAccess } = useAccess();

  const pendingFields = {
    title: pendingEvents?.some((e) => e.type === ProjectEventType.TitleChanged),
    description: pendingEvents?.some(
      (e) => e.type === ProjectEventType.DescriptionChanged,
    ),
    startDate: pendingEvents?.some(
      (e) => e.type === ProjectEventType.StartDateChanged,
    ),
    endDate: pendingEvents?.some(
      (e) => e.type === ProjectEventType.EndDateChanged,
    ),
    owningOrg: pendingEvents?.some(
      (e) => e.type === ProjectEventType.OwningOrgNodeChanged,
    ),
  };

  const disabledFields = {
    title: !hasAccess(ProjectEventType.TitleChanged) || pendingFields.title,
    description:
      !hasAccess(ProjectEventType.DescriptionChanged) ||
      pendingFields.description,
    startDate:
      !hasAccess(ProjectEventType.StartDateChanged) || pendingFields.startDate,
    endDate:
      !hasAccess(ProjectEventType.EndDateChanged) || pendingFields.endDate,
    owningOrg:
      !hasAccess(ProjectEventType.OwningOrgNodeChanged) ||
      pendingFields.owningOrg,
  };

  const oneFieldEditable =
    !disabledFields.title ||
    !disabledFields.description ||
    !disabledFields.startDate ||
    !disabledFields.endDate ||
    !disabledFields.owningOrg;

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <div className="grid gap-8 lg:grid-cols-3">
          <div className="lg:col-span-2 space-y-8">
            <Card>
              <CardHeader>
                <CardTitle>Project Details</CardTitle>
                <CardDescription>
                  Update the core information about your project.
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <ProjectFields
                  form={form}
                  disabledFields={disabledFields}
                  pendingFields={pendingFields}
                />

                {project?.owning_org_node?.id && (
                  <CustomFieldsRenderer
                    form={form}
                    organisationId={project.owning_org_node.id}
                  />
                )}

                <div className="flex justify-start">
                  <Button
                    type="submit"
                    disabled={isUpdating || !oneFieldEditable}
                  >
                    {isUpdating ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Saving...
                      </>
                    ) : (
                      <>
                        <Save className="mr-2 h-4 w-4" />
                        Save Changes
                      </>
                    )}
                  </Button>
                </div>
              </CardContent>
            </Card>
          </div>

          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="text-sm font-medium">
                  Organization
                  {pendingFields.owningOrg && (
                    <Badge variant="secondary" className="ml-2 h-5 text-[10px]">
                      Pending Approval
                    </Badge>
                  )}
                </CardTitle>
              </CardHeader>
              <CardContent>
                <FormField
                  control={form.control}
                  name="organisationID"
                  render={({ field }) => (
                    <FormItem>
                      <FormControl>
                        <OrganisationSearchSelect
                          value={field.value}
                          onSelect={(orgId) => {
                            field.onChange(orgId);
                          }}
                          disabled={isUpdating || disabledFields.owningOrg}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </CardContent>
            </Card>
          </div>
        </div>
      </form>
    </Form>
  );
}
