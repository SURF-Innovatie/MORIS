import { useState } from "react";
import { useGetAuthOrcidUrl, usePostAuthOrcidUnlink } from "@api/moris";
import { UserResponse } from "@api/model";
import { useAuth } from "@/hooks/useAuth";
import { useToast } from "@/hooks/use-toast";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
    DialogClose,
} from "@/components/ui/dialog";
import { Link, Unlink } from "lucide-react";

interface OrcidConnectionProps {
    user: UserResponse;
    refetchProfile: () => Promise<any>;
}

export function OrcidConnection({ user, refetchProfile }: OrcidConnectionProps) {
    const { updateUser } = useAuth();
    const { toast } = useToast();
    const [isUnlinkDialogOpen, setIsUnlinkDialogOpen] = useState(false);

    const { refetch: getAuthURL, isFetching: isGettingURL } = useGetAuthOrcidUrl({
        query: {
            enabled: false,
        },
    });

    const { mutateAsync: unlinkORCID, isPending: isUnlinking } =
        usePostAuthOrcidUnlink();

    const handleLinkORCID = async () => {
        try {
            const result = await getAuthURL();
            if (result.data?.url) {
                window.location.href = result.data.url;
            }
        } catch (error) {
            toast({
                title: "Error",
                description: "Failed to get ORCID authorization URL",
                variant: "destructive",
            });
        }
    };

    const handleUnlinkORCID = async () => {
        try {
            await unlinkORCID();
            toast({
                title: "Success",
                description: "ORCID account unlinked successfully",
            });
            setIsUnlinkDialogOpen(false);
            // Refetch profile to update UI
            const { data: updatedUser } = await refetchProfile();
            if (updatedUser) {
                // Map UserResponse to UserAccount
                updateUser(updatedUser);
            }
        } catch (error) {
            toast({
                title: "Error",
                description: "Failed to unlink ORCID account",
                variant: "destructive",
            });
        }
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>Integrations</CardTitle>
                <CardDescription>
                    Manage your external account connections
                </CardDescription>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    <div className="flex items-start justify-between gap-4 rounded-lg border p-4 bg-muted/30">
                        <div className="space-y-1">
                            <div className="flex items-center gap-2">
                                <h3 className="font-semibold text-sm">ORCID iD</h3>
                                {user.orcid && (
                                    <Badge
                                        variant="secondary"
                                        className="h-5 px-1.5 text-[10px] bg-green-500/10 text-green-600 hover:bg-green-500/20 border-green-500/20"
                                    >
                                        Verified
                                    </Badge>
                                )}
                            </div>
                            <p className="text-xs text-muted-foreground">
                                {user.orcid ? (
                                    <span className="font-mono">{user.orcid}</span>
                                ) : (
                                    "Connect your ORCID iD to your account"
                                )}
                            </p>
                        </div>
                        {user.orcid ? (
                            <Dialog
                                open={isUnlinkDialogOpen}
                                onOpenChange={setIsUnlinkDialogOpen}
                            >
                                <DialogTrigger asChild>
                                    <Button
                                        variant="ghost"
                                        size="icon"
                                        className="h-8 w-8 text-muted-foreground hover:text-destructive"
                                    >
                                        <Unlink className="h-4 w-4" />
                                    </Button>
                                </DialogTrigger>
                                <DialogContent>
                                    <DialogHeader>
                                        <DialogTitle>Unlink ORCID?</DialogTitle>
                                        <DialogDescription>
                                            Are you sure you want to unlink your ORCID account?
                                            This action cannot be undone easily.
                                        </DialogDescription>
                                    </DialogHeader>
                                    <DialogFooter>
                                        <DialogClose asChild>
                                            <Button variant="outline">Cancel</Button>
                                        </DialogClose>
                                        <Button
                                            variant="destructive"
                                            onClick={handleUnlinkORCID}
                                            disabled={isUnlinking}
                                        >
                                            {isUnlinking ? "Unlinking..." : "Unlink"}
                                        </Button>
                                    </DialogFooter>
                                </DialogContent>
                            </Dialog>
                        ) : (
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={handleLinkORCID}
                                disabled={isGettingURL}
                                className="h-8"
                            >
                                <Link className="mr-2 h-3.5 w-3.5" />
                                {isGettingURL ? "Connecting..." : "Connect"}
                            </Button>
                        )}
                    </div>
                </div>
            </CardContent>
        </Card>
    );
}
