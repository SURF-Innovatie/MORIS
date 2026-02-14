import { format } from "date-fns";
import { ProjectResponse } from "@api/model";

export function getProjectStatus(project: ProjectResponse) {
  if (!project.start_date || !project.end_date)
    return { label: "Unknown", variant: "secondary" as const };

  const now = new Date();
  const start = new Date(project.start_date);
  const end = new Date(project.end_date);

  if (now < start) return { label: "Upcoming", variant: "secondary" as const };
  if (now > end) return { label: "Completed", variant: "outline" as const };
  return { label: "Active", variant: "default" as const };
}

export function formatDateRange(project: ProjectResponse) {
  if (!project.start_date && !project.end_date) return "Timeline not set";
  const start = project.start_date
    ? format(new Date(project.start_date), "MMM yyyy")
    : "N/A";
  const end = project.end_date
    ? format(new Date(project.end_date), "MMM yyyy")
    : "Present";
  return `${start} · ${end}`;
}
