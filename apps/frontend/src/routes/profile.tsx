import { useState } from "react";
import {
  useGetAuthOrcidUrl,
  useGetProfile,
  usePostAuthOrcidUnlink,
  usePutPeopleId,
  useGetUsersIdEventsApproved,
} from "@api/moris";
import { format } from "date-fns";
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
import { Link, Unlink, Pencil } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";

const ProfileRoute = () => {
  const { updateUser } = useAuth();
  const { toast } = useToast();
  const [isUnlinkDialogOpen, setIsUnlinkDialogOpen] = useState(false);
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);

  const { data: user, isLoading, refetch: refetchProfile } = useGetProfile();

  const { refetch: getAuthURL, isFetching: isGettingURL } = useGetAuthOrcidUrl({
    query: {
      enabled: false,
    },
  });

  const { mutateAsync: unlinkORCID, isPending: isUnlinking } =
    usePostAuthOrcidUnlink();

  const { mutateAsync: updatePerson, isPending: isUpdating } = usePutPeopleId();

  const { data: eventsData, isLoading: isLoadingEvents } =
    useGetUsersIdEventsApproved(user?.id || "", {
      query: {
        enabled: !!user?.id,
      },
    });

  const [editForm, setEditForm] = useState({
    name: "",
    avatarUrl: "",
    description: "",
  });

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

  const handleEditOpen = () => {
    if (user) {
      setEditForm({
        name: user.name || "",
        avatarUrl: user.avatarUrl || "",
        description: user.description || "",
      });
      setIsEditDialogOpen(true);
    }
  };

  const handleEditSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!user) return;

    try {
      await updatePerson({
        id: user.person_id!,
        data: {
          name: editForm.name,
          email: user.email, // Email cannot be changed
          avatar_url: editForm.avatarUrl || undefined,
          description: editForm.description || undefined,
          user_id: user.id,
        },
      });

      toast({
        title: "Success",
        description: "Profile updated successfully",
      });
      setIsEditDialogOpen(false);
      refetchProfile();
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to update profile",
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
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <div className="space-y-1">
              <CardTitle>Personal Information</CardTitle>
              <CardDescription>Manage your personal details</CardDescription>
            </div>
            <Button variant="ghost" size="icon" onClick={handleEditOpen}>
              <Pencil className="h-4 w-4" />
            </Button>
          </CardHeader>
          <CardContent className="space-y-6 pt-4">
            <div className="flex flex-col items-center gap-4">
              <Avatar className="h-24 w-24">
                <AvatarImage src={user.avatarUrl || ""} alt={user.name} />
                <AvatarFallback className="text-2xl">
                  {user.name
                    ?.split(" ")
                    .map((n) => n[0])
                    .join("")
                    .toUpperCase()
                    .slice(0, 2)}
                </AvatarFallback>
              </Avatar>
            </div>

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
            {user.description && (
              <div className="flex flex-col gap-1">
                <label className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                  Description
                </label>
                <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                  {user.description}
                </p>
              </div>
            )}
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
                      <Badge
                        variant="secondary"
                        className="h-5 px-1.5 text-[10px] bg-green-500/10 text-green-600 hover:bg-green-500/20 border-green-500/20"
                      >
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
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 text-muted-foreground hover:text-destructive"
                      >
                        <Unlink className="h-4 w-4" />
                      </Button>
                    </DialogTrigger>
                    <DialogContent>
                      <DialogHeader>
                        <DialogTitle>Unlink ORCID?</DialogTitle>
                        <DialogDescription>
                          Are you sure you want to unlink your ORCID account?
                          This action cannot be undone easily.
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
          <CardContent className="flex flex-col">
            {isLoadingEvents ? (
              <div className="py-12 text-center text-muted-foreground">
                Loading activity...
              </div>
            ) : eventsData?.events?.length ? (
              <div className="space-y-6">
                {eventsData.events.map((event) => (
                  <div
                    key={event.id}
                    className="flex flex-col gap-1 border-b pb-4 last:border-0 last:pb-0"
                  >
                    <div className="flex items-center justify-between">
                      <Badge variant="outline" className="font-normal text-xs">
                        {event.type}
                      </Badge>
                      <span className="text-xs text-muted-foreground">
                        {event.at ? format(new Date(event.at), "PPP") : "N/A"}
                      </span>
                    </div>
                    <p className="text-sm mt-1">{event.details}</p>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground">
                <div className="h-12 w-12 rounded-full bg-muted/50 flex items-center justify-center mb-4">
                  <Link className="h-6 w-6 opacity-20" />
                </div>
                <p className="font-medium">No recent activity</p>
                <p className="text-sm mt-1 max-w-xs mx-auto">
                  Once you start working on projects or publishing research,
                  your activity will appear here.
                </p>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit Profile</DialogTitle>
            <DialogDescription>
              Update your personal information.
            </DialogDescription>
          </DialogHeader>
          <form onSubmit={handleEditSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">Full Name</Label>
              <Input
                id="name"
                value={editForm.name}
                onChange={(e) =>
                  setEditForm({ ...editForm, name: e.target.value })
                }
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input id="email" value={user.email} disabled />
              <p className="text-[0.8rem] text-muted-foreground">
                Email cannot be changed at this time.
              </p>
            </div>
            <div className="space-y-2">
              <Label htmlFor="avatarUrl">Avatar URL</Label>
              <Input
                id="avatarUrl"
                value={editForm.avatarUrl}
                onChange={(e) =>
                  setEditForm({ ...editForm, avatarUrl: e.target.value })
                }
                placeholder="https://example.com/avatar.jpg"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                value={editForm.description}
                onChange={(e) =>
                  setEditForm({ ...editForm, description: e.target.value })
                }
                placeholder="Tell us about yourself..."
              />
            </div>
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setIsEditDialogOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isUpdating}>
                {isUpdating ? "Saving..." : "Save Changes"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export default ProfileRoute;
