import { useState } from "react";
import {
  getGetOrganisationNodesIdMembershipsEffectiveQueryKey,
  useDeleteOrganisationMembershipsId,
} from "@api/moris";
import { Button } from "@/components/ui/button";
import { Trash2 } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";

interface RemoveMemberButtonProps {
  membershipId: string;
  nodeId: string;
}

export function RemoveMemberButton({
  membershipId,
  nodeId,
}: RemoveMemberButtonProps) {
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
}
