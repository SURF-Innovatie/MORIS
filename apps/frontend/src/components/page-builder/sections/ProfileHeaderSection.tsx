import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";

interface ProfileHeaderData {
  name?: string;
  role?: string;
  bio?: string;
  avatarUrl?: string;
}

interface ProfileHeaderEditorProps {
  data: ProfileHeaderData;
  onChange: (data: ProfileHeaderData) => void;
}

export function ProfileHeaderEditor({
  data,
  onChange,
}: ProfileHeaderEditorProps) {
  return (
    <div className="space-y-4 p-4 border rounded-md bg-white">
      <h4 className="font-semibold text-sm text-slate-500 uppercase tracking-wider mb-2">
        Profile Header Settings
      </h4>
      <div className="space-y-2">
        <Label>Name</Label>
        <Input
          value={data.name || ""}
          onChange={(e) => onChange({ ...data, name: e.target.value })}
          placeholder="e.g. Dr. Jane Doe"
        />
      </div>
      <div className="space-y-2">
        <Label>Role / Title</Label>
        <Input
          value={data.role || ""}
          onChange={(e) => onChange({ ...data, role: e.target.value })}
          placeholder="e.g. Senior Researcher"
        />
      </div>
      <div className="space-y-2">
        <Label>Bio</Label>
        <Textarea
          value={data.bio || ""}
          onChange={(e) => onChange({ ...data, bio: e.target.value })}
          placeholder="Short biography..."
        />
      </div>
      <div className="space-y-2">
        <Label>Avatar URL</Label>
        <Input
          value={data.avatarUrl || ""}
          onChange={(e) => onChange({ ...data, avatarUrl: e.target.value })}
          placeholder="https://..."
        />
      </div>
    </div>
  );
}

export function ProfileHeaderViewer({ data }: { data: ProfileHeaderData }) {
  return (
    <div className="flex flex-col items-center text-center p-8 bg-slate-50 rounded-xl">
      <div className="w-32 h-32 rounded-full bg-slate-200 overflow-hidden mb-4 border-4 border-white shadow-sm">
        {data.avatarUrl ? (
          <img
            src={data.avatarUrl}
            alt={data.name}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center text-slate-400 text-4xl">
            {data.name?.charAt(0) || "?"}
          </div>
        )}
      </div>
      <h1 className="text-3xl font-bold text-slate-900 mb-1">
        {data.name || "Name"}
      </h1>
      <p className="text-lg text-primary font-medium mb-4">
        {data.role || "Role"}
      </p>
      <p className="max-w-xl text-slate-600 leading-relaxed">
        {data.bio || "No biography provided."}
      </p>
    </div>
  );
}
