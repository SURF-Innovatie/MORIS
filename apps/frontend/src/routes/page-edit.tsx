import { useParams, useNavigate, useSearchParams } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { PageBuilder } from "@/components/page-builder/PageBuilder";
import {
  useGetPagesSlug,
  usePostPages,
  usePutPagesId,
  getGetProjectsProjectIdPagesQueryKey,
  getGetUsersUserIdPagesQueryKey,
} from "@/api/generated-orval/moris";
import { Section } from "@/components/page-builder/types";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { toast } from "sonner";
import { useEffect, useState } from "react";

export default function PageEditRoute() {
  const { slug } = useParams<{ slug: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const type = searchParams.get("type") || "project"; // 'project' or 'profile'
  const projectId = searchParams.get("projectId");
  const userId = searchParams.get("userId");

  // Fetch existing page data if we are editing an existing slug
  // For new pages, we might need a different flow or just start empty
  const pageSlug = slug ?? "";
  const { data: pageData, isLoading } = useGetPagesSlug(pageSlug, {
    query: {
      enabled: !!slug && slug !== "new",
    },
  });

  const [initialSections, setInitialSections] = useState<Section[]>([]);

  useEffect(() => {
    if (pageData?.content) {
      setInitialSections(pageData.content as Section[]);
    } else if (!slug || slug === "new") {
      // Apply templates based on type
      if (type === "profile") {
        setInitialSections([
          {
            id: crypto.randomUUID(),
            type: "profile_header",
            data: {
              name: "Your Name",
              role: "Researcher",
              bio: "Tell us about yourself...",
            },
          },
          {
            id: crypto.randomUUID(),
            type: "rich_text",
            data: {
              content:
                "<h2>About Me</h2><p>Share your background and interests.</p>",
            },
          },
          {
            id: crypto.randomUUID(),
            type: "links",
            data: {
              title: "Connect",
              links: [],
            },
          },
        ]);
      } else {
        // Default Project Template
        setInitialSections([
          {
            id: crypto.randomUUID(),
            type: "hero",
            data: {
              title: "Project Title",
              subtitle: "A brief tagline for your project",
            },
          },
          {
            id: crypto.randomUUID(),
            type: "statistics",
            data: {
              title: "Impact",
              stats: [
                { id: "1", value: "10+", label: "Publications" },
                { id: "2", value: "5", label: "Partners" },
              ],
            },
          },
          {
            id: crypto.randomUUID(),
            type: "rich_text",
            data: {
              content:
                "<h2>Overview</h2><p>Describe your project goals and methodology here.</p>",
            },
          },
        ]);
      }
    }
  }, [slug, pageData, type]);

  const createPage = usePostPages();
  const updatePage = usePutPagesId();

  const handleSave = async (sections: Section[]) => {
    try {
      if (!slug || slug === "new") {
        // Create new page
        // TODO: Get title/slug from a pre-step or form
        // For now, generating a random one, but in real app we'd likely have a modal or separate field
        const newTitle = type === "profile" ? "My Profile" : "New Project Page";
        const newSlug = `${type}-${Date.now()}`;

        await createPage.mutateAsync({
          data: {
            title: newTitle,
            slug: newSlug,
            type: type,
            content: sections,
            is_published: true,
            project_id: projectId || undefined,
            user_id: userId || undefined,
          },
        });

        // Invalidate related queries so buttons update
        if (projectId) {
          queryClient.invalidateQueries({
            queryKey: getGetProjectsProjectIdPagesQueryKey(projectId),
          });
        }
        if (userId) {
          queryClient.invalidateQueries({
            queryKey: getGetUsersUserIdPagesQueryKey(userId),
          });
        }

        toast.success("Page created successfully");
        // Navigate to the newly created page's view
        navigate(`/pages/${newSlug}`, { replace: true });
      } else if (pageData && pageData.id) {
        // Update existing page
        await updatePage.mutateAsync({
          id: pageData.id,
          data: {
            content: sections,
            // title: ...
          },
        });
        toast.success("Page updated successfully");
        queryClient.invalidateQueries({ queryKey: ["/pages", slug] });
      }
    } catch (error) {
      console.error(error);
      toast.error("Failed to save page");
    }
  };

  if (isLoading) {
    return <div>Loading editor...</div>;
  }

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col">
      <header className="bg-white border-b px-6 py-4 flex items-center gap-4 sticky top-0 z-50">
        <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div>
          <h1 className="text-xl font-bold text-slate-900">
            {pageData?.title || "Create New Page"}
          </h1>
          <p className="text-sm text-slate-500">Page Builder</p>
        </div>
      </header>

      <main className="flex-1 p-6 overflow-auto">
        <PageBuilder initialSections={initialSections} onSave={handleSave} />
      </main>
    </div>
  );
}
