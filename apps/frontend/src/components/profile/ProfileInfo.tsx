import { useState } from "react";
import { usePutPeopleId } from "@api/moris";
import { UserResponse } from "@api/model";
import { useToast } from "@/hooks/use-toast";
import { Button } from "@/components/ui/button";
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
} from "@/components/ui/dialog";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Pencil } from "lucide-react";

interface ProfileInfoProps {
    user: UserResponse;
    refetchProfile: () => void;
}

export function ProfileInfo({ user, refetchProfile }: ProfileInfoProps) {
    const { toast } = useToast();
    const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
    const [editForm, setEditForm] = useState({
        name: "",
        avatarUrl: "",
        description: "",
    });

    const { mutateAsync: updatePerson, isPending: isUpdating } = usePutPeopleId();

    const handleEditOpen = () => {
        if (user) {
            setEditForm({
                name: user.name || "",
                avatarUrl: user.avatarUrl || "",
                description: user.description || "",
            });
            setIsEditDialogOpen(true);
        }
    };

    const handleEditSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user) return;

        try {
            await updatePerson({
                id: user.person_id!,
                data: {
                    name: editForm.name,
                    email: user.email,
                    avatarUrl: editForm.avatarUrl || undefined,
                    description: editForm.description || undefined,
                    user_id: user.id,
                },
            });

            toast({
                title: "Success",
                description: "Profile updated successfully",
            });
            setIsEditDialogOpen(false);
            refetchProfile();
        } catch (error) {
            toast({
                title: "Error",
                description: "Failed to update profile",
                variant: "destructive",
            });
        }
    };

    return (
        <>
            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <div className="space-y-1">
                        <CardTitle>Personal Information</CardTitle>
                        <CardDescription>Manage your personal details</CardDescription>
                    </div>
                    <Button variant="ghost" size="icon" onClick={handleEditOpen}>
                        <Pencil className="h-4 w-4" />
                    </Button>
                </CardHeader>
                <CardContent className="space-y-6 pt-4">
                    <div className="flex flex-col items-center gap-4">
                        <Avatar className="h-24 w-24">
                            <AvatarImage src={user.avatarUrl || ""} alt={user.name} />
                            <AvatarFallback className="text-2xl">
                                {user.name
                                    ?.split(" ")
                                    .map((n: string) => n[0])
                                    .join("")
                                    .toUpperCase()
                                    .slice(0, 2)}
                            </AvatarFallback>
                        </Avatar>
                    </div>

                    <div className="flex flex-col gap-1">
                        <label className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                            Full Name
                        </label>
                        <p className="font-medium">{user.name || "N/A"}</p>
                    </div>
                    <div className="flex flex-col gap-1">
                        <label className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                            Email
                        </label>
                        <p className="font-medium">{user.email}</p>
                    </div>
                    {user.description && (
                        <div className="flex flex-col gap-1">
                            <label className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                                Description
                            </label>
                            <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                                {user.description}
                            </p>
                        </div>
                    )}
                </CardContent>
            </Card>

            <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Edit Profile</DialogTitle>
                        <DialogDescription>
                            Update your personal information.
                        </DialogDescription>
                    </DialogHeader>
                    <form onSubmit={handleEditSubmit} className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="name">Full Name</Label>
                            <Input
                                id="name"
                                value={editForm.name}
                                onChange={(e) =>
                                    setEditForm({ ...editForm, name: e.target.value })
                                }
                                required
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="email">Email</Label>
                            <Input id="email" value={user.email} disabled />
                            <p className="text-[0.8rem] text-muted-foreground">
                                Email cannot be changed at this time.
                            </p>
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="avatarUrl">Avatar URL</Label>
                            <Input
                                id="avatarUrl"
                                value={editForm.avatarUrl}
                                onChange={(e) =>
                                    setEditForm({ ...editForm, avatarUrl: e.target.value })
                                }
                                placeholder="https://example.com/avatar.jpg"
                            />
                        </div>
                        <div className="space-y-2">
                            <Label htmlFor="description">Description</Label>
                            <Textarea
                                id="description"
                                value={editForm.description}
                                onChange={(e) =>
                                    setEditForm({ ...editForm, description: e.target.value })
                                }
                                placeholder="Tell us about yourself..."
                            />
                        </div>
                        <DialogFooter>
                            <Button
                                type="button"
                                variant="outline"
                                onClick={() => setIsEditDialogOpen(false)}
                            >
                                Cancel
                            </Button>
                            <Button type="submit" disabled={isUpdating}>
                                {isUpdating ? "Saving..." : "Save Changes"}
                            </Button>
                        </DialogFooter>
                    </form>
                </DialogContent>
            </Dialog>
        </>
    );
}
