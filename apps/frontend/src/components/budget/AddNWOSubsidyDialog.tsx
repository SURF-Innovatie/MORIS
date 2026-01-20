import { useState, useEffect } from "react";
import { Search, ExternalLink, Loader2 } from "lucide-react";
import NwoIcon from "@/components/icons/nwoIcon";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { useGetNwoProjects } from "@api/moris";
import type { Project } from "@api/model";

interface AddNWOSubsidyDialogProps {
  onSelect: (project: Project) => void;
  disabled?: boolean;
  isSubmitting?: boolean;
}

export function AddNWOSubsidyDialog({
  onSelect,
  disabled,
  isSubmitting,
}: AddNWOSubsidyDialogProps) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [debouncedQuery, setDebouncedQuery] = useState("");

  // Debounce the search query
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQuery(query);
    }, 400);
    return () => clearTimeout(timer);
  }, [query]);

  // Query NWO projects
  const {
    data: searchResults,
    isLoading,
    isFetching,
  } = useGetNwoProjects(
    { title: debouncedQuery, per_page: 10 },
    {
      query: {
        enabled: debouncedQuery.length >= 2,
      },
    },
  );

  const results = searchResults?.projects || [];

  const handleSelect = (project: Project) => {
    onSelect(project);
    setOpen(false);
    setQuery("");
  };

  // Reset query when dialog closes
  useEffect(() => {
    if (!open) {
      setQuery("");
      setDebouncedQuery("");
    }
  }, [open]);

  const formatCurrency = (amount?: number) => {
    if (!amount) return "";
    return new Intl.NumberFormat("nl-NL", {
      style: "currency",
      currency: "EUR",
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const showLoading = isLoading || isFetching;

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          type="button"
          variant="outline"
          size="sm"
          disabled={disabled || isSubmitting}
        >
          <NwoIcon width={16} height={16} className="mr-2" />
          Add NWO Subsidy
          {isSubmitting && <Loader2 className="ml-2 h-4 w-4 animate-spin" />}
        </Button>
      </DialogTrigger>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <NwoIcon width={20} height={20} />
            Add NWO Subsidy as Line Item
          </DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <p className="text-sm text-muted-foreground">
            Search for an NWO grant to create a new budget line item with
            pre-filled information.
          </p>

          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search by project title..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              className="pl-10"
              autoFocus
            />
            {showLoading && (
              <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 animate-spin text-muted-foreground" />
            )}
          </div>

          <ScrollArea className="h-[400px] border rounded-md">
            {/* Initial state - no search yet */}
            {debouncedQuery.length < 2 && (
              <div className="p-8 text-center text-muted-foreground">
                <NwoIcon
                  width={32}
                  height={32}
                  className="mx-auto mb-3 opacity-50"
                />
                <p>Enter at least 2 characters to search NWO grants</p>
              </div>
            )}

            {/* Loading state */}
            {showLoading && debouncedQuery.length >= 2 && (
              <div className="p-8 text-center text-muted-foreground">
                <Loader2 className="h-8 w-8 animate-spin mx-auto mb-3" />
                <p>Searching NWO database...</p>
              </div>
            )}

            {/* No results */}
            {results.length === 0 &&
              debouncedQuery.length >= 2 &&
              !showLoading && (
                <div className="p-8 text-center text-muted-foreground">
                  <p>No grants found for "{debouncedQuery}"</p>
                </div>
              )}

            {/* Results list */}
            {!showLoading &&
              results.map((project) => (
                <button
                  key={project.project_id}
                  type="button"
                  onClick={() => handleSelect(project)}
                  className={cn(
                    "w-full text-left p-4 border-b hover:bg-muted/50 transition-colors",
                    "focus:outline-none focus:bg-muted/50",
                  )}
                >
                  <div className="space-y-1">
                    <div className="font-medium">{project.title}</div>
                    <div className="flex flex-wrap gap-2 text-xs text-muted-foreground">
                      <span className="font-mono">{project.project_id}</span>
                      {project.funding_scheme && (
                        <Badge variant="outline" className="text-xs">
                          {project.funding_scheme}
                        </Badge>
                      )}
                      {project.award_amount && (
                        <Badge
                          variant="secondary"
                          className="text-xs font-semibold"
                        >
                          {formatCurrency(project.award_amount)}
                        </Badge>
                      )}
                    </div>
                    {project.project_members &&
                      project.project_members.length > 0 && (
                        <div className="text-xs text-muted-foreground">
                          {project.project_members
                            .slice(0, 2)
                            .map(
                              (m) =>
                                `${m.first_name || ""} ${m.last_name || ""}`,
                            )
                            .join(", ")}
                          {project.project_members.length > 2 &&
                            ` +${project.project_members.length - 2} more`}
                        </div>
                      )}
                  </div>
                </button>
              ))}
          </ScrollArea>

          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <span>Data from NWO Open API</span>
            <a
              href="https://data.nwo.nl/en/how-to-use-the-nwopen-api"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-1 hover:underline"
            >
              data.nwo.nl
              <ExternalLink className="h-3 w-3" />
            </a>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
