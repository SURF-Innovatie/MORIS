import { Users } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { ProjectResponse } from "@api/model";

interface ProjectTeamListProps {
    project: ProjectResponse;
}

export function ProjectTeamList({ project }: ProjectTeamListProps) {
    return (
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
                                                {person.name}
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
    );
}
