import { BookOpen, ExternalLink } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { ProjectResponse } from "@api/model";

interface ProjectProductListProps {
    project: ProjectResponse;
}

export function ProjectProductList({ project }: ProjectProductListProps) {
    return (
        <div className="space-y-4">
            <div className="flex items-center gap-2">
                <BookOpen className="h-5 w-5 text-muted-foreground" />
                <h2 className="text-lg font-semibold">Products</h2>
                <Badge variant="secondary" className="ml-2">
                    {project.products?.length || 0}
                </Badge>
            </div>
            <Card>
                <CardContent className="p-0">
                    {project.products && project.products.length > 0 ? (
                        <div className="divide-y">
                            {project.products.map((product) => (
                                <div
                                    key={product.id}
                                    className="flex items-center justify-between p-4"
                                >
                                    <div className="space-y-1">
                                        <p className="text-sm font-medium leading-none">
                                            {product.name}
                                        </p>
                                        <div className="flex items-center gap-2 text-xs text-muted-foreground">
                                            <span className="capitalize">
                                                Product
                                            </span>
                                            {product.doi && (
                                                <>
                                                    <span>•</span>
                                                    <a
                                                        href={`https://doi.org/${product.doi}`}
                                                        target="_blank"
                                                        rel="noopener noreferrer"
                                                        className="hover:underline inline-flex items-center gap-1"
                                                    >
                                                        {product.doi}
                                                        <ExternalLink className="h-3 w-3" />
                                                    </a>
                                                </>
                                            )}
                                        </div>
                                    </div>
                                </div>
                            ))}
                        </div>
                    ) : (
                        <div className="flex flex-col items-center justify-center py-8 text-center">
                            <BookOpen className="mb-2 h-8 w-8 text-muted-foreground/30" />
                            <p className="text-sm text-muted-foreground">
                                No products added yet.
                            </p>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
