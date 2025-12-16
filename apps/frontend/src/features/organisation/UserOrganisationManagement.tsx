import { useState } from "react";
import { useGetOrganisationMembershipsMine, usePostOrganisationNodesIdChildren } from "@api/moris";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Plus } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";

export const UserOrganisationManagement = () => {
    const { data: memberships, isLoading } = useGetOrganisationMembershipsMine();

    if (isLoading) return <div>Loading...</div>;

    return (
        <div className="p-6">
            <h1 className="text-2xl font-bold mb-6">My Organizations</h1>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {memberships?.map((membership: any) => (
                    <MembershipCard key={membership.membershipID} membership={membership} />
                ))}
                {memberships?.length === 0 && <div className="text-gray-500">You are not a member of any organization.</div>}
            </div>
        </div>
    );
};

const MembershipCard = ({ membership }: { membership: any }) => {
    // membership is EffectiveMembershipResponse
    // Fields: membershipID, personID, roleScopeID, scopeRootID, roleID, roleKey, hasAdminRights
    // We don't have Node Name here?
    // Backend ListEffectiveMemberships (and ListMyMemberships) returns EffectiveMembershipResponse struct.
    // It DOES NOT include Node Name!
    // This is a missing feature in backend response. I need the node name to display it.
    // I should update backend to include Node Name or Node details.

    // For now, I will display Role Key.
    const canManage = membership.hasAdminRights;

    return (
        <div className="border rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow">
            <div className="flex justify-between items-start">
                <div>
                    <h3 className="font-semibold text-lg">Organization {membership.scopeRootID}</h3>
                    <p className="text-sm text-gray-500">{membership.roleKey} Role</p>
                </div>
                {canManage && <CreateChildDialog parentId={membership.scopeRootID} />}
            </div>
        </div>
    );
};

const CreateChildDialog = ({ parentId }: { parentId: string }) => {
    const [open, setOpen] = useState(false);
    const [name, setName] = useState("");
    const queryClient = useQueryClient();
    const { mutate: createChild, isPending } = usePostOrganisationNodesIdChildren({
        mutation: {
            onSuccess: () => {
                queryClient.invalidateQueries({ queryKey: ['/organisation-memberships/mine'] });
                setOpen(false);
                setName("");
            }
        }
    });

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button variant="outline" size="sm"><Plus size={14} className="mr-1" /> New Unit</Button>
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
