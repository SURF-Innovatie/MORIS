import { useState, useEffect } from "react";
import { Search, ExternalLink, X, Loader2 } from "lucide-react";
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
import { useGetNwoProjects, useGetNwoProjectProjectId } from "@api/moris";
import type { Project } from "@api/model";

interface NWOGrantSearchProps {
  value?: string | null;
  onChange: (grantId: string | null, project?: Project) => void;
  disabled?: boolean;
}

export function NWOGrantSearch({
  value,
  onChange,
  disabled,
}: NWOGrantSearchProps) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [debouncedQuery, setDebouncedQuery] = useState("");
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);

  // Debounce the search query
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQuery(query);
    }, 400);
    return () => clearTimeout(timer);
  }, [query]);

  // Query NWO projects
  const { data: searchResults, isLoading } = useGetNwoProjects(
    { title: debouncedQuery, per_page: 10 },
    {
      query: {
        enabled: debouncedQuery.length >= 2,
      },
    },
  );

  // Load project details when value is set
  const { data: projectDetails } = useGetNwoProjectProjectId(value || "", {
    query: {
      enabled: !!value && !selectedProject,
    },
  });

  // Update selected project when details are loaded
  useEffect(() => {
    if (projectDetails && !selectedProject) {
      setSelectedProject(projectDetails);
    }
  }, [projectDetails, selectedProject]);

  const results = searchResults?.projects || [];

  const handleSelect = (project: Project) => {
    setSelectedProject(project);
    onChange(project.project_id || null, project);
    setOpen(false);
    setQuery("");
  };

  const handleClear = () => {
    setSelectedProject(null);
    onChange(null);
  };

  const formatCurrency = (amount?: number) => {
    if (!amount) return "";
    return new Intl.NumberFormat("nl-NL", {
      style: "currency",
      currency: "EUR",
      maximumFractionDigits: 0,
    }).format(amount);
  };

  return (
    <div className="space-y-2">
      {value && selectedProject ? (
        <div className="flex items-center gap-2 p-3 border rounded-md bg-muted/50">
          <div className="flex-1 min-w-0">
            <div className="font-medium truncate text-sm">
              {selectedProject.title}
            </div>
            <div className="text-xs text-muted-foreground flex gap-2 items-center">
              <span>{selectedProject.project_id}</span>
              {selectedProject.award_amount && (
                <Badge variant="secondary" className="text-xs">
                  {formatCurrency(selectedProject.award_amount)}
                </Badge>
              )}
            </div>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon"
            onClick={handleClear}
            disabled={disabled}
          >
            <X className="h-4 w-4" />
          </Button>
        </div>
      ) : (
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button
              type="button"
              variant="outline"
              className="w-full justify-start text-muted-foreground"
              disabled={disabled}
            >
              <NwoIcon width={16} height={16} className="mr-2" />
              Link NWO Grant...
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl">
            <DialogHeader>
              <DialogTitle>Search NWO Grants</DialogTitle>
            </DialogHeader>
            <div className="space-y-4">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Search by project title..."
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  className="pl-10"
                  autoFocus
                />
                {isLoading && (
                  <Loader2 className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 animate-spin text-muted-foreground" />
                )}
              </div>

              <ScrollArea className="h-[400px] border rounded-md">
                {results.length === 0 &&
                  debouncedQuery.length >= 2 &&
                  !isLoading && (
                    <div className="p-8 text-center text-muted-foreground">
                      No projects found for "{debouncedQuery}"
                    </div>
                  )}
                {results.length === 0 && debouncedQuery.length < 2 && (
                  <div className="p-8 text-center text-muted-foreground">
                    Enter at least 2 characters to search
                  </div>
                )}
                {results.map((project) => (
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
                          <Badge variant="secondary" className="text-xs">
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
                  href="https://nwopen-api.nwo.nl"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-1 hover:underline"
                >
                  nwopen-api.nwo.nl
                  <ExternalLink className="h-3 w-3" />
                </a>
              </div>
            </div>
          </DialogContent>
        </Dialog>
      )}
    </div>
  );
}
