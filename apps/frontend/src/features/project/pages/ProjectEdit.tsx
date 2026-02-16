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
  ProjectEventType,
} from "@/api/events";

import { GeneralTab } from "@/components/project-edit/GeneralTab";
import { PeopleTab } from "@/components/project-edit/PeopleTab";
import { ChangelogTab } from "@/components/project-edit/ChangelogTab";
import { ProductsTab } from "@/components/project-edit/ProductsTab";
import { ProjectEventPoliciesTab } from "@/components/project-edit/ProjectEventPoliciesTab";
import { projectFormSchema } from "@/lib/schemas/project";
import { EMPTY_UUID } from "@/lib/constants";
import { ProjectAccessProvider } from "@/contexts/ProjectAccessContext";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

// Helper function to apply pending events (same as in original file)
function applyPendingEvents(project: any, events: any[]): any {
  if (!events || events.length === 0) return project;

  const p = JSON.parse(JSON.stringify(project));

  for (const e of events) {
    if (e.status !== "pending") continue;

    switch (e.type) {
      case ProjectEventType.TitleChanged:
        if (e.data?.title) p.title = e.data.title;
        break;
      case ProjectEventType.DescriptionChanged:
        if (e.data?.description) p.description = e.data.description;
        break;
      case ProjectEventType.StartDateChanged:
        if (e.data?.start_date) p.start_date = e.data.start_date;
        break;
      case ProjectEventType.EndDateChanged:
        if (e.data?.end_date) p.end_date = e.data.end_date;
        break;
      case ProjectEventType.OwningOrgNodeChanged:
        if (e.data?.owning_org_node_id) {
          p.owning_org_node = {
            ...(p.owning_org_node || {}),
            id: e.data.owning_org_node_id,
          };
        }
        break;
      case ProjectEventType.CustomFieldValueSet:
        if (e.data?.definition_id) {
          p.custom_fields = p.custom_fields || {};
          p.custom_fields[e.data.definition_id] = e.data.value;
        }
        break;
      case ProjectEventType.ProductAdded:
        if (e.product && e.product.id) {
          p.products = p.products || [];
          p.products.push({ ...e.product, pending: true });
        }
        break;
      case ProjectEventType.ProductRemoved:
        if (e.data?.product_id) {
          p.products = (p.products || []).filter(
            (prod: any) => prod.id !== e.data.product_id,
          );
        }
        break;
      case ProjectEventType.ProjectRoleAssigned:
        if (e.person && e.projectRole) {
          p.members = p.members || [];
          const newMember = {
            id: `pending-${e.person.id}-${e.projectRole.id}`,
            user_id: e.person.id,
            name: `${e.person.givenName} ${e.person.familyName}`.trim(),
            email: e.person.email,
            avatarUrl: e.person.avatarUrl,
            role: e.projectRole.slug,
            role_id: e.projectRole.id,
            role_name: e.projectRole.name,
            pending: true,
          };
          p.members.push(newMember);
        }
        break;
      case ProjectEventType.ProjectRoleUnassigned:
        if (e.data?.person_id && e.data?.project_role_id) {
          p.members = (p.members || []).filter(
            (m: any) =>
              !(
                m.user_id === e.data.person_id &&
                m.role_id === e.data.project_role_id
              ),
          );
        }
        break;
    }
  }
  return p;
}

export const ProjectEdit = () => {
  return (
    <ProjectAccessProvider>
      <ProjectEditForm />
    </ProjectAccessProvider>
  );
};

function ProjectEditForm() {
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
    // ... (Same submission logic as original, copied for brevity but ideally shared)
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
      if (values.startDate.toISOString() !== currentStartDate) {
        promises.push(
          createStartDateChangedEvent(id!, {
            start_date: values.startDate.toISOString(),
          }),
        );
      }

      const currentEndDate = projectedProject.end_date
        ? new Date(projectedProject.end_date).toISOString()
        : null;
      if (values.endDate.toISOString() !== currentEndDate) {
        promises.push(
          createEndDateChangedEvent(id!, {
            end_date: values.endDate.toISOString(),
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

  return (
    <div className="space-y-6">
      <Tabs
        value={activeTab}
        onValueChange={setActiveTab}
        className="space-y-6"
      >
        {/* Recreated Tabs List inside the page content if helpful, or rely on ProjectLayout tabs? 
               ProjectLayout tabs are for top-level navigation (Overview, Edit, Team).
               Inside "Edit", we might have sub-tabs like General, Policies, Changelog. 
           */}
        <TabsList className="grid w-full max-w-xl grid-cols-5">
          <TabsTrigger value="general">General</TabsTrigger>
          <TabsTrigger value="people">People</TabsTrigger>
          <TabsTrigger value="products">Products</TabsTrigger>
          <TabsTrigger value="policies">Policies</TabsTrigger>
          <TabsTrigger value="changelog">Changelog</TabsTrigger>
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

        <TabsContent value="people">
          <PeopleTab
            projectId={id!}
            members={projectedProject?.members || []}
            onRefresh={refetchProject}
          />
        </TabsContent>

        <TabsContent value="products">
          <ProductsTab
            projectId={id!}
            products={projectedProject?.products || []}
            onRefresh={refetchProject}
          />
        </TabsContent>

        <TabsContent value="changelog">
          <ChangelogTab projectId={id!} />
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
