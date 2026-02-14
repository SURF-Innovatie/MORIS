import { useState } from "react";
import {
  getGetOrganisationNodesIdMembershipsEffectiveQueryKey,
  usePostOrganisationScopes,
  usePostOrganisationMemberships,
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
import { UserSearchSelect } from "@/components/user/UserSearchSelect";
import { OrganisationEffectiveMembershipResponse } from "@/api/generated-orval/model";
import { ErrorModal } from "@/components/ui/error-modal";

interface AddMemberDialogProps {
  nodeId: string;
  roles: any[];
  members: OrganisationEffectiveMembershipResponse[];
  disabled?: boolean;
}

export function AddMemberDialog({
  nodeId,
  roles,
  members,
  disabled,
}: AddMemberDialogProps) {
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
      {!disabled && (
        <DialogTrigger asChild>
          <Button>Add Member</Button>
        </DialogTrigger>
      )}
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
}
