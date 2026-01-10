import { useGetProfile } from "@api/moris";
import { ProfileInfo } from "@/components/profile/ProfileInfo";
import { OrcidConnection } from "@/components/profile/OrcidConnection";
import { ZenodoConnection } from "@/components/profile/ZenodoConnection";
import { ProfileActivity } from "@/components/profile/ProfileActivity";

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
        <OrcidConnection user={user} refetchProfile={refetchProfile} />
        <ZenodoConnection user={user} refetchProfile={refetchProfile} />
      </div>

      {/* Right Column: Recent Activity */}
      <div className="lg:col-span-2 space-y-8">
        <ProfileActivity userId={user.id!} />
      </div>
    </div>
  );
};

export default ProfileRoute;
