import { useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import {
  useGetOrganisationNodesIdMembershipsEffective,
  getGetOrganisationNodesIdMembershipsEffectiveQueryKey,
  useGetOrganisationNodesIdCustomFields,
  usePutOrganisationNodesIdMembersPersonIdCustomFields,
  useGetOrganisationNodesIdOrganisationRoles,
  usePostOrganisationScopes,
  usePostOrganisationMemberships,
  useDeleteOrganisationMembershipsId,
  useGetOrganisationNodesIdPermissionsMine,
} from "@api/moris";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { useToast } from "@/hooks/use-toast";
import { useForm } from "react-hook-form";
import { Loader2, Pencil, Trash2 } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";

import { UserSearchSelect } from "@/components/user/UserSearchSelect";
import { OrganisationEffectiveMembershipResponse } from "@/api/generated-orval/model";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";
import { ErrorModal } from "@/components/ui/error-modal";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ProjectRolesList } from "./components/ProjectRolesList";
import { CustomFieldDefinitionsList } from "./components/CustomFieldDefinitionsList";
import { RoleList } from "./components/RoleList"; // New component

import { OrganisationEditTab } from "./components/OrganisationEditTab";

export const RoleManagement = () => {
  const { nodeId } = useParams<{ nodeId: string }>();
  const { data: members, isLoading: isLoadingMembers } =
    useGetOrganisationNodesIdMembershipsEffective(nodeId!);
  const { data: roles } = useGetOrganisationNodesIdOrganisationRoles(nodeId!);
  const { data: myPermissions, isLoading: isLoadingPerms } =
    useGetOrganisationNodesIdPermissionsMine(nodeId!);

  const canManageMembers = myPermissions?.includes("manage_members") ?? false;
  const canManageCustomFields =
    myPermissions?.includes("manage_custom_fields") ?? false;

  const [searchParams, setSearchParams] = useSearchParams();
  const currentTab = searchParams.get("tab") || "members";

  if (isLoadingMembers || isLoadingPerms) return <div>Loading...</div>;

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Organisation Settings</h1>
      </div>

      <Tabs
        value={currentTab}
        onValueChange={(val) => setSearchParams({ tab: val })}
        className="w-full"
      >
        <TabsList className="mb-4">
          <TabsTrigger value="members">Members & Permissions</TabsTrigger>
          <TabsTrigger value="roles">Roles</TabsTrigger>
          <TabsTrigger value="project-roles">Project Roles</TabsTrigger>
          <TabsTrigger value="custom-fields">Custom Fields</TabsTrigger>
          <TabsTrigger value="edit">Edit</TabsTrigger>
        </TabsList>

        <TabsContent value="members">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">Members</h2>
            {nodeId && (
              <AddMemberDialog
                nodeId={nodeId}
                roles={roles || []}
                members={members || []}
                disabled={!canManageMembers}
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
                    <td className="px-4 py-3 capitalize">
                      {m.scopeRootOrganisation?.name}
                    </td>
                    <td className="px-4 py-3 text-right flex gap-2 justify-end">
                      <EditMemberCustomFieldsButton
                        nodeId={nodeId!}
                        membership={m}
                        canEdit={canManageCustomFields}
                      />
                      {canManageMembers && (
                        <RemoveMemberButton
                          membershipId={m.membershipId!}
                          nodeId={nodeId!}
                        />
                      )}
                    </td>
                  </tr>
                ))}
                {members?.length === 0 && (
                  <tr>
                    <td
                      colSpan={4}
                      className="px-4 py-8 text-center text-gray-500"
                    >
                      No members found.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </TabsContent>

        <TabsContent value="roles">
          {nodeId ? <RoleList nodeId={nodeId} /> : <div>Invalid Node ID</div>}
        </TabsContent>

        <TabsContent value="project-roles">
          {nodeId ? (
            <ProjectRolesList nodeId={nodeId} />
          ) : (
            <div>Invalid Node ID</div>
          )}
        </TabsContent>

        <TabsContent value="custom-fields">
          {nodeId ? (
            <CustomFieldDefinitionsList nodeId={nodeId} />
          ) : (
            <div>Invalid Node ID</div>
          )}
        </TabsContent>

        <TabsContent value="edit">
          {nodeId ? (
            <OrganisationEditTab nodeId={nodeId} />
          ) : (
            <div>Invalid Node ID</div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
};

function EditMemberCustomFieldsButton({
  nodeId,
  membership,
  canEdit,
}: {
  nodeId: string;
  membership: OrganisationEffectiveMembershipResponse;
  canEdit: boolean;
}) {
  const [open, setOpen] = useState(false);
  const { data: fields, isLoading: isLoadingFields } =
    useGetOrganisationNodesIdCustomFields(nodeId, { category: "PERSON" });
  const { mutateAsync: updateFields, isPending } =
    usePutOrganisationNodesIdMembersPersonIdCustomFields();
  const { toast } = useToast();
  const queryClient = useQueryClient();

  // Prepare default values from membership custom fields
  const defaultValues: Record<string, any> = {};
  if (fields && membership.customFields) {
    fields.forEach((f: any) => {
      if (membership.customFields && membership.customFields[f.id]) {
        defaultValues[f.id] = membership.customFields[f.id];
      }
    });
  }

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    defaultValues,
    values: defaultValues, // Reactive update when fields load
  });

  const onSubmit = async (data: any) => {
    if (!canEdit) return;
    try {
      await updateFields({
        id: nodeId,
        personId: membership.person!.id!,
        data: { values: data },
      });
      toast({ title: "Custom fields updated" });
      queryClient.invalidateQueries({
        queryKey: getGetOrganisationNodesIdMembershipsEffectiveQueryKey(nodeId),
      });
      setOpen(false);
    } catch (error) {
      console.error("Failed to update custom fields", error);
      toast({
        title: "Failed to update custom fields",
        variant: "destructive",
      });
    }
  };

  // Debug check
  if (!membership.person?.id) {
    console.warn("EditMemberCustomFieldsButton: Missing person ID", membership);
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      {canEdit && (
        <DialogTrigger asChild>
          <Button
            variant="ghost"
            size="icon"
            title="Edit Custom Fields"
            disabled={!membership.person?.id}
          >
            <Pencil className="h-4 w-4 text-gray-500" />
          </Button>
        </DialogTrigger>
      )}
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>
            Edit Custom Fields for {membership.person?.name}
          </DialogTitle>
        </DialogHeader>

        {isLoadingFields ? (
          <div className="flex justify-center p-4">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : fields?.length === 0 ? (
          <div className="text-center p-4 text-muted-foreground">
            No person custom fields defined for this organisation.
          </div>
        ) : (
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {fields?.map((field: any) => (
              <div key={field.id} className="space-y-2">
                <label className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
                  {field.name}
                  {field.required && (
                    <span className="text-destructive ml-1">*</span>
                  )}
                </label>
                {field.type === "BOOLEAN" ? (
                  <div className="flex items-center space-x-2">
                    <Checkbox
                      id={field.id}
                      // Handling checkbox with register is tricky without Controller, but simple native behavior works usually if value prop is managed.
                      // For simplicity in this non-controlled form, we use register.
                      {...register(field.id)}
                    />
                    <label
                      htmlFor={field.id}
                      className="text-sm text-muted-foreground"
                    >
                      {field.description}
                    </label>
                  </div>
                ) : (
                  <Input
                    {...register(field.id, { required: field.required })}
                    type={
                      field.type === "NUMBER"
                        ? "number"
                        : field.type === "DATE"
                          ? "date"
                          : "text"
                    }
                    placeholder={field.example_value}
                  />
                )}
                {field.description && field.type !== "BOOLEAN" && (
                  <p className="text-[0.8rem] text-muted-foreground">
                    {field.description}
                  </p>
                )}
              </div>
            ))}
            <DialogFooter>
              <Button type="submit" disabled={isPending}>
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Save Changes
              </Button>
            </DialogFooter>
          </form>
        )}
      </DialogContent>
    </Dialog>
  );
}

function AddMemberDialog({
  nodeId,
  roles,
  members,
  disabled,
}: {
  nodeId: string;
  roles: any[];
  members: OrganisationEffectiveMembershipResponse[];
  disabled?: boolean;
}) {
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

function RemoveMemberButton({
  membershipId,
  nodeId,
}: {
  membershipId: string;
  nodeId: string;
}) {
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
