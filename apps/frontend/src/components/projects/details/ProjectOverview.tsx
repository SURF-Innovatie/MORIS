import { format } from "date-fns";
import { Calendar, Building2 } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ProjectResponse, CustomFieldDefinitionResponse } from "@api/model";
import { useGetOrganisationNodesIdCustomFields } from "@api/moris";

interface ProjectOverviewProps {
    project: ProjectResponse;
}

const getProjectStatus = (project: ProjectResponse) => {
    if (!project.start_date || !project.end_date)
        return { label: "Unknown", variant: "secondary" as const };

    const now = new Date();
    const start = new Date(project.start_date);
    const end = new Date(project.end_date);

    if (now < start) return { label: "Upcoming", variant: "secondary" as const };
    if (now > end) return { label: "Completed", variant: "outline" as const };
    return { label: "Active", variant: "default" as const };
};

export function ProjectOverview({ project }: ProjectOverviewProps) {
    const status = getProjectStatus(project);

    const { data: customFields } = useGetOrganisationNodesIdCustomFields(
        project.owning_org_node?.id!,
        { query: { enabled: !!project.owning_org_node?.id } }
    );

    return (
        <Card>
            <CardHeader>
                <div className="flex items-start justify-between">
                    <div className="space-y-1">
                        <CardTitle className="text-2xl">{project.title}</CardTitle>
                        <div className="flex items-center gap-2 pt-2">
                            <Badge variant={status.variant}>{status.label}</Badge>
                            {project.owning_org_node && (
                                <Badge variant="outline" className="flex items-center gap-1">
                                    <Building2 className="h-3 w-3" />
                                    {project.owning_org_node.name}
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
                                    {project.start_date
                                        ? format(new Date(project.start_date), "MMMM d, yyyy")
                                        : "N/A"}{" "}
                                    -{" "}
                                    {project.end_date
                                        ? format(new Date(project.end_date), "MMMM d, yyyy")
                                        : "N/A"}
                                </span>
                            </div>
                        </div>

                        {/* Custom Fields Display */}
                        {project.custom_fields && Object.keys(project.custom_fields).length > 0 && customFields && (
                            <div>
                                <h3 className="text-sm font-medium text-muted-foreground mb-2 mt-4">
                                    Additional Information
                                </h3>
                                <div className="grid grid-cols-2 gap-4">
                                    {Object.entries(project.custom_fields).map(([id, value]) => {
                                        const def = customFields.find(f => f.id === id);
                                        if (!def) return null;
                                        
                                        let displayValue = String(value);
                                        if (def.type === "BOOLEAN") displayValue = value ? "Yes" : "No";
                                        if (def.type === "DATE" && value) {
                                            try {
                                                displayValue = format(new Date(value as string), "MMMM d, yyyy");
                                            } catch (e) {
                                                displayValue = String(value);
                                            }
                                        }

                                        return (
                                            <div key={id}>
                                                <p className="text-xs text-muted-foreground">{def.name}</p>
                                                <p className="text-sm font-medium">{displayValue}</p>
                                            </div>
                                        );
                                    })}
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}
