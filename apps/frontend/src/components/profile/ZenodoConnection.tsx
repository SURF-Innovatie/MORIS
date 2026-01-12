import { useState } from "react";
import {
  useGetZenodoAuthUrl,
  useDeleteZenodoUnlink,
  useGetZenodoStatus,
} from "@api/moris";
import { UserResponse } from "@api/model";
import { useToast } from "@/hooks/use-toast";
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
import { Link, Unlink, ExternalLink } from "lucide-react";
import ZenodoIcon from "@/components/icons/zenodoIcon";

interface ZenodoConnectionProps {
  user: UserResponse;
  refetchProfile: () => Promise<any>;
}

export function ZenodoConnection({ refetchProfile }: ZenodoConnectionProps) {
  const { toast } = useToast();
  const [isUnlinkDialogOpen, setIsUnlinkDialogOpen] = useState(false);

  const { data: statusData, refetch: refetchStatus } = useGetZenodoStatus();

  const { refetch: getAuthURL, isFetching: isGettingURL } = useGetZenodoAuthUrl(
    {
      query: {
        enabled: false,
      },
    }
  );

  const { mutateAsync: unlinkZenodo, isPending: isUnlinking } =
    useDeleteZenodoUnlink();

  const handleLinkZenodo = async () => {
    try {
      const result = await getAuthURL();
      if (result.data?.auth_url) {
        window.location.href = result.data.auth_url;
      }
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to get Zenodo authorization URL",
        variant: "destructive",
      });
    }
  };

  const handleUnlinkZenodo = async () => {
    try {
      await unlinkZenodo();
      toast({
        title: "Success",
        description: "Zenodo account unlinked successfully",
      });
      setIsUnlinkDialogOpen(false);
      await refetchStatus();
      await refetchProfile();
    } catch (error) {
      toast({
        title: "Error",
        description: "Failed to unlink Zenodo account",
        variant: "destructive",
      });
    }
  };

  const isLinked = statusData?.linked ?? false;

  return (
    <div className="flex items-start gap-4 py-5">
      <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-muted/50 overflow-hidden mt-0.5">
        <ZenodoIcon width={24} height={24} />
      </div>
      <div className="flex-1 space-y-1">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <h3 className="font-medium text-sm">Zenodo</h3>
            {isLinked && (
              <Badge
                variant="secondary"
                className="h-5 px-1.5 text-[10px] bg-blue-500/10 text-blue-600 hover:bg-blue-500/20 border-blue-500/20"
              >
                Connected
              </Badge>
            )}
          </div>
          {isLinked ? (
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
                  <DialogTitle>Unlink Zenodo?</DialogTitle>
                  <DialogDescription>
                    Are you sure you want to unlink your Zenodo account? You
                    will no longer be able to upload products to Zenodo from
                    MORIS.
                  </DialogDescription>
                </DialogHeader>
                <DialogFooter>
                  <DialogClose asChild>
                    <Button variant="outline">Cancel</Button>
                  </DialogClose>
                  <Button
                    variant="destructive"
                    onClick={handleUnlinkZenodo}
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
              onClick={handleLinkZenodo}
              disabled={isGettingURL}
              className="h-8 px-3"
            >
              <Link className="mr-2 h-3.5 w-3.5" />
              {isGettingURL ? "Connecting..." : "Connect"}
            </Button>
          )}
        </div>
        <p className="text-xs text-muted-foreground leading-relaxed max-w-[400px]">
          {isLinked ? (
            <span className="flex items-center gap-1">
              Upload research products directly to Zenodo
              <a
                href="https://zenodo.org"
                target="_blank"
                rel="noopener noreferrer"
                className="text-primary hover:underline inline-flex items-center"
              >
                <ExternalLink className="h-3 w-3 ml-0.5" />
              </a>
            </span>
          ) : (
            "Connect your Zenodo account to upload research products"
          )}
        </p>
      </div>
    </div>
  );
}
