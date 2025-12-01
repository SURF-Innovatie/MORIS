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

  return (
    <Card>
      <CardHeader>
        <CardTitle>Project History</CardTitle>
        <CardDescription>
          View the recent activity and changes in this project.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="relative space-y-8 pl-8 before:absolute before:left-3.5 before:top-2 before:h-full before:w-px before:bg-border">
          {entries.map((log, index) => (
            <div key={index} className="relative">
              <div className="absolute -left-8 top-1 flex h-7 w-7 items-center justify-center rounded-full border bg-background">
                <History className="h-3 w-3 text-muted-foreground" />
              </div>
              <div className="flex flex-col gap-1">
                <p className="text-sm font-medium leading-none">
                  <span className="font-semibold">System</span> {log.event}
                </p>
                <p className="text-xs text-muted-foreground">
                  {format(new Date(log.at!), "PPP p")}
                </p>
              </div>
            </div>
          ))}
          {entries.length === 0 && (
            <p className="text-sm text-muted-foreground">
              No history available.
            </p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
