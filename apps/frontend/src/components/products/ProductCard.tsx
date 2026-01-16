import { ExternalLink, Trash2 } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ProductResponse, ProductType } from "@api/model";
import ZenodoIcon from "@/components/icons/zenodoIcon";

interface ProductCardProps {
  product: ProductResponse;
  onRemove?: (id: string) => void;
  canRemove?: boolean;
}

export const getProductTypeLabel = (type?: ProductType) => {
  if (type === undefined) return "Unknown";

  const typeMap: Record<number, string> = {
    0: "Cartographic Material",
    1: "Dataset",
    2: "Image",
    3: "Interactive Resource",
    4: "Learning Object",
    5: "Other",
    6: "Software",
    7: "Sound",
    8: "Trademark",
    9: "Workflow",
  };

  return typeMap[type as number] || "Unknown";
};

export function ProductCard({
  product,
  onRemove,
  canRemove,
  pending,
}: ProductCardProps & { pending?: boolean }) {
  return (
    <Card className={pending ? "opacity-70 border-dashed" : ""}>
      <CardHeader className="pb-2">
        <div className="flex justify-between items-start">
          <CardTitle className="text-base font-medium line-clamp-2">
            {product.name}
          </CardTitle>
          {pending && (
            <span className="text-[10px] font-semibold bg-yellow-100 text-yellow-800 px-2 py-0.5 rounded-full border border-yellow-200">
              Pending
            </span>
          )}
        </div>
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
        <div className="flex justify-between items-center">
          {product.zenodo_deposition_id && (
            <a
              href={`https://sandbox.zenodo.org/records/${product.zenodo_deposition_id}`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-primary"
              title="View on Zenodo"
            >
              <ZenodoIcon width={16} height={16} />
              <span className="text-xs">Zenodo</span>
            </a>
          )}
          <div className="flex-1" />
          {canRemove && onRemove && !pending && (
            <Button
              variant="ghost"
              size="sm"
              className="text-destructive hover:text-destructive/90"
              onClick={() => onRemove(product.id!)}
            >
              <Trash2 className="mr-2 h-4 w-4" />
              Remove
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
