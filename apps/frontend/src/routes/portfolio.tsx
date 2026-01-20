import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { format } from "date-fns";
import {
  Building2,
  Calendar,
  ExternalLink,
  Folder,
  Mail,
  Package,
  Sparkles,
} from "lucide-react";

import {
  useGetOrganisationMembershipsMine,
  useGetPortfolioMe,
  useGetProductsMe,
  useGetProfile,
  useGetProjects,
  useGetUsersIdEventsApproved,
  usePutPortfolioMe,
  getGetPortfolioMeQueryKey,
} from "@api/moris";
import { ProjectResponse } from "@api/model";
import { ProjectEvent } from "@/api/events";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
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
import { ProductCard, getProductTypeLabel } from "@/components/products/ProductCard";
import { EventRenderer } from "@/components/events/EventRenderer";
import { useToast } from "@/hooks/use-toast";
import { useQueryClient } from "@tanstack/react-query";

const getProjectStatus = (project: ProjectResponse) => {
  if (!project.start_date || !project.end_date)
    return { label: "Unknown", variant: "secondary" as const };

  const now = new Date();
  const start = new Date(project.start_date);
  const end = new Date(project.end_date);

  if (now < start) return { label: "Upcoming", variant: "secondary" as const };
  if (now > end) return { label: "Completed", variant: "outline" as const };
  return { label: "Active", variant: "default" as const };
};

const formatDateRange = (project: ProjectResponse) => {
  if (!project.start_date && !project.end_date) return "Timeline not set";
  const start = project.start_date
    ? format(new Date(project.start_date), "MMM yyyy")
    : "N/A";
  const end = project.end_date
    ? format(new Date(project.end_date), "MMM yyyy")
    : "Present";
  return `${start} · ${end}`;
};

const PortfolioRoute = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { toast } = useToast();
  const { data: user, isLoading: isLoadingUser } = useGetProfile();
  const { data: portfolio, isLoading: isLoadingPortfolio } = useGetPortfolioMe();
  const { data: projects, isLoading: isLoadingProjects } = useGetProjects();
  const { data: products, isLoading: isLoadingProducts } = useGetProductsMe();
  const { data: memberships } = useGetOrganisationMembershipsMine();
  const { mutateAsync: updatePortfolio, isPending: isUpdatingPortfolio } =
    usePutPortfolioMe({
      mutation: {
        onSuccess: () => {
          queryClient.invalidateQueries({
            queryKey: getGetPortfolioMeQueryKey(),
          });
          toast({ title: "Portfolio updated" });
        },
        onError: () => {
          toast({
            title: "Failed to update portfolio",
            variant: "destructive",
          });
        },
      },
    });

  const [headline, setHeadline] = useState("");
  const [summary, setSummary] = useState("");
  const [website, setWebsite] = useState("");
  const [showEmail, setShowEmail] = useState(true);
  const [showOrcid, setShowOrcid] = useState(true);
  const [pinnedProjectIds, setPinnedProjectIds] = useState<string[]>([]);
  const [pinnedProductIds, setPinnedProductIds] = useState<string[]>([]);
  const [isEditing, setIsEditing] = useState(false);

  const { data: eventsData } = useGetUsersIdEventsApproved(user?.id ?? "", {
    query: {
      enabled: !!user?.id,
    },
  });

  const relevantProjects = useMemo(() => {
    if (!projects?.length) return [];
    if (!user?.id) return projects;

    const filtered = projects.filter((project) =>
      project.members?.some((member) => member.user_id === user.id)
    );

    return filtered.length ? filtered : projects;
  }, [projects, user?.id]);

  useEffect(() => {
    if (!portfolio) return;
    setHeadline(portfolio.headline ?? "");
    setSummary(portfolio.summary ?? "");
    setWebsite(portfolio.website ?? "");
    setShowEmail(portfolio.show_email ?? true);
    setShowOrcid(portfolio.show_orcid ?? true);
    setPinnedProjectIds(portfolio.pinned_project_ids ?? []);
    setPinnedProductIds(portfolio.pinned_product_ids ?? []);
  }, [portfolio]);

  const pinnedProjects = useMemo(() => {
    if (!pinnedProjectIds.length || !relevantProjects.length) return [];
    const projectMap = new Map(
      relevantProjects
        .filter((project) => project.id)
        .map((project) => [project.id!, project])
    );
    return pinnedProjectIds
      .map((id) => projectMap.get(id))
      .filter(Boolean) as ProjectResponse[];
  }, [pinnedProjectIds, relevantProjects]);

  const featuredProjects = useMemo(() => {
    const sorted = [...relevantProjects]
      .sort((a, b) => {
        const aDate = a.end_date || a.start_date || "";
        const bDate = b.end_date || b.start_date || "";
        return new Date(bDate).getTime() - new Date(aDate).getTime();
      });

    if (pinnedProjects.length === 0) {
      return sorted.slice(0, 4);
    }

    const pinnedIds = new Set(pinnedProjects.map((project) => project.id));
    const remaining = sorted.filter((project) => !pinnedIds.has(project.id));
    return [...pinnedProjects, ...remaining].slice(0, 4);
  }, [pinnedProjects, relevantProjects]);

  const productTypeCounts = useMemo(() => {
    if (!products?.length) return [];
    const counts = products.reduce<Record<string, number>>((acc, product) => {
      const label = getProductTypeLabel(product.type);
      acc[label] = (acc[label] ?? 0) + 1;
      return acc;
    }, {});

    return Object.entries(counts).sort((a, b) => b[1] - a[1]);
  }, [products]);

  const pinnedProducts = useMemo(() => {
    if (!pinnedProductIds.length || !products?.length) return [];
    const productMap = new Map(
      products
        .filter((product) => product.id)
        .map((product) => [product.id!, product])
    );
    return pinnedProductIds
      .map((id) => productMap.get(id))
      .filter(Boolean) as NonNullable<(typeof products)[number]>[];
  }, [pinnedProductIds, products]);

  const featuredProducts = useMemo(() => {
    if (!products?.length) return [];
    if (pinnedProducts.length === 0) return products.slice(0, 6);

    const pinnedIds = new Set(pinnedProducts.map((product) => product.id));
    const remaining = products.filter((product) => !pinnedIds.has(product.id));
    return [...pinnedProducts, ...remaining].slice(0, 6);
  }, [pinnedProducts, products]);

  const handleSavePortfolio = async () => {
    await updatePortfolio({
      data: {
        headline: headline.trim(),
        summary: summary.trim(),
        website: website.trim(),
        show_email: showEmail,
        show_orcid: showOrcid,
        pinned_project_ids: pinnedProjectIds,
        pinned_product_ids: pinnedProductIds,
      },
    });
  };

  const togglePinnedProject = (projectId: string) => {
    setPinnedProjectIds((prev) =>
      prev.includes(projectId)
        ? prev.filter((id) => id !== projectId)
        : [...prev, projectId]
    );
  };

  const togglePinnedProduct = (productId: string) => {
    setPinnedProductIds((prev) =>
      prev.includes(productId)
        ? prev.filter((id) => id !== productId)
        : [...prev, productId]
    );
  };

  const eventHighlights = (eventsData?.events ?? []).slice(0, 5) as ProjectEvent[];

  if (isLoadingUser) {
    return <div>Loading portfolio...</div>;
  }

  if (!user) {
    return <div>User not found</div>;
  }

  return (
    <div className="space-y-8">
      <div className="relative overflow-hidden rounded-3xl border bg-linear-to-br from-primary/10 via-background to-background p-8">
        <div className="absolute inset-0 opacity-40">
          <div className="h-full w-full bg-[radial-gradient(circle_at_top,var(--tw-gradient-stops))] from-primary/20 via-transparent to-transparent" />
        </div>
        <div className="relative flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
          <div className="flex items-start gap-6">
            <Avatar className="h-24 w-24 border">
              <AvatarImage src={user.avatarUrl || ""} alt={user.name} />
              <AvatarFallback className="text-2xl">
                {user.name
                  ?.split(" ")
                  .map((n) => n[0])
                  .join("")
                  .toUpperCase()
                  .slice(0, 2)}
              </AvatarFallback>
            </Avatar>
            <div className="space-y-3">
              <div>
                <p className="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">
                  Portfolio
                </p>
                <h1 className="text-3xl font-semibold tracking-tight">
                  {headline || user.name || "Anonymous Researcher"}
                </h1>
                <p className="text-sm text-muted-foreground">
                  {summary ||
                    user.description ||
                    "Showcasing projects, deliverables, and contributions across the research ecosystem."}
                </p>
              </div>
              <div className="flex flex-wrap items-center gap-3 text-sm text-muted-foreground">
                {showEmail && user.email && (
                  <span className="inline-flex items-center gap-2">
                    <Mail className="h-4 w-4" />
                    {user.email}
                  </span>
                )}
                {showOrcid && user.orcid && (
                  <a
                    href={`https://orcid.org/${user.orcid}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-2 text-primary hover:underline"
                  >
                    <ExternalLink className="h-4 w-4" />
                    ORCID {user.orcid}
                  </a>
                )}
                {website.trim() && (
                  <a
                    href={website.trim()}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-2 text-primary hover:underline"
                  >
                    <ExternalLink className="h-4 w-4" />
                    Website
                  </a>
                )}
              </div>
            </div>
          </div>
          <div className="flex flex-wrap gap-3">
            <Button variant="outline" onClick={() => navigate("/dashboard/projects")}
            >
              <Folder className="mr-2 h-4 w-4" />
              View Projects
            </Button>
            <Button variant="outline" onClick={() => navigate("/dashboard/settings")}>
              Edit Profile
            </Button>
            <Button onClick={() => setIsEditing((prev) => !prev)}>
              {isEditing ? "Close Editor" : "Edit Portfolio"}
            </Button>
          </div>
        </div>
      </div>

      {isEditing && (
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
                    onChange={(event) => setHeadline(event.target.value)}
                    placeholder="Lead with your strongest message"
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="summary">Summary</Label>
                  <Textarea
                    id="summary"
                    value={summary}
                    onChange={(event) => setSummary(event.target.value)}
                    placeholder="A short overview of your research focus"
                    rows={4}
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="website">Website</Label>
                  <Input
                    id="website"
                    value={website}
                    onChange={(event) => setWebsite(event.target.value)}
                    placeholder="https://your-site.example"
                  />
                </div>
                <div className="flex flex-wrap gap-6">
                  <div className="flex items-center gap-2">
                    <Switch
                      id="show-email"
                      checked={showEmail}
                      onCheckedChange={setShowEmail}
                    />
                    <Label htmlFor="show-email">Show email</Label>
                  </div>
                  <div className="flex items-center gap-2">
                    <Switch
                      id="show-orcid"
                      checked={showOrcid}
                      onCheckedChange={setShowOrcid}
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
                                togglePinnedProject(projectId)
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
                                togglePinnedProduct(productId)
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
                onClick={handleSavePortfolio}
                disabled={isUpdatingPortfolio || isLoadingPortfolio}
              >
                {isUpdatingPortfolio ? "Saving..." : "Save portfolio"}
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader>
            <CardDescription>Projects</CardDescription>
            <CardTitle className="text-3xl">
              {isLoadingProjects ? "…" : relevantProjects.length}
            </CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader>
            <CardDescription>Deliverables</CardDescription>
            <CardTitle className="text-3xl">
              {isLoadingProducts ? "…" : products?.length ?? 0}
            </CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader>
            <CardDescription>Affiliations</CardDescription>
            <CardTitle className="text-3xl">
              {memberships?.length ?? 0}
            </CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader>
            <CardDescription>Highlights</CardDescription>
            <CardTitle className="text-3xl">
              {eventHighlights.length}
            </CardTitle>
          </CardHeader>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Sparkles className="h-5 w-5 text-primary" />
              Project Highlights
            </CardTitle>
            <CardDescription>
              Featured work showing recent contributions and outcomes.
            </CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4 md:grid-cols-2">
            {featuredProjects.length === 0 ? (
              <div className="col-span-full rounded-lg border border-dashed p-6 text-center text-sm text-muted-foreground">
                No projects to showcase yet.
              </div>
            ) : (
              featuredProjects.map((project) => {
                const status = getProjectStatus(project);
                return (
                  <div
                    key={project.id}
                    className="rounded-xl border bg-background p-4 shadow-sm"
                  >
                    <div className="flex items-start justify-between gap-2">
                      <div>
                        <h3 className="text-base font-semibold">
                          {project.title || "Untitled Project"}
                        </h3>
                        <p className="text-xs text-muted-foreground">
                          {formatDateRange(project)}
                        </p>
                      </div>
                      <Badge variant={status.variant}>{status.label}</Badge>
                    </div>
                    <p className="mt-3 line-clamp-3 text-sm text-muted-foreground">
                      {project.description ||
                        "No description added yet for this project."}
                    </p>
                    <div className="mt-4 flex flex-wrap gap-2 text-xs text-muted-foreground">
                      {project.owning_org_node?.name && (
                        <span className="inline-flex items-center gap-1">
                          <Building2 className="h-3 w-3" />
                          {project.owning_org_node.name}
                        </span>
                      )}
                      <span className="inline-flex items-center gap-1">
                        <Package className="h-3 w-3" />
                        {project.products?.length ?? 0} deliverables
                      </span>
                      {project.start_date && (
                        <span className="inline-flex items-center gap-1">
                          <Calendar className="h-3 w-3" />
                          {format(new Date(project.start_date), "MMM yyyy")}
                        </span>
                      )}
                    </div>
                  </div>
                );
              })
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Deliverable Mix</CardTitle>
            <CardDescription>
              A snapshot of your most common output types.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {productTypeCounts.length === 0 ? (
              <div className="rounded-lg border border-dashed p-6 text-center text-sm text-muted-foreground">
                No deliverables found.
              </div>
            ) : (
              productTypeCounts.map(([label, count]) => (
                <div
                  key={label}
                  className="flex items-center justify-between rounded-lg border px-3 py-2 text-sm"
                >
                  <span>{label}</span>
                  <Badge variant="outline">{count}</Badge>
                </div>
              ))
            )}
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
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
            ) : products?.length ? (
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

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Building2 className="h-5 w-5 text-primary" />
              Affiliations
            </CardTitle>
            <CardDescription>
              Your current organizational memberships and roles.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            {memberships?.length ? (
              memberships.map((membership) => (
                <div
                  key={membership.membershipId}
                  className="rounded-lg border px-4 py-3"
                >
                  <p className="text-sm font-medium">
                    {membership.scopeRootOrganisation?.name ||
                      "Organisation unavailable"}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    Role: {membership.roleKey || "Member"}
                  </p>
                </div>
              ))
            ) : (
              <div className="rounded-lg border border-dashed p-6 text-center text-sm text-muted-foreground">
                No organisation memberships yet.
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Folder className="h-5 w-5 text-primary" />
            Activity Highlights
          </CardTitle>
          <CardDescription>
            Recent updates across projects, roles, and publications.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {eventHighlights.length ? (
            <div className="space-y-4">
              {eventHighlights.map((event) => (
                <div
                  key={event.id}
                  className="rounded-lg border px-4 py-3"
                >
                  <EventRenderer event={event} variant="compact" />
                </div>
              ))}
            </div>
          ) : (
            <div className="rounded-lg border border-dashed p-6 text-center text-sm text-muted-foreground">
              No recent activity to show yet.
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default PortfolioRoute;
