import { useState } from "react";
import { Outlet } from "react-router-dom";

import QueryIndicator from "@/components/status/query-indicator";
import { ExpandableNavbar } from "@/components/layout";

const RootLayout = () => {
  const [isNavExpanded, setIsNavExpanded] = useState(true);

  return (
    <div className="min-h-screen bg-background">
      <ExpandableNavbar onExpandChange={setIsNavExpanded} />

      {/* Main Content Area */}
      <div
        className="transition-all duration-300 lg:ml-20"
        style={{
          marginLeft:
            typeof window !== "undefined" && window.innerWidth >= 1024
              ? isNavExpanded
                ? "16rem"
                : "5rem"
              : "0",
        }}
      >
        <div className="min-h-screen flex flex-col pt-16 lg:pt-0">
          <main className="flex-1 px-6 py-10 lg:px-10 lg:py-12">
            <div className="mx-auto max-w-6xl">
              <Outlet />
            </div>
          </main>

          <footer className="mt-auto border-t border-white/10 px-6 py-6 lg:px-10">
            <div className="mx-auto max-w-6xl flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-muted-foreground">
              <span>Built with love and the shadcn UI design kit.</span>
              <span>&copy; {new Date().getFullYear()} MORIS</span>
            </div>
          </footer>
        </div>
      </div>

      <QueryIndicator />
    </div>
  );
};

export default RootLayout;
