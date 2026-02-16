import { useState } from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Pencil, MoreHorizontal, Trash, Ban, CheckCircle } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useQueryClient } from "@tanstack/react-query";
import {
  useGetAdminUsersList,
  useDeleteUsersId,
  usePostAdminUsersIdToggleActive,
  getGetAdminUsersListQueryKey,
} from "@api/moris";
import { toast } from "sonner";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";
import { UserDialog } from "@/components/user/UserDialog";
import { UserResponse } from "@api/model";

export const AdminUsersRoute = () => {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const {
    data: response,
    isLoading,
    error,
  } = useGetAdminUsersList({ page, page_size: pageSize });
  const { mutateAsync: deleteUser } = useDeleteUsersId();
  const { mutateAsync: toggleActive } = usePostAdminUsersIdToggleActive();
  const queryClient = useQueryClient();

  const users = response?.data;
  const totalPages = response?.total_pages || 1;

  const [userToDelete, setUserToDelete] = useState<string | null>(null);
  const [selectedUser, setSelectedUser] = useState<UserResponse | null>(null);
  const [isUserDialogOpen, setIsUserDialogOpen] = useState(false);

  const handleCreateUser = () => {
    setSelectedUser(null);
    setIsUserDialogOpen(true);
  };

  const handleEditUser = (user: UserResponse) => {
    setSelectedUser(user);
    setIsUserDialogOpen(true);
  };

  const handleToggleActive = async (userId: string, currentStatus: boolean) => {
    try {
      await toggleActive({ id: userId, data: { is_active: !currentStatus } });
      await queryClient.invalidateQueries({
        queryKey: getGetAdminUsersListQueryKey(),
      });
      toast.success("Success", {
        description: `User ${currentStatus ? "deactivated" : "activated"} successfully`,
      });
    } catch (e: any) {
      console.error("Failed to toggle active status", e);
      toast.error("Error", {
        description:
          e.response?.data?.message || "Failed to update user status",
      });
    }
  };

  const confirmDelete = async () => {
    if (!userToDelete) return;
    try {
      await deleteUser({ id: userToDelete });
      await queryClient.invalidateQueries({
        queryKey: getGetAdminUsersListQueryKey(),
      });
      toast.success("Success", {
        description: "User deleted successfully",
      });
    } catch (e: any) {
      console.error("Failed to delete user", e);
      toast.error("Error", {
        description: e.response?.data?.message || "Failed to delete user",
      });
    } finally {
      setUserToDelete(null);
    }
  };

  if (isLoading) {
    return <div>Loading users...</div>;
  }

  if (error) {
    return <div className="text-red-500">Error loading users</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold tracking-tight">User Management</h1>
        <Button onClick={handleCreateUser}>Create User</Button>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Email</TableHead>
              <TableHead>Role</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="w-[100px]">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {users?.map((user) => (
              <TableRow key={user.id}>
                <TableCell className="font-medium">
                  {user.name || "N/A"}
                </TableCell>
                <TableCell>{user.email}</TableCell>
                <TableCell>
                  {user.is_sys_admin ? (
                    <Badge variant="default">SysAdmin</Badge>
                  ) : (
                    <Badge variant="outline">User</Badge>
                  )}
                </TableCell>
                <TableCell>
                  {user.is_active ? (
                    <Badge
                      variant="secondary"
                      className="bg-green-100 text-green-800 hover:bg-green-100/80"
                    >
                      Active
                    </Badge>
                  ) : (
                    <Badge variant="destructive">Inactive</Badge>
                  )}
                </TableCell>
                <TableCell>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" className="h-8 w-8 p-0">
                        <span className="sr-only">Open menu</span>
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuLabel>Actions</DropdownMenuLabel>
                      <DropdownMenuItem onClick={() => handleEditUser(user)}>
                        <Pencil className="mr-2 h-4 w-4" />
                        Edit
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() =>
                          handleToggleActive(user.id!, user.is_active || false)
                        }
                      >
                        {user.is_active ? (
                          <>
                            <Ban className="mr-2 h-4 w-4" />
                            Deactivate
                          </>
                        ) : (
                          <>
                            <CheckCircle className="mr-2 h-4 w-4" />
                            Activate
                          </>
                        )}
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        onClick={() => setUserToDelete(user.id!)}
                        className="text-red-600 focus:text-red-600"
                      >
                        <Trash className="mr-2 h-4 w-4" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <p className="text-sm font-medium">Rows per page</p>
          <Select
            value={`${pageSize}`}
            onValueChange={(value: string) => {
              setPageSize(Number(value));
              setPage(1);
            }}
          >
            <SelectTrigger className="h-8 w-[70px]">
              <SelectValue placeholder={pageSize} />
            </SelectTrigger>
            <SelectContent side="top">
              {[10, 20, 30, 40, 50].map((pageSize) => (
                <SelectItem key={pageSize} value={`${pageSize}`}>
                  {pageSize}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <Pagination className="w-auto mx-0">
          <PaginationContent>
            <PaginationItem>
              {page > 1 ? (
                <PaginationPrevious
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  className="cursor-pointer"
                  size="default"
                />
              ) : (
                <PaginationPrevious
                  aria-disabled
                  className="pointer-events-none opacity-50"
                  size="default"
                />
              )}
            </PaginationItem>

            <PaginationItem>
              <span className="text-sm text-muted-foreground px-4">
                Page {page} of {totalPages}
              </span>
            </PaginationItem>

            <PaginationItem>
              {page < totalPages ? (
                <PaginationNext
                  onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                  className="cursor-pointer"
                  size="default"
                />
              ) : (
                <PaginationNext
                  aria-disabled
                  className="pointer-events-none opacity-50"
                  size="default"
                />
              )}
            </PaginationItem>
          </PaginationContent>
        </Pagination>
      </div>
      <ConfirmationModal
        isOpen={!!userToDelete}
        onClose={() => setUserToDelete(null)}
        onConfirm={confirmDelete}
        title="Delete User"
        description="Are you sure you want to delete this user? This action cannot be undone."
        confirmLabel="Delete"
        variant="destructive"
      />
      <UserDialog
        open={isUserDialogOpen}
        onOpenChange={setIsUserDialogOpen}
        user={selectedUser}
      />
    </div>
  );
};
export default AdminUsersRoute;
