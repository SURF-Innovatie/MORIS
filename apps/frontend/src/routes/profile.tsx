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
import { SidebarLayout } from "@/components/layout";

const ProfileRoute = () => {
  const { data: user, isLoading, refetch: refetchProfile } = useGetProfile();

  if (isLoading && !user) {
    return <div>Loading...</div>;
  }

  if (!user) {
    return <div>User not found</div>;
  }

  const sidebarGroups = [
    {
      label: "User Settings",
      items: [
        { label: "Profile", href: "/dashboard/settings", icon: User },
        // These would ideally be separate routes or query params, but for now we can render all or keep tabs.
        // Following the GitHub pattern, these are usually separate items in the sidebar.
        // For simplicity in this refactor, I'll keep the Tabs but wrapped in SidebarLayout for consistency with the Request.
      ],
    },
  ];

  // Actually, to truly mirror GitHub, "Settings" pages usually have a side nav.
  // I will use SidebarLayout and render the current Tabs content inside it, but potentially
  // we should split these into sub-routes eventually. For now, layout consistency is key.

  return (
    <SidebarLayout sidebarGroups={sidebarGroups}>
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h1 className="text-2xl font-semibold tracking-tight">
            Public Profile
          </h1>
        </div>

        <Tabs defaultValue="profile" className="space-y-6">
          <TabsList>
            <TabsTrigger value="profile" className="flex items-center gap-2">
              <User className="h-4 w-4" />
              Profile
            </TabsTrigger>
            <TabsTrigger
              value="connections"
              className="flex items-center gap-2"
            >
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
                  <OrcidConnection
                    user={user}
                    refetchProfile={refetchProfile}
                  />
                  <ZenodoConnection
                    user={user}
                    refetchProfile={refetchProfile}
                  />
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="security">
            <SecuritySettings />
          </TabsContent>
        </Tabs>
      </div>
    </SidebarLayout>
  );
};

export default ProfileRoute;
