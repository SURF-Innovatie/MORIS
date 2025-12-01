import { useState } from "react";
import {
  useGetAuthOrcidUrl,
  useGetProfile,
  usePostAuthOrcidUnlink,
} from "@api/moris";
import { useToast } from "@/hooks/use-toast";
import { useAuth } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogClose,
} from "@/components/ui/dialog";
import { Link, Unlink } from "lucide-react";

const ProfileRoute = () => {
  const { updateUser } = useAuth();
  const { toast } = useToast();
  const [isUnlinkDialogOpen, setIsUnlinkDialogOpen] = useState(false);

  const { data: user, isLoading, refetch: refetchProfile } = useGetProfile();

  const { refetch: getAuthURL, isFetching: isGettingURL } = useGetAuthOrcidUrl({
    query: {
      enabled: false,
    },
  });

  const { mutateAsync: unlinkORCID, isPending: isUnlinking } =
    usePostAuthOrcidUnlink();

  const handleLinkORCID = async () => {
    try {
      const result = await getAuthURL();
      if (result.data?.url) {
        window.location.href = result.data.url;
      }
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to get ORCID authorization URL",
        variant: "destructive",
      });
    }
  };

  const handleUnlinkORCID = async () => {
    try {
      await unlinkORCID();
      toast({
        title: "Success",
        description: "ORCID account unlinked successfully",
      });
      setIsUnlinkDialogOpen(false);
      // Refetch profile to update UI
      const { data: updatedUser } = await refetchProfile();
      if (updatedUser) {
        // Map UserResponse to UserAccount
        updateUser({
          user: {
            id: updatedUser.id,
            personID: updatedUser.person_id,
          },
          person: {
            id: updatedUser.person_id,
            email: updatedUser.email,
            name: updatedUser.name,
            givenName: updatedUser.givenName,
            familyName: updatedUser.familyName,
            orciD: updatedUser.orcid,
            userID: updatedUser.id,
          },
        });
      }
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to unlink ORCID account",
        variant: "destructive",
      });
    }
  };

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
        <Card>
          <CardHeader>
            <CardTitle>Personal Information</CardTitle>
            <CardDescription>Manage your personal details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex flex-col gap-1">
              <label className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                Full Name
              </label>
              <p className="font-medium">{user.name || "N/A"}</p>
            </div>
            <div className="flex flex-col gap-1">
              <label className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                Email
              </label>
              <p className="font-medium">{user.email}</p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Integrations</CardTitle>
            <CardDescription>
              Manage your external account connections
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-start justify-between gap-4 rounded-lg border p-4 bg-muted/30">
                <div className="space-y-1">
                  <div className="flex items-center gap-2">
                    <h3 className="font-semibold text-sm">ORCID iD</h3>
                    {user.orcid && (
                      <Badge variant="secondary" className="h-5 px-1.5 text-[10px] bg-green-500/10 text-green-600 hover:bg-green-500/20 border-green-500/20">
                        Verified
                      </Badge>
                    )}
                  </div>
                  <p className="text-xs text-muted-foreground">
                    {user.orcid ? (
                      <span className="font-mono">{user.orcid}</span>
                    ) : (
                      "Connect your ORCID iD to your account"
                    )}
                  </p>
                </div>
                {user.orcid ? (
                  <Dialog
                    open={isUnlinkDialogOpen}
                    onOpenChange={setIsUnlinkDialogOpen}
                  >
                    <DialogTrigger asChild>
                      <Button variant="ghost" size="icon" className="h-8 w-8 text-muted-foreground hover:text-destructive">
                        <Unlink className="h-4 w-4" />
                      </Button>
                    </DialogTrigger>
                    <DialogContent>
                      <DialogHeader>
                        <DialogTitle>Unlink ORCID?</DialogTitle>
                        <DialogDescription>
                          Are you sure you want to unlink your ORCID account? This
                          action cannot be undone easily.
                        </DialogDescription>
                      </DialogHeader>
                      <DialogFooter>
                        <DialogClose asChild>
                          <Button variant="outline">Cancel</Button>
                        </DialogClose>
                        <Button
                          variant="destructive"
                          onClick={handleUnlinkORCID}
                          disabled={isUnlinking}
                        >
                          {isUnlinking ? "Unlinking..." : "Unlink"}
                        </Button>
                      </DialogFooter>
                    </DialogContent>
                  </Dialog>
                ) : (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={handleLinkORCID}
                    disabled={isGettingURL}
                    className="h-8"
                  >
                    <Link className="mr-2 h-3.5 w-3.5" />
                    {isGettingURL ? "Connecting..." : "Connect"}
                  </Button>
                )}
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Right Column: Recent Activity (Placeholder) */}
      <div className="lg:col-span-2 space-y-8">
        <Card className="h-full border-dashed">
          <CardHeader>
            <CardTitle>Recent Activity</CardTitle>
            <CardDescription>
              Your recent publications and project updates.
            </CardDescription>
          </CardHeader>
          <CardContent className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground">
            <div className="h-12 w-12 rounded-full bg-muted/50 flex items-center justify-center mb-4">
              <Link className="h-6 w-6 opacity-20" />
            </div>
            <p className="font-medium">No recent activity</p>
            <p className="text-sm mt-1 max-w-xs mx-auto">
              Once you start working on projects or publishing research, your activity will appear here.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default ProfileRoute;
