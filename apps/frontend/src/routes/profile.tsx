import { useGetProfile } from "@api/moris";
import { ProfileInfo } from "@/components/profile/ProfileInfo";
import { OrcidConnection } from "@/components/profile/OrcidConnection";
import { ZenodoConnection } from "@/components/profile/ZenodoConnection";
import { ProfileActivity } from "@/components/profile/ProfileActivity";
import { SecuritySettings } from "@/components/profile/SecuritySettings";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { User, Link as LinkIcon, ShieldAlert } from "lucide-react";

const ProfileRoute = () => {
  const { data: user, isLoading, refetch: refetchProfile } = useGetProfile();

  if (isLoading && !user) {
    return <div>Loading...</div>;
  }

  if (!user) {
    return <div>User not found</div>;
  }

  return (
    <div className="space-y-6">
      <Tabs defaultValue="profile" className="space-y-6">
        <TabsList>
          <TabsTrigger value="profile" className="flex items-center gap-2">
            <User className="h-4 w-4" />
            Profile
          </TabsTrigger>
          <TabsTrigger value="connections" className="flex items-center gap-2">
            <LinkIcon className="h-4 w-4" />
            Connections
          </TabsTrigger>
          <TabsTrigger value="security" className="flex items-center gap-2">
            <ShieldAlert className="h-4 w-4" />
            Security
          </TabsTrigger>
        </TabsList>

        <TabsContent value="profile" className="space-y-6">
          <div className="grid gap-8 lg:grid-cols-3">
            {/* Left Column: Personal Info */}
            <div className="lg:col-span-1 space-y-8">
              <ProfileInfo user={user} refetchProfile={refetchProfile} />
            </div>

            {/* Right Column: Recent Activity */}
            <div className="lg:col-span-2 space-y-8">
              <ProfileActivity userId={user.id!} />
            </div>
          </div>
        </TabsContent>

        <TabsContent value="connections">
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
        </TabsContent>

        <TabsContent value="security">
          <SecuritySettings />
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default ProfileRoute;
