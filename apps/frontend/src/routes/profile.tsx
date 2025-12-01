import { useState } from "react";
import {
  useGetAuthOrcidUrl,
  useGetProfile,
  usePostAuthOrcidUnlink,
} from "@api/moris";
import { useToast } from "@/hooks/use-toast";
import { useAuth } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";
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
    <div className="container mx-auto max-w-2xl py-8">
      <h1 className="text-3xl font-bold mb-8">Profile</h1>

      <Card className="mb-8">
        <CardHeader>
          <CardTitle>Personal Information</CardTitle>
          <CardDescription>Manage your personal details</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="text-sm font-medium text-muted-foreground">
                Email
              </label>
              <p className="text-lg">{user.email}</p>
            </div>
            {/* Add Name if available in user object, currently AuthenticatedUser only has Email and ID */}
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
          <div className="flex items-center justify-between">
            <div>
              <h3 className="font-medium">ORCID</h3>
              <p className="text-sm text-muted-foreground">
                {user.orcid ? (
                  <>
                    Linked: <span className="font-mono">{user.orcid}</span>
                  </>
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
                  <Button variant="outline" size="sm">
                    <Unlink className="mr-2 h-4 w-4" /> Unlink ORCID
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
              >
                <Link className="mr-2 h-4 w-4" />
                {isGettingURL ? "Connecting..." : "Link ORCID"}
              </Button>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default ProfileRoute;
