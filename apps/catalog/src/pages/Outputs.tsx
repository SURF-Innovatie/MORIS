import { useState, useMemo, useCallback } from "react";
import { useCatalog } from "../components/CatalogProvider";

const TYPE_LABELS: Record<number, string> = {
  0: "Cartographic Material",
  1: "Data Publication",
  2: "Image",
  3: "Interactive Resource",
  4: "Learning Object",
  5: "Other",
  6: "Model/Code",
  7: "Sound",
  8: "Trademark",
  9: "Workflow",
};

type SortKey = "title" | "year" | "type" | "project";

interface FlatOutput {
  productId: string;
  name: string;
  doi: string;
  type: number;
  typeLabel: string;
  authorIds: string[];
  authorNames: string[];
  projectId: string;
  projectTitle: string;
  year: number | null;
}

export function Outputs() {
  const { data } = useCatalog();
  const catalog = data?.catalog;
  const projects = data?.projects ?? [];
  const people = data?.people ?? {};

  // --- Flatten all outputs across projects ---
  const outputs: FlatOutput[] = useMemo(() => {
    const seen = new Set<string>();
    const results: FlatOutput[] = [];

    for (const pd of projects) {
      const proj = pd.project;
      if (!proj) continue;

      for (const prod of pd.products ?? []) {
        if (seen.has(prod.Id)) continue;
        seen.add(prod.Id);

        const authorNames = (prod.AuthorPersonIDs ?? [])
          .map((id: string) => people[id]?.Name)
          .filter(Boolean) as string[];

        results.push({
          productId: prod.Id,
          name: prod.Name,
          doi: prod.DOI ?? "",
          type: prod.Type as number,
          typeLabel: TYPE_LABELS[prod.Type as number] ?? "Other",
          authorIds: prod.AuthorPersonIDs ?? [],
          authorNames,
          projectId: proj.id,
          projectTitle: proj.title,
          year: proj.start_date
            ? new Date(proj.start_date).getFullYear()
            : null,
        });
      }
    }

    return results;
  }, [projects, people]);

  // --- Derive filter options from data ---
  const projectOptions = useMemo(
    () =>
      projects
        .map((p) => ({
          id: p.project?.id ?? "",
          title: p.project?.title ?? "",
        }))
        .filter((p) => p.id),
    [projects],
  );

  const typeOptions = useMemo(() => {
    const types = new Set(outputs.map((o) => o.type));
    return [...types]
      .sort((a, b) => a - b)
      .map((t) => ({ value: t, label: TYPE_LABELS[t] ?? `Type ${t}` }));
  }, [outputs]);

  const authorOptions = useMemo(() => {
    const map = new Map<string, string>();
    for (const o of outputs) {
      for (let i = 0; i < o.authorIds.length; i++) {
        if (o.authorNames[i]) map.set(o.authorIds[i], o.authorNames[i]);
      }
    }
    return [...map.entries()]
      .sort((a, b) => a[1].localeCompare(b[1]))
      .map(([id, name]) => ({ id, name }));
  }, [outputs]);

  // --- Filter state ---
  const [search, setSearch] = useState("");
  const [selectedProjects, setSelectedProjects] = useState<Set<string>>(
    new Set(),
  );
  const [selectedTypes, setSelectedTypes] = useState<Set<number>>(new Set());
  const [selectedAuthors, setSelectedAuthors] = useState<Set<string>>(
    new Set(),
  );
  const [sortBy, setSortBy] = useState<SortKey>("title");
  const [sortAsc, setSortAsc] = useState(true);

  const toggleFilter = useCallback(
    <T,>(set: Set<T>, val: T, setter: (s: Set<T>) => void) => {
      const next = new Set(set);
      if (next.has(val)) next.delete(val);
      else next.add(val);
      setter(next);
    },
    [],
  );

  const activeFilterCount =
    selectedProjects.size + selectedTypes.size + selectedAuthors.size;

  const clearFilters = useCallback(() => {
    setSelectedProjects(new Set());
    setSelectedTypes(new Set());
    setSelectedAuthors(new Set());
    setSearch("");
  }, []);

  // --- Apply filters ---
  const filtered = useMemo(() => {
    let list = outputs;

    if (search.trim()) {
      const q = search.toLowerCase();
      list = list.filter(
        (o) =>
          o.name.toLowerCase().includes(q) ||
          o.doi.toLowerCase().includes(q) ||
          o.projectTitle.toLowerCase().includes(q) ||
          o.authorNames.some((a) => a.toLowerCase().includes(q)),
      );
    }

    if (selectedProjects.size > 0) {
      list = list.filter((o) => selectedProjects.has(o.projectId));
    }
    if (selectedTypes.size > 0) {
      list = list.filter((o) => selectedTypes.has(o.type));
    }
    if (selectedAuthors.size > 0) {
      list = list.filter((o) =>
        o.authorIds.some((id) => selectedAuthors.has(id)),
      );
    }

    const sorted = [...list].sort((a, b) => {
      let cmp = 0;
      switch (sortBy) {
        case "title":
          cmp = a.name.localeCompare(b.name);
          break;
        case "year":
          cmp = (a.year ?? 0) - (b.year ?? 0);
          break;
        case "type":
          cmp = a.typeLabel.localeCompare(b.typeLabel);
          break;
        case "project":
          cmp = a.projectTitle.localeCompare(b.projectTitle);
          break;
      }
      return sortAsc ? cmp : -cmp;
    });

    return sorted;
  }, [
    outputs,
    search,
    selectedProjects,
    selectedTypes,
    selectedAuthors,
    sortBy,
    sortAsc,
  ]);

  const handleSort = (key: SortKey) => {
    if (sortBy === key) {
      setSortAsc(!sortAsc);
    } else {
      setSortBy(key);
      setSortAsc(true);
    }
  };

  const SortIcon = ({ column }: { column: SortKey }) => {
    if (sortBy !== column)
      return <span className="text-gray-300 ml-1">&darr;</span>;
    return (
      <span className="text-(--color-primary) ml-1">
        {sortAsc ? "\u2191" : "\u2193"}
      </span>
    );
  };

  return (
    <>
      {/* Hero banner */}
      <div className="bg-(--color-primary) pb-28 pt-10">
        <div className="max-w-7xl mx-auto px-6">
          <h1 className="text-white font-condensed-extra text-4xl md:text-5xl font-semibold leading-tight max-w-lg">
            Output Catalogus
          </h1>
          <p className="text-white/80 mt-3 max-w-xl text-lg">
            All publications, datasets, models and code from{" "}
            {catalog?.title ?? "the programme"}.
          </p>
        </div>
      </div>

      {/* Main content — overlaps the hero */}
      <div className="max-w-7xl mx-auto px-6 -mt-20 pb-12">
        {/* Search bar card */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 mb-6">
          <div className="flex items-center gap-4">
            <div className="relative grow">
              <input
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Search publications by title, author, DOI or project..."
                className="w-full border border-gray-300 rounded-lg py-3 pl-4 pr-12 text-base focus:outline-none focus:ring-2 focus:ring-(--color-primary)/30 focus:border-(--color-primary) transition"
              />
              <svg
                className="absolute right-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
            </div>
            <div className="flex items-center gap-2 text-sm text-gray-500 shrink-0">
              <span className="bg-(--color-accent) text-white text-xs font-bold rounded-full h-6 min-w-6 flex items-center justify-center px-1.5">
                {filtered.length}
              </span>
              <span>
                {filtered.length === 1 ? "publication" : "publications"}
              </span>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-[260px_1fr] gap-6">
          {/* Sidebar: Filters */}
          <aside className="space-y-5">
            {/* Active filters summary */}
            {activeFilterCount > 0 && (
              <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-semibold text-gray-900">
                    Active filters ({activeFilterCount})
                  </span>
                  <button
                    onClick={clearFilters}
                    className="text-xs text-(--color-primary) hover:underline"
                  >
                    Clear all
                  </button>
                </div>
                <div className="flex flex-wrap gap-1.5">
                  {[...selectedProjects].map((id) => {
                    const label =
                      projectOptions.find((p) => p.id === id)?.title ?? id;
                    return (
                      <button
                        key={`p-${id}`}
                        onClick={() =>
                          toggleFilter(
                            selectedProjects,
                            id,
                            setSelectedProjects,
                          )
                        }
                        className="inline-flex items-center gap-1 bg-(--color-primary)/10 text-(--color-primary) text-xs rounded-full px-2.5 py-1 hover:bg-(--color-primary)/20 transition"
                      >
                        <span className="truncate max-w-[140px]">{label}</span>
                        <span>&times;</span>
                      </button>
                    );
                  })}
                  {[...selectedTypes].map((t) => (
                    <button
                      key={`t-${t}`}
                      onClick={() =>
                        toggleFilter(selectedTypes, t, setSelectedTypes)
                      }
                      className="inline-flex items-center gap-1 bg-(--color-primary)/10 text-(--color-primary) text-xs rounded-full px-2.5 py-1 hover:bg-(--color-primary)/20 transition"
                    >
                      {TYPE_LABELS[t]} <span>&times;</span>
                    </button>
                  ))}
                  {[...selectedAuthors].map((id) => {
                    const label =
                      authorOptions.find((a) => a.id === id)?.name ?? id;
                    return (
                      <button
                        key={`a-${id}`}
                        onClick={() =>
                          toggleFilter(
                            selectedAuthors,
                            id,
                            setSelectedAuthors,
                          )
                        }
                        className="inline-flex items-center gap-1 bg-(--color-primary)/10 text-(--color-primary) text-xs rounded-full px-2.5 py-1 hover:bg-(--color-primary)/20 transition"
                      >
                        <span className="truncate max-w-[140px]">{label}</span>
                        <span>&times;</span>
                      </button>
                    );
                  })}
                </div>
              </div>
            )}

            {/* Project filter */}
            <FilterGroup title="Project" count={selectedProjects.size}>
              {projectOptions.map((p) => (
                <FilterCheckbox
                  key={p.id}
                  label={p.title}
                  checked={selectedProjects.has(p.id)}
                  onChange={() =>
                    toggleFilter(selectedProjects, p.id, setSelectedProjects)
                  }
                />
              ))}
            </FilterGroup>

            {/* Type filter */}
            <FilterGroup title="Publication type" count={selectedTypes.size}>
              {typeOptions.map((t) => (
                <FilterCheckbox
                  key={t.value}
                  label={t.label}
                  checked={selectedTypes.has(t.value)}
                  onChange={() =>
                    toggleFilter(selectedTypes, t.value, setSelectedTypes)
                  }
                />
              ))}
            </FilterGroup>

            {/* Author filter */}
            <FilterGroup title="Authors" count={selectedAuthors.size}>
              {authorOptions.map((a) => (
                <FilterCheckbox
                  key={a.id}
                  label={a.name}
                  checked={selectedAuthors.has(a.id)}
                  onChange={() =>
                    toggleFilter(selectedAuthors, a.id, setSelectedAuthors)
                  }
                />
              ))}
            </FilterGroup>
          </aside>

          {/* Results */}
          <section>
            {/* Sort bar */}
            <div className="flex items-center gap-1 text-xs font-medium text-gray-500 mb-3 px-1">
              <span className="mr-2">Sort by:</span>
              {(
                [
                  ["title", "Title"],
                  ["year", "Year"],
                  ["type", "Type"],
                  ["project", "Project"],
                ] as const
              ).map(([key, label]) => (
                <button
                  key={key}
                  onClick={() => handleSort(key)}
                  className={`px-2 py-1 rounded transition ${sortBy === key ? "bg-(--color-primary)/10 text-(--color-primary)" : "hover:bg-gray-100"}`}
                >
                  {label}
                  <SortIcon column={key} />
                </button>
              ))}
            </div>

            {/* Output list */}
            {filtered.length === 0 ? (
              <div className="bg-white rounded-xl border border-gray-100 shadow-sm p-10 text-center text-gray-400">
                No publications found matching your criteria.
              </div>
            ) : (
              <div className="space-y-3">
                {filtered.map((o) => (
                  <div
                    key={o.productId}
                    className="bg-white rounded-lg border border-gray-100 shadow-sm p-5 hover:shadow-md transition-shadow"
                  >
                    <div className="flex items-start justify-between gap-4">
                      <div className="min-w-0">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="inline-block bg-(--color-primary)/10 text-(--color-primary) text-xs font-medium rounded-full px-2.5 py-0.5">
                            {o.typeLabel}
                          </span>
                          {o.year && (
                            <span className="text-xs text-gray-400">
                              {o.year}
                            </span>
                          )}
                        </div>
                        <h3 className="text-base font-semibold text-gray-900 leading-snug">
                          {o.doi ? (
                            <a
                              href={`https://doi.org/${o.doi}`}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="hover:text-(--color-primary) transition-colors"
                            >
                              {o.name}
                            </a>
                          ) : (
                            o.name
                          )}
                        </h3>
                        {o.authorNames.length > 0 && (
                          <p className="text-sm text-gray-500 mt-1">
                            {o.authorNames.join(", ")}
                          </p>
                        )}
                      </div>
                      {o.doi && (
                        <a
                          href={`https://doi.org/${o.doi}`}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="shrink-0 text-xs text-(--color-primary) border border-(--color-primary)/30 rounded-full px-3 py-1 hover:bg-(--color-primary)/5 transition"
                        >
                          DOI
                        </a>
                      )}
                    </div>
                    <div className="mt-2 text-xs text-gray-400">
                      <a
                        href={`/projects/${o.projectId}`}
                        className="hover:text-(--color-primary) transition-colors"
                      >
                        {o.projectTitle}
                      </a>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </section>
        </div>
      </div>
    </>
  );
}

// --- Filter UI components ---

function FilterGroup({
  title,
  count,
  children,
}: {
  title: string;
  count: number;
  children: React.ReactNode;
}) {
  const [expanded, setExpanded] = useState(true);

  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm">
      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full flex items-center justify-between p-4 text-sm font-semibold text-gray-900"
      >
        <span className="flex items-center gap-2">
          {title}
          {count > 0 && (
            <span className="bg-(--color-accent) text-white text-[10px] font-bold rounded-full h-5 min-w-5 flex items-center justify-center px-1">
              {count}
            </span>
          )}
        </span>
        <span className="text-gray-400 text-xs">
          {expanded ? "\u2212" : "+"}
        </span>
      </button>
      {expanded && (
        <div className="px-4 pb-4 space-y-1.5 max-h-52 overflow-y-auto">
          {children}
        </div>
      )}
    </div>
  );
}

function FilterCheckbox({
  label,
  checked,
  onChange,
}: {
  label: string;
  checked: boolean;
  onChange: () => void;
}) {
  return (
    <label className="flex items-center gap-2 cursor-pointer group">
      <input
        type="checkbox"
        checked={checked}
        onChange={onChange}
        className="rounded border-gray-300 text-(--color-primary) focus:ring-(--color-primary)/30 h-3.5 w-3.5"
      />
      <span
        className={`text-sm truncate ${checked ? "text-(--color-primary) font-medium" : "text-gray-600 group-hover:text-gray-900"} transition-colors`}
      >
        {label}
      </span>
    </label>
  );
}
