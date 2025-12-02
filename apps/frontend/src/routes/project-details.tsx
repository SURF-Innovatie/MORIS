import { useNavigate, useParams } from "react-router-dom";
import { format } from "date-fns";
import {
  ArrowLeft,
  Calendar,
  Building2,
  Users,
  Pencil,
  Loader2,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useGetProjectsId } from "@api/moris";
import { ProjectResponse } from "@api/model";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

const getProjectStatus = (project: ProjectResponse) => {
  if (!project.startDate || !project.endDate)
    return { label: "Unknown", variant: "secondary" as const };

  const now = new Date();
  const start = new Date(project.startDate);
  const end = new Date(project.endDate);

  if (now < start) return { label: "Upcoming", variant: "secondary" as const };
  if (now > end) return { label: "Completed", variant: "outline" as const };
  return { label: "Active", variant: "default" as const };
};

export default function ProjectDetailsRoute() {
  const { id } = useParams();
  const navigate = useNavigate();

  const {
    data: project,
    isLoading,
    error,
  } = useGetProjectsId(id!, {
    query: {
      enabled: !!id,
    },
  });

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="flex h-screen flex-col items-center justify-center gap-4">
        <p className="text-destructive">Failed to load project details.</p>
        <Button variant="outline" onClick={() => navigate("/dashboard")}>
          Back to Dashboard
        </Button>
      </div>
    );
  }

  const status = getProjectStatus(project);

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
                Project Details
              </h1>
            </div>
          </div>
          <Button onClick={() => navigate(`/projects/${id}/edit`)}>
            <Pencil className="mr-2 h-4 w-4" />
            Edit Project
          </Button>
        </div>
      </header>

      <main className="container flex-1 py-8 space-y-8">
        {/* Project Header Card */}
        <Card>
          <CardHeader>
            <div className="flex items-start justify-between">
              <div className="space-y-1">
                <CardTitle className="text-2xl">{project.title}</CardTitle>
                <div className="flex items-center gap-2 pt-2">
                  <Badge variant={status.variant}>{status.label}</Badge>
                  {project.organization && (
                    <Badge
                      variant="outline"
                      className="flex items-center gap-1"
                    >
                      <Building2 className="h-3 w-3" />
                      {project.organization.name}
                    </Badge>
                  )}
                </div>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <div className="grid gap-6 md:grid-cols-2">
              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium text-muted-foreground mb-2">
                    Description
                  </h3>
                  <p className="text-sm leading-relaxed">
                    {project.description || "No description provided."}
                  </p>
                </div>
              </div>
              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-medium text-muted-foreground mb-2">
                    Timeline
                  </h3>
                  <div className="flex items-center gap-2 text-sm">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <span>
                      {project.startDate
                        ? format(new Date(project.startDate), "MMMM d, yyyy")
                        : "N/A"}{" "}
                      -{" "}
                      {project.endDate
                        ? format(new Date(project.endDate), "MMMM d, yyyy")
                        : "N/A"}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Team Section */}
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <Users className="h-5 w-5 text-muted-foreground" />
            <h2 className="text-lg font-semibold">Team Members</h2>
            <Badge variant="secondary" className="ml-2">
              {project.people?.length || 0}
            </Badge>
          </div>
          <Card>
            <CardContent className="p-0">
              {project.people && project.people.length > 0 ? (
                <div className="divide-y">
                  {project.people.map((person) => (
                    <div
                      key={person.id}
                      className="flex items-center justify-between p-4"
                    >
                      <div className="flex items-center gap-3">
                        <Avatar className="h-9 w-9">
                          <AvatarImage
                            src={person.avatar_url || ""}
                            alt={person.email}
                          />
                          <AvatarFallback className="text-xs font-medium uppercase">
                            {person.email?.charAt(0) || "?"}
                          </AvatarFallback>
                        </Avatar>
                        <div className="space-y-0.5">
                          <p className="text-sm font-medium leading-none">
                            {person.givenName} {person.familyName}
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {person.email}
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        {project.projectAdmin === person.id && (
                          <Badge className="text-xs">Admin</Badge>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="flex flex-col items-center justify-center py-8 text-center">
                  <Users className="mb-2 h-8 w-8 text-muted-foreground/30" />
                  <p className="text-sm text-muted-foreground">
                    No team members added yet.
                  </p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  );
}
