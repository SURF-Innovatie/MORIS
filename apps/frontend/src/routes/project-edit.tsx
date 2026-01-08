import { useEffect, useState } from "react";
import { useParams, useNavigate, useSearchParams } from "react-router-dom";
import { useForm } from "react-hook-form";
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema";
import { z } from "zod";
import { Loader2, ArrowLeft, Eye } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useToast } from "@/hooks/use-toast";
import { useGetProjectsId } from "@api/moris";
import {
  createTitleChangedEvent,
  createDescriptionChangedEvent,
  createStartDateChangedEvent,
  createEndDateChangedEvent,
  createOwningOrgNodeChangedEvent,
  createCustomFieldValueSetEvent,
} from "@/api/events";

import { GeneralTab } from "@/components/project-edit/GeneralTab";
import { PeopleTab } from "@/components/project-edit/PeopleTab";
import { ChangelogTab } from "@/components/project-edit/ChangelogTab";
import { ProductsTab } from "@/components/project-edit/ProductsTab";
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
    if (project) {
      form.reset({
        title: project.title || "",
        description: project.description || "",
        startDate: project.start_date
          ? new Date(project.start_date)
          : undefined,
        endDate: project.end_date ? new Date(project.end_date) : undefined,
        organisationID: project.owning_org_node?.id || EMPTY_UUID,
        customFields: project.custom_fields || {},
      });
    }
  }, [project, form]);

  async function onSubmit(values: z.infer<typeof projectFormSchema>) {
    if (!project) return;

    setIsSaving(true);
    try {
      const promises: Promise<any>[] = [];

      // Compare and emit events for changed fields
      if (values.title !== project.title) {
        promises.push(createTitleChangedEvent(id!, { title: values.title }));
      }

      if (values.description !== project.description) {
        promises.push(
          createDescriptionChangedEvent(id!, {
            description: values.description,
          })
        );
      }

      const currentStartDate = project.start_date
        ? new Date(project.start_date).toISOString()
        : null;
      if (values.startDate.toISOString() !== currentStartDate) {
        promises.push(
          createStartDateChangedEvent(id!, {
            start_date: values.startDate.toISOString(),
          })
        );
      }

      const currentEndDate = project.end_date
        ? new Date(project.end_date).toISOString()
        : null;
      if (values.endDate.toISOString() !== currentEndDate) {
        promises.push(
          createEndDateChangedEvent(id!, {
            end_date: values.endDate.toISOString(),
          })
        );
      }

      const currentOrgNodeId = project.owning_org_node?.id || EMPTY_UUID;
      if (values.organisationID !== currentOrgNodeId) {
        promises.push(
          createOwningOrgNodeChangedEvent(id!, {
            owning_org_node_id: values.organisationID,
          })
        );
      }

      // Handle Custom Fields updates
      if (values.customFields) {
        const currentFields = project.custom_fields || {};
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
                 project_id: id!,
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
      <header className="sticky top-0 z-10 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
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
                {project?.title || "Project Settings"}
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
            <TabsList className="grid w-full max-w-md grid-cols-4">
              <TabsTrigger value="general">General</TabsTrigger>
              <TabsTrigger value="people">People</TabsTrigger>
              <TabsTrigger value="products">Products</TabsTrigger>
              <TabsTrigger value="changelog">Changelog</TabsTrigger>
            </TabsList>
          </div>

          <TabsContent value="general" className="space-y-4">
            <GeneralTab
              form={form}
              onSubmit={onSubmit}
              isUpdating={isSaving}
              // project pass-through might need check if GeneralTab uses project props
              // But looking at previous code, it just passed `project`.
              // GeneralTab probably needs updating if it uses snake_case props?
              project={project}
            />
          </TabsContent>

          <TabsContent value="people">
            <PeopleTab
              projectId={id!}
              members={project?.members || []}
              onRefresh={refetchProject}
            />
          </TabsContent>

          <TabsContent value="products">
            <ProductsTab
              projectId={id!}
              products={project?.products || []}
              onRefresh={refetchProject}
            />
          </TabsContent>

          <TabsContent value="changelog">
            <ChangelogTab projectId={id!} />
          </TabsContent>
        </Tabs>
      </main>
    </div>
  );
}
