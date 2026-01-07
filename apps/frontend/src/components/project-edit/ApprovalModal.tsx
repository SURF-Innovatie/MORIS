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
  useGetEventsId,
  usePostEventsIdApprove,
  usePostEventsIdReject,
} from "@/api/generated-orval/moris";
import { useQueryClient } from "@tanstack/react-query";
import { EventRenderer } from "@/components/events/EventRenderer";
import { ProjectEvent } from "@/api/events";

interface ApprovalModalProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  eventId: string;
  notificationMessage?: string;
}

export function ApprovalModal({
  isOpen,
  onClose,
  projectId,
  eventId,
  notificationMessage,
}: ApprovalModalProps) {
  const queryClient = useQueryClient();

  const { data: eventData, isLoading: isLoadingEvent } = useGetEventsId(
    eventId,
    {
      query: {
        enabled: isOpen && !!eventId,
      },
    }
  );

  const { mutate: approve, isPending: isApproving } = usePostEventsIdApprove({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: [`/projects/${projectId}/pending-events`],
        });
        queryClient.invalidateQueries({ queryKey: [`/notifications`] });
        queryClient.invalidateQueries({ queryKey: [`/events/${eventId}`] });
        onClose();
      },
    },
  });

  const { mutate: reject, isPending: isRejecting } = usePostEventsIdReject({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: [`/projects/${projectId}/pending-events`],
        });
        queryClient.invalidateQueries({ queryKey: [`/notifications`] });
        queryClient.invalidateQueries({ queryKey: [`/events/${eventId}`] });
        onClose();
      },
    },
  });

  const event = eventData;
  const isPending = event?.status === "pending";

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
            {isPending
              ? "Review the details of this request before approving or rejecting."
              : "This request has already been processed."}
          </DialogDescription>
        </DialogHeader>

        {isLoadingEvent ? (
          <div className="flex justify-center py-8">
            <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
          </div>
        ) : event ? (
          <div className="py-4 space-y-4">
            <EventRenderer
              event={event as ProjectEvent}
              className="border rounded-lg p-4 bg-gray-50/50"
            />

            {!isPending && (
              <div className="text-center text-sm text-muted-foreground bg-muted/50 p-2 rounded">
                Status:{" "}
                <span className="font-medium capitalize">{event.status}</span>
              </div>
            )}
          </div>
        ) : (
          <div className="py-6 space-y-4">
            <div className="rounded-md bg-muted p-4">
              <p className="text-sm font-medium">Notification Message:</p>
              <p className="text-sm text-muted-foreground mt-1">
                {notificationMessage}
              </p>
            </div>
            <p className="text-center text-sm text-muted-foreground">
              Unable to load event details.
            </p>
          </div>
        )}

        <DialogFooter>
          <Button
            variant="outline"
            onClick={onClose}
            disabled={isApproving || isRejecting}
          >
            Close
          </Button>
          {event && isPending && (
            <>
              <Button
                variant="destructive"
                onClick={handleReject}
                disabled={isApproving || isRejecting}
              >
                {isRejecting && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Reject
              </Button>
              <Button
                onClick={handleApprove}
                disabled={isApproving || isRejecting}
              >
                {isApproving && (
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                )}
                Approve
              </Button>
            </>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
