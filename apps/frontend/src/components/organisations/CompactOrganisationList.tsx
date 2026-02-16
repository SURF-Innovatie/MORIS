import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { useGetOrganisationMembershipsMine } from "@api/moris";
import { Skeleton } from "@/components/ui/skeleton";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";

export const CompactOrganisationList = () => {
  const navigate = useNavigate();
  const { data: memberships, isLoading } = useGetOrganisationMembershipsMine();

  if (isLoading) {
    return (
      <div className="space-y-2 px-2">
        <Skeleton className="h-8 w-8 rounded-md" />
        <Skeleton className="h-8 w-8 rounded-md" />
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center justify-between px-2">
        <h4 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
          Your Teams
        </h4>
      </div>

      <div className="space-y-1">
        {memberships && memberships.length === 0 ? (
          <p className="px-2 text-xs text-muted-foreground italic">
            No organisations.
          </p>
        ) : (
          memberships?.map((membership) => {
            const org = membership.scopeRootOrganisation;
            if (!org) return null;

            return (
              <Button
                key={membership.membershipId}
                variant="ghost"
                className="w-full justify-start gap-2 px-2 py-1.5 h-auto text-sm font-normal text-muted-foreground hover:text-foreground"
                onClick={() =>
                  navigate(`/dashboard/organisations/${org.id}/members`)
                }
              >
                <Avatar className="h-5 w-5 rounded-md">
                  <AvatarFallback className="rounded-md text-[9px] bg-primary/10 text-primary">
                    {org.name?.substring(0, 2).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                <span className="truncate">{org.name}</span>
              </Button>
            );
          })
        )}
      </div>
    </div>
  );
};
