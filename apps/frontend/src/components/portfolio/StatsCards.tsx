import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface StatsCardsProps {
  projectCount: number;
  productCount: number;
  affiliationCount: number;
  highlightCount: number;
  isLoadingProjects: boolean;
  isLoadingProducts: boolean;
}

export const StatsCards = ({
  projectCount,
  productCount,
  affiliationCount,
  highlightCount,
  isLoadingProjects,
  isLoadingProducts,
}: StatsCardsProps) => {
  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      <Card>
        <CardHeader>
          <CardDescription>Projects</CardDescription>
          <CardTitle className="text-3xl">
            {isLoadingProjects ? "..." : projectCount}
          </CardTitle>
        </CardHeader>
      </Card>
      <Card>
        <CardHeader>
          <CardDescription>Deliverables</CardDescription>
          <CardTitle className="text-3xl">
            {isLoadingProducts ? "..." : productCount}
          </CardTitle>
        </CardHeader>
      </Card>
      <Card>
        <CardHeader>
          <CardDescription>Affiliations</CardDescription>
          <CardTitle className="text-3xl">{affiliationCount}</CardTitle>
        </CardHeader>
      </Card>
      <Card>
        <CardHeader>
          <CardDescription>Highlights</CardDescription>
          <CardTitle className="text-3xl">{highlightCount}</CardTitle>
        </CardHeader>
      </Card>
    </div>
  );
};
