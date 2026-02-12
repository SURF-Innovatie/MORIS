import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { TipTapEditor } from "@/components/ui/tiptap-editor";
import { useToast } from "@/hooks/use-toast";
import { ArrowLeft, Loader2, Save } from "lucide-react";
import {
  useGetCatalogsId,
  usePostCatalogs,
  usePutCatalogsId,
  getGetCatalogsQueryKey,
  getGetCatalogsIdQueryKey,
} from "@api/moris";
import type { CreateRequest, UpdateRequest } from "@api/model";

interface CatalogFormData {
  name: string;
  title: string;
  description: string;
  rich_description: string;
  logo_url: string;
  primary_color: string;
  secondary_color: string;
  accent_color: string;
  font_family: string;
  favicon: string;
}

const EMPTY_FORM: CatalogFormData = {
  name: "",
  title: "",
  description: "",
  rich_description: "",
  logo_url: "",
  primary_color: "#008094",
  secondary_color: "#004c5a",
  accent_color: "#f59e0b",
  font_family: "",
  favicon: "",
};

export const AdminCatalogFormRoute = () => {
  const { id } = useParams<{ id: string }>();
  const isEdit = !!id;
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { toast } = useToast();
  const [form, setForm] = useState<CatalogFormData>(EMPTY_FORM);

  const { data: details, isLoading: isLoadingCatalog } = useGetCatalogsId(
    id!,
    { query: { enabled: !!id } },
  );

  useEffect(() => {
    const catalog = details?.catalog;
    if (catalog) {
      setForm({
        name: catalog.name || "",
        title: catalog.title || "",
        description: catalog.description || "",
        rich_description: catalog.rich_description || "",
        logo_url: catalog.logo_url || "",
        primary_color: catalog.primary_color || "#008094",
        secondary_color: catalog.secondary_color || "#004c5a",
        accent_color: catalog.accent_color || "#f59e0b",
        font_family: catalog.font_family || "",
        favicon: catalog.favicon || "",
      });
    }
  }, [details]);

  const { mutate: createCatalog, isPending: isCreating } = usePostCatalogs({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetCatalogsQueryKey(),
        });
        toast({
          title: "Success",
          description: "Catalog created successfully",
        });
        navigate("/dashboard/admin/catalogs");
      },
      onError: (e: any) => {
        toast({
          variant: "destructive",
          title: "Error",
          description: e.message || "Failed to create catalog",
        });
      },
    },
  });

  const { mutate: updateCatalog, isPending: isUpdating } = usePutCatalogsId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: getGetCatalogsQueryKey(),
        });
        if (id) {
          queryClient.invalidateQueries({
            queryKey: getGetCatalogsIdQueryKey(id),
          });
        }
        toast({
          title: "Success",
          description: "Catalog updated successfully",
        });
        navigate("/dashboard/admin/catalogs");
      },
      onError: (e: any) => {
        toast({
          variant: "destructive",
          title: "Error",
          description: e.message || "Failed to update catalog",
        });
      },
    },
  });

  const toPayload = (
    data: CatalogFormData,
  ): CreateRequest | UpdateRequest => ({
    name: data.name,
    title: data.title,
    description: data.description || undefined,
    rich_description: data.rich_description || undefined,
    logo_url: data.logo_url || undefined,
    primary_color: data.primary_color || undefined,
    secondary_color: data.secondary_color || undefined,
    accent_color: data.accent_color || undefined,
    font_family: data.font_family || undefined,
    favicon: data.favicon || undefined,
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name || !form.title) {
      toast({
        variant: "destructive",
        title: "Validation",
        description: "Name and title are required",
      });
      return;
    }
    if (isEdit) {
      updateCatalog({ id: id!, data: toPayload(form) as UpdateRequest });
    } else {
      createCatalog({ data: toPayload(form) as CreateRequest });
    }
  };

  const isSaving = isCreating || isUpdating;

  if (isEdit && isLoadingCatalog) {
    return <div>Loading catalog...</div>;
  }

  const updateField = (field: keyof CatalogFormData, value: string) =>
    setForm((prev) => ({ ...prev, [field]: value }));

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => navigate("/dashboard/admin/catalogs")}
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back
        </Button>
        <h1 className="text-2xl font-bold tracking-tight">
          {isEdit ? "Edit Catalog" : "New Catalog"}
        </h1>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Basic Info */}
        <Card>
          <CardHeader>
            <CardTitle>Basic Information</CardTitle>
            <CardDescription>
              Name is used as the URL slug, title is the display name.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="name">Name *</Label>
                <Input
                  id="name"
                  value={form.name}
                  onChange={(e) => updateField("name", e.target.value)}
                  placeholder="e.g. deepnl"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="title">Title *</Label>
                <Input
                  id="title"
                  value={form.title}
                  onChange={(e) => updateField("title", e.target.value)}
                  placeholder="e.g. DeepNL"
                />
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="description">Short Description</Label>
              <Textarea
                id="description"
                value={form.description}
                onChange={(e) => updateField("description", e.target.value)}
                placeholder="A brief description of the catalog"
                rows={3}
              />
            </div>
          </CardContent>
        </Card>

        {/* Rich Description */}
        <Card>
          <CardHeader>
            <CardTitle>Rich Description</CardTitle>
            <CardDescription>
              A detailed description with rich text formatting.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <TipTapEditor
              content={form.rich_description}
              onChange={(html) => updateField("rich_description", html)}
              placeholder="Write a detailed description..."
            />
          </CardContent>
        </Card>

        {/* Branding */}
        <Card>
          <CardHeader>
            <CardTitle>Branding</CardTitle>
            <CardDescription>
              Customize the look and feel of the catalog viewer.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="logo_url">Logo URL</Label>
                <Input
                  id="logo_url"
                  value={form.logo_url}
                  onChange={(e) => updateField("logo_url", e.target.value)}
                  placeholder="https://example.com/logo.png"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="favicon">Favicon URL</Label>
                <Input
                  id="favicon"
                  value={form.favicon}
                  onChange={(e) => updateField("favicon", e.target.value)}
                  placeholder="https://example.com/favicon.ico"
                />
              </div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="space-y-2">
                <Label htmlFor="primary_color">Primary Color</Label>
                <div className="flex items-center gap-2">
                  <input
                    type="color"
                    id="primary_color"
                    value={form.primary_color}
                    onChange={(e) =>
                      updateField("primary_color", e.target.value)
                    }
                    className="h-10 w-10 rounded border cursor-pointer"
                  />
                  <Input
                    value={form.primary_color}
                    onChange={(e) =>
                      updateField("primary_color", e.target.value)
                    }
                    className="flex-1"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="secondary_color">Secondary Color</Label>
                <div className="flex items-center gap-2">
                  <input
                    type="color"
                    id="secondary_color"
                    value={form.secondary_color}
                    onChange={(e) =>
                      updateField("secondary_color", e.target.value)
                    }
                    className="h-10 w-10 rounded border cursor-pointer"
                  />
                  <Input
                    value={form.secondary_color}
                    onChange={(e) =>
                      updateField("secondary_color", e.target.value)
                    }
                    className="flex-1"
                  />
                </div>
              </div>
              <div className="space-y-2">
                <Label htmlFor="accent_color">Accent Color</Label>
                <div className="flex items-center gap-2">
                  <input
                    type="color"
                    id="accent_color"
                    value={form.accent_color}
                    onChange={(e) =>
                      updateField("accent_color", e.target.value)
                    }
                    className="h-10 w-10 rounded border cursor-pointer"
                  />
                  <Input
                    value={form.accent_color}
                    onChange={(e) =>
                      updateField("accent_color", e.target.value)
                    }
                    className="flex-1"
                  />
                </div>
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="font_family">Font Family</Label>
              <Input
                id="font_family"
                value={form.font_family}
                onChange={(e) => updateField("font_family", e.target.value)}
                placeholder="e.g. Inter, Roboto, Open Sans"
              />
              <p className="text-xs text-muted-foreground">
                Google Fonts name. Leave empty for the default system font.
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Color Preview */}
        <Card>
          <CardHeader>
            <CardTitle>Preview</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-4">
              <div
                className="h-12 flex-1 rounded-md flex items-center justify-center text-white text-sm font-medium"
                style={{ backgroundColor: form.primary_color }}
              >
                Primary
              </div>
              <div
                className="h-12 flex-1 rounded-md flex items-center justify-center text-white text-sm font-medium"
                style={{ backgroundColor: form.secondary_color }}
              >
                Secondary
              </div>
              <div
                className="h-12 flex-1 rounded-md flex items-center justify-center text-sm font-medium"
                style={{ backgroundColor: form.accent_color }}
              >
                Accent
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Submit */}
        <div className="flex justify-end gap-3">
          <Button
            type="button"
            variant="outline"
            onClick={() => navigate("/dashboard/admin/catalogs")}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={isSaving}>
            {isSaving && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            <Save className="mr-2 h-4 w-4" />
            {isEdit ? "Update Catalog" : "Create Catalog"}
          </Button>
        </div>
      </form>
    </div>
  );
};

export default AdminCatalogFormRoute;
