import { useGetOrganisationMembershipsMine } from "@api/moris";
import { Button } from "@/components/ui/button";
import { Plus, Settings } from "lucide-react";
import { Link } from "react-router-dom";
import { OrganisationNode } from "./components/OrganisationNode";
import { CreateChildDialog } from "./components/CreateChildDialog";
import { OrganisationEffectiveMembershipResponse } from "@/api/generated-orval/model";
import { OrganisationListLayout } from "./components/OrganisationListLayout";

export const UserOrganisationManagement = () => {
  const { data: memberships, isLoading } = useGetOrganisationMembershipsMine();

  return (
    <OrganisationListLayout
      title="My Organizations"
      isLoading={isLoading}
      isEmpty={memberships?.length === 0}
      emptyMessage="You are not a member of any organization."
    >
      {memberships?.map(
        (membership: OrganisationEffectiveMembershipResponse) => (
          <UserOrganisationNode
            key={membership.membershipId}
            membership={membership}
          />
        )
      )}
    </OrganisationListLayout>
  );
};

const UserOrganisationNode = ({
  membership,
}: {
  membership: OrganisationEffectiveMembershipResponse;
}) => {
  const rootNode = {
    id: membership.scopeRootOrganisation?.id,
    name: membership.scopeRootOrganisation?.name,
  };

  const renderActions = (node: any) => {
    const isRoot = node.id === membership.scopeRootOrganisation?.id;
    const canManage = membership.hasAdminRights;

    return (
      <>
        {isRoot && (
          <span className="text-sm text-gray-500 mr-2 self-center hidden md:inline">
            ({membership.roleKey})
          </span>
        )}
        {canManage && (
          <>
            <Button
              variant="outline"
              size="sm"
              asChild
              className="h-7 text-xs px-2"
            >
              <Link to={`/dashboard/organisations/${node.id}/members`}>
                <Settings size={14} className="mr-1" /> Settings
              </Link>
            </Button>
            <CreateChildDialog
              parentId={node.id!}
              trigger={
                <Button
                  variant="outline"
                  size="sm"
                  className="h-7 text-xs px-2"
                >
                  <Plus size={14} className="mr-1" /> New Unit
                </Button>
              }
            />
          </>
        )}
      </>
    );
  };

  return (
    <OrganisationNode
      node={rootNode}
      renderActions={renderActions}
      defaultExpanded={true}
    />
  );
};
