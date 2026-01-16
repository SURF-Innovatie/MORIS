import { useEffect, useState, useMemo } from "react";
import { useParams, useNavigate, useSearchParams } from "react-router-dom";
import { useForm } from "react-hook-form";
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema";
import { z } from "zod";
import { Loader2, ArrowLeft, Eye } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useToast } from "@/hooks/use-toast";
import {
  useGetProjectsId,
  useGetProjectsIdPendingEvents,
  usePostProjectsProjectIdBudget,
} from "@api/moris";
import { useQueryClient } from "@tanstack/react-query";
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
import { BudgetOverview } from "@/components/budget/BudgetOverview";
import { BudgetEditor } from "@/components/budget/BudgetEditor";
import { projectFormSchema } from "@/lib/schemas/project";
import { EMPTY_UUID } from "@/lib/constants";
import { ProjectAccessProvider } from "@/context/ProjectAccessContext";

export default function ProjectEditRoute() {
  return (
    <ProjectAccessProvider>
      <ProjectEditForm />
    </ProjectAccessProvider>
  );
}

function ProjectEditForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { toast } = useToast();
  const queryClient = useQueryClient();

  const [searchParams, setSearchParams] = useSearchParams();
  const activeTab = searchParams.get("tab") || "general";
  const setActiveTab = (tab: string) => {
    setSearchParams(
      (prev) => {
        prev.set("tab", tab);
        return prev;
      },
      { replace: true }
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
  const [isEditingBudget, setIsEditingBudget] = useState(false);

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

      // Compare and emit events for changed fields
      if (values.title !== projectedProject.title) {
        promises.push(createTitleChangedEvent(id!, { title: values.title }));
      }

      if (values.description !== projectedProject.description) {
        promises.push(
          createDescriptionChangedEvent(id!, {
            description: values.description,
          })
        );
      }

      const currentStartDate = projectedProject.start_date
        ? new Date(projectedProject.start_date).toISOString()
        : null;
      if (values.startDate.toISOString() !== currentStartDate) {
        promises.push(
          createStartDateChangedEvent(id!, {
            start_date: values.startDate.toISOString(),
          })
        );
      }

      const currentEndDate = projectedProject.end_date
        ? new Date(projectedProject.end_date).toISOString()
        : null;
      if (values.endDate.toISOString() !== currentEndDate) {
        promises.push(
          createEndDateChangedEvent(id!, {
            end_date: values.endDate.toISOString(),
          })
        );
      }

      const currentOrgNodeId =
        projectedProject.owning_org_node?.id || EMPTY_UUID;
      if (values.organisationID !== currentOrgNodeId) {
        promises.push(
          createOwningOrgNodeChangedEvent(id!, {
            owning_org_node_id: values.organisationID,
          })
        );
      }

      // Handle Custom Fields updates
      if (values.customFields) {
        const currentFields = projectedProject.custom_fields || {};
        Object.entries(values.customFields).map(([defId, value]) => {
          let valStr = String(value);
          if (value instanceof Date) valStr = value.toISOString();

          let currentValStr = String(currentFields[defId]);
          // Handle Date comparison if needed, or if current is undefined
          if (currentFields[defId] === undefined) currentValStr = "";

          if (valStr !== currentValStr) {
            // For Boolean: form true -> "true", current true (bool) -> "true"
            promises.push(
              createCustomFieldValueSetEvent(id!, {
                definition_id: defId,
                value: valStr,
              })
            );
          }
        });
      }

      if (promises.length === 0) {
        toast({
          title: "No changes",
          description: "No changes were detected to save.",
        });
        setIsSaving(false);
        return;
      }

      await Promise.all(promises);
      await refetchProject();
      await refetchPending();

      toast({
        title: "Project updated",
        description: "The project details have been successfully saved.",
      });
    } catch (error) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to update project. Please try again.",
      });
    } finally {
      setIsSaving(false);
    }
  }

  const createBudgetMutation = usePostProjectsProjectIdBudget();

  const handleCreateBudget = async () => {
    try {
      if (!id) return;
      await createBudgetMutation.mutateAsync({
        projectId: id,
        data: {
          title: `${projectedProject?.title} Budget`,
          description: "Initial budget draft",
        },
      });
      await queryClient.invalidateQueries({ queryKey: ["budget", id] });
      toast({
        title: "Budget Created",
        description: "A new budget has been initialized for this project.",
      });
    } catch (error) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to create budget. Please try again.",
      });
    }
  };

  if (isLoadingProject) {
    return (
      <div className="flex h-screen items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      {/* Header */}
      <header className="sticky top-0 z-10 border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
        <div className="container flex h-16 items-center justify-between py-4">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => navigate("/dashboard")}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div className="flex flex-col">
              <h1 className="text-lg font-semibold leading-none tracking-tight">
                {projectedProject?.title || "Project Settings"}
              </h1>
              <p className="text-sm text-muted-foreground">
                Manage your project settings and team
              </p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={() => navigate(`/projects/${id}`)}
            >
              <Eye className="mr-2 h-4 w-4" />
              View Page
            </Button>
          </div>
        </div>
      </header>

      <main className="container flex-1 py-8">
        <Tabs
          value={activeTab}
          onValueChange={setActiveTab}
          className="space-y-8"
        >
          <div className="flex items-center justify-between">
            <TabsList className="grid w-full max-w-2xl grid-cols-6">
              <TabsTrigger value="general">General</TabsTrigger>
              <TabsTrigger value="people">People</TabsTrigger>
              <TabsTrigger value="products">Products</TabsTrigger>
              <TabsTrigger value="budget">Budget</TabsTrigger>
              <TabsTrigger value="policies">Policies</TabsTrigger>
              <TabsTrigger value="changelog">Changelog</TabsTrigger>
            </TabsList>
          </div>

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

          <TabsContent value="budget">
            {isEditingBudget ? (
              <BudgetEditor
                projectId={id!}
                onDone={() => setIsEditingBudget(false)}
              />
            ) : (
              <BudgetOverview
                projectId={id!}
                onCreateBudget={handleCreateBudget}
                onEditBudget={() => setIsEditingBudget(true)}
              />
            )}
          </TabsContent>
        </Tabs>
      </main>
    </div>
  );
}

function applyPendingEvents(
  project: any, // Using any here to allow augmentation with pending flags easily
  events: any[]
): any {
  if (!events || events.length === 0) return project;

  // Deep clone to avoid mutating original
  const p = JSON.parse(JSON.stringify(project));

  // Sort events by date? usually they come sorted? Assuming they are sorted chronologically or we trust the order.
  // Actually the order matters.

  for (const e of events) {
    if (e.status !== "pending") continue; // should only be pending events here anyway based on hook

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
        // We only have ID, so we patch it partial
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
            (prod: any) => prod.id !== e.data.product_id
          );
        }
        break;
      case ProjectEventType.ProjectRoleAssigned:
        if (e.person && e.projectRole) {
          p.members = p.members || [];
          // Construct member object
          const newMember = {
            id: `pending-${e.person.id}-${e.projectRole.id}`, // Temporary ID
            user_id: e.person.id,
            name: `${e.person.givenName} ${e.person.familyName}`.trim(),
            email: e.person.email,
            avatarUrl: e.person.avatarUrl, // Assuming it might be here
            role: e.projectRole.slug, // Checking role slug
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
              )
          );
        }
        break;
    }
  }

  return p;
}
