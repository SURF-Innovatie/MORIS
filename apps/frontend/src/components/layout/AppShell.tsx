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
import { Star } from "lucide-react";

export const AppShell = () => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);
  // Better handled by clicking item, but this is a fail-safe.

  return (
    <NotificationProvider>
      <div className="flex min-h-screen flex-col bg-background">
        <GlobalHeader onMenuClick={() => setIsSidebarOpen(true)} />

        <main className="flex-1 flex flex-col">
          <Outlet />
        </main>

        <Sheet open={isSidebarOpen} onOpenChange={setIsSidebarOpen}>
          <SheetContent side="left" className="w-[300px] sm:w-[350px] p-0">
            <SheetHeader className="p-4 border-b">
              <SheetTitle className="flex items-center gap-2">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
                  <Star className="h-5 w-5 fill-current" />
                </div>
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
