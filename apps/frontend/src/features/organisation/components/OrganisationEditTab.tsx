import { useState, useEffect } from "react";
import {
  useGetOrganisationNodesId,
  usePatchOrganisationNodesId,
} from "@api/moris";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Loader2 } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { RorSearchSelect } from "@/components/organisation/RorSearchSelect";
import { toast } from "sonner";
import { getGetOrganisationNodesIdQueryKey } from "@api/moris";

interface OrganisationEditTabProps {
  nodeId: string;
}

export const OrganisationEditTab = ({ nodeId }: OrganisationEditTabProps) => {
  const { data: node, isLoading } = useGetOrganisationNodesId(nodeId);

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [avatarUrl, setAvatarUrl] = useState("");
  const [rorId, setRorId] = useState<string | undefined>(undefined);
  const queryClient = useQueryClient();

  useEffect(() => {
    if (node) {
      setName(node.name || "");
      setDescription(node.description || "");
      setAvatarUrl(node.avatarUrl || "");
      setRorId(node.rorId || undefined);
    }
  }, [node]);

  const { mutate: updateNode, isPending } = usePatchOrganisationNodesId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetOrganisationNodesIdQueryKey(nodeId),
        });
        queryClient.invalidateQueries({
          queryKey: ["/organisation-nodes/roots"],
        }); // Invalidate roots in case name changed
        toast.success("Organisation updated", {
          description: "Your changes have been saved successfully.",
        });
      },
      onError: (error: any) => {
        toast.error("Failed to update organisation", {
          description: error?.message || "An unknown error occurred",
        });
      },
    },
  });

  if (isLoading) {
    return (
      <div className="flex justify-center p-8">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!node) {
    return (
      <div className="p-4 text-center text-muted-foreground">
        Organisation not found
      </div>
    );
  }

  const handleSave = () => {
    updateNode({
      id: nodeId,
      data: {
        name,
        description: description || undefined, // Send undefined if empty to avoid clearing if backend treats empty string as "clear" or to be consistent
        avatarUrl: avatarUrl || undefined,
        rorId,
      },
    });
  };

  const hasChanges =
    name !== node.name ||
    (description || "") !== (node.description || "") ||
    (avatarUrl || "") !== (node.avatarUrl || "") ||
    rorId !== node.rorId;

  return (
    <div className="space-y-6 max-w-2xl py-4">
      <div className="space-y-4">
        <div className="grid w-full items-center gap-1.5">
          <Label htmlFor="name">Name</Label>
          <Input
            id="name"
            placeholder="Organisation Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>

        <div className="grid w-full items-center gap-1.5">
          <Label htmlFor="description">Description (Optional)</Label>
          <Textarea
            id="description"
            placeholder="Enter a brief description..."
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="min-h-[100px]"
          />
        </div>

        <div className="grid w-full items-center gap-1.5">
          <Label htmlFor="avatarUrl">Avatar URL (Optional)</Label>
          <Input
            id="avatarUrl"
            placeholder="https://example.com/logo.png"
            value={avatarUrl}
            onChange={(e) => setAvatarUrl(e.target.value)}
          />
          {avatarUrl && (
            <div className="mt-2">
              <p className="text-xs text-muted-foreground mb-1">Preview:</p>
              <img
                src={avatarUrl}
                alt="Avatar Preview"
                className="h-16 w-16 object-contain rounded-md border bg-muted/50"
                onError={(e) => {
                  (e.target as HTMLImageElement).src = ""; // Clear on error or show placeholder
                  // You might want to show a broken image icon or text here
                }}
              />
            </div>
          )}
        </div>

        <div className="grid w-full items-center gap-1.5">
          <Label>ROR</Label>
          <div className="flex flex-col gap-1">
            <RorSearchSelect
              value={rorId}
              onSelect={(id, item) => {
                setRorId(id);
                if (!name) setName(item.name || "");
              }}
            />
            <p className="text-[0.8rem] text-muted-foreground">
              Linking a Research Organization Registry (ROR) ID helps verify
              this organization.
            </p>
          </div>
        </div>

        <div className="flex justify-end pt-4">
          <Button
            onClick={handleSave}
            disabled={isPending || !name || !hasChanges}
          >
            {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Save Changes
          </Button>
        </div>
      </div>
    </div>
  );
};
