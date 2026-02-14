import { ExternalLink, Folder, Mail } from "lucide-react";

import { UserResponse } from "@api/model";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";

interface PortfolioHeaderProps {
  user: UserResponse;
  headline: string;
  summary: string;
  website: string;
  showEmail: boolean;
  showOrcid: boolean;
  isEditing: boolean;
  onNavigateProjects: () => void;
  onNavigateSettings: () => void;
  onToggleEditing: () => void;
}

export const PortfolioHeader = ({
  user,
  headline,
  summary,
  website,
  showEmail,
  showOrcid,
  isEditing,
  onNavigateProjects,
  onNavigateSettings,
  onToggleEditing,
}: PortfolioHeaderProps) => {
  return (
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
          <Button variant="outline" onClick={onNavigateProjects}>
            <Folder className="mr-2 h-4 w-4" />
            View Projects
          </Button>
          <Button variant="outline" onClick={onNavigateSettings}>
            Edit Profile
          </Button>
          <Button onClick={onToggleEditing}>
            {isEditing ? "Close Editor" : "Edit Portfolio"}
          </Button>
        </div>
      </div>
    </div>
  );
};
