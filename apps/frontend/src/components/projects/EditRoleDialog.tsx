import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { usePutProjectsIdPeoplePersonId, useGetProjectsRoles } from "@api/moris";
import { MemberResponse } from "@api/model";
import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogFooter,
} from "@/components/ui/dialog";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/hooks/use-toast";
import { Loader2 } from "lucide-react";

interface EditRoleDialogProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    member: MemberResponse;
    projectId: string;
    onSuccess?: () => void;
}

export function EditRoleDialog({
    open,
    onOpenChange,
    member,
    projectId,
    onSuccess,
}: EditRoleDialogProps) {
    const [role, setRole] = useState(member.role || "");
    const queryClient = useQueryClient();
    const { toast } = useToast();

    const { data: roles, isLoading: isLoadingRoles } = useGetProjectsRoles();

    const { mutate, isPending } = usePutProjectsIdPeoplePersonId({
        mutation: {
            onSuccess: () => {
                queryClient.invalidateQueries({ queryKey: ["/projects", projectId] });
                toast({ title: "Role updated" });
                onOpenChange(false);
                onSuccess?.();
            },
            onError: () => {
                toast({
                    title: "Failed to update role",
                    variant: "destructive"
                });
            }
        },
    });

    const handleSave = () => {
        if (!member.id) return;
        mutate({
            id: projectId,
            personId: member.id,
            data: { role },
        });
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent>
                <DialogHeader>
                    <DialogTitle>Edit Role</DialogTitle>
                </DialogHeader>
                <div className="py-4">
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Role</label>
                        {isLoadingRoles ? (
                            <div className="flex bg-muted h-10 w-full items-center px-3 rounded-md">
                                <Loader2 className="h-4 w-4 animate-spin text-muted-foreground mr-2" />
                                <span className="text-sm text-muted-foreground">Loading roles...</span>
                            </div>
                        ) : (
                            <Select value={role} onValueChange={setRole}>
                                <SelectTrigger>
                                    <SelectValue placeholder="Select a role" />
                                </SelectTrigger>
                                <SelectContent>
                                    {roles?.map((r) => (
                                        <SelectItem key={r.key} value={r.key || ""}>
                                            {r.name}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                        )}
                    </div>
                </div>
                <DialogFooter>
                    <Button variant="outline" onClick={() => onOpenChange(false)}>
                        Cancel
                    </Button>
                    <Button onClick={handleSave} disabled={isPending || isLoadingRoles}>
                        {isPending ? "Saving..." : "Save Changes"}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
