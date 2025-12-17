import { useState } from "react";
import {
  useGetOrganisationNodesRoots,
  usePostOrganisationNodes,
} from "@api/moris";
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
import { Link } from "react-router-dom";
import { Plus, Settings } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { OrganisationNode } from "./components/OrganisationNode";
import { CreateChildDialog } from "./components/CreateChildDialog";

export const AdminOrganisationPanel = () => {
  const { data: roots, isLoading } = useGetOrganisationNodesRoots();

  if (isLoading) return <div>Loading...</div>;

  const renderActions = (node: OrganisationResponse) => {
    return (
      <>
        <CreateChildDialog parentId={node.id!} />
        <Link to={`/dashboard/admin/organisations/${node.id}/roles`}>
          <Button variant="ghost" size="sm">
            <Settings size={14} />
          </Button>
        </Link>
      </>
    );
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Organisation Management</h1>
        <CreateRootDialog />
      </div>
      <div className="space-y-2">
        {roots?.map((node: OrganisationResponse) => (
          <OrganisationNode
            key={node.id}
            node={node}
            renderActions={renderActions}
          />
        ))}
      </div>
    </div>
  );
};

const CreateRootDialog = () => {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const queryClient = useQueryClient();
  const { mutate: createRoot, isPending } = usePostOrganisationNodes({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: ["/organisation-nodes/roots"],
        });
        setOpen(false);
        setName("");
      },
    },
  });

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button>
          <Plus size={16} className="mr-2" /> New Root Organisation
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Create Root Organisation</DialogTitle>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <Input
            placeholder="Organisation Name"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
          <Button
            onClick={() => createRoot({ data: { name } })}
            disabled={isPending || !name}
          >
            {isPending ? "Creating..." : "Create"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};
