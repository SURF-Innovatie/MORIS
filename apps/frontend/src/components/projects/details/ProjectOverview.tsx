import { format } from "date-fns";
import { Calendar, Building2 } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ProjectResponse } from "@api/model";

interface ProjectOverviewProps {
    project: ProjectResponse;
}

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

export function ProjectOverview({ project }: ProjectOverviewProps) {
    const status = getProjectStatus(project);

    return (
        <Card>
            <CardHeader>
                <div className="flex items-start justify-between">
                    <div className="space-y-1">
                        <CardTitle className="text-2xl">{project.title}</CardTitle>
                        <div className="flex items-center gap-2 pt-2">
                            <Badge variant={status.variant}>{status.label}</Badge>
                            {project.organization && (
                                <Badge variant="outline" className="flex items-center gap-1">
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
    );
}
