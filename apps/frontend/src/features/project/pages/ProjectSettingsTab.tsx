import { useEffect, useState, useMemo } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import { useForm } from "react-hook-form";
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema";
import { z } from "zod";
import { Loader2 } from "lucide-react";

import { toast } from "sonner";
import { useGetProjectsId, useGetProjectsIdPendingEvents } from "@api/moris";
import {
  createTitleChangedEvent,
  createDescriptionChangedEvent,
  createStartDateChangedEvent,
  createEndDateChangedEvent,
  createOwningOrgNodeChangedEvent,
  createCustomFieldValueSetEvent,
} from "@/api/events";

import { GeneralTab } from "@/components/project-edit/GeneralTab";
import { ProjectEventPoliciesTab } from "@/components/project-edit/ProjectEventPoliciesTab";
import { projectFormSchema } from "@/lib/schemas/project";
import { EMPTY_UUID } from "@/lib/constants";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { applyPendingEvents } from "@/lib/events/projection";

/**
 * ProjectSettingsTab - Admin settings for project configuration
 *
 * This component provides admin controls for project settings with two sub-tabs:
 * - General: Basic project information (title, description, dates, organization, custom fields)
 * - Policies: Event approval policies configuration
 *
 * Uses URL params for sub-tab navigation (?tab=general or ?tab=policies)
 * Only accessible to users with edit permissions (enforced at route level)
 */
export function ProjectSettingsTab() {
  const { id } = useParams();

  const [searchParams, setSearchParams] = useSearchParams();
  const activeTab = searchParams.get("tab") || "general";
  const setActiveTab = (tab: string) => {
    setSearchParams(
      (prev) => {
        prev.set("tab", tab);
        return prev;
      },
      { replace: true },
    );
  };

  const {
    data: project,
    isLoading: isLoadingProject,
    refetch: refetchProject,
  } = useGetProjectsId(id!, {
    query: {
      enabled: !!id,
    },
  });

  const { data: pendingEventsData, refetch: refetchPending } =
    useGetProjectsIdPendingEvents(id!, {
      query: {
        enabled: !!id,
      },
    });

  const projectedProject = useMemo(() => {
    if (!project) return undefined;
    if (!pendingEventsData?.events) return project;
    return applyPendingEvents(project, pendingEventsData.events);
  }, [project, pendingEventsData]);

  const [isSaving, setIsSaving] = useState(false);

  const form = useForm<z.infer<typeof projectFormSchema>>({
    resolver: standardSchemaResolver(projectFormSchema),
    defaultValues: {
      title: "",
      description: "",
      organisationID: EMPTY_UUID,
    },
  });

  useEffect(() => {
    if (projectedProject) {
      form.reset({
        title: projectedProject.title || "",
        description: projectedProject.description || "",
        startDate: projectedProject.start_date
          ? new Date(projectedProject.start_date)
          : undefined,
        endDate: projectedProject.end_date
          ? new Date(projectedProject.end_date)
          : undefined,
        organisationID: projectedProject.owning_org_node?.id || EMPTY_UUID,
        customFields: projectedProject.custom_fields || {},
      });
    }
  }, [projectedProject, form]);

  useEffect(() => {
    const handleRefresh = () => {
      refetchPending();
    };

    window.addEventListener("notifications:should-refresh", handleRefresh);

    return () => {
      window.removeEventListener("notifications:should-refresh", handleRefresh);
    };
  }, [refetchPending]);

  async function onSubmit(values: z.infer<typeof projectFormSchema>) {
    if (!projectedProject) return;

    setIsSaving(true);
    try {
      const promises: Promise<any>[] = [];

      if (values.title !== projectedProject.title) {
        promises.push(createTitleChangedEvent(id!, { title: values.title }));
      }

      if (values.description !== projectedProject.description) {
        promises.push(
          createDescriptionChangedEvent(id!, {
            description: values.description,
          }),
        );
      }

      const currentStartDate = projectedProject.start_date
        ? new Date(projectedProject.start_date).toISOString()
        : null;
      if (values.startDate?.toISOString() !== currentStartDate) {
        promises.push(
          createStartDateChangedEvent(id!, {
            start_date: values.startDate!.toISOString(),
          }),
        );
      }

      const currentEndDate = projectedProject.end_date
        ? new Date(projectedProject.end_date).toISOString()
        : null;
      if (values.endDate?.toISOString() !== currentEndDate) {
        promises.push(
          createEndDateChangedEvent(id!, {
            end_date: values.endDate!.toISOString(),
          }),
        );
      }

      const currentOrgNodeId =
        projectedProject.owning_org_node?.id || EMPTY_UUID;
      if (values.organisationID !== currentOrgNodeId) {
        promises.push(
          createOwningOrgNodeChangedEvent(id!, {
            owning_org_node_id: values.organisationID,
          }),
        );
      }

      // Handle Custom Fields updates
      if (values.customFields) {
        const currentFields = projectedProject.custom_fields || {};
        Object.entries(values.customFields).map(([defId, value]) => {
          let valStr = String(value);
          if (value instanceof Date) valStr = value.toISOString();

          let currentValStr = String(currentFields[defId]);
          if (currentFields[defId] === undefined) currentValStr = "";

          if (valStr !== currentValStr) {
            promises.push(
              createCustomFieldValueSetEvent(id!, {
                definition_id: defId,
                value: valStr,
              }),
            );
          }
        });
      }

      if (promises.length === 0) {
        toast("No changes", {
          description: "No changes were detected to save.",
        });
        setIsSaving(false);
        return;
      }

      await Promise.all(promises);
      await refetchProject();
      await refetchPending();

      toast.success("Project updated", {
        description: "The project details have been successfully saved.",
      });
    } catch (error) {
      toast.error("Error", {
        description: "Failed to update project. Please try again.",
      });
    } finally {
      setIsSaving(false);
    }
  }

  if (isLoadingProject) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!projectedProject) {
    return (
      <div className="flex h-64 items-center justify-center text-muted-foreground">
        Project not found
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <Tabs
        value={activeTab}
        onValueChange={setActiveTab}
        className="space-y-6"
      >
        <TabsList className="grid w-full max-w-md grid-cols-2">
          <TabsTrigger value="general">General</TabsTrigger>
          <TabsTrigger value="policies">Policies</TabsTrigger>
        </TabsList>

        <TabsContent value="general" className="space-y-4">
          <GeneralTab
            form={form}
            onSubmit={onSubmit}
            isUpdating={isSaving}
            project={projectedProject}
            pendingEvents={pendingEventsData?.events as any}
          />
        </TabsContent>

        <TabsContent value="policies">
          <ProjectEventPoliciesTab
            projectId={id!}
            orgNodeId={projectedProject?.owning_org_node?.id || ""}
          />
        </TabsContent>
      </Tabs>
    </div>
  );
}
