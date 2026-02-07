import { Section } from "./types";
import { HeroSectionEditor, HeroSectionViewer } from "./sections/HeroSection";
import { RichTextEditor, RichTextViewer } from "./sections/RichTextSection";
import {
  ProjectListEditor,
  ProjectListViewer,
} from "./sections/ProjectListSection";
import {
  ProfileHeaderEditor,
  ProfileHeaderViewer,
} from "./sections/ProfileHeaderSection";
import {
  LinksSectionEditor,
  LinksSectionViewer,
} from "./sections/LinksSection";
import {
  StatisticsSectionEditor,
  StatisticsSectionViewer,
} from "./sections/StatisticsSection";

interface SectionRegistryProps {
  section: Section;
  mode: "edit" | "view";
  onChange?: (data: any) => void;
}

export function SectionRegistry({
  section,
  mode,
  onChange,
}: SectionRegistryProps) {
  const isEdit = mode === "edit";

  switch (section.type) {
    case "hero":
      return isEdit ? (
        <HeroSectionEditor data={section.data} onChange={onChange!} />
      ) : (
        <HeroSectionViewer data={section.data} />
      );
    case "rich_text":
      return isEdit ? (
        <RichTextEditor data={section.data} onChange={onChange!} />
      ) : (
        <RichTextViewer data={section.data} />
      );
    case "project_list":
      return isEdit ? (
        <ProjectListEditor data={section.data} onChange={onChange!} />
      ) : (
        <ProjectListViewer data={section.data} />
      );
    case "profile_header":
      return isEdit ? (
        <ProfileHeaderEditor data={section.data} onChange={onChange!} />
      ) : (
        <ProfileHeaderViewer data={section.data} />
      );
    case "links":
      return isEdit ? (
        <LinksSectionEditor data={section.data} onChange={onChange!} />
      ) : (
        <LinksSectionViewer data={section.data} />
      );
    case "statistics":
      return isEdit ? (
        <StatisticsSectionEditor data={section.data} onChange={onChange!} />
      ) : (
        <StatisticsSectionViewer data={section.data} />
      );
    default:
      return (
        <div className="p-4 text-red-500">
          Unknown section type: {section.type}
        </div>
      );
  }
}
