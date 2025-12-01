import { Loader2 } from "lucide-react";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import {
    useGetProjectsIdPendingEvents,
    usePostEventsIdApprove,
    usePostEventsIdReject,
} from "@/api/generated-orval/moris";
import { useQueryClient } from "@tanstack/react-query";

interface ApprovalModalProps {
    isOpen: boolean;
    onClose: () => void;
    projectId: string;
    eventId: string;
}

export function ApprovalModal({
    isOpen,
    onClose,
    projectId,
    eventId,
}: ApprovalModalProps) {
    const queryClient = useQueryClient();

    const { data: events, isLoading } = useGetProjectsIdPendingEvents(projectId, {
        query: {
            enabled: isOpen && !!projectId,
        },
    });

    const { mutate: approve, isPending: isApproving } = usePostEventsIdApprove({
        mutation: {
            onSuccess: () => {
                queryClient.invalidateQueries({ queryKey: [`/projects/${projectId}/pending-events`] });
                queryClient.invalidateQueries({ queryKey: [`/notifications`] });
                onClose();
            },
        },
    });

    const { mutate: reject, isPending: isRejecting } = usePostEventsIdReject({
        mutation: {
            onSuccess: () => {
                queryClient.invalidateQueries({ queryKey: [`/projects/${projectId}/pending-events`] });
                queryClient.invalidateQueries({ queryKey: [`/notifications`] });
                onClose();
            },
        },
    });

    const event = events?.events?.find((e) => e.id === eventId);

    const handleApprove = () => {
        approve({ id: eventId });
    };

    const handleReject = () => {
        reject({ id: eventId });
    };

    if (!isOpen) return null;

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Approval Request</DialogTitle>
                    <DialogDescription>
                        Review the details of this request before approving or rejecting.
                    </DialogDescription>
                </DialogHeader>

                {isLoading ? (
                    <div className="flex justify-center py-8">
                        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
                    </div>
                ) : event ? (
                    <div className="grid gap-4 py-4">
                        <div className="grid grid-cols-4 items-center gap-4">
                            <span className="font-medium text-right">Type:</span>
                            <span className="col-span-3 capitalize">
                                {event.type?.replace(/_/g, " ")}
                            </span>
                        </div>
                        <div className="grid grid-cols-4 items-center gap-4">
                            <span className="font-medium text-right">Details:</span>
                            <span className="col-span-3 text-sm text-muted-foreground">
                                {JSON.stringify(event.details, null, 2)}
                            </span>
                        </div>
                        <div className="grid grid-cols-4 items-center gap-4">
                            <span className="font-medium text-right">Requested by:</span>
                            <span className="col-span-3 text-sm text-muted-foreground">
                                {event.createdBy}
                            </span>
                        </div>
                    </div>
                ) : (
                    <div className="py-8 text-center text-muted-foreground">
                        Event not found or already processed.
                    </div>
                )}

                <DialogFooter>
                    <Button variant="outline" onClick={onClose} disabled={isApproving || isRejecting}>
                        Cancel
                    </Button>
                    <Button
                        variant="destructive"
                        onClick={handleReject}
                        disabled={isApproving || isRejecting || !event}
                    >
                        {isRejecting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                        Reject
                    </Button>
                    <Button
                        onClick={handleApprove}
                        disabled={isApproving || isRejecting || !event}
                    >
                        {isApproving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                        Approve
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
