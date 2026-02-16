import { useEffect, useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";

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
import { getProductTypeLabel } from "@/components/products/ProductCard";
import { PortfolioHeader } from "@/components/portfolio/PortfolioHeader";
import { PortfolioEditor } from "@/components/portfolio/PortfolioEditor";
import { StatsCards } from "@/components/portfolio/StatsCards";
import { ProjectHighlights } from "@/components/portfolio/ProjectHighlights";
import { DeliverableMix } from "@/components/portfolio/DeliverableMix";
import { DeliverablesSection } from "@/components/portfolio/DeliverablesSection";
import { AffiliationsCard } from "@/components/portfolio/AffiliationsCard";
import { ActivityHighlights } from "@/components/portfolio/ActivityHighlights";
import { useToast } from "@/hooks/use-toast";
import { useQueryClient } from "@tanstack/react-query";

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
      <PortfolioHeader
        user={user}
        headline={headline}
        summary={summary}
        website={website}
        showEmail={showEmail}
        showOrcid={showOrcid}
        isEditing={isEditing}
        onNavigateProjects={() => navigate("/dashboard/projects")}
        onNavigateSettings={() => navigate("/dashboard/settings")}
        onToggleEditing={() => setIsEditing((prev) => !prev)}
      />

      {isEditing && (
        <PortfolioEditor
          headline={headline}
          summary={summary}
          website={website}
          showEmail={showEmail}
          showOrcid={showOrcid}
          pinnedProjectIds={pinnedProjectIds}
          pinnedProductIds={pinnedProductIds}
          relevantProjects={relevantProjects}
          products={products}
          isLoadingProjects={isLoadingProjects}
          isLoadingProducts={isLoadingProducts}
          isLoadingPortfolio={isLoadingPortfolio}
          isUpdatingPortfolio={isUpdatingPortfolio}
          onHeadlineChange={setHeadline}
          onSummaryChange={setSummary}
          onWebsiteChange={setWebsite}
          onShowEmailChange={setShowEmail}
          onShowOrcidChange={setShowOrcid}
          onTogglePinnedProject={togglePinnedProject}
          onTogglePinnedProduct={togglePinnedProduct}
          onSave={handleSavePortfolio}
        />
      )}

      <StatsCards
        projectCount={relevantProjects.length}
        productCount={products?.length ?? 0}
        affiliationCount={memberships?.length ?? 0}
        highlightCount={eventHighlights.length}
        isLoadingProjects={isLoadingProjects}
        isLoadingProducts={isLoadingProducts}
      />

      <div className="grid gap-6 lg:grid-cols-3">
        <ProjectHighlights featuredProjects={featuredProjects} />
        <DeliverableMix productTypeCounts={productTypeCounts} />
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <DeliverablesSection
          featuredProducts={featuredProducts}
          isLoadingProducts={isLoadingProducts}
          hasProducts={!!products?.length}
        />
        <AffiliationsCard memberships={memberships} />
      </div>

      <ActivityHighlights eventHighlights={eventHighlights} />
    </div>
  );
};

export default PortfolioRoute;
