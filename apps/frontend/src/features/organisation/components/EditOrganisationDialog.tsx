import { useState, useEffect } from "react";
import { usePatchOrganisationNodesId } from "@api/moris";
import { OrganisationResponse } from "@api/model";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Pencil } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { RorSearchSelect } from "@/components/organisation/RorSearchSelect";
import { Label } from "@/components/ui/label";

interface EditOrganisationDialogProps {
  node: OrganisationResponse;
  trigger?: React.ReactNode;
}

export const EditOrganisationDialog = ({
  node,
  trigger,
}: EditOrganisationDialogProps) => {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState(node.name);
  const [rorId, setRorId] = useState<string | undefined>(node.rorId);

  // Sync state when node changes (e.g. if parent list updates)
  useEffect(() => {
    if (open) {
      setName(node.name);
      setRorId(node.rorId);
    }
  }, [open, node]);

  const queryClient = useQueryClient();
  const { mutate: updateNode, isPending } = usePatchOrganisationNodesId({
    mutation: {
      onSuccess: () => {
        // Invalidate root list and specific node
        // Ideally we invalidate the parent's children list too, but we might not know parent ID easily if it's not passed.
        // Invalidate everything organisation related for safety/simplicity or use smarter cache updates.
        queryClient.invalidateQueries({ queryKey: ["/organisation-nodes"] });
        setOpen(false);
      },
    },
  });

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger || (
          <Button variant="ghost" size="sm">
            <Pencil size={14} />
          </Button>
        )}
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Organisation</DialogTitle>
        </DialogHeader>
        <div className="space-y-4 pt-4">
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
            <Label>ROR</Label>
            <RorSearchSelect
              value={rorId}
              onSelect={(id, item) => {
                setRorId(id);
                if (!name) setName(item.name);
              }}
            />
          </div>

          <div className="flex justify-end gap-2">
            <Button variant="outline" onClick={() => setOpen(false)}>
              Cancel
            </Button>
            <Button
              onClick={() =>
                updateNode({ id: node.id!, data: { name, rorId } })
              }
              disabled={isPending || !name}
            >
              {isPending ? "Saving..." : "Save Changes"}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
};
