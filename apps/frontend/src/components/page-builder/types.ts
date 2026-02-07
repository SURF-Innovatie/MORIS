export type SectionType =
  | "rich_text"
  | "hero"
  | "project_list"
  | "profile_header"
  | "links"
  | "statistics";

export interface Section {
  id: string;
  type: SectionType;
  data: Record<string, any>;
}

export interface PageContent {
  sections: Section[];
}
