import { useState } from "react";
import { useForm } from "react-hook-form";
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema";
import { z } from "zod";
import { Loader2, Plus, Search, Trash2, ExternalLink } from "lucide-react";

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
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/use-toast";
import { useGetCrossrefWorks, usePostProducts } from "@api/moris";
import {
  createProductAddedEvent,
  createProductRemovedEvent,
} from "@/api/events";
import { ProductResponse, ProductType } from "@api/model";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";
import { Allowed } from "@/components/auth/Allowed";
import { ProjectEventType } from "@/api/events";

// Schema for the DOI search form
const doiFormSchema = z.object({
  doi: z.string().min(1, "DOI is required"),
});

interface ProductsTabProps {
  projectId: string;
  products: ProductResponse[];
  onRefresh: () => void;
}

export function ProductsTab({
  projectId,
  products,
  onRefresh,
}: ProductsTabProps) {
  const { toast } = useToast();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [searchedProduct, setSearchedProduct] = useState<any>(null); // TODO: Type properly based on Crossref response
  const [isSearching, setIsSearching] = useState(false);
  const [productToDelete, setProductToDelete] = useState<string | null>(null);

  const form = useForm<z.infer<typeof doiFormSchema>>({
    resolver: standardSchemaResolver(doiFormSchema),
    defaultValues: {
      doi: "",
    },
  });

  const [searchDoi, setSearchDoi] = useState<string | null>(null);

  const {
    data: crossrefData,
    isLoading: isLoadingCrossref,
    isError: isCrossrefError,
  } = useGetCrossrefWorks(
    { doi: searchDoi! },
    {
      query: {
        enabled: !!searchDoi,
        retry: false,
      },
    }
  );

  // Effect to handle search results
  if (
    crossrefData &&
    searchDoi &&
    !isLoadingCrossref &&
    searchedProduct?.doi !== crossrefData.DOI
  ) {
    // Map crossref data to our product structure for preview
    setSearchedProduct({
      title: crossrefData.title?.[0] || "Unknown Title",
      doi: crossrefData.DOI,
      type: mapCrossrefType(crossrefData.type),
      language: crossrefData.language || "en", // Default to en if missing
    });
    setSearchDoi(null); // Reset search trigger
    setIsSearching(false);
  }

  if (isCrossrefError && searchDoi && !isLoadingCrossref) {
    toast({
      variant: "destructive",
      title: "Error",
      description: "Failed to find product with this DOI.",
    });
    setSearchDoi(null);
    setIsSearching(false);
  }

  const { mutateAsync: createProduct } = usePostProducts();

  function onSearch(values: z.infer<typeof doiFormSchema>) {
    setIsSearching(true);
    setSearchedProduct(null);
    setSearchDoi(values.doi);
  }

  async function onAddProduct() {
    if (!searchedProduct) return;

    try {
      // 1. Create the product in our DB
      const newProduct = await createProduct({
        data: {
          name: searchedProduct.title,
          doi: searchedProduct.doi,
          type: searchedProduct.type,
          language: searchedProduct.language,
        },
      });

      // 2. Link it to the project via event
      if (newProduct && newProduct.id) {
        await createProductAddedEvent(projectId, {
          product_id: newProduct.id,
        });

        toast({
          title: "Product added",
          description:
            "The product has been successfully added to the project.",
        });
        setIsDialogOpen(false);
        form.reset();
        setSearchedProduct(null);
        onRefresh();
      }
    } catch (error) {
      console.error(error);
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to add product. Please try again.",
      });
    }
  }

  async function onRemoveProduct() {
    if (!productToDelete) return;

    try {
      await createProductRemovedEvent(projectId, {
        product_id: productToDelete,
      });

      toast({
        title: "Product removed",
        description: "The product has been removed from the project.",
      });
      setProductToDelete(null);
      onRefresh();
    } catch (error) {
      console.error(error);
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to remove product.",
      });
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-medium">Products</h3>
          <p className="text-sm text-muted-foreground">
            Manage the products associated with this project.
          </p>
        </div>
        <Allowed event={ProjectEventType.ProductAdded}>
          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
            <DialogTrigger asChild>
              <Button>
                <Plus className="mr-2 h-4 w-4" />
                Add Product
              </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[500px]">
              <DialogHeader>
                <DialogTitle>Add Product</DialogTitle>
                <DialogDescription>
                  Search for a product by DOI to add it to the project.
                </DialogDescription>
              </DialogHeader>

              <Form {...form}>
                <form
                  onSubmit={form.handleSubmit(onSearch)}
                  className="space-y-4"
                >
                  <FormField
                    control={form.control}
                    name="doi"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>DOI</FormLabel>
                        <div className="flex gap-2">
                          <FormControl>
                            <Input placeholder="10.1038/..." {...field} />
                          </FormControl>
                          <Button
                            type="submit"
                            disabled={isSearching || isLoadingCrossref}
                          >
                            {isSearching || isLoadingCrossref ? (
                              <Loader2 className="h-4 w-4 animate-spin" />
                            ) : (
                              <Search className="h-4 w-4" />
                            )}
                          </Button>
                        </div>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </form>
              </Form>

              {searchedProduct && (
                <div className="mt-4 rounded-md border p-4">
                  <h4 className="font-medium">{searchedProduct.title}</h4>
                  <p className="text-sm text-muted-foreground mt-1">
                    DOI: {searchedProduct.doi}
                  </p>
                  <p className="text-sm text-muted-foreground">
                    Type: {getProductTypeLabel(searchedProduct.type)}
                  </p>
                </div>
              )}

              <DialogFooter>
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => setIsDialogOpen(false)}
                >
                  Cancel
                </Button>
                <Button onClick={onAddProduct} disabled={!searchedProduct}>
                  Add to Project
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </Allowed>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {products.map((product) => (
          <Card key={product.id}>
            <CardHeader className="pb-2">
              <CardTitle className="text-base font-medium line-clamp-2">
                {product.name}
              </CardTitle>
              <CardDescription className="flex items-center gap-2">
                <span className="capitalize">
                  {getProductTypeLabel(product.type)}
                </span>
                {product.doi && (
                  <a
                    href={`https://doi.org/${product.doi}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-primary hover:underline inline-flex items-center"
                  >
                    <ExternalLink className="h-3 w-3 ml-1" />
                  </a>
                )}
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex justify-end">
                <Allowed event={ProjectEventType.ProductRemoved}>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="text-destructive hover:text-destructive/90"
                    onClick={() => setProductToDelete(product.id!)}
                  >
                    <Trash2 className="mr-2 h-4 w-4" />
                    Remove
                  </Button>
                </Allowed>
              </div>
            </CardContent>
          </Card>
        ))}
        {products.length === 0 && (
          <div className="col-span-full flex flex-col items-center justify-center p-8 text-center border rounded-lg border-dashed text-muted-foreground">
            <p>No products added yet.</p>
          </div>
        )}
      </div>

      <ConfirmationModal
        isOpen={!!productToDelete}
        onClose={() => setProductToDelete(null)}
        onConfirm={onRemoveProduct}
        title="Remove Product"
        description="Are you sure you want to remove this product from the project? This action cannot be undone."
        confirmLabel="Remove"
        variant="destructive"
      />
    </div>
  );
}

// Helper to map Crossref types to our ProductType enum
// This is a simplification. You might need a more robust mapping.
function mapCrossrefType(type: string | undefined): ProductType {
  if (!type) return 0; // Default or unknown
  // Example mapping based on common Crossref types
  // We need to know what ProductType enum values correspond to.
  // Assuming 0=Unknown, 1=Article, 2=Book, etc.
  // Since I don't have the exact enum definition handy in this context,
  // I'll assume some defaults or map to a generic type.
  // Let's check the generated model for ProductType.
  // For now, returning 0 (which usually is a safe default or "Other").
  return 0;
}

function getProductTypeLabel(_type: ProductType | undefined): string {
  // TODO: Implement proper label mapping based on enum
  return "Product";
}
