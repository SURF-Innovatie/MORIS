import { createContext, useContext, useEffect, type ReactNode } from "react";
import { useGetCatalogsId } from "@/api/generated-orval/catalogs/catalogs";
import type { InternalAppCatalogCatalogDetails } from "@/api/generated-orval/model";

// Default NWO Colors
const DEFAULT_PRIMARY = "#008094";
const DEFAULT_SECONDARY = "#004c5a";
const DEFAULT_ACCENT = "#f59e0b";

interface CatalogContextType {
  data: InternalAppCatalogCatalogDetails | undefined;
  isLoading: boolean;
  error: unknown;
}

const CatalogContext = createContext<CatalogContextType | undefined>(undefined);

export function CatalogProvider({ children }: { children: ReactNode }) {
  const catalogId = import.meta.env.VITE_CATALOG_ID;

  const {
    data: axiosData,
    isLoading,
    error,
  } = useGetCatalogsId(catalogId, {
    query: {
      enabled: !!catalogId,
      // staleTime: 0, // Default is 0, which means "stale immediately".
      // This combined with persistence means: serve from cache (if exists), then refetch in background.
    },
  });

  const data = axiosData?.data;

  useEffect(() => {
    if (data && data.catalog) {
      const root = document.documentElement;
      root.style.setProperty(
        "--color-primary",
        data.catalog.primary_color || DEFAULT_PRIMARY,
      );
      root.style.setProperty(
        "--color-secondary",
        data.catalog.secondary_color || DEFAULT_SECONDARY,
      );
      root.style.setProperty(
        "--color-accent",
        data.catalog.accent_color || DEFAULT_ACCENT,
      );
      // Load heading fonts (Saira Condensed family)
      const headingLink = document.createElement("link");
      headingLink.href =
        "https://fonts.googleapis.com/css2?family=Saira+Condensed:wght@300;400;500;600;700&family=Saira+Extra+Condensed:wght@400;500;600;700&display=swap";
      headingLink.rel = "stylesheet";
      document.head.appendChild(headingLink);

      // Load the catalog's configured body font
      if (data.catalog.font_family) {
        const fontName = data.catalog.font_family.replace(/ /g, "+");
        const bodyLink = document.createElement("link");
        bodyLink.href = `https://fonts.googleapis.com/css2?family=${fontName}:ital,wght@0,300;0,400;0,500;0,600;0,700;1,400&display=swap`;
        bodyLink.rel = "stylesheet";
        document.head.appendChild(bodyLink);
        root.style.setProperty("--font-sans", `'${data.catalog.font_family}', sans-serif`);
      }
      if (data.catalog.title) {
        document.title = data.catalog.title;
      }
    }
  }, [data]);

  if (!catalogId) {
    return (
      <div className="min-h-screen flex items-center justify-center text-red-500">
        VITE_CATALOG_ID is not set
      </div>
    );
  }

  return (
    <CatalogContext.Provider value={{ data, isLoading, error }}>
      {children}
    </CatalogContext.Provider>
  );
}

export const useCatalog = () => {
  const context = useContext(CatalogContext);
  if (context === undefined) {
    throw new Error("useCatalog must be used within a CatalogProvider");
  }
  return context;
};
