import { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Bell,
  Inbox,
  LayoutGrid,
  Table as TableIcon,
  Users as UsersIcon,
  Calendar,
  Building2,
  ExternalLink,
  CheckCircle2,
} from "lucide-react";

import { Badge } from "../components/ui/badge";
import { Button } from "../components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "../components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../components/ui/table";
import { useGetProjects } from "../api/generated-orval/moris";

const FAKE_NOTIFICATIONS = [
  {
    id: "1",
    type: "project",
    title: "New project assignment",
    description: 'You have been added to "Website Redesign 2024"',
    timestamp: "2 hours ago",
    read: false,
  },
  {
    id: "2",
    type: "message",
    title: "Team message from Sarah Chen",
    description: "Can you review the latest design mockups?",
    timestamp: "5 hours ago",
    read: false,
  },
  {
    id: "3",
    type: "update",
    title: "Project milestone completed",
    description: "Mobile App Development reached 75% completion",
    timestamp: "1 day ago",
    read: true,
  },
  {
    id: "4",
    type: "project",
    title: "Project deadline approaching",
    description: "API Integration project is due in 3 days",
    timestamp: "1 day ago",
    read: true,
  },
];

const DashboardRoute = () => {
  const navigate = useNavigate();
  const [viewMode, setViewMode] = useState<"cards" | "table">("cards");
  const { data: projects, isLoading, error } = useGetProjects();

  const formatDate = (dateString?: string) => {
    if (!dateString) return "N/A";
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  return (
    <div className="flex flex-col gap-8">
      {/* Notifications Inbox Section */}
      <section>
        <div className="mb-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Inbox className="h-5 w-5 text-muted-foreground" />
            <h2 className="text-2xl font-semibold tracking-tight">Inbox</h2>
            <Badge variant="default" className="ml-2">
              {FAKE_NOTIFICATIONS.filter((n) => !n.read).length}
            </Badge>
          </div>
          <Button variant="ghost" size="sm">
            <CheckCircle2 className="mr-2 h-4 w-4" />
            Mark all as read
          </Button>
        </div>

        <Card>
          <CardContent className="p-0">
            {FAKE_NOTIFICATIONS.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <Bell className="mb-4 h-12 w-12 text-muted-foreground/50" />
                <p className="text-sm text-muted-foreground">
                  No new notifications
                </p>
              </div>
            ) : (
              <div className="divide-y">
                {FAKE_NOTIFICATIONS.map((notification) => (
                  <div
                    key={notification.id}
                    className={`flex items-start gap-4 p-4 transition-colors hover:bg-muted/50 ${
                      !notification.read ? "bg-primary/5" : ""
                    }`}
                  >
                    <div className="mt-1">
                      {!notification.read && (
                        <div className="h-2 w-2 rounded-full bg-primary" />
                      )}
                    </div>
                    <div className="flex-1 space-y-1">
                      <p className="font-medium leading-none">
                        {notification.title}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {notification.description}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        {notification.timestamp}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </section>

      {/* Projects Section */}
      <section>
        <div className="mb-4 flex items-center justify-between">
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
            <Button
              variant={viewMode === "cards" ? "default" : "ghost"}
              size="sm"
              onClick={() => setViewMode("cards")}
            >
              <LayoutGrid className="mr-2 h-4 w-4" />
              Cards
            </Button>
            <Button
              variant={viewMode === "table" ? "default" : "ghost"}
              size="sm"
              onClick={() => setViewMode("table")}
            >
              <TableIcon className="mr-2 h-4 w-4" />
              Table
            </Button>
            <Button
              size="sm"
              onClick={() => navigate("/dashboard/projects/new")}
            >
              Create Project
            </Button>
          </div>
        </div>

        {isLoading && (
          <Card>
            <CardContent className="flex items-center justify-center py-12">
              <p className="text-sm text-muted-foreground">
                Loading projects...
              </p>
            </CardContent>
          </Card>
        )}

        {error && (
          <Card>
            <CardContent className="flex items-center justify-center py-12">
              <p className="text-sm text-destructive">
                Failed to load projects
              </p>
            </CardContent>
          </Card>
        )}

        {projects && projects.length === 0 && (
          <Card>
            <CardContent className="flex flex-col items-center justify-center py-12 text-center">
              <Building2 className="mb-4 h-12 w-12 text-muted-foreground/50" />
              <p className="text-sm text-muted-foreground">No projects found</p>
              <Button
                className="mt-4"
                size="sm"
                onClick={() => navigate("/dashboard/projects/new")}
              >
                Create your first project
              </Button>
            </CardContent>
          </Card>
        )}

        {projects && projects.length > 0 && viewMode === "cards" && (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {projects.map((project) => (
              <Card key={project.id} className="transition-all hover:shadow-lg">
                <CardHeader>
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <CardTitle className="line-clamp-1">
                        {project.title || "Untitled Project"}
                      </CardTitle>
                      <CardDescription className="mt-2 line-clamp-2">
                        {project.description || "No description available"}
                      </CardDescription>
                    </div>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8"
                      onClick={() => navigate(`/projects/${project.id}/edit`)}
                    >
                      <ExternalLink className="h-4 w-4" />
                    </Button>
                  </div>
                </CardHeader>
                <CardContent className="space-y-3">
                  {/* {project.organisation && (
                    <div className="flex items-center gap-2 text-sm">
                      <Building2 className="h-4 w-4 text-muted-foreground" />
                      <span className="text-muted-foreground">
                        {project.organisation}
                      </span>
                    </div>
                  )} */}
                  <div className="flex items-center gap-2 text-sm">
                    <Calendar className="h-4 w-4 text-muted-foreground" />
                    <span className="text-muted-foreground">
                      {formatDate(project.startDate)} -{" "}
                      {formatDate(project.endDate)}
                    </span>
                  </div>
                  {project.people && project.people.length > 0 && (
                    <div className="flex items-center gap-2 text-sm">
                      <UsersIcon className="h-4 w-4 text-muted-foreground" />
                      <span className="text-muted-foreground">
                        {project.people.length}{" "}
                        {project.people.length === 1 ? "member" : "members"}
                      </span>
                    </div>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {projects && projects.length > 0 && viewMode === "table" && (
          <Card>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-[250px]">Title</TableHead>
                    <TableHead>Description</TableHead>
                    <TableHead>Organisation</TableHead>
                    <TableHead>Start Date</TableHead>
                    <TableHead>End Date</TableHead>
                    <TableHead className="text-right">Members</TableHead>
                    <TableHead className="w-[50px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {projects.map((project) => (
                    <TableRow key={project.id}>
                      <TableCell className="font-medium">
                        {project.title || "Untitled Project"}
                      </TableCell>
                      <TableCell className="max-w-[300px] truncate">
                        {project.description || "No description"}
                      </TableCell>
                      <TableCell>
                        {/* {project.organisation || "N/A"} */}
                        N/A
                      </TableCell>
                      <TableCell>{formatDate(project.startDate)}</TableCell>
                      <TableCell>{formatDate(project.endDate)}</TableCell>
                      <TableCell className="text-right">
                        {project.people?.length || 0}
                      </TableCell>
                      <TableCell>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8"
                          onClick={() =>
                            navigate(`/projects/${project.id}/edit`)
                          }
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
    </div>
  );
};

export default DashboardRoute;
