import { useState } from "react";
import { useGetOrganisationMembershipsMine, usePostOrganisationNodesIdChildren } from "@api/moris";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Plus, Users } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { EffectiveMembershipResponse } from "@/api/generated-orval/model";

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

const MembershipCard = ({ membership }: { membership: EffectiveMembershipResponse }) => {
    const canManage = membership.hasAdminRights;

    return (
        <div className="border rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow">
            <div className="flex justify-between items-start">
                <div>
                    <h3 className="font-semibold text-lg">{membership.organisationName || `Organization ${membership.scopeRootID}`}</h3>
                    <p className="text-sm text-gray-500">{membership.roleKey} Role</p>
                </div>
                <div className="flex gap-2">
                    {canManage && (
                        <Button variant="outline" size="sm" asChild>
                            <Link to={`/dashboard/organisations/${membership.scopeRootId}/members`}>
                                <Users size={14} className="mr-1" /> Members
                            </Link>
                        </Button>
                    )}
                    {canManage && <CreateChildDialog parentId={membership.scopeRootId!} />}
                </div>
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
