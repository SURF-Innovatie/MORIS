import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";
import { useToast } from "@/hooks/use-toast";
import {
  Pencil,
  MoreHorizontal,
  Trash,
  Plus,
  ExternalLink,
} from "lucide-react";
import {
  useGetCatalogs,
  useDeleteCatalogsId,
  getGetCatalogsQueryKey,
} from "@api/moris";

export const AdminCatalogsRoute = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const [catalogToDelete, setCatalogToDelete] = useState<string | null>(null);

  const { data: catalogs, isLoading, error } = useGetCatalogs();

  const { mutateAsync: deleteCatalog, isPending: isDeleting } =
    useDeleteCatalogsId({
      mutation: {
        onSuccess: () => {
          queryClient.invalidateQueries({
            queryKey: getGetCatalogsQueryKey(),
          });
          toast({
            title: "Success",
            description: "Catalog deleted successfully",
          });
        },
        onError: (error: any) => {
          toast({
            variant: "destructive",
            title: "Error",
            description: error?.message || "Failed to delete catalog",
          });
        },
      },
    });

  const confirmDelete = async () => {
    if (!catalogToDelete) return;
    try {
      await deleteCatalog({ id: catalogToDelete });
    } finally {
      setCatalogToDelete(null);
    }
  };

  if (isLoading) return <div>Loading catalogs...</div>;
  if (error)
    return <div className="text-red-500">Error loading catalogs</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold tracking-tight">
          Catalog Management
        </h1>
        <Button onClick={() => navigate("/dashboard/admin/catalogs/new")}>
          <Plus className="mr-2 h-4 w-4" />
          New Catalog
        </Button>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Title</TableHead>
              <TableHead>Projects</TableHead>
              <TableHead>Colors</TableHead>
              <TableHead className="w-[100px]">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {catalogs && catalogs.length > 0 ? (
              catalogs.map((catalog) => (
                <TableRow key={catalog.id}>
                  <TableCell className="font-medium">
                    {catalog.name}
                  </TableCell>
                  <TableCell>{catalog.title}</TableCell>
                  <TableCell>
                    <Badge variant="secondary">
                      {catalog.project_ids?.length ?? 0} projects
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-1.5">
                      {catalog.primary_color && (
                        <div
                          className="h-5 w-5 rounded-full border border-black/10"
                          style={{ backgroundColor: catalog.primary_color }}
                          title={`Primary: ${catalog.primary_color}`}
                        />
                      )}
                      {catalog.secondary_color && (
                        <div
                          className="h-5 w-5 rounded-full border border-black/10"
                          style={{ backgroundColor: catalog.secondary_color }}
                          title={`Secondary: ${catalog.secondary_color}`}
                        />
                      )}
                      {catalog.accent_color && (
                        <div
                          className="h-5 w-5 rounded-full border border-black/10"
                          style={{ backgroundColor: catalog.accent_color }}
                          title={`Accent: ${catalog.accent_color}`}
                        />
                      )}
                    </div>
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
                        <DropdownMenuItem
                          onClick={() =>
                            navigate(
                              `/dashboard/admin/catalogs/${catalog.id}/edit`,
                            )
                          }
                        >
                          <Pencil className="mr-2 h-4 w-4" />
                          Edit
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() =>
                            window.open(`/catalog/${catalog.name}`, "_blank")
                          }
                        >
                          <ExternalLink className="mr-2 h-4 w-4" />
                          View
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          onClick={() => setCatalogToDelete(catalog.id!)}
                          className="text-red-600 focus:text-red-600"
                        >
                          <Trash className="mr-2 h-4 w-4" />
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className="h-24 text-center text-muted-foreground"
                >
                  No catalogs found. Create your first catalog to get started.
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <ConfirmationModal
        isOpen={!!catalogToDelete}
        onClose={() => setCatalogToDelete(null)}
        onConfirm={confirmDelete}
        title="Delete Catalog"
        description="Are you sure you want to delete this catalog? This action cannot be undone."
        confirmLabel="Delete"
        variant="destructive"
        isLoading={isDeleting}
      />
    </div>
  );
};

export default AdminCatalogsRoute;
