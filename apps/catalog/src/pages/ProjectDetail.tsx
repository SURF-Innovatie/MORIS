import { useParams } from "react-router-dom";
import { useCatalog } from "../components/CatalogProvider";

const PRODUCT_TYPE_LABELS: Record<number, string> = {
  0: "Cartographic Material",
  1: "Dataset",
  2: "Image",
  3: "Interactive Resource",
  4: "Learning Object",
  5: "Other",
  6: "Software",
  7: "Sound",
  8: "Trademark",
  9: "Workflow",
};

function formatDate(dateStr?: string) {
  if (!dateStr) return null;
  const d = new Date(dateStr);
  if (isNaN(d.getTime())) return null;
  return d.toLocaleDateString("en-GB", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
}

export function ProjectDetail() {
  const { id } = useParams<{ id: string }>();
  const { data } = useCatalog();
  const catalog = data?.catalog;
  const people = data?.people ?? {};
  const organisations = data?.organisations ?? {};

  const projectDetail = data?.projects?.find((p) => p.project?.id === id);
  const project = projectDetail?.project;
  const members = projectDetail?.members ?? [];
  const products = projectDetail?.products ?? [];
  const org = project ? organisations[project.owning_org_node_id] : undefined;

  if (!project) {
    return (
      <div className="max-w-7xl mx-auto px-6 py-12">
        <a
          href="/"
          className="text-sm text-gray-500 hover:text-(--color-primary) transition-colors"
        >
          &larr; Back to overview
        </a>
        <h1 className="font-condensed text-3xl font-semibold text-(--color-primary) mt-4">
          Project Not Found
        </h1>
        <p className="text-gray-600 mt-2">Could not find project with ID: {id}</p>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-6 py-8">
      {/* Breadcrumb */}
      <nav className="text-sm text-gray-500 mb-6">
        <a href="/" className="hover:text-(--color-primary) transition-colors">
          {catalog?.title ?? "Projects"}
        </a>
        <span className="mx-2">&rsaquo;</span>
        <span className="text-gray-700">{project.title}</span>
      </nav>

      <div className="grid grid-cols-1 lg:grid-cols-[1fr_340px] gap-10">
        {/* Main content */}
        <div>
          <h1 className="font-condensed text-4xl md:text-5xl font-semibold text-gray-900 leading-tight mb-6">
            {project.title}
          </h1>

          {project.description && (
            <p className="text-lg font-medium text-gray-800 leading-relaxed mb-8">
              {project.description}
            </p>
          )}

          {/* Timeline bar */}
          {(project.start_date || project.end_date) && (
            <div className="flex items-center gap-4 mb-10">
              {project.start_date && (
                <div className="flex items-center gap-2">
                  <span className="flex items-center justify-center w-8 h-8 rounded-full bg-(--color-primary) text-white text-sm font-bold">
                    1
                  </span>
                  <span className="text-sm font-medium text-gray-700">
                    Start: {formatDate(project.start_date)}
                  </span>
                </div>
              )}
              {project.start_date && project.end_date && (
                <div className="grow h-0.5 bg-(--color-primary)/20" />
              )}
              {project.end_date && (
                <div className="flex items-center gap-2">
                  <span className="flex items-center justify-center w-8 h-8 rounded-full bg-(--color-primary) text-white text-sm font-bold">
                    2
                  </span>
                  <span className="text-sm font-medium text-gray-700">
                    End: {formatDate(project.end_date)}
                  </span>
                </div>
              )}
            </div>
          )}

          {/* Members */}
          {members.length > 0 && (
            <section className="mb-10">
              <h2 className="font-condensed text-2xl font-semibold text-(--color-primary) mb-4">
                Team
              </h2>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                {members.map((m) => (
                  <div
                    key={m.Person.ID}
                    className="flex items-center gap-3 bg-white border border-gray-100 rounded-lg p-3"
                  >
                    {m.Person.AvatarUrl ? (
                      <img
                        src={m.Person.AvatarUrl}
                        alt=""
                        className="w-10 h-10 rounded-full object-cover shrink-0"
                      />
                    ) : (
                      <div className="w-10 h-10 rounded-full bg-(--color-primary)/10 text-(--color-primary) flex items-center justify-center text-sm font-semibold shrink-0">
                        {(m.Person.GivenName?.[0] ?? m.Person.Name[0]).toUpperCase()}
                      </div>
                    )}
                    <div className="min-w-0">
                      <p className="text-sm font-medium text-gray-900 truncate">
                        {m.Person.Name}
                      </p>
                      <p className="text-xs text-gray-500 truncate">
                        {m.Role.Name}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </section>
          )}

          {/* Products / Outputs */}
          {products.length > 0 && (
            <section>
              <h2 className="font-condensed text-2xl font-semibold text-(--color-primary) mb-4">
                Outputs
              </h2>
              <div className="flex flex-wrap gap-2">
                {products.map((prod) => (
                  <span
                    key={prod.Id}
                    className="inline-flex items-center gap-1.5 border border-gray-300 rounded-full px-4 py-1.5 text-sm text-gray-700"
                  >
                    {prod.DOI ? (
                      <a
                        href={`https://doi.org/${prod.DOI}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="hover:text-(--color-primary) transition-colors"
                      >
                        {prod.Name}
                      </a>
                    ) : (
                      prod.Name
                    )}
                    <span className="text-xs text-gray-400">
                      {PRODUCT_TYPE_LABELS[prod.Type as number] ?? ""}
                    </span>
                  </span>
                ))}
              </div>
            </section>
          )}
        </div>

        {/* Sidebar */}
        <aside className="h-fit">
          <div className="bg-gray-50 rounded-xl border border-gray-100 p-6 space-y-5">
            <h2 className="font-condensed text-2xl font-semibold text-(--color-primary)">
              Details
            </h2>

            {org && (
              <div>
                <dt className="text-sm font-semibold text-gray-900">Organisation</dt>
                <dd className="text-sm text-(--color-primary) mt-0.5">
                  {org.name}
                </dd>
              </div>
            )}

            {project.start_date && (
              <div>
                <dt className="text-sm font-semibold text-gray-900">Start date</dt>
                <dd className="text-sm text-gray-600 mt-0.5">
                  {formatDate(project.start_date)}
                </dd>
              </div>
            )}

            {project.end_date && (
              <div>
                <dt className="text-sm font-semibold text-gray-900">End date</dt>
                <dd className="text-sm text-gray-600 mt-0.5">
                  {formatDate(project.end_date)}
                </dd>
              </div>
            )}

            {members.length > 0 && (
              <div>
                <dt className="text-sm font-semibold text-gray-900">Team members</dt>
                <dd className="text-sm text-(--color-primary) mt-0.5">
                  {members.length} member{members.length !== 1 && "s"}
                </dd>
              </div>
            )}

            {products.length > 0 && (
              <div>
                <dt className="text-sm font-semibold text-gray-900">Outputs</dt>
                <dd className="text-sm text-(--color-primary) mt-0.5">
                  {products.length} output{products.length !== 1 && "s"}
                </dd>
              </div>
            )}

            {catalog && (
              <div>
                <dt className="text-sm font-semibold text-gray-900">Catalog</dt>
                <dd className="text-sm mt-0.5">
                  <a
                    href="/"
                    className="text-(--color-primary) hover:underline"
                  >
                    {catalog.title}
                  </a>
                </dd>
              </div>
            )}
          </div>

          {/* Back link */}
          <a
            href="/"
            className="inline-flex items-center gap-2 mt-5 bg-(--color-primary) text-white text-sm font-medium px-5 py-2.5 rounded-full hover:opacity-90 transition-opacity"
          >
            <span>&rarr;</span>
            View all projects
          </a>
        </aside>
      </div>
    </div>
  );
}
