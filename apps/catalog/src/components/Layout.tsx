import { type ReactNode, useEffect } from "react";
import { useCatalog } from "./CatalogProvider";

export function Layout({ children }: { children: ReactNode }) {
  const { data, isLoading } = useCatalog();
  const catalog = data?.catalog;

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        Loading...
      </div>
    );
  }

  useEffect(() => {
    if (catalog?.favicon) {
      const link = document.querySelector(
        "link[rel~='icon']",
      ) as HTMLLinkElement;
      if (link) {
        link.href = catalog.favicon;
      } else {
        const newLink = document.createElement("link");
        newLink.rel = "icon";
        newLink.href = catalog.favicon;
        document.head.appendChild(newLink);
      }
    }
  }, [catalog?.favicon]);

  return (
    <div className="min-h-screen flex flex-col font-sans bg-gray-50 text-gray-900">
      <header className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
          <a href="/" className="flex items-center gap-3">
            {catalog?.logo_url && (
              <img src={catalog.logo_url} alt="" className="h-10" />
            )}
            <span className="text-xl font-semibold font-condensed text-(--color-primary)">
              {catalog?.title || "Catalog"}
            </span>
          </a>
          <nav className="flex items-center gap-6 text-sm font-medium text-gray-700">
            <a
              href="/"
              className="hover:text-(--color-primary) transition-colors"
            >
              Home
            </a>
            <a
              href="/outputs"
              className="hover:text-(--color-primary) transition-colors"
            >
              Outputs
            </a>
          </nav>
        </div>
      </header>

      <main className="grow">{children}</main>

      <footer className="bg-(--color-secondary) text-white py-10">
        <div className="max-w-7xl mx-auto px-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div>
              {catalog?.logo_url && (
                <img
                  src={catalog.logo_url}
                  alt=""
                  className="h-10 mb-4 brightness-0 invert"
                />
              )}
              <p className="text-sm text-white/70">
                {catalog?.description}
              </p>
            </div>
            <div>
              <h3 className="font-semibold font-condensed text-lg mb-3">
                Menu
              </h3>
              <ul className="space-y-2 text-sm text-white/80">
                <li>
                  <a href="/" className="hover:text-white transition-colors">
                    Home
                  </a>
                </li>
                <li>
                  <a
                    href="/outputs"
                    className="hover:text-white transition-colors"
                  >
                    Outputs
                  </a>
                </li>
              </ul>
            </div>
            <div>
              <h3 className="font-semibold font-condensed text-lg mb-3">
                Contact
              </h3>
              <p className="text-sm text-white/70">
                {catalog?.title}
              </p>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
