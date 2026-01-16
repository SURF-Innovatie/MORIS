import { useParams } from "react-router-dom";
import { OrgAnalyticsDashboard } from "@/components/analytics/OrgAnalyticsDashboard";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";

export default function OrgAnalyticsRoute() {
  const { orgId } = useParams();
  const navigate = useNavigate();

  if (!orgId) {
    return <div>Organisation ID is required</div>;
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <header className="sticky top-0 z-10 border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
        <div className="container flex h-16 items-center flex-row justify-between py-4">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-lg font-semibold leading-none tracking-tight">
                Organisation Analytics
              </h1>
              <p className="text-sm text-muted-foreground">
                Financial overview and budget tracking
              </p>
            </div>
          </div>
        </div>
      </header>

      <main className="container flex-1 py-8">
        <OrgAnalyticsDashboard orgId={orgId} />
      </main>
    </div>
  );
}
