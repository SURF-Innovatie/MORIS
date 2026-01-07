import { useState } from "react";
import { useGetUsersSearch } from "@api/moris"; // Generated hook
import { PersonResponse } from "@api/model"; // Generated model
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Search, Check, User as UserIcon } from "lucide-react";
import { cn } from "@/lib/utils";

interface UserSearchSelectProps {
    value?: string;
    onSelect: (personId: string, person: PersonResponse) => void;
    disabled?: boolean;
}

export const UserSearchSelect = ({ value, onSelect, disabled }: UserSearchSelectProps) => {
    const [open, setOpen] = useState(false);
    const [query, setQuery] = useState("");
    const [selectedPerson, setSelectedPerson] = useState<PersonResponse | null>(null);

    // Debounce query or just let React Query handle it with keepPreviousData?
    // Let's use a simple debounce hook or just rely on typing speed. 
    // For simplicity, direct query for now.
    const { data: results, isLoading } = useGetUsersSearch(
        { q: query },
        { query: { enabled: open && query.length > 0 } }
    );

    // Initial fetch for selected value if needed? 
    // Ideally, parent passes the person object or we fetch it. 
    // For "Add Member", we usually start empty.

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
                    {selectedPerson ? (
                        <div className="flex items-center gap-2">
                            <UserAvatar person={selectedPerson} />
                            <span>{selectedPerson.name || selectedPerson.email}</span>
                        </div>
                    ) : (
                        <span className="text-muted-foreground">{value ? "Selected (Details loading...)" : "Select user..."}</span>
                    )}
                    <Search className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-[300px] p-0" align="start">
                <div className="flex items-center border-b px-3">
                    <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
                    <Input
                        className="flex h-11 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground border-none focus-visible:ring-0"
                        placeholder="Search by name or email..."
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                    />
                </div>
                <div className="max-h-[300px] overflow-y-auto p-1">
                    {isLoading && <div className="py-6 text-center text-sm text-muted-foreground">Searching...</div>}

                    {!isLoading && results?.length === 0 && (
                        <div className="py-6 text-center text-sm text-muted-foreground">No user found.</div>
                    )}

                    {!isLoading && results?.map((person) => (
                        <div
                            key={person.id}
                            className={cn(
                                "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground",
                                value === person.id && "bg-accent text-accent-foreground"
                            )}
                            onClick={() => {
                                onSelect(person.id!, person);
                                setSelectedPerson(person);
                                setOpen(false);
                            }}
                        >
                            <Check
                                className={cn(
                                    "mr-2 h-4 w-4",
                                    value === person.id ? "opacity-100" : "opacity-0"
                                )}
                            />
                            <div className="flex items-center gap-2">
                                <UserAvatar person={person} />
                                <div className="flex flex-col">
                                    <span>{person.name}</span>
                                    <span className="text-xs text-muted-foreground">{person.email}</span>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </PopoverContent>
        </Popover>
    );
};

const UserAvatar = ({ person }: { person: PersonResponse }) => {
    if (person.avatarUrl) {
        return <img src={person.avatarUrl} alt={person.name} className="h-6 w-6 rounded-full object-cover" />;
    }
    return <UserIcon className="h-6 w-6 p-1 rounded-full bg-secondary text-secondary-foreground" />;
};
