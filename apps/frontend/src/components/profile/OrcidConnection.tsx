import { useState } from "react";
import { useGetAuthOrcidUrl, usePostAuthOrcidUnlink } from "@api/moris";
import { UserResponse } from "@api/model";
import { useAuth } from "@/hooks/useAuth";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

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
import OrcidIcon from "../icons/orcidIcon";

interface OrcidConnectionProps {
  user: UserResponse;
  refetchProfile: () => Promise<any>;
}

export function OrcidConnection({
  user,
  refetchProfile,
}: OrcidConnectionProps) {
  const { updateUser } = useAuth();
  const [isUnlinkDialogOpen, setIsUnlinkDialogOpen] = useState(false);

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
      toast.error("Error", {
        description: "Failed to get ORCID authorization URL",
      });
    }
  };

  const handleUnlinkORCID = async () => {
    try {
      await unlinkORCID();
      toast.success("Success", {
        description: "ORCID account unlinked successfully",
      });
      setIsUnlinkDialogOpen(false);
      // Refetch profile to update UI
      const { data: updatedUser } = await refetchProfile();
      if (updatedUser) {
        // Map UserResponse to UserAccount
        updateUser(updatedUser);
      }
    } catch (error) {
      toast.error("Error", {
        description: "Failed to unlink ORCID account",
      });
    }
  };

  return (
    <div className="flex items-start gap-4 py-5">
      <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-muted/50 overflow-hidden mt-0.5">
        <OrcidIcon width={24} height={24} />
      </div>
      <div className="flex-1 space-y-1">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <h3 className="font-medium text-sm">ORCiD</h3>
            {user.orcid && (
              <Badge
                variant="secondary"
                className="h-5 px-1.5 text-[10px] bg-green-500/10 text-green-600 hover:bg-green-500/20 border-green-500/20"
              >
                Verified
              </Badge>
            )}
          </div>
          {user.orcid ? (
            <Dialog
              open={isUnlinkDialogOpen}
              onOpenChange={setIsUnlinkDialogOpen}
            >
              <DialogTrigger asChild>
                <Button
                  variant="ghost"
                  size="sm"
                  className="text-muted-foreground hover:text-destructive h-8 px-2"
                >
                  <Unlink className="mr-2 h-3.5 w-3.5" />
                  Disconnect
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
              className="h-8 px-3"
            >
              <Link className="mr-2 h-3.5 w-3.5" />
              {isGettingURL ? "Connecting..." : "Connect"}
            </Button>
          )}
        </div>
        <p className="text-xs text-muted-foreground leading-relaxed max-w-[400px]">
          {user.orcid ? (
            <span className="font-mono">{user.orcid}</span>
          ) : (
            "Connect your ORCID iD to your account"
          )}
        </p>
      </div>
    </div>
  );
}
