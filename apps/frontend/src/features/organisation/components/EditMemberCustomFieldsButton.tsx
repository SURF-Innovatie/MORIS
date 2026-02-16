import { useState } from "react";
import {
  useGetOrganisationNodesIdCustomFields,
  usePutOrganisationNodesIdMembersPersonIdCustomFields,
  getGetOrganisationNodesIdMembershipsEffectiveQueryKey,
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
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { useToast } from "@/hooks/use-toast";
import { useForm } from "react-hook-form";
import { Loader2, Pencil } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { OrganisationEffectiveMembershipResponse } from "@/api/generated-orval/model";

interface EditMemberCustomFieldsButtonProps {
  nodeId: string;
  membership: OrganisationEffectiveMembershipResponse;
  canEdit: boolean;
}

export function EditMemberCustomFieldsButton({
  nodeId,
  membership,
  canEdit,
}: EditMemberCustomFieldsButtonProps) {
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

  const { register, handleSubmit } = useForm({
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
