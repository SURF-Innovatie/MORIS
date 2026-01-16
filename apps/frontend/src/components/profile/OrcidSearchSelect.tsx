import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Search, Check, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";
import OrcidIcon from "@/components/icons/orcidIcon";
import { useGetOrcidSearch } from "@api/moris";
import { OrcidPerson } from "@api/model";

interface OrcidSearchSelectProps {
  value?: string; // This is the ORCID ID
  onSelect: (item: OrcidPerson) => void;
  disabled?: boolean;
}

export const OrcidSearchSelect = ({
  value,
  onSelect,
  disabled,
}: OrcidSearchSelectProps) => {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [selectedItem, setSelectedItem] = useState<OrcidPerson | null>(null);

  const { data: results, isLoading } = useGetOrcidSearch(
    { q: query },
    { query: { enabled: open && query.length > 2 } }
  );

  const getDisplayName = (item: OrcidPerson) => {
    if (item.credit_name) return item.credit_name;
    if (item.first_name && item.last_name)
      return `${item.first_name} ${item.last_name}`;
    return item.first_name || item.last_name || "Unknown Name";
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between"
          disabled={disabled}
        >
          {selectedItem || value ? (
            <div className="flex items-center gap-2 overflow-hidden">
              <OrcidIcon width={16} height={16} className="shrink-0" />
              <span className="truncate">
                {selectedItem ? getDisplayName(selectedItem) : value}
                <span className="text-muted-foreground font-normal ml-1">
                  {selectedItem?.orcid || value}
                </span>
              </span>
            </div>
          ) : (
            <span className="text-muted-foreground">
              Search ORCID (optional)...
            </span>
          )}
          <Search className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        className="w-[400px] p-0"
        align="start"
        side="bottom"
        avoidCollisions={false}
      >
        <div className="flex items-center border-b px-3">
          <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
          <Input
            className="flex h-11 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground border-none focus-visible:ring-0 shadow-none"
            placeholder="Search by name (min 3 chars)..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
        </div>
        <div className="max-h-[300px] overflow-y-auto p-1">
          {isLoading && (
            <div className="py-6 flex justify-center text-sm text-muted-foreground">
              <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Searching...
            </div>
          )}

          {!isLoading && query.length > 2 && results?.length === 0 && (
            <div className="py-6 text-center text-sm text-muted-foreground">
              No person found.
            </div>
          )}

          {!isLoading && query.length <= 2 && (
            <div className="py-6 text-center text-sm text-muted-foreground">
              Type at least 3 characters to search.
            </div>
          )}

          {!isLoading &&
            results?.map((item) => (
              <div
                key={item.orcid}
                className={cn(
                  "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground",
                  value === item.orcid && "bg-accent text-accent-foreground"
                )}
                onClick={() => {
                  onSelect(item);
                  setSelectedItem(item);
                  setOpen(false);
                }}
              >
                <Check
                  className={cn(
                    "mr-2 h-4 w-4",
                    value === item.orcid ? "opacity-100" : "opacity-0"
                  )}
                />
                <div className="flex items-center gap-2 overflow-hidden bg-transparent">
                  <div className="flex flex-col">
                    <span className="font-medium">{getDisplayName(item)}</span>
                    <span className="text-xs text-muted-foreground">
                      {item.orcid}
                    </span>
                  </div>
                </div>
              </div>
            ))}
        </div>
      </PopoverContent>
    </Popover>
  );
};
