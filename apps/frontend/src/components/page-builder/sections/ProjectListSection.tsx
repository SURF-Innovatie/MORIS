import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface ProjectListData {
  limit?: number;
  filterTag?: string; // Example: "AI", "Quantum"
}

interface ProjectListEditorProps {
  data: ProjectListData;
  onChange: (data: ProjectListData) => void;
}

export function ProjectListEditor({ data, onChange }: ProjectListEditorProps) {
  return (
    <div className="space-y-4 p-4 border rounded-md bg-white">
      <h4 className="font-semibold text-sm text-slate-500 uppercase tracking-wider mb-2">
        Project List Settings
      </h4>
      <div className="grid grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label>Number of Items</Label>
          <Input
            type="number"
            value={data.limit || 5}
            onChange={(e) =>
              onChange({ ...data, limit: parseInt(e.target.value) })
            }
          />
        </div>
        <div className="space-y-2">
          <Label>Filter by Tag (Optional)</Label>
          <Input
            value={data.filterTag || ""}
            onChange={(e) => onChange({ ...data, filterTag: e.target.value })}
            placeholder="e.g. AI"
          />
        </div>
      </div>
    </div>
  );
}

export function ProjectListViewer({ data }: { data: ProjectListData }) {
  // In a real implementation, we would use a React Query hook here to fetch projects
  // const { data: projects } = useSearchProjects({ query: data.filterTag, limit: data.limit });

  // Mock data for display
  const projects = [];
  const limit = data.limit || 3;
  for (let i = 0; i < limit; i++) {
    projects.push({
      id: i,
      title: `Research Project ${i + 1} ${data.filterTag ? `(${data.filterTag})` : ""}`,
      description:
        "This is a sample project description showing how the list item would look.",
    });
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {projects.map((project) => (
        <div
          key={project.id}
          className="border rounded-xl p-6 bg-white shadow-sm hover:shadow-md transition-shadow"
        >
          <div className="h-40 bg-slate-100 rounded-lg mb-4 flex items-center justify-center text-slate-400">
            Project Image
          </div>
          <h3 className="font-bold text-lg text-slate-900 mb-2">
            {project.title}
          </h3>
          <p className="text-slate-600 line-clamp-3">{project.description}</p>
        </div>
      ))}
    </div>
  );
}
