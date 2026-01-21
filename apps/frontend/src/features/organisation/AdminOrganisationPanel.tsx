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
import { Plus, Settings, Network } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { OrganisationNode } from "./components/OrganisationNode";
import { CreateChildDialog } from "./components/CreateChildDialog";
import { RorSearchSelect } from "@/components/organisation/RorSearchSelect";
import { EditOrganisationDialog } from "./components/EditOrganisationDialog";
import { OrganisationListLayout } from "./components/OrganisationListLayout";
import { OrganisationTreeView } from "./components/OrganisationTreeView";

export const AdminOrganisationPanel = () => {
  const { data: roots, isLoading } = useGetOrganisationNodesRoots();

  const renderActions = (node: OrganisationResponse) => {
    return (
      <>
        <Button
          variant="outline"
          size="sm"
          asChild
          className="h-7 text-xs px-2"
        >
          <Link to={`/dashboard/admin/organisations/${node.id}/roles`}>
            <Settings size={14} className="mr-1" /> Settings
          </Link>
        </Button>
        <CreateChildDialog
          parentId={node.id!}
          trigger={
            <Button variant="outline" size="sm" className="h-7 text-xs px-2">
              <Plus size={14} className="mr-1" /> New Unit
            </Button>
          }
        />
        <EditOrganisationDialog node={node} />
      </>
    );
  };

  return (
    <OrganisationListLayout
      title="Organisation Management"
      headerActions={
        <div className="flex gap-2">
          <ViewTreeDialog />
          <CreateRootDialog />
        </div>
      }
      isLoading={isLoading}
      isEmpty={roots?.length === 0}
    >
      {roots?.map((node: OrganisationResponse) => (
        <OrganisationNode
          key={node.id}
          node={node}
          renderActions={renderActions}
          defaultExpanded={true}
        />
      ))}
    </OrganisationListLayout>
  );
};

const ViewTreeDialog = () => {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">
          <Network size={16} className="mr-2" /> View Tree
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-[90vw] h-[90vh] flex flex-col">
        <DialogHeader>
          <DialogTitle>Organisation Structure</DialogTitle>
        </DialogHeader>
        <div className="flex-1 min-h-0">
          <OrganisationTreeView height={window.innerHeight * 0.8} />
        </div>
      </DialogContent>
    </Dialog>
  );
}; // End ViewTreeDialog

// Add imports for ViewTreeDialog dependencies
// The user imports are handled at the top, I will add them in a separate step or try to merge imports if multi replace works better.
// But for now I'm replacing the body of AdminOrganisationPanel + ViewTreeDialog definition.
// Wait, I should not define ViewTreeDialog inside AdminOrganisationPanel or duplicate it if I'm not careful.
// I'll place ViewTreeDialog at the bottom or imported. But since I'm editing the file, I can define it in the same file for simplicity as CreateRootDialog is there.

const CreateRootDialog = () => {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [rorId, setRorId] = useState<string | undefined>(undefined);
  const queryClient = useQueryClient();
  const { mutate: createRoot, isPending } = usePostOrganisationNodes({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: ["/organisation-nodes/roots"],
        });
        setOpen(false);
        setName("");
        setRorId(undefined);
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
          <RorSearchSelect
            value={rorId}
            onSelect={(id, item) => {
              setRorId(id);
              if (!name) setName(item.name || "");
            }}
          />
          <Button
            onClick={() => createRoot({ data: { name, rorId } })}
            disabled={isPending || !name}
          >
            {isPending ? "Creating..." : "Create"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};
