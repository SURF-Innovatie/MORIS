import { useState } from "react";
import { usePostOrganisationNodesIdChildren } from "@api/moris";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Plus } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { RorSearchSelect } from "@/components/organisation/RorSearchSelect";

interface CreateChildDialogProps {
  parentId: string;
  trigger?: React.ReactNode;
  onSuccess?: () => void;
}

export const CreateChildDialog = ({
  parentId,
  trigger,
  onSuccess,
}: CreateChildDialogProps) => {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [rorId, setRorId] = useState<string | undefined>(undefined);
  const queryClient = useQueryClient();
  const { mutate: createChild, isPending } = usePostOrganisationNodesIdChildren(
    {
      mutation: {
        onSuccess: () => {
          // Invalidate children of the parent
          queryClient.invalidateQueries({
            queryKey: [`/organisation-nodes/${parentId}/children`],
          });
          queryClient.invalidateQueries({
            queryKey: ["/organisation-memberships/mine"],
          });

          if (onSuccess) onSuccess();
          setOpen(false);
          setName("");
          setRorId(undefined);
        },
      },
    }
  );

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger || (
          <Button variant="ghost" size="sm">
            <Plus size={14} />
          </Button>
        )}
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Child Unit</DialogTitle>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <Input
            placeholder="Unit Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <RorSearchSelect
            value={rorId}
            onSelect={(id, item) => {
              setRorId(id);
              if (!name) setName(item.name || ""); // Auto-fill name if empty
            }}
          />
          <Button
            onClick={() => createChild({ id: parentId, data: { name, rorId } })}
            disabled={isPending || !name}
          >
            {isPending ? "Creating..." : "Create"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};
