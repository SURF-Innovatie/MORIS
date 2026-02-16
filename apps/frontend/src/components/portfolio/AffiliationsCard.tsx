import { Building2 } from "lucide-react";

import { OrganisationEffectiveMembershipResponse } from "@api/model";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface AffiliationsCardProps {
  memberships: OrganisationEffectiveMembershipResponse[] | undefined;
}

export const AffiliationsCard = ({ memberships }: AffiliationsCardProps) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Building2 className="h-5 w-5 text-primary" />
          Affiliations
        </CardTitle>
        <CardDescription>
          Your current organizational memberships and roles.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-3">
        {memberships?.length ? (
          memberships.map((membership) => (
            <div
              key={membership.membershipId}
              className="rounded-lg border px-4 py-3"
            >
              <p className="text-sm font-medium">
                {membership.scopeRootOrganisation?.name ||
                  "Organisation unavailable"}
              </p>
              <p className="text-xs text-muted-foreground">
                Role: {membership.roleKey || "Member"}
              </p>
            </div>
          ))
        ) : (
          <div className="rounded-lg border border-dashed p-6 text-center text-sm text-muted-foreground">
            No organisation memberships yet.
          </div>
        )}
      </CardContent>
    </Card>
  );
};
