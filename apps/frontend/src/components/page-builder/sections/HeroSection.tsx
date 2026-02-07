import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

interface HeroData {
  title?: string;
  subtitle?: string;
  imageUrl?: string;
}

interface HeroSectionEditorProps {
  data: HeroData;
  onChange: (data: HeroData) => void;
}

export function HeroSectionEditor({ data, onChange }: HeroSectionEditorProps) {
  return (
    <div className="space-y-4 p-4 border rounded-md bg-white">
      <h4 className="font-semibold text-sm text-slate-500 uppercase tracking-wider mb-2">
        Hero Banner Settings
      </h4>
      <div className="space-y-2">
        <Label>Title</Label>
        <Input
          value={data.title || ""}
          onChange={(e) => onChange({ ...data, title: e.target.value })}
          placeholder="Enter a catchy title"
        />
      </div>
      <div className="space-y-2">
        <Label>Subtitle</Label>
        <Input
          value={data.subtitle || ""}
          onChange={(e) => onChange({ ...data, subtitle: e.target.value })}
          placeholder="Enter a brief description"
        />
      </div>
      <div className="space-y-2">
        <Label>Background Image URL</Label>
        <Input
          value={data.imageUrl || ""}
          onChange={(e) => onChange({ ...data, imageUrl: e.target.value })}
          placeholder="https://example.com/image.jpg"
        />
      </div>
    </div>
  );
}

export function HeroSectionViewer({ data }: { data: HeroData }) {
  return (
    <div
      className="relative rounded-xl overflow-hidden bg-slate-100 min-h-[300px] flex items-center justify-center text-center p-8"
      style={{
        backgroundImage: data.imageUrl ? `url(${data.imageUrl})` : undefined,
        backgroundSize: "cover",
        backgroundPosition: "center",
      }}
    >
      {data.imageUrl && <div className="absolute inset-0 bg-black/40" />}
      <div className="relative z-10 max-w-2xl mx-auto text-white">
        <h1 className="text-4xl md:text-5xl font-bold mb-4 font-serif tracking-tight drop-shadow-sm">
          {data.title || "Untitled Page"}
        </h1>
        <p className="text-lg md:text-xl text-white/90 font-light drop-shadow-sm">
          {data.subtitle || "Add a subtitle to describe this page."}
        </p>
      </div>
    </div>
  );
}
