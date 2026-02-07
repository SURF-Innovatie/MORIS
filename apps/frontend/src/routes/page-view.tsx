import { useParams, useNavigate } from "react-router-dom";
import { useGetPagesSlug } from "@/api/generated-orval/moris";
import { PageViewer } from "@/components/page-builder/PageViewer";
import { Section } from "@/components/page-builder/types";
import { useAuth } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Pencil } from "lucide-react";

export default function PageViewRoute() {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();

  const { data: pageData, isLoading, error } = useGetPagesSlug(slug!);

  // Determine if current user can edit this page
  const canEdit = (() => {
    if (!user || !pageData) return false;
    // User pages: check if page belongs to this user
    if (pageData.user_id && pageData.user_id === user.id) return true;
    // Project pages: for now, allow any authenticated user to edit
    // TODO: Add proper project membership check
    if (pageData.project_id) return true;
    return false;
  })();

  if (isLoading) {
    return <div className="flex justify-center py-20">Loading page...</div>;
  }

  if (error || !pageData) {
    return (
      <div className="flex flex-col items-center justify-center py-20 text-center">
        <h1 className="text-2xl font-bold text-slate-800 mb-2">
          Page Not Found
        </h1>
        <p className="text-slate-500">
          The page you are looking for does not exist.
        </p>
        <Button variant="outline" className="mt-4" onClick={() => navigate(-1)}>
          <ArrowLeft className="mr-2 h-4 w-4" /> Go Back
        </Button>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white">
      {/* Header with Edit Button */}
      <header className="sticky top-0 z-10 bg-white/95 backdrop-blur border-b">
        <div className="container flex h-14 items-center justify-between">
          <div className="flex items-center gap-3">
            <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <h1 className="font-semibold text-lg">{pageData.title}</h1>
          </div>
          {canEdit && (
            <Button
              variant="outline"
              onClick={() => navigate(`/pages/${slug}/edit`)}
            >
              <Pencil className="mr-2 h-4 w-4" /> Edit Page
            </Button>
          )}
        </div>
      </header>

      <PageViewer sections={pageData.content as Section[]} />
    </div>
  );
}
