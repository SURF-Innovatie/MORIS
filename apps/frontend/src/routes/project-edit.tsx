import { useEffect } from "react";
import { useParams, useNavigate, useSearchParams } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2, ArrowLeft, Eye } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useToast } from "@/hooks/use-toast";
import { useGetProjectsId, usePutProjectsId } from "@api/moris";

import { GeneralTab } from "@/components/project-edit/GeneralTab";
import { PeopleTab } from "@/components/project-edit/PeopleTab";
import { ChangelogTab } from "@/components/project-edit/ChangelogTab";
import { ProductsTab } from "@/components/project-edit/ProductsTab";
import { projectFormSchema } from "@/components/project-edit/schema";

export default function ProjectEditRoute() {
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

  const { mutateAsync: updateProject, isPending: isUpdating } =
    usePutProjectsId();

  const form = useForm<z.infer<typeof projectFormSchema>>({
    resolver: zodResolver(projectFormSchema),
    defaultValues: {
      title: "",
      description: "",
      organisationID: "00000000-0000-0000-0000-000000000000",
    },
  });

  useEffect(() => {
    if (project) {
      form.reset({
        title: project.title || "",
        description: project.description || "",
        startDate: project.startDate ? new Date(project.startDate) : undefined,
        endDate: project.endDate ? new Date(project.endDate) : undefined,
        organisationID:
          project.organization?.id || "00000000-0000-0000-0000-000000000000",
      });
    }
  }, [project, form]);

  async function onSubmit(values: z.infer<typeof projectFormSchema>) {
    try {
      await updateProject({
        id: id!,
        data: {
          title: values.title,
          description: values.description,
          startDate: values.startDate.toISOString(),
          endDate: values.endDate.toISOString(),
          organisationID: values.organisationID,
        },
      });
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
              isUpdating={isUpdating}
              project={project}
            />
          </TabsContent>

          <TabsContent value="people">
            <PeopleTab
              projectId={id!}
              people={project?.people || []}
              adminId={project?.projectAdmin}
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
