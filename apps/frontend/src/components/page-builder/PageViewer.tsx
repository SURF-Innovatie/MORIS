import { Section } from "./types";
import { SectionRegistry } from "./SectionRegistry";

interface PageViewerProps {
  sections: Section[];
}

export function PageViewer({ sections }: PageViewerProps) {
  if (!sections || sections.length === 0) {
    return (
      <div className="text-center py-20 text-slate-400">
        This page has no content yet.
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-0 w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {sections.map((section) => (
        <div key={section.id} className="w-full">
          <SectionRegistry section={section} mode="view" />
        </div>
      ))}
    </div>
  );
}
