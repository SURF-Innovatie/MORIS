import { ProjectResponse, ProductResponse } from "@api/model";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import { getProductTypeLabel } from "@/components/products/ProductCard";
import { formatDateRange } from "@/lib/format";

interface PortfolioEditorProps {
  headline: string;
  summary: string;
  website: string;
  showEmail: boolean;
  showOrcid: boolean;
  pinnedProjectIds: string[];
  pinnedProductIds: string[];
  relevantProjects: ProjectResponse[];
  products: ProductResponse[] | undefined;
  isLoadingProjects: boolean;
  isLoadingProducts: boolean;
  isLoadingPortfolio: boolean;
  isUpdatingPortfolio: boolean;
  onHeadlineChange: (value: string) => void;
  onSummaryChange: (value: string) => void;
  onWebsiteChange: (value: string) => void;
  onShowEmailChange: (value: boolean) => void;
  onShowOrcidChange: (value: boolean) => void;
  onTogglePinnedProject: (projectId: string) => void;
  onTogglePinnedProduct: (productId: string) => void;
  onSave: () => void;
}

export const PortfolioEditor = ({
  headline,
  summary,
  website,
  showEmail,
  showOrcid,
  pinnedProjectIds,
  pinnedProductIds,
  relevantProjects,
  products,
  isLoadingProjects,
  isLoadingProducts,
  isLoadingPortfolio,
  isUpdatingPortfolio,
  onHeadlineChange,
  onSummaryChange,
  onWebsiteChange,
  onShowEmailChange,
  onShowOrcidChange,
  onTogglePinnedProject,
  onTogglePinnedProduct,
  onSave,
}: PortfolioEditorProps) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Customize your portfolio</CardTitle>
        <CardDescription>
          Pin standout projects and deliverables, add a headline, and
          control visibility for contact info.
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid gap-6 md:grid-cols-2">
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="headline">Headline</Label>
              <Input
                id="headline"
                value={headline}
                onChange={(event) => onHeadlineChange(event.target.value)}
                placeholder="Lead with your strongest message"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="summary">Summary</Label>
              <Textarea
                id="summary"
                value={summary}
                onChange={(event) => onSummaryChange(event.target.value)}
                placeholder="A short overview of your research focus"
                rows={4}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="website">Website</Label>
              <Input
                id="website"
                value={website}
                onChange={(event) => onWebsiteChange(event.target.value)}
                placeholder="https://your-site.example"
              />
            </div>
            <div className="flex flex-wrap gap-6">
              <div className="flex items-center gap-2">
                <Switch
                  id="show-email"
                  checked={showEmail}
                  onCheckedChange={onShowEmailChange}
                />
                <Label htmlFor="show-email">Show email</Label>
              </div>
              <div className="flex items-center gap-2">
                <Switch
                  id="show-orcid"
                  checked={showOrcid}
                  onCheckedChange={onShowOrcidChange}
                />
                <Label htmlFor="show-orcid">Show ORCID</Label>
              </div>
            </div>
          </div>

          <div className="space-y-4">
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <Label>Pinned projects</Label>
                <span className="text-xs text-muted-foreground">
                  {pinnedProjectIds.length} selected
                </span>
              </div>
              <div className="max-h-56 space-y-2 overflow-y-auto rounded-lg border p-3">
                {isLoadingProjects ? (
                  <div className="text-sm text-muted-foreground">
                    Loading projects...
                  </div>
                ) : relevantProjects.length ? (
                  relevantProjects.map((project) => {
                    const projectId = project.id;
                    if (!projectId) return null;
                    return (
                      <div
                        key={projectId}
                        className="flex items-start gap-3"
                      >
                        <Checkbox
                          checked={pinnedProjectIds.includes(projectId)}
                          onCheckedChange={() =>
                            onTogglePinnedProject(projectId)
                          }
                        />
                        <div>
                          <p className="text-sm font-medium">
                            {project.title || "Untitled Project"}
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {formatDateRange(project)}
                          </p>
                        </div>
                      </div>
                    );
                  })
                ) : (
                  <div className="text-sm text-muted-foreground">
                    No projects available yet.
                  </div>
                )}
              </div>
            </div>

            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <Label>Pinned deliverables</Label>
                <span className="text-xs text-muted-foreground">
                  {pinnedProductIds.length} selected
                </span>
              </div>
              <div className="max-h-56 space-y-2 overflow-y-auto rounded-lg border p-3">
                {isLoadingProducts ? (
                  <div className="text-sm text-muted-foreground">
                    Loading deliverables...
                  </div>
                ) : products?.length ? (
                  products.map((product) => {
                    const productId = product.id;
                    if (!productId) return null;
                    return (
                      <div
                        key={productId}
                        className="flex items-start gap-3"
                      >
                        <Checkbox
                          checked={pinnedProductIds.includes(productId)}
                          onCheckedChange={() =>
                            onTogglePinnedProduct(productId)
                          }
                        />
                        <div>
                          <p className="text-sm font-medium">
                            {product.name || "Untitled deliverable"}
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {getProductTypeLabel(product.type)}
                          </p>
                        </div>
                      </div>
                    );
                  })
                ) : (
                  <div className="text-sm text-muted-foreground">
                    No deliverables available yet.
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>

        <div className="flex justify-end">
          <Button
            onClick={onSave}
            disabled={isUpdatingPortfolio || isLoadingPortfolio}
          >
            {isUpdatingPortfolio ? "Saving..." : "Save portfolio"}
          </Button>
        </div>
      </CardContent>
    </Card>
  );
};
