import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { useState } from "react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

import { v4 as uuidv4 } from "uuid";
import { Button } from "@/components/ui/button";
import {
  Plus,
  Save,
  LayoutTemplate,
  Type,
  LayoutList,
  BarChart3,
  UserCircle,
  Link as LinkIcon,
  Eye,
} from "lucide-react";
import { Section, SectionType } from "./types";
import { SortableSection } from "./SortableSection";
import { SectionRegistry } from "./SectionRegistry";
import { PageViewer } from "./PageViewer";

interface PageBuilderProps {
  initialSections?: Section[];
  onSave: (sections: Section[]) => void;
}

export function PageBuilder({
  initialSections = [],
  onSave,
}: PageBuilderProps) {
  const [sections, setSections] = useState<Section[]>(initialSections);
  const [previewOpen, setPreviewOpen] = useState(false);

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;

    if (over && active.id !== over.id) {
      setSections((items) => {
        const oldIndex = items.findIndex((item) => item.id === active.id);
        const newIndex = items.findIndex((item) => item.id === over.id);
        return arrayMove(items, oldIndex, newIndex);
      });
    }
  }

  const addSection = (type: SectionType) => {
    const newSection: Section = {
      id: uuidv4(),
      type,
      data: {},
    };
    setSections([...sections, newSection]);
  };

  const updateSection = (id: string, data: any) => {
    setSections(sections.map((s) => (s.id === id ? { ...s, data } : s)));
  };

  const removeSection = (id: string) => {
    setSections(sections.filter((s) => s.id !== id));
  };

  return (
    <>
      <div className="flex gap-6">
        {/* Sidebar / Toolbox */}
        <div className="w-64 flex flex-col gap-2 p-4 border rounded-lg bg-slate-50 sticky top-4 h-fit">
          <h3 className="font-semibold text-slate-700 mb-2">Add Section</h3>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button className="w-full justify-start">
                <Plus className="mr-2 h-4 w-4" /> Add Section
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="w-56">
              <DropdownMenuLabel>Content Blocks</DropdownMenuLabel>
              <DropdownMenuItem onClick={() => addSection("hero")}>
                <LayoutTemplate className="w-4 h-4 mr-2" /> Hero Banner
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => addSection("rich_text")}>
                <Type className="w-4 h-4 mr-2" /> Rich Text
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuLabel>Dynamic Data</DropdownMenuLabel>
              <DropdownMenuItem onClick={() => addSection("project_list")}>
                <LayoutList className="w-4 h-4 mr-2" /> Project List
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => addSection("statistics")}>
                <BarChart3 className="w-4 h-4 mr-2" /> Statistics / KPIs
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuLabel>Profile</DropdownMenuLabel>
              <DropdownMenuItem onClick={() => addSection("profile_header")}>
                <UserCircle className="w-4 h-4 mr-2" /> Profile Header
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => addSection("links")}>
                <LinkIcon className="w-4 h-4 mr-2" /> Links / Socials
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <div className="mt-auto pt-4 border-t space-y-2">
            <Button
              variant="outline"
              onClick={() => setPreviewOpen(true)}
              className="w-full"
            >
              <Eye className="w-4 h-4 mr-2" /> Preview
            </Button>
            <Button onClick={() => onSave(sections)} className="w-full">
              <Save className="w-4 h-4 mr-2" /> Save Page
            </Button>
          </div>
        </div>

        {/* Canvas */}
        <div className="flex-1 min-h-[500px] border-2 border-dashed border-slate-200 rounded-lg p-6 bg-white">
          <DndContext
            sensors={sensors}
            collisionDetection={closestCenter}
            onDragEnd={handleDragEnd}
          >
            <SortableContext
              items={sections}
              strategy={verticalListSortingStrategy}
            >
              <div className="flex flex-col gap-4">
                {sections.map((section) => (
                  <SortableSection
                    key={section.id}
                    id={section.id}
                    onDelete={() => removeSection(section.id)}
                  >
                    <SectionRegistry
                      section={section}
                      mode="edit"
                      onChange={(data: Record<string, any>) =>
                        updateSection(section.id, data)
                      }
                    />
                  </SortableSection>
                ))}
              </div>
            </SortableContext>
          </DndContext>

          {sections.length === 0 && (
            <div className="h-full flex items-center justify-center text-slate-400">
              Drag items or click buttons to add sections
            </div>
          )}
        </div>
      </div>

      {/* Preview Dialog */}
      <Dialog open={previewOpen} onOpenChange={setPreviewOpen}>
        <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Page Preview</DialogTitle>
          </DialogHeader>
          <div className="border rounded-lg overflow-hidden bg-white">
            <PageViewer sections={sections} />
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}
