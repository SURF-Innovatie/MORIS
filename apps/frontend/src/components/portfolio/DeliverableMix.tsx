import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

interface DeliverableMixProps {
  productTypeCounts: [string, number][];
}

export const DeliverableMix = ({
  productTypeCounts,
}: DeliverableMixProps) => {
  return (
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
  );
};
