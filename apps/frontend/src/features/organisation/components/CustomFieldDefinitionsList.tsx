import { Loader2, Trash2, Info } from "lucide-react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import {
  useGetOrganisationNodesIdCustomFields,
  useDeleteOrganisationNodesIdCustomFieldsFieldId,
} from "@api/moris";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { AddCustomFieldDialog } from "./AddCustomFieldDialog";

interface CustomFieldDefinitionsListProps {
  nodeId: string;
}

export const CustomFieldDefinitionsList = ({
  nodeId,
}: CustomFieldDefinitionsListProps) => {
  const {
    data: fields,
    isLoading,
    refetch,
  } = useGetOrganisationNodesIdCustomFields(nodeId);
  const { mutateAsync: deleteField } =
    useDeleteOrganisationNodesIdCustomFieldsFieldId();

  const handleDelete = async (fieldId: string) => {
    if (
      !confirm(
        "Are you sure you want to delete this field? Data in projects or member profiles might be lost.",
      )
    ) {
      return;
    }
    try {
      await deleteField({ id: nodeId, fieldId });
      toast.success("Custom field deleted");
      refetch();
    } catch (error) {
      console.error("Failed to delete field", error);
      toast.error("Error deleting field");
    }
  };

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="text-lg font-medium">Custom Fields</h3>
        <AddCustomFieldDialog nodeId={nodeId} onSuccess={refetch} />
      </div>

      <div className="border rounded-lg overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[200px]">Name</TableHead>
              <TableHead>Category</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Required</TableHead>
              <TableHead>Description</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={6} className="h-24 text-center">
                  <div className="flex items-center justify-center">
                    <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                  </div>
                </TableCell>
              </TableRow>
            ) : fields?.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={6}
                  className="h-24 text-center text-muted-foreground"
                >
                  No custom fields defined.
                </TableCell>
              </TableRow>
            ) : (
              fields?.map((field: any) => (
                <TableRow key={field.id}>
                  <TableCell className="font-medium">
                    <div className="flex items-center gap-2">
                      {field.name}
                      {field.validation_regex && (
                        <TooltipProvider>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Info className="h-3 w-3 text-muted-foreground cursor-help" />
                            </TooltipTrigger>
                            <TooltipContent>
                              <p className="text-xs">
                                Regex: {field.validation_regex}
                              </p>
                            </TooltipContent>
                          </Tooltip>
                        </TooltipProvider>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <span
                      className={`text-xs font-mono rounded px-2 py-1 inline-block ${field.category === "PERSON" ? "bg-blue-100 text-blue-800" : "bg-gray-100 text-gray-800"}`}
                    >
                      {field.category || "PROJECT"}
                    </span>
                  </TableCell>
                  <TableCell className="text-xs font-mono">
                    {field.type}
                  </TableCell>
                  <TableCell>{field.required ? "Yes" : "No"}</TableCell>
                  <TableCell
                    className="text-muted-foreground max-w-[300px] truncate"
                    title={field.description}
                  >
                    {field.description}
                  </TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() => handleDelete(field.id)}
                      title="Delete Field"
                    >
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
      <p className="text-xs text-muted-foreground mt-2">
        Note: These fields will be available within this organisation and its
        sub-organisations.
      </p>
    </div>
  );
};
