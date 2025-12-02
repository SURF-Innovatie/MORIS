import { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  LayoutGrid,
  Table as TableIcon,
  Building2,
  ExternalLink,
  Calendar,
} from "lucide-react";
import { format } from "date-fns";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useGetProjects } from "@api/moris";
import { ProjectResponse } from "@api/model";

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

const formatDate = (dateString?: string) => {
  if (!dateString) return "N/A";
  return format(new Date(dateString), "MMM d, yyyy");
};

interface ProjectListProps {
  showCreateButton?: boolean;
}

export const ProjectList = ({ showCreateButton = true }: ProjectListProps) => {
  const navigate = useNavigate();
  const [viewMode, setViewMode] = useState<"cards" | "table">("cards");

  const {
    data: projects,
    isLoading: isLoadingProjects,
    error: projectsError,
  } = useGetProjects();

  return (
    <section>
      <div className="mb-6 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Building2 className="h-5 w-5 text-muted-foreground" />
          <h2 className="text-2xl font-semibold tracking-tight">Projects</h2>
          {projects && (
            <Badge variant="outline" className="ml-2">
              {projects.length}
            </Badge>
          )}
        </div>
        <div className="flex items-center gap-2">
          <div className="flex items-center rounded-lg border bg-muted/50 p-1">
            <Button
              variant={viewMode === "cards" ? "secondary" : "ghost"}
              size="sm"
              className="h-7 px-2"
              onClick={() => setViewMode("cards")}
            >
              <LayoutGrid className="h-4 w-4" />
            </Button>
            <Button
              variant={viewMode === "table" ? "secondary" : "ghost"}
              size="sm"
              className="h-7 px-2"
              onClick={() => setViewMode("table")}
            >
              <TableIcon className="h-4 w-4" />
            </Button>
          </div>
          {showCreateButton && (
            <Button
              size="sm"
              onClick={() => navigate("/dashboard/projects/new")}
            >
              Create Project
            </Button>
          )}
        </div>
      </div>

      {isLoadingProjects && (
        <Card>
          <CardContent className="flex items-center justify-center py-12">
            <p className="text-sm text-muted-foreground">Loading projects...</p>
          </CardContent>
        </Card>
      )}

      {projectsError && (
        <Card>
          <CardContent className="flex items-center justify-center py-12">
            <p className="text-sm text-destructive">Failed to load projects</p>
          </CardContent>
        </Card>
      )}

      {projects && projects.length === 0 && (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12 text-center">
            <Building2 className="mb-4 h-12 w-12 text-muted-foreground/50" />
            <p className="text-sm text-muted-foreground">No projects found</p>
            {showCreateButton && (
              <Button
                className="mt-4"
                size="sm"
                onClick={() => navigate("/dashboard/projects/new")}
              >
                Create your first project
              </Button>
            )}
          </CardContent>
        </Card>
      )}

      {projects && projects.length > 0 && viewMode === "cards" && (
        <div className="grid gap-6 sm:grid-cols-2">
          {projects.map((project) => {
            const status = getProjectStatus(project);

            return (
              <Card
                key={project.id}
                className="group flex flex-col transition-all hover:shadow-md hover:border-primary/20"
                onClick={() => navigate(`/projects/${project.id}/edit`)}
              >
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between gap-4">
                    <div className="space-y-1">
                      <CardTitle className="line-clamp-1 text-base">
                        {project.title || "Untitled Project"}
                      </CardTitle>
                      <Badge
                        variant={status.variant}
                        className="text-[10px] px-1.5 py-0 h-5"
                      >
                        {status.label}
                      </Badge>
                    </div>
                  </div>
                  <CardDescription className="line-clamp-2 mt-2 text-xs">
                    {project.description ||
                      "No description available for this project."}
                  </CardDescription>
                </CardHeader>
                <CardContent className="pb-3 flex-1">
                  <div className="flex items-center gap-2 text-xs text-muted-foreground mb-4">
                    <Calendar className="h-3.5 w-3.5" />
                    <span>
                      {formatDate(project.startDate)} -{" "}
                      {formatDate(project.endDate)}
                    </span>
                  </div>

                  <div className="flex items-center -space-x-2 overflow-hidden">
                    {project.people && project.people.length > 0 ? (
                      <>
                        {project.people.slice(0, 4).map((person, i) => (
                          <Avatar
                            key={i}
                            className="h-6 w-6 ring-2 ring-background"
                          >
                            <AvatarImage
                              src={person.avatar_url || ""}
                              alt={person.email}
                            />
                            <AvatarFallback className="text-[10px] bg-muted text-muted-foreground">
                              {person.email?.charAt(0).toUpperCase() || "?"}
                            </AvatarFallback>
                          </Avatar>
                        ))}
                        {project.people.length > 4 && (
                          <div className="h-6 w-6 rounded-full ring-2 ring-background bg-muted flex items-center justify-center text-[10px] font-medium text-muted-foreground text-center">
                            +{project.people.length - 4}
                          </div>
                        )}
                      </>
                    ) : (
                      <span className="text-xs text-muted-foreground italic">
                        No members
                      </span>
                    )}
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>
      )}

      {projects && projects.length > 0 && viewMode === "table" && (
        <Card>
          <CardContent className="p-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[250px]">Title</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Dates</TableHead>
                  <TableHead className="text-right">Members</TableHead>
                  <TableHead className="w-[50px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {projects.map((project) => (
                  <TableRow key={project.id}>
                    <TableCell className="font-medium">
                      <div className="flex flex-col">
                        <span>{project.title || "Untitled Project"}</span>
                        <span className="text-xs text-muted-foreground truncate max-w-[200px]">
                          {project.description}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={getProjectStatus(project).variant}
                        className="text-[10px] h-5"
                      >
                        {getProjectStatus(project).label}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">
                      {formatDate(project.startDate)} -{" "}
                      {formatDate(project.endDate)}
                    </TableCell>
                    <TableCell className="text-right">
                      {project.people?.length || 0}
                    </TableCell>
                    <TableCell>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                        onClick={() => navigate(`/projects/${project.id}/edit`)}
                      >
                        <ExternalLink className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      )}
    </section>
  );
};
