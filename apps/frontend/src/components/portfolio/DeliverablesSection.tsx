import { Package } from "lucide-react";

import { ProductResponse } from "@api/model";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ProductCard } from "@/components/products/ProductCard";

interface DeliverablesSectionProps {
  featuredProducts: ProductResponse[];
  isLoadingProducts: boolean;
  hasProducts: boolean;
}

export const DeliverablesSection = ({
  featuredProducts,
  isLoadingProducts,
  hasProducts,
}: DeliverablesSectionProps) => {
  return (
    <Card className="lg:col-span-2">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Package className="h-5 w-5 text-primary" />
          Deliverables
        </CardTitle>
        <CardDescription>
          Publications, software, datasets, and other research outputs.
        </CardDescription>
      </CardHeader>
      <CardContent>
        {isLoadingProducts ? (
          <div className="py-6 text-center text-sm text-muted-foreground">
            Loading deliverables...
          </div>
        ) : hasProducts ? (
          <div className="grid gap-4 md:grid-cols-2">
            {featuredProducts.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        ) : (
          <div className="rounded-lg border border-dashed p-6 text-center text-sm text-muted-foreground">
            Add your first deliverable to start building your portfolio.
          </div>
        )}
      </CardContent>
    </Card>
  );
};
