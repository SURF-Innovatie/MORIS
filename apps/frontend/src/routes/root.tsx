import { Outlet } from "react-router-dom";

import QueryIndicator from "@/components/status/query-indicator";
import { ExpandableNavbar } from "@/components/layout";
import { Breadcrumbs } from "@/components/breadcrumbs";
import { NotificationProvider } from "@/context/NotificationContext";

const RootLayout = () => {
  return (
    <NotificationProvider>
      <div className="min-h-screen bg-background">
        <ExpandableNavbar />

        {/* Main Content Area */}
        <div className="transition-all duration-300 lg:ml-64">
          <div className="min-h-screen flex flex-col pt-16 lg:pt-0">
            <main className="flex-1 px-6 py-10 lg:px-10 lg:py-12">
              <div className="mx-auto max-w-7xl w-full p-6">
                <Breadcrumbs />
                <Outlet />
              </div>
            </main>

            <footer className="mt-auto border-t border-white/10 px-6 py-6 lg:px-10">
              <div className="mx-auto max-w-6xl flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-muted-foreground">
                <span>&copy; {new Date().getFullYear()} SURF</span>
              </div>
            </footer>
          </div>
        </div>

        <QueryIndicator />
      </div>
    </NotificationProvider>
  );
};

export default RootLayout;
