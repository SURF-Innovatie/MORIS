import { useState } from "react";
import { useParams } from "react-router-dom";
import {
    useGetOrganisationNodesIdMembershipsEffective,
    getGetOrganisationNodesIdMembershipsEffectiveQueryKey,
    useGetOrganisationRoles,
    usePostOrganisationScopes,
    usePostOrganisationMemberships,
    useDeleteOrganisationMembershipsId
} from "@api/moris";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useQueryClient } from "@tanstack/react-query";
import { Trash2 } from "lucide-react";
import { UserSearchSelect } from "@/components/user/UserSearchSelect";
import { EffectiveMembershipResponse } from "@/api/generated-orval/model";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";

export const RoleManagement = () => {
    const { nodeId } = useParams<{ nodeId: string }>();
    const { data: members, isLoading } = useGetOrganisationNodesIdMembershipsEffective(nodeId!);
    const { data: roles } = useGetOrganisationRoles();

    if (isLoading) return <div>Loading...</div>;

    return (
        <div className="p-6">
            <div className="flex justify-between items-center mb-6">
                <h2 className="text-xl font-bold">Members & Roles</h2>
                {nodeId && <AddMemberDialog nodeId={nodeId} roles={roles || []} />}
            </div>

            <div className="border rounded-lg overflow-hidden">
                <table className="w-full text-sm text-left">
                    <thead className="bg-gray-50 text-gray-700 font-medium">
                        <tr>
                            <th className="px-4 py-3">User</th>
                            <th className="px-4 py-3">Role</th>
                            <th className="px-4 py-3 text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y">
                        {members?.map((m: EffectiveMembershipResponse) => (
                            <tr key={m.membershipId} className="bg-white">
                                <td className="px-4 py-3">{m.person?.name}</td>
                                <td className="px-4 py-3 capitalize">{m.roleKey}</td>
                                <td className="px-4 py-3 text-right">
                                    <RemoveMemberButton membershipId={m.membershipId!} nodeId={nodeId!} />
                                </td>
                            </tr>
                        ))}
                        {members?.length === 0 && (
                            <tr><td colSpan={3} className="px-4 py-8 text-center text-gray-500">No members found.</td></tr>
                        )}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

const AddMemberDialog = ({ nodeId, roles }: { nodeId: string, roles: any[] }) => {
    const [open, setOpen] = useState(false);
    const [personId, setPersonId] = useState("");
    const [roleKey, setRoleKey] = useState("");

    const queryClient = useQueryClient();

    const { mutateAsync: createScope } = usePostOrganisationScopes();
    const { mutateAsync: addMembership } = usePostOrganisationMemberships();

    const handleAdd = async () => {
        try {
            const scopeRes = await createScope({ data: { roleKey, rootNodeId: nodeId } });
            const scopeId = (scopeRes as any).id || (scopeRes as any).data?.id;

            await addMembership({ data: { personId: personId, roleScopeId: scopeId } });

            queryClient.invalidateQueries({ queryKey: getGetOrganisationNodesIdMembershipsEffectiveQueryKey(nodeId) });
            setOpen(false);
            setPersonId("");
            setRoleKey("");
        } catch (e) {
            console.error("Failed to add member", e);
            alert("Failed to add member (scope might already exist?)");
        }
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button>Add Member</Button>
            </DialogTrigger>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Add Member</DialogTitle>
                </DialogHeader>
                <div className="space-y-4 pt-4">
                    <UserSearchSelect
                        value={personId}
                        onSelect={(id) => setPersonId(id)}
                    />
                    <Select onValueChange={setRoleKey}>
                        <SelectTrigger><SelectValue placeholder="Select Role" /></SelectTrigger>
                        <SelectContent>
                            {roles.map((r: any) => (
                                <SelectItem key={r.key} value={r.key}>{r.key}</SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                    <Button onClick={handleAdd} disabled={!personId || !roleKey}>
                        Add
                    </Button>
                </div>
            </DialogContent>
        </Dialog>
    );
};

const RemoveMemberButton = ({ membershipId, nodeId }: { membershipId: string, nodeId: string }) => {
    const [showConfirm, setShowConfirm] = useState(false);
    const queryClient = useQueryClient();
    const { mutate: remove, isPending } = useDeleteOrganisationMembershipsId({
        mutation: {
            onSuccess: () => {
                queryClient.invalidateQueries({ queryKey: getGetOrganisationNodesIdMembershipsEffectiveQueryKey(nodeId) });
                setShowConfirm(false);
            }
        }
    });

    return (
        <>
            <Button variant="ghost" size="icon" onClick={() => setShowConfirm(true)}>
                <Trash2 size={16} className="text-red-500" />
            </Button>

            <ConfirmationModal
                isOpen={showConfirm}
                onClose={() => setShowConfirm(false)}
                onConfirm={() => remove({ id: membershipId })}
                title="Remove Member"
                description="Are you sure you want to remove this member? This action cannot be undone."
                confirmLabel="Remove"
                variant="destructive"
                isLoading={isPending}
            />
        </>
    );
};
