import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Plus, Trash2 } from "lucide-react";

interface StatItem {
  id: string;
  value: string;
  label: string;
}

interface StatisticsData {
  title?: string;
  stats?: StatItem[];
}

interface StatisticsSectionEditorProps {
  data: StatisticsData;
  onChange: (data: StatisticsData) => void;
}

export function StatisticsSectionEditor({
  data,
  onChange,
}: StatisticsSectionEditorProps) {
  const addStat = () => {
    const newStat = {
      id: crypto.randomUUID(),
      value: "0",
      label: "New Stat",
    };
    onChange({ ...data, stats: [...(data.stats || []), newStat] });
  };

  const updateStat = (id: string, field: keyof StatItem, value: string) => {
    const newStats = (data.stats || []).map((stat) =>
      stat.id === id ? { ...stat, [field]: value } : stat,
    );
    onChange({ ...data, stats: newStats });
  };

  const removeStat = (id: string) => {
    const newStats = (data.stats || []).filter((stat) => stat.id !== id);
    onChange({ ...data, stats: newStats });
  };

  return (
    <div className="space-y-4 p-4 border rounded-md bg-white">
      <h4 className="font-semibold text-sm text-slate-500 uppercase tracking-wider mb-2">
        Statistics Settings
      </h4>
      <div className="space-y-2">
        <Label>Section Title (Optional)</Label>
        <Input
          value={data.title || ""}
          onChange={(e) => onChange({ ...data, title: e.target.value })}
          placeholder="e.g. Key Metrics"
        />
      </div>

      <div className="space-y-3 mt-4">
        <Label>Statistics</Label>
        {(data.stats || []).map((stat) => (
          <div key={stat.id} className="flex gap-2 items-start">
            <div className="grid grid-cols-2 gap-2 flex-1">
              <Input
                value={stat.value}
                onChange={(e) => updateStat(stat.id, "value", e.target.value)}
                placeholder="Value (e.g. 50+)"
              />
              <Input
                value={stat.label}
                onChange={(e) => updateStat(stat.id, "label", e.target.value)}
                placeholder="Label (e.g. Projects)"
              />
            </div>
            <Button
              variant="ghost"
              size="icon"
              className="text-red-500 hover:text-red-600 hover:bg-red-50"
              onClick={() => removeStat(stat.id)}
            >
              <Trash2 className="w-4 h-4" />
            </Button>
          </div>
        ))}
        <Button
          onClick={addStat}
          variant="outline"
          size="sm"
          className="w-full"
        >
          <Plus className="w-4 h-4 mr-2" /> Add Statistic
        </Button>
      </div>
    </div>
  );
}

export function StatisticsSectionViewer({ data }: { data: StatisticsData }) {
  if (!data.stats || data.stats.length === 0) return null;

  return (
    <div className="py-8">
      {data.title && (
        <h3 className="text-xl font-bold text-slate-900 mb-6 text-center">
          {data.title}
        </h3>
      )}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
        {data.stats.map((stat) => (
          <div key={stat.id} className="text-center p-4 bg-slate-50 rounded-lg">
            <div className="text-4xl font-black text-primary mb-2">
              {stat.value}
            </div>
            <div className="text-sm uppercase tracking-wide text-slate-500 font-semibold">
              {stat.label}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
