import { useGetProjectsProjectIdBudget } from "@api/moris";
import { BudgetResponse as Budget } from "@api/model";
import { statusLabels } from "@/lib/constants";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Loader2,
  Plus,
  TrendingUp,
  DollarSign,
  AlertTriangle,
} from "lucide-react";

interface BudgetOverviewProps {
  projectId: string;
  onCreateBudget?: () => void;
  onEditBudget?: (budget: Budget) => void;
}

export function BudgetOverview({
  projectId,
  onCreateBudget,
  onEditBudget,
}: BudgetOverviewProps) {
  const {
    data: budget,
    isLoading,
    error,
  } = useGetProjectsProjectIdBudget(projectId, {
    query: {
      retry: false,
    },
  });

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-48">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !budget) {
    return (
      <Card className="border-dashed">
        <CardContent className="flex flex-col items-center justify-center py-12">
          <DollarSign className="h-12 w-12 text-muted-foreground mb-4" />
          <h3 className="text-lg font-semibold mb-2">No Budget Yet</h3>
          <p className="text-muted-foreground text-center mb-4">
            Create a budget to start tracking project expenditures
          </p>
          {onCreateBudget && (
            <Button onClick={onCreateBudget}>
              <Plus className="h-4 w-4 mr-2" />
              Create Budget
            </Button>
          )}
        </CardContent>
      </Card>
    );
  }

  const burnRate = budget.burnRate ?? 0;
  const burnRateColor =
    burnRate > 90
      ? "text-red-500"
      : burnRate > 75
        ? "text-yellow-500"
        : "text-green-500";

  return (
    <div className="space-y-6">
      {/* Budget Header */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              {budget.title || "Untitled Budget"}
              <StatusBadge status={budget.status || "draft"} />
            </CardTitle>
            <CardDescription>{budget.description}</CardDescription>
          </div>
          {onEditBudget && budget.status === "draft" && (
            <Button variant="outline" onClick={() => onEditBudget(budget)}>
              Edit Budget
            </Button>
          )}
        </CardHeader>
      </Card>

      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <SummaryCard
          title="Total Budgeted"
          value={formatCurrency(
            budget.totalBudgeted || 0,
            budget.currency || "EUR"
          )}
          icon={<DollarSign className="h-4 w-4" />}
        />
        <SummaryCard
          title="Total Spent"
          value={formatCurrency(
            budget.totalActuals || 0,
            budget.currency || "EUR"
          )}
          icon={<TrendingUp className="h-4 w-4" />}
        />
        <SummaryCard
          title="Remaining"
          value={formatCurrency(
            budget.remaining || 0,
            budget.currency || "EUR"
          )}
          icon={<DollarSign className="h-4 w-4" />}
          variant={(budget.remaining || 0) < 0 ? "destructive" : undefined}
        />
        <SummaryCard
          title="Burn Rate"
          value={`${burnRate.toFixed(1)}%`}
          icon={<TrendingUp className={`h-4 w-4 ${burnRateColor}`} />}
        />
      </div>

      {/* Progress Bar */}
      <Card>
        <CardHeader>
          <CardTitle className="text-sm font-medium">
            Budget Consumption
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            <Progress
              value={Math.min(burnRate, 100)}
              className={burnRate > 100 ? "bg-red-200" : undefined}
            />
            <div className="flex justify-between text-sm text-muted-foreground">
              <span>€{(budget.totalActuals || 0).toLocaleString()}</span>
              <span>€{(budget.totalBudgeted || 0).toLocaleString()}</span>
            </div>
          </div>
          {burnRate > 90 && (
            <div className="flex items-center gap-2 mt-4 text-yellow-600">
              <AlertTriangle className="h-4 w-4" />
              <span className="text-sm">Budget is nearly exhausted</span>
            </div>
          )}
          {burnRate > 100 && (
            <div className="flex items-center gap-2 mt-4 text-red-600">
              <AlertTriangle className="h-4 w-4" />
              <span className="text-sm">Budget has been exceeded!</span>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Line Items Summary */}
      <Card>
        <CardHeader>
          <CardTitle>Line Items ({(budget.lineItems || []).length})</CardTitle>
          <CardDescription>Budget breakdown by category</CardDescription>
        </CardHeader>
        <CardContent>
          {!budget.lineItems?.length ? (
            <p className="text-muted-foreground text-center py-4">
              No line items yet
            </p>
          ) : (
            <div className="space-y-4">
              {(budget.lineItems || []).slice(0, 5).map((item) => (
                <div
                  key={item.id}
                  className="flex items-center justify-between"
                >
                  <div>
                    <p className="font-medium">{item.description}</p>
                    <p className="text-sm text-muted-foreground">
                      {item.category} • {item.year}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="font-medium">
                      €{(item.budgetedAmount || 0).toLocaleString()}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      €{(item.totalActuals || 0).toLocaleString()} spent
                    </p>
                  </div>
                </div>
              ))}
              {(budget.lineItems || []).length > 5 && (
                <p className="text-sm text-muted-foreground text-center">
                  +{(budget.lineItems?.length || 0) - 5} more items
                </p>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

// Helper Components
function StatusBadge({ status }: { status: string }) {
  const variant =
    status === "approved"
      ? "default"
      : status === "locked"
        ? "secondary"
        : status === "submitted"
          ? "outline"
          : "outline";

  return <Badge variant={variant}>{statusLabels[status] || status}</Badge>;
}

interface SummaryCardProps {
  title: string;
  value: string;
  icon?: React.ReactNode;
  variant?: "default" | "destructive";
}

function SummaryCard({ title, value, icon, variant }: SummaryCardProps) {
  return (
    <Card
      className={
        variant === "destructive" ? "border-red-200 bg-red-50" : undefined
      }
    >
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div
          className={`text-2xl font-bold ${variant === "destructive" ? "text-red-600" : ""}`}
        >
          {value}
        </div>
      </CardContent>
    </Card>
  );
}

function formatCurrency(amount: number, currency: string = "EUR"): string {
  return new Intl.NumberFormat("nl-NL", {
    style: "currency",
    currency: currency,
  }).format(amount);
}
