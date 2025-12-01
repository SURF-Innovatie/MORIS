import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2, ArrowLeft } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useToast } from "@/hooks/use-toast";
import { useGetProjectsId, usePutProjectsId } from "@api/moris";

import { GeneralTab } from "@/components/project-edit/GeneralTab";
import { PeopleTab } from "@/components/project-edit/PeopleTab";
import { ChangelogTab } from "@/components/project-edit/ChangelogTab";
import { projectFormSchema } from "@/components/project-edit/schema";

// Mock data for People
const MOCK_PEOPLE = [
  {
    id: "1",
    name: "Alice Johnson",
    email: "alice@example.com",
    role: "Admin",
    avatar: null,
  },
  {
    id: "2",
    name: "Bob Smith",
    email: "bob@example.com",
    role: "Member",
    avatar: null,
  },
  {
    id: "3",
    name: "Charlie Brown",
    email: "charlie@example.com",
    role: "Viewer",
    avatar: null,
  },
];

export default function ProjectEditRoute() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { toast } = useToast();

  // State for mock features
  const [activeTab, setActiveTab] = useState("general");
  const [people, setPeople] = useState(MOCK_PEOPLE);

  const handleAddMember = () => {
    const newMember = {
      id: Math.random().toString(),
      name: "New Member",
      email: "new@example.com",
      role: "Viewer",
      avatar: null,
    };
    setPeople([...people, newMember]);
  };

  const { data: project, isLoading: isLoadingProject } = useGetProjectsId(id!, {
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
          project.organisation || "00000000-0000-0000-0000-000000000000",
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
            {/* Status badge placeholder */}
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
            <TabsList className="grid w-full max-w-md grid-cols-3">
              <TabsTrigger value="general">General</TabsTrigger>
              <TabsTrigger value="people">People</TabsTrigger>
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
            <PeopleTab people={people} onAddMember={handleAddMember} />
          </TabsContent>

          <TabsContent value="changelog">
            <ChangelogTab projectId={id!} />
          </TabsContent>
        </Tabs>
      </main>
    </div>
  );
}
