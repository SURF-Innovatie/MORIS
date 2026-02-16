import { useParams } from "react-router-dom";
import {
  useGetOrganisationNodesIdMembershipsEffective,
  useGetOrganisationNodesIdOrganisationRoles,
  useGetOrganisationNodesIdPermissionsMine,
} from "@api/moris";
import { OrganisationEffectiveMembershipResponse } from "@/api/generated-orval/model";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useUrlTabs } from "@/hooks/useUrlTabs";
import { ProjectRolesList } from "./components/ProjectRolesList";
import { CustomFieldDefinitionsList } from "./components/CustomFieldDefinitionsList";
import { RoleList } from "./components/RoleList";
import { OrganisationEditTab } from "./components/OrganisationEditTab";
import { EventPoliciesTab } from "./components/EventPoliciesTab";
import { EditMemberCustomFieldsButton } from "./components/EditMemberCustomFieldsButton";
import { AddMemberDialog } from "./components/AddMemberDialog";
import { RemoveMemberButton } from "./components/RemoveMemberButton";

export const RoleManagement = () => {
  const { nodeId } = useParams<{ nodeId: string }>();
  const { data: members, isLoading: isLoadingMembers } =
    useGetOrganisationNodesIdMembershipsEffective(nodeId!);
  const { data: roles } = useGetOrganisationNodesIdOrganisationRoles(nodeId!);
  const { data: myPermissions, isLoading: isLoadingPerms } =
    useGetOrganisationNodesIdPermissionsMine(nodeId!);

  const canManageMembers = myPermissions?.includes("manage_members") ?? false;
  const canManageCustomFields =
    myPermissions?.includes("manage_custom_fields") ?? false;

  const [currentTab, setTab] = useUrlTabs("members");

  if (isLoadingMembers || isLoadingPerms) return <div>Loading...</div>;

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Organisation Settings</h1>
      </div>

      <Tabs
        value={currentTab}
        onValueChange={setTab}
        className="w-full"
      >
        <TabsList className="mb-4">
          <TabsTrigger value="members">Members & Permissions</TabsTrigger>
          <TabsTrigger value="roles">Roles</TabsTrigger>
          <TabsTrigger value="project-roles">Project Roles</TabsTrigger>
          <TabsTrigger value="custom-fields">Custom Fields</TabsTrigger>
          <TabsTrigger value="event-policies">Event Policies</TabsTrigger>
          <TabsTrigger value="edit">Edit</TabsTrigger>
        </TabsList>

        <TabsContent value="members">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">Members</h2>
            {nodeId && (
              <AddMemberDialog
                nodeId={nodeId}
                roles={roles || []}
                members={members || []}
                disabled={!canManageMembers}
              />
            )}
          </div>

          <div className="border rounded-lg overflow-hidden">
            <table className="w-full text-sm text-left">
              <thead className="bg-gray-50 text-gray-700 font-medium">
                <tr>
                  <th className="px-4 py-3">User</th>
                  <th className="px-4 py-3">Role</th>
                  <th className="px-4 py-3">Owning Organisation</th>
                  <th className="px-4 py-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {members?.map((m: OrganisationEffectiveMembershipResponse) => (
                  <tr key={m.membershipId} className="bg-white">
                    <td className="px-4 py-3">{m.person?.name}</td>
                    <td className="px-4 py-3 capitalize">{m.roleKey}</td>
                    <td className="px-4 py-3 capitalize">
                      {m.scopeRootOrganisation?.name}
                    </td>
                    <td className="px-4 py-3 text-right flex gap-2 justify-end">
                      <EditMemberCustomFieldsButton
                        nodeId={nodeId!}
                        membership={m}
                        canEdit={canManageCustomFields}
                      />
                      {canManageMembers && (
                        <RemoveMemberButton
                          membershipId={m.membershipId!}
                          nodeId={nodeId!}
                        />
                      )}
                    </td>
                  </tr>
                ))}
                {members?.length === 0 && (
                  <tr>
                    <td
                      colSpan={4}
                      className="px-4 py-8 text-center text-gray-500"
                    >
                      No members found.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </TabsContent>

        <TabsContent value="roles">
          {nodeId ? <RoleList nodeId={nodeId} /> : <div>Invalid Node ID</div>}
        </TabsContent>

        <TabsContent value="project-roles">
          {nodeId ? (
            <ProjectRolesList nodeId={nodeId} />
          ) : (
            <div>Invalid Node ID</div>
          )}
        </TabsContent>

        <TabsContent value="custom-fields">
          {nodeId ? (
            <CustomFieldDefinitionsList nodeId={nodeId} />
          ) : (
            <div>Invalid Node ID</div>
          )}
        </TabsContent>

        <TabsContent value="edit">
          {nodeId ? (
            <OrganisationEditTab nodeId={nodeId} />
          ) : (
            <div>Invalid Node ID</div>
          )}
        </TabsContent>

        <TabsContent value="event-policies">
          {nodeId ? (
            <EventPoliciesTab nodeId={nodeId} />
          ) : (
            <div>Invalid Node ID</div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
};
