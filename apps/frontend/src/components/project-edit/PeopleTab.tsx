import { useState } from "react";
import { MoreHorizontal, Crown } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { AddPersonDialog } from "./AddPersonDialog";
import { EditRoleDialog } from "@/components/project-edit/EditRoleDialog";
import { ProjectMemberResponse } from "@/api/generated-orval/model";
import { useAccess } from "@/context/AccessContext";
import { Allowed } from "@/components/auth/Allowed";
import { ProjectEventType } from "@/api/events";

interface PeopleTabProps {
  projectId: string;
  members: ProjectMemberResponse[];
  onRefresh: () => void;
}

export function PeopleTab({ projectId, members, onRefresh }: PeopleTabProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle>Team Members</CardTitle>
          <CardDescription>
            Manage who has access to this project.
          </CardDescription>
        </div>
        <Allowed event={ProjectEventType.ProjectRoleAssigned}>
          <AddPersonDialog projectId={projectId} onPersonAdded={onRefresh} />
        </Allowed>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {members.map((member) => (
            <MemberRow
              key={member.id}
              member={member}
              allMembers={members}
              projectId={projectId}
              onRefresh={onRefresh}
            />
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

function MemberRow({
  member,
  allMembers,
  projectId,
  onRefresh,
}: {
  member: ProjectMemberResponse;
  allMembers: ProjectMemberResponse[];
  projectId: string;
  onRefresh: () => void;
}) {
  // Get all member entries for the same person (to see all their current roles)
  const allMembersForPerson = allMembers.filter((m) => m.id === member.id);
  const [editOpen, setEditOpen] = useState(false);
  const { hasAccess } = useAccess();
  const canEditRole = hasAccess(ProjectEventType.ProjectRoleAssigned);
  const canRemove = hasAccess(ProjectEventType.ProjectRoleUnassigned);
  const pending = (member as any).pending;

  if (pending || (!canEditRole && !canRemove)) {
    return (
      <div className="flex items-center justify-between rounded-lg border p-4 hover:bg-muted/50 transition-colors">
        <MemberInfo member={member} allMembersForPerson={allMembersForPerson} />
      </div>
    );
  }

  return (
    <>
      <div className="flex items-center justify-between rounded-lg border p-4 hover:bg-muted/50 transition-colors">
        <MemberInfo member={member} allMembersForPerson={allMembersForPerson} />
        <div className="flex items-center gap-2">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>Actions</DropdownMenuLabel>
              <Allowed event={ProjectEventType.ProjectRoleAssigned}>
                <DropdownMenuItem onClick={() => setEditOpen(true)}>
                  Edit Role
                </DropdownMenuItem>
              </Allowed>
              <DropdownMenuSeparator />
              <Allowed event={ProjectEventType.ProjectRoleUnassigned}>
                <DropdownMenuItem className="text-destructive">
                  Remove
                </DropdownMenuItem>
              </Allowed>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <EditRoleDialog
        open={editOpen}
        onOpenChange={setEditOpen}
        member={member}
        allMembersForPerson={allMembersForPerson}
        projectId={projectId}
        onSuccess={onRefresh}
      />
    </>
  );
}

function MemberInfo({
  member,
  allMembersForPerson,
}: {
  member: ProjectMemberResponse & { pending?: boolean };
  allMembersForPerson?: ProjectMemberResponse[];
}) {
  // Get all unique roles for this person
  const roles = allMembersForPerson
    ? allMembersForPerson
        .map((m) => ({ id: m.role_id, name: m.role_name || m.role }))
        .filter(
          (r, index, self) =>
            r.id && self.findIndex((s) => s.id === r.id) === index,
        )
    : member.role_id
      ? [{ id: member.role_id, name: member.role_name || member.role }]
      : [];

  const hasLeadRole = roles.some((r) => r.name?.toLowerCase() === "lead");

  return (
    <div
      className={
        member.pending
          ? "flex items-center gap-4 opacity-70"
          : "flex items-center gap-4"
      }
    >
      <Avatar className="h-10 w-10 border">
        <AvatarImage src={member.avatarUrl || ""} />
        <AvatarFallback className="font-semibold text-primary">
          {(member.name || "Unknown")
            .split(" ")
            .map((n) => n[0])
            .join("")
            .toUpperCase()
            .slice(0, 2)}
        </AvatarFallback>
      </Avatar>
      <div>
        <div className="flex items-center gap-2 flex-wrap">
          <p className="font-semibold leading-none">
            {member.name || "Unknown"}
          </p>
          {member.pending && (
            <Badge
              variant="outline"
              className="text-[10px] h-5 px-1.5 border-yellow-500 text-yellow-600 bg-yellow-50"
            >
              Pending
            </Badge>
          )}
          {hasLeadRole && (
            <Crown className="h-3.5 w-3.5 text-yellow-500 fill-yellow-500" />
          )}
          {roles.map((role) => (
            <Badge
              key={role.id}
              variant="secondary"
              className="text-[10px] h-5 px-1.5 font-normal capitalize"
            >
              {role.name}
            </Badge>
          ))}
        </div>
        <p className="text-sm text-muted-foreground mt-1">{member.email}</p>
      </div>
    </div>
  );
}
