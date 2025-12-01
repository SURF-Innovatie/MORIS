import { useState } from "react";
import {
    LayoutGrid,
    Table as TableIcon,
    Package,
    ExternalLink,
} from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { useGetProductsMe } from "@api/moris";
import { ProductType } from "@api/model";

const getProductTypeLabel = (type?: ProductType) => {
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

export const ProductList = () => {
    const [viewMode, setViewMode] = useState<"cards" | "table">("cards");

    const { data: products, isLoading, error } = useGetProductsMe();

    return (
        <section>
            <div className="mb-6 flex items-center justify-between">
                <div className="flex items-center gap-2">
                    <Package className="h-5 w-5 text-muted-foreground" />
                    <h2 className="text-2xl font-semibold tracking-tight">Products</h2>
                    {products && (
                        <Badge variant="outline" className="ml-2">
                            {products.length}
                        </Badge>
                    )}
                </div>
                <div className="flex items-center gap-2">
                    <div className="flex items-center rounded-lg border bg-muted/50 p-1">
                        <Button
                            variant={viewMode === "cards" ? "secondary" : "ghost"}
                            size="sm"
                            className="h-7 px-2"
                            onClick={() => setViewMode("cards")}
                        >
                            <LayoutGrid className="h-4 w-4" />
                        </Button>
                        <Button
                            variant={viewMode === "table" ? "secondary" : "ghost"}
                            size="sm"
                            className="h-7 px-2"
                            onClick={() => setViewMode("table")}
                        >
                            <TableIcon className="h-4 w-4" />
                        </Button>
                    </div>
                </div>
            </div>

            {isLoading && (
                <Card>
                    <CardContent className="flex items-center justify-center py-12">
                        <p className="text-sm text-muted-foreground">
                            Loading products...
                        </p>
                    </CardContent>
                </Card>
            )}

            {error && (
                <Card>
                    <CardContent className="flex items-center justify-center py-12">
                        <p className="text-sm text-destructive">
                            Failed to load products
                        </p>
                    </CardContent>
                </Card>
            )}

            {products && products.length === 0 && (
                <Card>
                    <CardContent className="flex flex-col items-center justify-center py-12 text-center">
                        <Package className="mb-4 h-12 w-12 text-muted-foreground/50" />
                        <p className="text-sm text-muted-foreground">No products found</p>
                    </CardContent>
                </Card>
            )}

            {products && products.length > 0 && viewMode === "cards" && (
                <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
                    {products.map((product) => (
                        <Card key={product.id} className="group flex flex-col transition-all hover:shadow-md hover:border-primary/20">
                            <CardHeader className="pb-3">
                                <div className="flex items-start justify-between gap-4">
                                    <div className="space-y-1">
                                        <CardTitle className="line-clamp-1 text-base">
                                            {product.name || "Untitled Product"}
                                        </CardTitle>
                                        <Badge variant="secondary" className="text-[10px] px-1.5 py-0 h-5">
                                            {getProductTypeLabel(product.type)}
                                        </Badge>
                                    </div>
                                    {product.doi && (
                                        <Button
                                            variant="ghost"
                                            size="icon"
                                            className="h-8 w-8 -mt-1 -mr-2 text-muted-foreground"
                                            asChild
                                        >
                                            <a href={`https://doi.org/${product.doi}`} target="_blank" rel="noreferrer">
                                                <ExternalLink className="h-4 w-4" />
                                            </a>
                                        </Button>
                                    )}
                                </div>
                                <CardDescription className="line-clamp-1 mt-2 text-xs">
                                    {product.doi || "No DOI available"}
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="pb-3 flex-1">
                                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                                    <span className="uppercase font-medium text-[10px] bg-muted px-1.5 py-0.5 rounded">
                                        {product.language || "Unknown Language"}
                                    </span>
                                </div>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}

            {products && products.length > 0 && viewMode === "table" && (
                <Card>
                    <CardContent className="p-0">
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>Name</TableHead>
                                    <TableHead>Type</TableHead>
                                    <TableHead>DOI</TableHead>
                                    <TableHead>Language</TableHead>
                                    <TableHead className="w-[50px]"></TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {products.map((product) => (
                                    <TableRow key={product.id}>
                                        <TableCell className="font-medium">
                                            {product.name || "Untitled Product"}
                                        </TableCell>
                                        <TableCell>
                                            <Badge variant="secondary" className="text-[10px] h-5">
                                                {getProductTypeLabel(product.type)}
                                            </Badge>
                                        </TableCell>
                                        <TableCell className="text-xs text-muted-foreground">
                                            {product.doi || "-"}
                                        </TableCell>
                                        <TableCell className="text-xs text-muted-foreground">
                                            {product.language || "-"}
                                        </TableCell>
                                        <TableCell>
                                            {product.doi && (
                                                <Button
                                                    variant="ghost"
                                                    size="icon"
                                                    className="h-8 w-8"
                                                    asChild
                                                >
                                                    <a href={`https://doi.org/${product.doi}`} target="_blank" rel="noreferrer">
                                                        <ExternalLink className="h-4 w-4" />
                                                    </a>
                                                </Button>
                                            )}
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    </CardContent>
                </Card>
            )}
        </section>
    );
};
