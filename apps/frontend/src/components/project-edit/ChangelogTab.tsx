import { format } from "date-fns";
import { History, Loader2 } from "lucide-react";
import { useGetProjectsIdChangelog } from "@api/moris";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface ChangelogTabProps {
  projectId: string;
}

export function ChangelogTab({ projectId }: ChangelogTabProps) {
  const {
    data: changelog,
    isLoading,
    error,
  } = useGetProjectsIdChangelog(projectId);

  if (isLoading) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="flex justify-center">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center text-sm text-destructive">
            Failed to load changelog
          </div>
        </CardContent>
      </Card>
    );
  }

  const entries = changelog?.entries || [];

  // Mock helper to resolve UUIDs to names
  const resolveName = (id: string) => {
    // In a real app, we would look this up from a user cache or API
    // For now, we'll generate a deterministic fake name or return "System"
    if (!id) return "System";
    if (id.startsWith("user-")) return "Dr. Elaine Carter"; // Example
    return "System User";
  };

  // Group entries by date
  const groupedEntries = entries.reduce((acc, log) => {
    const date = format(new Date(log.at!), "yyyy-MM-dd");
    if (!acc[date]) acc[date] = [];
    acc[date].push(log);
    return acc;
  }, {} as Record<string, typeof entries>);

  return (
    <Card>
      <CardHeader>
        <CardTitle>Project History</CardTitle>
        <CardDescription>
          View the recent activity and changes in this project.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-8">
          {Object.entries(groupedEntries).map(([date, logs]) => (
            <div key={date}>
              <h4 className="mb-4 text-sm font-medium text-muted-foreground sticky top-0 bg-background py-2">
                {format(new Date(date), "EEEE, MMMM d, yyyy")}
              </h4>
              <div className="relative space-y-6 pl-6 before:absolute before:left-2 before:top-2 before:h-full before:w-px before:bg-border">
                {logs.map((log, index) => (
                  <div key={index} className="relative">
                    <div className="absolute -left-6 top-0.5 flex h-4 w-4 items-center justify-center rounded-full border bg-background ring-4 ring-background">
                      <History className="h-2.5 w-2.5 text-muted-foreground" />
                    </div>
                    <div className="flex flex-col gap-1">
                      <p className="text-sm leading-none">
                        <span className="font-semibold text-foreground">
                          {resolveName((log as any).by || "")}
                        </span>{" "}
                        <span className="text-muted-foreground">
                          {log.event?.toLowerCase().replace(/_/g, " ")}
                        </span>
                      </p>
                      <p className="text-xs text-muted-foreground/70">
                        {format(new Date(log.at!), "p")}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}
          {entries.length === 0 && (
            <div className="flex flex-col items-center justify-center py-8 text-center text-muted-foreground">
              <History className="mb-2 h-8 w-8 opacity-20" />
              <p>No history available.</p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
