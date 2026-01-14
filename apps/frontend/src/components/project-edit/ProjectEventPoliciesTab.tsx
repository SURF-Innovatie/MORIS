import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  useGetProjectsIdPolicies,
  useGetEventTypes,
  getGetProjectsIdPoliciesQueryKey,
  useGetProjectsIdRoles,
  useGetOrganisationNodesIdOrganisationRoles,
} from "@api/moris";
import { EventPolicyResponse } from "@api/model";
import { createEventPolicyRemovedEvent, ProjectEventType } from "@/api/events";
import { useAccess } from "@/context/AccessContext";
import { Button } from "@/components/ui/button";
import { Dialog, DialogTrigger } from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Loader2, Plus, Lock } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import { Card, CardContent } from "@/components/ui/card";
import { ProjectPolicyCard } from "./ProjectPolicyCard";
import { ProjectPolicyFormDialog } from "./ProjectPolicyFormDialog";

interface ProjectEventPoliciesTabProps {
  projectId: string;
  orgNodeId: string;
}

export function ProjectEventPoliciesTab({
  projectId,
  orgNodeId,
}: ProjectEventPoliciesTabProps) {
  const [createDialogOpen, setCreateDialogOpen] = useState(false);
  const [editPolicy, setEditPolicy] = useState<EventPolicyResponse | null>(
    null
  );
  const [isDeleting, setIsDeleting] = useState<string | null>(null);
  const queryClient = useQueryClient();
  const { toast } = useToast();
  const { hasAccess } = useAccess();

  const canAddPolicy = hasAccess(ProjectEventType.EventPolicyAdded);
  const canRemovePolicy = hasAccess(ProjectEventType.EventPolicyRemoved);
  const canUpdatePolicy = hasAccess(ProjectEventType.EventPolicyUpdated);

  // Fetch event types from backend
  const { data: eventTypes = [] } = useGetEventTypes();

  // Fetch project roles
  const { data: projectRoles = [] } = useGetProjectsIdRoles(projectId);

  // Fetch org roles
  const { data: orgRoles = [] } =
    useGetOrganisationNodesIdOrganisationRoles(orgNodeId);

  const { data: policies, isLoading } = useGetProjectsIdPolicies(projectId, {
    org_node_id: orgNodeId,
    inherited: true,
  });

  const handleDeletePolicy = async (policy: EventPolicyResponse) => {
    if (!policy.id) return;
    setIsDeleting(policy.id);
    try {
      await createEventPolicyRemovedEvent(projectId, {
        policy_id: policy.id,
        name: policy.name || "",
      });
      queryClient.invalidateQueries({
        queryKey: getGetProjectsIdPoliciesQueryKey(projectId),
      });
      toast({ title: "Policy removal requested" });
    } catch (error) {
      toast({ title: "Failed to remove policy", variant: "destructive" });
    } finally {
      setIsDeleting(null);
    }
  };

  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-12">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const projectPolicies =
    policies?.filter((p: EventPolicyResponse) => p.project_id === projectId) ||
    [];
  const inheritedPolicies =
    policies?.filter((p: EventPolicyResponse) => p.project_id !== projectId) ||
    [];

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-xl font-semibold">Event Policies</h2>
          <p className="text-sm text-muted-foreground">
            Configure automatic notifications and approval workflows for this
            project
          </p>
        </div>
        {canAddPolicy ? (
          <Dialog open={createDialogOpen} onOpenChange={setCreateDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="h-4 w-4 mr-2" />
                New Policy
              </Button>
            </DialogTrigger>
            <ProjectPolicyFormDialog
              projectId={projectId}
              eventTypes={eventTypes}
              projectRoles={projectRoles}
              orgRoles={orgRoles}
              onClose={() => setCreateDialogOpen(false)}
            />
          </Dialog>
        ) : (
          <Button disabled>
            <Lock className="h-4 w-4 mr-2" />
            No Permission
          </Button>
        )}
      </div>

      {/* Project Policies */}
      <div className="space-y-4">
        <h3 className="text-lg font-medium">Project Policies</h3>
        {projectPolicies.length === 0 ? (
          <Card>
            <CardContent className="py-8 text-center text-muted-foreground">
              No project-specific policies configured.
            </CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {projectPolicies.map((policy: EventPolicyResponse) => (
              <ProjectPolicyCard
                key={policy.id}
                policy={policy}
                eventTypes={eventTypes}
                onEdit={
                  canUpdatePolicy ? () => setEditPolicy(policy) : undefined
                }
                onDelete={
                  canRemovePolicy ? () => handleDeletePolicy(policy) : undefined
                }
                isDeleting={isDeleting === policy.id}
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
              From Organisation
            </Badge>
          </h3>
          <div className="grid gap-4">
            {inheritedPolicies.map((policy: EventPolicyResponse) => (
              <ProjectPolicyCard
                key={policy.id}
                policy={policy}
                eventTypes={eventTypes}
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
          <ProjectPolicyFormDialog
            projectId={projectId}
            eventTypes={eventTypes}
            projectRoles={projectRoles}
            orgRoles={orgRoles}
            policy={editPolicy}
            onClose={() => setEditPolicy(null)}
          />
        )}
      </Dialog>
    </div>
  );
}
