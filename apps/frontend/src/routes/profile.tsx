import { useGetProfile } from "@api/moris";
import { ProfileInfo } from "@/components/profile/ProfileInfo";
import { OrcidConnection } from "@/components/profile/OrcidConnection";
import { ZenodoConnection } from "@/components/profile/ZenodoConnection";
import { ProfileActivity } from "@/components/profile/ProfileActivity";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

const ProfileRoute = () => {
  const { data: user, isLoading, refetch: refetchProfile } = useGetProfile();

  if (isLoading && !user) {
    return <div>Loading...</div>;
  }

  if (!user) {
    return <div>User not found</div>;
  }

  return (
    <div className="grid gap-8 lg:grid-cols-3">
      {/* Left Column: Personal Info & Integrations */}
      <div className="lg:col-span-1 space-y-8">
        <ProfileInfo user={user} refetchProfile={refetchProfile} />
        <Card>
          <CardHeader>
            <CardTitle>Integrations</CardTitle>
            <CardDescription>
              Manage your external account connections
            </CardDescription>
          </CardHeader>
          <CardContent className="px-6 py-0">
            <div className="divide-y">
              <OrcidConnection user={user} refetchProfile={refetchProfile} />
              <ZenodoConnection user={user} refetchProfile={refetchProfile} />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Right Column: Recent Activity */}
      <div className="lg:col-span-2 space-y-8">
        <ProfileActivity userId={user.id!} />
      </div>
    </div>
  );
};

export default ProfileRoute;
