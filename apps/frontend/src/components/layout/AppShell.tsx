import { useState } from "react";
import { Outlet } from "react-router-dom";
import { GlobalHeader } from "./GlobalHeader";
import { NotificationProvider } from "@/contexts/NotificationContext";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { AppNavigation } from "./AppNavigation";

export const AppShell = () => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  // Better handled by clicking item, but this is a fail-safe.

  return (
    <NotificationProvider>
      <div className="flex min-h-screen flex-col bg-background">
        <GlobalHeader onMenuClick={() => setIsSidebarOpen(true)} />

        <main className="flex-1 flex flex-col container mx-auto px-4 py-6 md:px-6 lg:px-8">
          <Outlet />
        </main>

        <Sheet open={isSidebarOpen} onOpenChange={setIsSidebarOpen}>
          <SheetContent side="left" className="w-[300px] sm:w-[350px] p-0">
            <SheetHeader className="p-4 border-b">
              <SheetTitle className="flex items-center gap-2">
                <span className="font-display text-lg tracking-tight">
                  MORIS
                </span>
              </SheetTitle>
            </SheetHeader>
            <div className="h-full overflow-y-auto pb-6">
              <AppNavigation onItemClick={() => setIsSidebarOpen(false)} />
            </div>
          </SheetContent>
        </Sheet>
      </div>
    </NotificationProvider>
  );
};
