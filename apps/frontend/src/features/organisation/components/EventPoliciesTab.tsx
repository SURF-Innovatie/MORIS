import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  useGetOrganisationsIdPolicies,
  useDeletePoliciesId,
  useGetEventTypes,
  getGetOrganisationsIdPoliciesQueryKey,
  useGetOrganisationNodesIdOrganisationRoles,
  useGetOrganisationNodesIdRoles,
} from "@api/moris";
import { EventPolicyResponse } from "@api/model";
import { Button } from "@/components/ui/button";
import { Dialog, DialogTrigger } from "@/components/ui/dialog";
import { Loader2, Plus } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { PolicyCard } from "./PolicyCard";
import { PolicyFormDialog } from "./PolicyFormDialog";

interface EventPoliciesTabProps {
  nodeId: string;
}

export function EventPoliciesTab({ nodeId }: EventPoliciesTabProps) {
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editPolicy, setEditPolicy] = useState<EventPolicyResponse | null>(
    null
  );
  const queryClient = useQueryClient();
  const { toast } = useToast();

  // Fetch event types from backend
  const { data: eventTypes = [] } = useGetEventTypes();

  // Fetch org roles for this node
  const { data: orgRoles = [] } =
    useGetOrganisationNodesIdOrganisationRoles(nodeId);

  // Fetch project roles available at this node (includes inherited)
  const { data: projectRoles = [] } = useGetOrganisationNodesIdRoles(nodeId);

  const { data: policies, isLoading } = useGetOrganisationsIdPolicies(nodeId, {
    inherited: true,
  });

  const deleteMutation = useDeletePoliciesId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetOrganisationsIdPoliciesQueryKey(nodeId),
        });
        toast({ title: "Policy deleted" });
      },
      onError: () => {
        toast({ title: "Failed to delete policy", variant: "destructive" });
      },
    },
  });

  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const ownPolicies =
    policies?.filter((p: EventPolicyResponse) => !p.inherited) || [];
  const inheritedPolicies =
    policies?.filter((p: EventPolicyResponse) => p.inherited) || [];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-xl font-semibold">Event Policies</h2>
          <p className="text-sm text-muted-foreground">
            Configure automatic notifications and approval workflows
          </p>
        </div>
        <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              New Policy
            </Button>
          </DialogTrigger>
          <PolicyFormDialog
            nodeId={nodeId}
            eventTypes={eventTypes}
            orgRoles={orgRoles}
            projectRoles={projectRoles}
            onClose={() => setCreateDialogOpen(false)}
          />
        </Dialog>
      </div>

      {/* Own Policies */}
      <div className="space-y-4">
        <h3 className="text-lg font-medium">Organisation Policies</h3>
        {ownPolicies.length === 0 ? (
          <Card>
            <CardContent className="py-8 text-center text-muted-foreground">
              No policies configured. Create one to get started.
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {ownPolicies.map((policy: EventPolicyResponse) => (
              <PolicyCard
                key={policy.id}
                policy={policy}
                eventTypes={eventTypes}
                projectRoles={projectRoles}
                orgRoles={orgRoles}
                onEdit={() => setEditPolicy(policy)}
                onDelete={() => deleteMutation.mutate({ id: policy.id! })}
                isDeleting={deleteMutation.isPending}
              />
            ))}
          </div>
        )}
      </div>

      {/* Inherited Policies */}
      {inheritedPolicies.length > 0 && (
        <div className="space-y-4">
          <h3 className="text-lg font-medium flex items-center gap-2">
            Inherited Policies
            <Badge variant="outline" className="font-normal">
              Read-only
            </Badge>
          </h3>
          <div className="grid gap-4">
            {inheritedPolicies.map((policy: EventPolicyResponse) => (
              <PolicyCard
                key={policy.id}
                policy={policy}
                eventTypes={eventTypes}
                projectRoles={projectRoles}
                orgRoles={orgRoles}
                inherited
              />
            ))}
          </div>
        </div>
      )}

      {/* Edit Dialog */}
      <Dialog
        open={!!editPolicy}
        onOpenChange={(open) => !open && setEditPolicy(null)}
      >
        {editPolicy && (
          <PolicyFormDialog
            nodeId={nodeId}
            eventTypes={eventTypes}
            orgRoles={orgRoles}
            projectRoles={projectRoles}
            policy={editPolicy}
            onClose={() => setEditPolicy(null)}
          />
        )}
      </Dialog>
    </div>
  );
}
