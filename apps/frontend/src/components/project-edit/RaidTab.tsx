import { useState } from "react";
import { Loader2, Plus, RefreshCw, ExternalLink } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useToast } from "@/hooks/use-toast";
import { ProjectEventType } from "@/api/events";
import { useAccess } from "@/context/AccessContext";
import {
  useGetProjectsIdRaid,
  usePostProjectsIdRaid,
  usePutProjectsIdRaid,
} from "@/api/generated-orval/moris";

interface RaidTabProps {
  projectId: string;
}

export function RaidTab({ projectId }: RaidTabProps) {
  const { toast } = useToast();
  const { hasAccess } = useAccess();
  const [isMinting, setIsMinting] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);

  const {
    data: raid,
    isLoading,
    error,
    refetch,
  } = useGetProjectsIdRaid(projectId, {
    query: {
      retry: false, // Don't retry on 404
    },
  });

  const { mutateAsync: mintRaid } = usePostProjectsIdRaid();
  const { mutateAsync: updateRaid } = usePutProjectsIdRaid();

  const handleMint = async () => {
    setIsMinting(true);
    try {
      await mintRaid({ id: projectId });
      toast({
        title: "RAiD Minted",
        description: "Successfully minted a new RAiD for this project.",
      });
      refetch();
    } catch (e: any) {
      toast({
        variant: "destructive",
        title: "Minting Failed",
        description:
          e.response?.data?.message || e.message || "Failed to mint RAiD",
      });
    } finally {
      setIsMinting(false);
    }
  };

  const handleUpdate = async () => {
    setIsUpdating(true);
    try {
      await updateRaid({ id: projectId });
      toast({
        title: "RAiD Updated",
        description: "Successfully updated RAiD metadata.",
      });
      refetch();
    } catch (e: any) {
      toast({
        variant: "destructive",
        title: "Update Failed",
        description:
          e.response?.data?.message || e.message || "Failed to update RAiD",
      });
    } finally {
      setIsUpdating(false);
    }
  };

  if (isLoading) {
    return (
      <div className="flex h-40 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  // Check for 404 (Not Found) -> Show Mint UI
  const isNotFound = error && (error as any).status === 404;

  if (isNotFound) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>RAiD Integration</CardTitle>
          <CardDescription>
            This project does not have a RAiD (Research Activity Identifier)
            yet.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground mb-4">
            Minting a RAiD will assign a persistent identifier to this project
            and register it with the RAiD service.
          </p>
        </CardContent>
        <CardFooter>
          <Button
            onClick={handleMint}
            disabled={isMinting || !hasAccess(ProjectEventType.RaidLinked)}
          >
            {isMinting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {!isMinting && <Plus className="mr-2 h-4 w-4" />}
            Mint RAiD
          </Button>
          {!hasAccess(ProjectEventType.RaidLinked) && (
            <p className="ml-4 text-xs text-red-500">
              You do not have permission to mint RAiDs.
            </p>
          )}
        </CardFooter>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="border-red-200">
        <CardHeader>
          <CardTitle className="text-red-500">Error</CardTitle>
        </CardHeader>
        <CardContent>
          <p>Failed to load RAiD information.</p>
        </CardContent>
      </Card>
    );
  }

  if (!raid) return null;

  return (
    <Card>
      <CardHeader>
        <CardTitle>RAiD Information</CardTitle>
        <CardDescription>
          Persistent identifier details for this project.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-2">
          <Label>RAiD Handle</Label>
          <div className="flex items-center gap-2">
            <Input value={raid.identifier?.id || ""} readOnly />
            <Button variant="outline" size="icon" asChild title="Open handle">
              <a
                href={
                  raid.identifier?.id
                    ? `https://handle.net/${raid.identifier.id}`
                    : "#"
                }
                target="_blank"
                rel="noreferrer"
              >
                <ExternalLink className="h-4 w-4" />
              </a>
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div className="grid gap-2">
            <Label>Owner</Label>
            <Input value={raid.identifier?.owner?.id || ""} readOnly />
          </div>
          <div className="grid gap-2">
            <Label>Registration Agency</Label>
            <Input
              value={raid.identifier?.registrationAgency?.id || ""}
              readOnly
            />
          </div>
        </div>

        {/* Display more metadata if needed, e.g. Title, Dates */}
        {raid.title && raid.title.length > 0 && (
          <div className="grid gap-2">
            <Label>Title (in RAiD)</Label>
            <Input value={raid.title[0].text} readOnly />
          </div>
        )}
      </CardContent>
      <CardFooter className="flex justify-between">
        <div className="text-xs text-muted-foreground">
          {/* Show sync status or similar if available */}
        </div>
        <div className="flex gap-4 items-center">
          {!hasAccess(ProjectEventType.RaidUpdated) && (
            <span className="text-xs text-red-500">Read only</span>
          )}
          <Button
            variant="outline"
            onClick={handleUpdate}
            disabled={isUpdating || !hasAccess(ProjectEventType.RaidUpdated)}
          >
            {isUpdating && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {!isUpdating && <RefreshCw className="mr-2 h-4 w-4" />}
            Update Metadata
          </Button>
        </div>
      </CardFooter>
    </Card>
  );
}
