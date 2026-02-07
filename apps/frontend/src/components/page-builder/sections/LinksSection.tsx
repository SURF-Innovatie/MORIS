import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Plus, Trash2, ExternalLink } from "lucide-react";

interface LinkItem {
  id: string;
  label: string;
  url: string;
}

interface LinksData {
  title?: string;
  links?: LinkItem[];
}

interface LinksSectionEditorProps {
  data: LinksData;
  onChange: (data: LinksData) => void;
}

export function LinksSectionEditor({
  data,
  onChange,
}: LinksSectionEditorProps) {
  const addLink = () => {
    const newLink = {
      id: crypto.randomUUID(),
      label: "",
      url: "",
    };
    onChange({ ...data, links: [...(data.links || []), newLink] });
  };

  const updateLink = (id: string, field: keyof LinkItem, value: string) => {
    const newLinks = (data.links || []).map((link) =>
      link.id === id ? { ...link, [field]: value } : link,
    );
    onChange({ ...data, links: newLinks });
  };

  const removeLink = (id: string) => {
    const newLinks = (data.links || []).filter((link) => link.id !== id);
    onChange({ ...data, links: newLinks });
  };

  return (
    <div className="space-y-4 p-4 border rounded-md bg-white">
      <h4 className="font-semibold text-sm text-slate-500 uppercase tracking-wider mb-2">
        Links Block Settings
      </h4>
      <div className="space-y-2">
        <Label>Section Title</Label>
        <Input
          value={data.title || ""}
          onChange={(e) => onChange({ ...data, title: e.target.value })}
          placeholder="e.g. Social Media, Resources"
        />
      </div>

      <div className="space-y-3 mt-4">
        <Label>Links</Label>
        {(data.links || []).map((link) => (
          <div key={link.id} className="flex gap-2 items-start">
            <div className="grid grid-cols-2 gap-2 flex-1">
              <Input
                value={link.label}
                onChange={(e) => updateLink(link.id, "label", e.target.value)}
                placeholder="Label (e.g. Twitter)"
              />
              <Input
                value={link.url}
                onChange={(e) => updateLink(link.id, "url", e.target.value)}
                placeholder="URL (https://...)"
              />
            </div>
            <Button
              variant="ghost"
              size="icon"
              className="text-red-500 hover:text-red-600 hover:bg-red-50"
              onClick={() => removeLink(link.id)}
            >
              <Trash2 className="w-4 h-4" />
            </Button>
          </div>
        ))}
        <Button
          onClick={addLink}
          variant="outline"
          size="sm"
          className="w-full"
        >
          <Plus className="w-4 h-4 mr-2" /> Add Link
        </Button>
      </div>
    </div>
  );
}

export function LinksSectionViewer({ data }: { data: LinksData }) {
  if (!data.links || data.links.length === 0) return null;

  return (
    <div className="py-6">
      {data.title && (
        <h3 className="text-xl font-bold text-slate-900 mb-4">{data.title}</h3>
      )}
      <div className="flex flex-wrap gap-3">
        {data.links.map((link) => (
          <a
            key={link.id}
            href={link.url}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center px-4 py-2 rounded-lg bg-white border border-slate-200 text-slate-700 hover:border-primary hover:text-primary transition-colors shadow-sm"
          >
            {link.label}
            <ExternalLink className="w-3 h-3 ml-2 opacity-50" />
          </a>
        ))}
      </div>
    </div>
  );
}
