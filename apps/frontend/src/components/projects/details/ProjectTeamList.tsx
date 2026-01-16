import { useState } from "react";
import { Users, MoreVertical, Pencil } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { ProjectMemberResponse, ProjectResponse } from "@api/model";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { EditRoleDialog } from "../../project-edit/EditRoleDialog";

interface ProjectTeamListProps {
  project: ProjectResponse;
  onRefresh?: () => void;
}

export function ProjectTeamList({ project, onRefresh }: ProjectTeamListProps) {
  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2">
        <Users className="h-5 w-5 text-muted-foreground" />
        <h2 className="text-lg font-semibold">Team Members</h2>
        <Badge variant="secondary" className="ml-2">
          {project.members?.length || 0}
        </Badge>
      </div>
      <Card>
        <CardContent className="p-0">
          {project.members && project.members.length > 0 ? (
            <div className="divide-y">
              {project.members.map((member) => (
                <MemberRow
                  key={member.id}
                  member={member}
                  projectId={project.id!}
                  onRefresh={onRefresh}
                />
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

function MemberRow({
  member,
  projectId,
  onRefresh,
}: {
  member: ProjectMemberResponse;
  projectId: string;
  onRefresh?: () => void;
}) {
  const [editOpen, setEditOpen] = useState(false);

  return (
    <>
      <div className="flex items-center justify-between p-4">
        <div className="flex items-center gap-3">
          <Avatar className="h-9 w-9">
            <AvatarImage src={member.avatarUrl || ""} alt={member.email} />
            <AvatarFallback className="text-xs font-medium uppercase">
              {member.email?.charAt(0) || "?"}
            </AvatarFallback>
          </Avatar>
          <div className="space-y-0.5">
            <p className="text-sm font-medium leading-none">{member.name}</p>
            <p className="text-xs text-muted-foreground">{member.email}</p>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <Badge variant="outline" className="capitalize">
            {member.role_name || member.role}
          </Badge>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem onClick={() => setEditOpen(true)}>
                <Pencil className="mr-2 h-4 w-4" />
                Edit Role
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <EditRoleDialog
        open={editOpen}
        onOpenChange={setEditOpen}
        member={member}
        projectId={projectId}
        onSuccess={onRefresh}
      />
    </>
  );
}
