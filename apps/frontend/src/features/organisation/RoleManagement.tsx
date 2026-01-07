import { useState } from "react";
import { useParams } from "react-router-dom";
import {
  useGetOrganisationNodesIdMembershipsEffective,
  getGetOrganisationNodesIdMembershipsEffectiveQueryKey,
  useGetOrganisationRoles,
  usePostOrganisationScopes,
  usePostOrganisationMemberships,
  useDeleteOrganisationMembershipsId,
} from "@api/moris";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useQueryClient } from "@tanstack/react-query";
import { Trash2 } from "lucide-react";
import { UserSearchSelect } from "@/components/user/UserSearchSelect";
import { OrganisationEffectiveMembershipResponse } from "@/api/generated-orval/model";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";
import { ErrorModal } from "@/components/ui/error-modal";

export const RoleManagement = () => {
  const { nodeId } = useParams<{ nodeId: string }>();
  const { data: members, isLoading } =
    useGetOrganisationNodesIdMembershipsEffective(nodeId!);
  const { data: roles } = useGetOrganisationRoles();

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-xl font-bold">Members & Roles</h2>
        {nodeId && (
          <AddMemberDialog
            nodeId={nodeId}
            roles={roles || []}
            members={members || []}
          />
        )}
      </div>

      <div className="border rounded-lg overflow-hidden">
        <table className="w-full text-sm text-left">
          <thead className="bg-gray-50 text-gray-700 font-medium">
            <tr>
              <th className="px-4 py-3">User</th>
              <th className="px-4 py-3">Role</th>
              <th className="px-4 py-3">Owning Organisation</th>
              <th className="px-4 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y">
            {members?.map((m: OrganisationEffectiveMembershipResponse) => (
              <tr key={m.membershipId} className="bg-white">
                <td className="px-4 py-3">{m.person?.name}</td>
                <td className="px-4 py-3 capitalize">{m.roleKey}</td>
                <td className="px-4 py-3 capitalize">{m.scopeRootOrganisation?.name}</td>
                <td className="px-4 py-3 text-right">
                  <RemoveMemberButton
                    membershipId={m.membershipId!}
                    nodeId={nodeId!}
                  />
                </td>
              </tr>
            ))}
            {members?.length === 0 && (
              <tr>
                <td colSpan={3} className="px-4 py-8 text-center text-gray-500">
                  No members found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

const AddMemberDialog = ({
  nodeId,
  roles,
  members,
}: {
  nodeId: string;
  roles: any[];
  members: OrganisationEffectiveMembershipResponse[];
}) => {
  const [open, setOpen] = useState(false);
  const [personId, setPersonId] = useState("");
  const [roleKey, setRoleKey] = useState("");

  // Error state
  const [error, setError] = useState<string | null>(null);

  const queryClient = useQueryClient();

  const { mutateAsync: createScope } = usePostOrganisationScopes();
  const { mutateAsync: addMembership } = usePostOrganisationMemberships();

  const handleAdd = async () => {
    // Double check validation
    const exists = members.some(
      (m) => m.person?.id === personId && m.roleKey === roleKey
    );
    if (exists) {
      setError("This user already has this role.");
      return;
    }

    try {
      const scopeRes = await createScope({
        data: { roleKey, rootNodeId: nodeId },
      });
      const scopeId = (scopeRes as any).id || (scopeRes as any).data?.id;

      await addMembership({
        data: { personId: personId, roleScopeId: scopeId },
      });

      queryClient.invalidateQueries({
        queryKey: getGetOrganisationNodesIdMembershipsEffectiveQueryKey(nodeId),
      });
      setOpen(false);
      setPersonId("");
      setRoleKey("");
    } catch (e: any) {
      console.error("Failed to add member", e);
      // Show error modal instead of alert
      const msg =
        e?.response?.data?.message || e?.message || "Failed to add member";
      setError(msg);
    }
  };

  // Derived state for validation warning
  const isDuplicate = members.some(
    (m) => m.person?.id === personId && m.roleKey === roleKey
  );

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
            <SelectTrigger>
              <SelectValue placeholder="Select Role" />
            </SelectTrigger>
            <SelectContent>
              {roles.map((r: any) => (
                <SelectItem key={r.key} value={r.key}>
                  {r.key}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          {isDuplicate && (
            <div className="text-sm text-red-500">
              This user already has the '{roleKey}' role.
            </div>
          )}

          <Button
            onClick={handleAdd}
            disabled={!personId || !roleKey || isDuplicate}
          >
            Add
          </Button>
        </div>
      </DialogContent>

      <ErrorModal
        isOpen={!!error}
        onClose={() => setError(null)}
        description={error || "An unknown error occurred"}
        title="Cannot Add Member"
      />
    </Dialog>
  );
};

const RemoveMemberButton = ({
  membershipId,
  nodeId,
}: {
  membershipId: string;
  nodeId: string;
}) => {
  const [showConfirm, setShowConfirm] = useState(false);
  const queryClient = useQueryClient();
  const { mutate: remove, isPending } = useDeleteOrganisationMembershipsId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey:
            getGetOrganisationNodesIdMembershipsEffectiveQueryKey(nodeId),
        });
        setShowConfirm(false);
      },
    },
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
