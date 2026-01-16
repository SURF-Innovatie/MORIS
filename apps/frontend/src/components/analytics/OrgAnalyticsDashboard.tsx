// import { useQuery } from "@tanstack/react-query"; // Unused now
import {
  useGetOrganisationsOrgIdAnalyticsSummary,
  useGetOrganisationsOrgIdAnalyticsByCategory,
  useGetOrganisationsOrgIdAnalyticsBurnRate,
  useGetOrganisationsOrgIdAnalyticsByProject,
  useGetOrganisationsOrgIdAnalyticsByFunding,
} from "@api/moris";
import { healthStatusEmoji, healthStatusLabels } from "@/lib/constants";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Loader2,
  TrendingUp,
  DollarSign,
  AlertTriangle,
  CheckCircle,
} from "lucide-react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  LineChart,
  Line,
  PieChart,
  Pie,
  Cell,
} from "recharts";

interface OrgAnalyticsDashboardProps {
  orgId: string;
}

const COLORS = [
  "#3b82f6",
  "#10b981",
  "#f59e0b",
  "#ef4444",
  "#8b5cf6",
  "#6b7280",
];

export function OrgAnalyticsDashboard({ orgId }: OrgAnalyticsDashboardProps) {
  const { data: summary, isLoading: summaryLoading } =
    useGetOrganisationsOrgIdAnalyticsSummary(orgId);

  const { data: categoryData } =
    useGetOrganisationsOrgIdAnalyticsByCategory(orgId);

  const { data: burnRateData } =
    useGetOrganisationsOrgIdAnalyticsBurnRate(orgId);

  const { data: projectHealth } =
    useGetOrganisationsOrgIdAnalyticsByProject(orgId);

  const { data: fundingData } =
    useGetOrganisationsOrgIdAnalyticsByFunding(orgId);

  if (summaryLoading) {
    return (
      <div className="flex items-center justify-center h-96">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <SummaryCard
          title="Total Projects"
          value={summary?.totalProjects ?? 0}
          icon={<CheckCircle className="h-4 w-4 text-green-500" />}
        />
        <SummaryCard
          title="Total Budgeted"
          value={formatCurrency(summary?.totalBudgeted ?? 0)}
          icon={<DollarSign className="h-4 w-4 text-blue-500" />}
        />
        <SummaryCard
          title="Total Spent"
          value={formatCurrency(summary?.totalActuals ?? 0)}
          icon={<TrendingUp className="h-4 w-4 text-green-500" />}
        />
        <SummaryCard
          title="Projects At Risk"
          value={summary?.projectsAtRisk ?? 0}
          icon={<AlertTriangle className="h-4 w-4 text-red-500" />}
          variant={summary?.projectsAtRisk ? "warning" : undefined}
        />
      </div>

      {/* Charts Row 1 */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Budget vs Actuals by Category */}
        <Card>
          <CardHeader>
            <CardTitle>Budget vs Actuals by Category</CardTitle>
            <CardDescription>
              Spending breakdown across expense categories
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-80">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={categoryData || []}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="category" />
                  <YAxis
                    tickFormatter={(v: number) => `€${(v / 1000).toFixed(0)}k`}
                  />
                  <Tooltip
                    formatter={(v: number) => formatCurrency(Number(v))}
                  />
                  <Legend />
                  <Bar dataKey="budgeted" fill="#3b82f6" name="Budgeted" />
                  <Bar dataKey="actuals" fill="#10b981" name="Actuals" />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        {/* Burn Rate Chart */}
        <Card>
          <CardHeader>
            <CardTitle>Cumulative Spend Over Time</CardTitle>
            <CardDescription>Budget consumption trend</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-80">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={burnRateData || []}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis
                    tickFormatter={(v: number) => `€${(v / 1000).toFixed(0)}k`}
                  />
                  <Tooltip
                    formatter={(v: number) => formatCurrency(Number(v))}
                  />
                  <Legend />
                  <Line
                    type="monotone"
                    dataKey="budgeted"
                    stroke="#3b82f6"
                    strokeDasharray="5 5"
                    name="Ideal Burn"
                  />
                  <Line
                    type="monotone"
                    dataKey="actual"
                    stroke="#10b981"
                    name="Actual Spend"
                    strokeWidth={2}
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Charts Row 2 */}
      <div className="grid gap-6 lg:grid-cols-2">
        {/* Funding Source Breakdown */}
        <Card>
          <CardHeader>
            <CardTitle>Funding Sources</CardTitle>
            <CardDescription>
              Budget allocation by funding source
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="h-80">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={fundingData || []}
                    dataKey="budgeted"
                    nameKey="fundingSource"
                    cx="50%"
                    cy="50%"
                    outerRadius={100}
                    label={({
                      name,
                      percent,
                    }: {
                      name: string;
                      percent: number;
                    }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                  >
                    {(fundingData || []).map((_, index) => (
                      <Cell
                        key={`cell-${index}`}
                        fill={COLORS[index % COLORS.length]}
                      />
                    ))}
                  </Pie>
                  <Tooltip
                    formatter={(v: number) => formatCurrency(Number(v))}
                  />
                  <Legend />
                </PieChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        {/* Project Health Table */}
        <Card>
          <CardHeader>
            <CardTitle>Project Health</CardTitle>
            <CardDescription>Status of all projects</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3 max-h-80 overflow-y-auto">
              {(projectHealth || []).length === 0 ? (
                <p className="text-center text-muted-foreground py-4">
                  No projects to display
                </p>
              ) : (
                projectHealth?.map((project) => (
                  <div
                    key={project.projectId}
                    className="flex items-center justify-between p-3 border rounded-lg hover:bg-muted/50 cursor-pointer"
                  >
                    <div>
                      <p className="font-medium">
                        {project.projectName ||
                          `Project ${(project.projectId || "").slice(0, 8)}`}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {formatCurrency(project.spent || 0)} /{" "}
                        {formatCurrency(project.budgeted || 0)}
                      </p>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="text-sm font-medium">
                        {(project.burnRate || 0).toFixed(1)}%
                      </span>
                      <span title={healthStatusLabels[project.status || ""]}>
                        {project.status && healthStatusEmoji[project.status]}
                      </span>
                    </div>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

interface SummaryCardProps {
  title: string;
  value: string | number;
  icon?: React.ReactNode;
  variant?: "default" | "warning";
}

function SummaryCard({ title, value, icon, variant }: SummaryCardProps) {
  return (
    <Card
      className={
        variant === "warning" ? "border-yellow-200 bg-yellow-50" : undefined
      }
    >
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        {icon}
      </CardHeader>
      <CardContent>
        <div
          className={`text-2xl font-bold ${variant === "warning" ? "text-yellow-600" : ""}`}
        >
          {value}
        </div>
      </CardContent>
    </Card>
  );
}

function formatCurrency(amount: number): string {
  return new Intl.NumberFormat("nl-NL", {
    style: "currency",
    currency: "EUR",
    maximumFractionDigits: 0,
  }).format(amount);
}
