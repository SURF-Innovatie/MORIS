import { format } from "date-fns";
import { Building2, Calendar, Package, Sparkles } from "lucide-react";

import { ProjectResponse } from "@api/model";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { getProjectStatus, formatDateRange } from "@/lib/format";

interface ProjectHighlightsProps {
  featuredProjects: ProjectResponse[];
}

export const ProjectHighlights = ({
  featuredProjects,
}: ProjectHighlightsProps) => {
  return (
    <Card className="lg:col-span-2">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Sparkles className="h-5 w-5 text-primary" />
          Project Highlights
        </CardTitle>
        <CardDescription>
          Featured work showing recent contributions and outcomes.
        </CardDescription>
      </CardHeader>
      <CardContent className="grid gap-4 md:grid-cols-2">
        {featuredProjects.length === 0 ? (
          <div className="col-span-full rounded-lg border border-dashed p-6 text-center text-sm text-muted-foreground">
            No projects to showcase yet.
          </div>
        ) : (
          featuredProjects.map((project) => {
            const status = getProjectStatus(project);
            return (
              <div
                key={project.id}
                className="rounded-xl border bg-background p-4 shadow-sm"
              >
                <div className="flex items-start justify-between gap-2">
                  <div>
                    <h3 className="text-base font-semibold">
                      {project.title || "Untitled Project"}
                    </h3>
                    <p className="text-xs text-muted-foreground">
                      {formatDateRange(project)}
                    </p>
                  </div>
                  <Badge variant={status.variant}>{status.label}</Badge>
                </div>
                <p className="mt-3 line-clamp-3 text-sm text-muted-foreground">
                  {project.description ||
                    "No description added yet for this project."}
                </p>
                <div className="mt-4 flex flex-wrap gap-2 text-xs text-muted-foreground">
                  {project.owning_org_node?.name && (
                    <span className="inline-flex items-center gap-1">
                      <Building2 className="h-3 w-3" />
                      {project.owning_org_node.name}
                    </span>
                  )}
                  <span className="inline-flex items-center gap-1">
                    <Package className="h-3 w-3" />
                    {project.products?.length ?? 0} deliverables
                  </span>
                  {project.start_date && (
                    <span className="inline-flex items-center gap-1">
                      <Calendar className="h-3 w-3" />
                      {format(new Date(project.start_date), "MMM yyyy")}
                    </span>
                  )}
                </div>
              </div>
            );
          })
        )}
      </CardContent>
    </Card>
  );
};
