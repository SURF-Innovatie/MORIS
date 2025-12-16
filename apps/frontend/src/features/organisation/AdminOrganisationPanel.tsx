import { useState } from "react";
import {
    useGetOrganisationNodesRoots,
    usePostOrganisationNodes,
    useGetOrganisationNodesIdChildren,
    usePostOrganisationNodesIdChildren,
} from "@api/moris";
import { OrganisationResponse } from "@api/model";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Link } from "react-router-dom";
import { ChevronRight, ChevronDown, Plus, Settings } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";

export const AdminOrganisationPanel = () => {
    const { data: roots, isLoading } = useGetOrganisationNodesRoots();

    if (isLoading) return <div>Loading...</div>;

    return (
        <div className="p-6">
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-bold">Organisation Management</h1>
                <CreateRootDialog />
            </div>
            <div className="space-y-2">
                {roots?.map((node) => (
                    <OrganisationNodeItem key={node.id} node={node} />
                ))}
            </div>
        </div>
    );
};

const OrganisationNodeItem = ({ node }: { node: OrganisationResponse }) => {
    const [isExpanded, setIsExpanded] = useState(false);
    const { data: children } = useGetOrganisationNodesIdChildren(node.id!, { query: { enabled: isExpanded } });

    return (
        <div className="ml-4 border-l pl-4">
            <div className="flex items-center gap-2 py-2">
                <button onClick={() => setIsExpanded(!isExpanded)} className="p-1 hover:bg-gray-100 rounded">
                    {isExpanded ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
                </button>
                <span className="font-medium">{node.name}</span>
                <div className="ml-auto flex gap-2">
                    <CreateChildDialog parentId={node.id!} />
                    <Link to={`/dashboard/admin/organisations/${node.id}/roles`}>
                        <Button variant="ghost" size="sm"><Settings size={14} /></Button>
                    </Link>
                </div>
            </div>
            {isExpanded && (
                <div className="ml-4">
                    {children?.map((child) => (
                        <OrganisationNodeItem key={child.id} node={child} />
                    ))}
                    {children?.length === 0 && <div className="text-gray-500 italic text-sm">No children</div>}
                </div>
            )}
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
                queryClient.invalidateQueries({ queryKey: ['/organisation-nodes/roots'] });
                setOpen(false);
                setName("");
            }
        }
    });

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button><Plus size={16} className="mr-2" /> New Root Organisation</Button>
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
                    <Button onClick={() => createRoot({ data: { name } })} disabled={isPending || !name}>
                        {isPending ? "Creating..." : "Create"}
                    </Button>
                </div>
            </DialogContent>
        </Dialog>
    );
};

const CreateChildDialog = ({ parentId }: { parentId: string }) => {
    const [open, setOpen] = useState(false);
    const [name, setName] = useState("");
    const queryClient = useQueryClient();
    const { mutate: createChild, isPending } = usePostOrganisationNodesIdChildren({
        mutation: {
            onSuccess: () => {
                queryClient.invalidateQueries({
                    predicate: (query) =>
                        query.queryKey.includes(parentId) &&
                        query.queryKey.includes('children')
                });
                setOpen(false);
                setName("");
            }
        }
    });

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button variant="ghost" size="sm"><Plus size={14} /></Button>
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
                    <Button
                        onClick={() => createChild({ id: parentId, data: { name } })}
                        disabled={isPending || !name}
                    >
                        {isPending ? "Creating..." : "Create"}
                    </Button>
                </div>
            </DialogContent>
        </Dialog>
    );
};
