import { useState, useMemo } from "react";
import { useCatalog } from "../components/CatalogProvider";

export function Home() {
  const { data } = useCatalog();
  const catalog = data?.catalog;
  const projects = data?.projects ?? [];
  const organisations = data?.organisations ?? {};

  const [search, setSearch] = useState("");

  const filtered = useMemo(() => {
    if (!search.trim()) return [];
    const q = search.toLowerCase();
    return projects.filter(
      (item) =>
        item.project?.title.toLowerCase().includes(q) ||
        item.project?.description?.toLowerCase().includes(q),
    );
  }, [search, projects]);

  const showResults = search.trim().length > 0;

  return (
    <>
      {/* Hero banner */}
      <div className="bg-(--color-primary) pb-28 pt-10">
        <div className="max-w-7xl mx-auto px-6">
          <h1 className="text-white font-condensed-extra text-4xl md:text-5xl font-semibold leading-tight max-w-lg">
            {catalog?.title}
          </h1>
          {catalog?.description && (
            <p className="text-white/80 mt-3 max-w-xl text-lg">
              {catalog.description}
            </p>
          )}
        </div>
      </div>

      {/* Main content — overlaps the hero */}
      <div className="max-w-7xl mx-auto px-6 -mt-20 pb-12">
        <div className="grid grid-cols-1 lg:grid-cols-[1fr_320px] gap-6">
          {/* Left: Search section */}
          <section className="bg-white rounded-xl shadow-sm border border-gray-100 p-8">
            <h2 className="font-condensed text-2xl md:text-3xl font-semibold text-gray-900 mb-1">
              Find Projects
            </h2>
            <p className="text-sm text-gray-500 mb-5">
              Search across all projects
            </p>

            {/* Search input */}
            <div className="relative">
              <input
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Search projects..."
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

            {/* Search results */}
            {showResults && (
              <div className="mt-4 space-y-2">
                {filtered.length === 0 ? (
                  <p className="text-gray-400 text-sm py-2">
                    No projects found for "{search}"
                  </p>
                ) : (
                  filtered.map((item) => {
                    const p = item.project;
                    if (!p) return null;
                    return (
                      <a
                        key={p.id}
                        href={`/projects/${p.id}`}
                        className="flex items-start gap-3 p-3 rounded-lg hover:bg-gray-50 transition-colors group"
                      >
                        <span className="text-(--color-primary) mt-0.5 shrink-0">
                          &rarr;
                        </span>
                        <div>
                          <span className="font-medium text-gray-900 group-hover:text-(--color-primary) transition-colors">
                            {p.title}
                          </span>
                          {p.description && (
                            <p className="text-sm text-gray-500 line-clamp-1 mt-0.5">
                              {p.description}
                            </p>
                          )}
                        </div>
                      </a>
                    );
                  })
                )}
              </div>
            )}

            {/* Quick links */}
            {!showResults && projects.length > 0 && (
              <div className="mt-6">
                <p className="text-sm font-medium text-gray-700 mb-3">
                  or go directly to
                </p>
                <div className="flex flex-wrap gap-2">
                  {projects.slice(0, 8).map((item) => {
                    const p = item.project;
                    if (!p) return null;
                    return (
                      <a
                        key={p.id}
                        href={`/projects/${p.id}`}
                        className="inline-flex items-center gap-1.5 border border-gray-300 rounded-full px-4 py-1.5 text-sm text-gray-700 hover:border-(--color-primary) hover:text-(--color-primary) transition-colors"
                      >
                        <span className="text-xs">&rarr;</span>
                        <span className="truncate max-w-[250px]">
                          {p.title}
                        </span>
                      </a>
                    );
                  })}
                </div>
              </div>
            )}
          </section>

          {/* Right: Projects sidebar */}
          <aside className="bg-white/80 backdrop-blur rounded-xl shadow-sm border border-gray-100 p-6 h-fit">
            <div className="flex items-center gap-2 mb-4">
              <h2 className="font-condensed text-2xl font-semibold text-gray-900">
                Projects
              </h2>
              <span className="bg-(--color-accent) text-white text-xs font-bold rounded-full h-6 min-w-6 flex items-center justify-center px-1.5">
                {projects.length}
              </span>
            </div>

            <ul className="space-y-3">
              {projects.slice(0, 5).map((item) => {
                const p = item.project;
                if (!p) return null;
                const org = organisations[p.owning_org_node_id];
                return (
                  <li key={p.id}>
                    <a
                      href={`/projects/${p.id}`}
                      className="text-(--color-primary) text-sm font-medium hover:underline leading-snug block"
                    >
                      {p.title}
                    </a>
                    {org && (
                      <p className="text-xs text-gray-400 mt-0.5">
                        {org.name}
                      </p>
                    )}
                  </li>
                );
              })}
            </ul>

            {projects.length > 5 && (
              <a
                href="/projects"
                className="inline-flex items-center gap-2 mt-5 bg-(--color-primary) text-white text-sm font-medium px-5 py-2 rounded-full hover:opacity-90 transition-opacity"
              >
                <span>&rarr;</span>
                View all projects
              </a>
            )}
          </aside>
        </div>

        {/* Project cards grid */}
        {projects.length > 0 && (
          <section className="mt-10">
            <h2 className="font-condensed text-2xl font-semibold text-gray-900 mb-5">
              All Projects
            </h2>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
              {projects.map((item) => {
                const p = item.project;
                if (!p) return null;
                const org = organisations[p.owning_org_node_id];
                const memberCount = item.members?.length ?? 0;
                const productCount = item.products?.length ?? 0;
                return (
                  <a
                    key={p.id}
                    href={`/projects/${p.id}`}
                    className="bg-white rounded-lg border border-gray-100 shadow-sm hover:shadow-md transition-shadow p-5 group"
                  >
                    {org && (
                      <span className="text-xs text-gray-400 uppercase tracking-wide">
                        {org.name}
                      </span>
                    )}
                    <h3 className="text-lg font-semibold text-gray-900 mt-1 mb-2 group-hover:text-(--color-primary) transition-colors leading-snug">
                      {p.title}
                    </h3>
                    {p.description && (
                      <p className="text-sm text-gray-500 line-clamp-2 mb-3">
                        {p.description}
                      </p>
                    )}
                    <div className="flex items-center gap-4 text-xs text-gray-400">
                      {memberCount > 0 && (
                        <span>
                          {memberCount} member{memberCount !== 1 && "s"}
                        </span>
                      )}
                      {productCount > 0 && (
                        <span>
                          {productCount} output{productCount !== 1 && "s"}
                        </span>
                      )}
                      {p.start_date && (
                        <span>
                          {new Date(p.start_date).getFullYear()}
                        </span>
                      )}
                    </div>
                  </a>
                );
              })}
            </div>
          </section>
        )}
      </div>
    </>
  );
}
