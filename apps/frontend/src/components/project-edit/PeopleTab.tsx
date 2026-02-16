import { useState } from "react";
import { Crown, Edit, Trash } from "lucide-react";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { AddPersonDialog } from "./AddPersonDialog";
import { EditRoleDialog } from "@/components/project-edit/EditRoleDialog";
import { ProjectMemberResponse } from "@/api/generated-orval/model";
import { useAccess } from "@/contexts/AccessContext";
import { Allowed } from "@/components/auth/Allowed";
import { ProjectEventType } from "@/api/events";
import { ListItem, ActionMenu } from "@/components/composition";
import type { ActionMenuItem } from "@/components/composition";

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

  // Get all unique roles for this person
  const roles = allMembersForPerson
    .map((m) => ({ id: m.role_id, name: m.role_name || m.role }))
    .filter(
      (r, index, self) =>
        r.id && self.findIndex((s) => s.id === r.id) === index,
    );

  const hasLeadRole = roles.some((r) => r.name?.toLowerCase() === "lead");

  // Build badges array
  const badges = roles.map((role) => ({
    label: role.name || "",
    variant: "secondary" as const,
  }));

  // Add crown icon for lead role (prepend to badges display)
  const subtitle = (
    <div className="flex items-center gap-1.5">
      {hasLeadRole && (
        <Crown className="h-3.5 w-3.5 text-yellow-500 fill-yellow-500" />
      )}
      <span>{member.email}</span>
    </div>
  );

  // Build action menu items
  const menuItems: ActionMenuItem[] = [];
  if (canEditRole) {
    menuItems.push({
      label: "Edit Role",
      icon: Edit,
      onClick: () => setEditOpen(true),
    });
  }
  if (canRemove) {
    menuItems.push({
      label: "Remove",
      icon: Trash,
      onClick: () => {
        // TODO: Implement remove
        console.log("Remove member", member.id);
      },
      destructive: true,
    });
  }

  // Show action menu only if user has permissions and member is not pending
  const showActions = !pending && menuItems.length > 0;

  const avatarFallback = (member.name || "Unknown")
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);

  return (
    <>
      <ListItem
        title={member.name || "Unknown"}
        subtitle={subtitle as any}
        avatarUrl={member.avatarUrl || ""}
        avatarFallback={avatarFallback}
        badges={badges}
        pending={pending}
        action={
          showActions ? (
            <ActionMenu
              sections={[
                {
                  label: "Actions",
                  items: menuItems,
                },
              ]}
            />
          ) : undefined
        }
      />

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
