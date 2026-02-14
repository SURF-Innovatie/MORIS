import { useState } from "react";
import { useForm } from "react-hook-form";
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema";
import { z } from "zod";
import { Loader2, Plus, Search, Upload, ChevronDown } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/use-toast";
import {
  useGetDoiResolve,
  usePostProducts,
  useGetZenodoStatus,
} from "@api/moris";
import {
  createProductAddedEvent,
  createProductRemovedEvent,
} from "@/api/events";
import { ProductResponse, ProductType, UploadType, Work } from "@api/model";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";
import { Allowed } from "@/components/auth/Allowed";
import { ProjectEventType } from "@/api/events";
import { ZenodoUploadDialog } from "@/components/products/ZenodoUploadDialog";
import { ProductCard, getProductTypeLabel } from "../products/ProductCard";
import { useAccess } from "@/contexts/AccessContext";
import { Doi } from "@/lib/doi";

// Schema for the DOI search form
const doiFormSchema = z.object({
  doi: z
    .string()
    .min(1, "DOI is required")
    .refine((val) => Doi.tryParse(val) !== null, {
      message: "Invalid DOI format",
    })
    .transform((val) => Doi.tryParse(val)?.toString() ?? val),
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
  const { hasAccess } = useAccess();
  const { mutateAsync: addProduct } = usePostProducts();
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [isZenodoDialogOpen, setIsZenodoDialogOpen] = useState(false);
  const [searchedProduct, setSearchedProduct] = useState<Work | null>(null);
  const [isSearching, setIsSearching] = useState(false);
  const [productToDelete, setProductToDelete] = useState<string | null>(null);

  const { data: zenodoStatus } = useGetZenodoStatus();

  const form = useForm<z.infer<typeof doiFormSchema>>({
    resolver: standardSchemaResolver(doiFormSchema),
    defaultValues: {
      doi: "",
    },
  });

  const [searchDoi, setSearchDoi] = useState<string | null>(null);

  const {
    data: doiData,
    isLoading: isLoadingDoi,
    isError: isDoiError,
  } = useGetDoiResolve(
    { doi: searchDoi! },
    {
      query: {
        enabled: !!searchDoi,
        retry: false,
      },
    },
  );

  // Effect to handle search results
  if (
    doiData &&
    searchDoi &&
    !isLoadingDoi &&
    searchedProduct?.doi !== doiData.doi
  ) {
    setSearchedProduct(doiData);
    setSearchDoi(null); // Reset search trigger
    setIsSearching(false);
  }

  if (isDoiError && searchDoi && !isLoadingDoi) {
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
          language: "en",
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

  async function handleZenodoUploadSuccess(
    doi: string,
    _zenodoUrl: string,
    depositionId: number,
    title: string,
    uploadType: UploadType,
  ) {
    try {
      // Create the product in our DB with the DOI from Zenodo
      const newProduct = await addProduct({
        data: {
          name: title, // Use user-provided title from upload dialog
          doi: doi,
          type: mapZenodoType(uploadType),
          language: "en",
          zenodo_deposition_id: depositionId,
        },
      });

      if (newProduct && newProduct.id) {
        await createProductAddedEvent(projectId, {
          product_id: newProduct.id,
        });

        toast({
          title: "Product uploaded to Zenodo",
          description:
            "The product has been published and added to the project.",
        });
        setIsZenodoDialogOpen(false);
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
          {zenodoStatus?.linked ? (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button>
                  <Plus className="mr-2 h-4 w-4" />
                  Add Product
                  <ChevronDown className="ml-2 h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => setIsDialogOpen(true)}>
                  <Search className="mr-2 h-4 w-4" />
                  Import from DOI
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => setIsZenodoDialogOpen(true)}>
                  <Upload className="mr-2 h-4 w-4" />
                  Upload to Zenodo
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          ) : (
            <Button onClick={() => setIsDialogOpen(true)}>
              <Plus className="mr-2 h-4 w-4" />
              Add Product
            </Button>
          )}

          <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
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
                            disabled={isSearching || isLoadingDoi}
                          >
                            {isSearching || isLoadingDoi ? (
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
        <ZenodoUploadDialog
          open={isZenodoDialogOpen}
          onOpenChange={setIsZenodoDialogOpen}
          onSuccess={handleZenodoUploadSuccess}
        />
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {products.map((product) => (
          <ProductCard
            key={product.id}
            product={product}
            onRemove={(id) => setProductToDelete(id)}
            canRemove={hasAccess(ProjectEventType.ProductRemoved)}
            pending={(product as any).pending}
          />
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

// Helper to map Zenodo UploadType to ProductType
function mapZenodoType(uploadType: UploadType): ProductType {
  switch (uploadType) {
    case UploadType.UploadTypeDataset:
      return ProductType.Dataset;
    case UploadType.UploadTypeSoftware:
      return ProductType.Software;
    case UploadType.UploadTypeImage:
      return ProductType.Image;
    case UploadType.UploadTypeVideo:
      return ProductType.Sound; // Closest match, could also be Other
    case UploadType.UploadTypeLesson:
      return ProductType.LearningObject;
    case UploadType.UploadTypePublication:
    case UploadType.UploadTypePoster:
    case UploadType.UploadTypePresentation:
    case UploadType.UploadTypePhysicalObject:
    case UploadType.UploadTypeOther:
    default:
      return ProductType.Other;
  }
}
