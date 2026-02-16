import { Link, useNavigate } from "react-router-dom";
import { Search, Plus, Bell, Menu, LogOut, Settings, User } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useAuth } from "@/hooks/useAuth";
import { useLogout } from "@/hooks/useLogout";
import { useNotifications } from "@/contexts/NotificationContext";

export const GlobalHeader = ({ onMenuClick }: { onMenuClick?: () => void }) => {
  const { user } = useAuth();
  const { mutate: logout } = useLogout();
  const navigate = useNavigate();
  const { unreadCount } = useNotifications();

  const handleLogout = () => {
    logout(undefined, {
      onSuccess: () => {
        navigate("/", { replace: true });
      },
    });
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
      <div className="flex h-14 items-center gap-4 px-4 lg:px-6">
        {/* Mobile Menu Trigger (Visible on small screens) */}
        <Button
          variant="ghost"
          size="icon"
          className="md:hidden"
          onClick={onMenuClick}
        >
          <Menu className="h-5 w-5" />
          <span className="sr-only">Toggle menu</span>
        </Button>

        {/* Desktop Menu Trigger (Visible on large screens too per user request for global sidebar) */}
        <Button
          variant="ghost"
          size="icon"
          className="hidden md:flex bg-muted/20"
          onClick={onMenuClick}
        >
          <Menu className="h-5 w-5" />
          <span className="sr-only">Toggle menu</span>
        </Button>

        {/* Logo */}
        <Link to="/dashboard" className="flex items-center gap-2 font-semibold">
          <span className="hidden font-display text-lg tracking-tight sm:inline-block">
            MORIS
          </span>
        </Link>

        {/* Search Bar */}
        <div className="flex flex-1 items-center gap-2 md:gap-4 lg:ml-4">
          <Button
            variant="outline"
            className="relative h-9 w-full justify-start rounded-[0.5rem] bg-background text-sm font-normal text-muted-foreground shadow-none sm:pr-12 md:w-64 lg:w-80"
            onClick={() => {
              // TODO: Trigger Command Palette
              console.log("Open Command Palette");
            }}
          >
            <span className="inline-flex gap-2 items-center">
              <Search className="h-4 w-4" />
              Type{" "}
              <kbd className="pointer-events-none inline-flex h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
                <span className="text-xs">/</span>
              </kbd>{" "}
              to search...
            </span>
          </Button>
        </div>

        {/* Right Actions */}
        <div className="flex items-center gap-2">
          {/* Create New Dropdown */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-9 w-9">
                <Plus className="h-5 w-5" />
                <span className="sr-only">Create new...</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem
                onClick={() => navigate("/dashboard/projects/new")}
              >
                New Project
              </DropdownMenuItem>
              <DropdownMenuItem
                onClick={() => navigate("/dashboard/organisations/new")}
              >
                New Organisation
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          {/* Notifications */}
          <Button
            variant="ghost"
            size="icon"
            className="h-9 w-9 relative"
            onClick={() => navigate("/dashboard/inbox")}
          >
            <Bell className="h-5 w-5" />
            {unreadCount > 0 && (
              <span className="absolute top-1.5 right-1.5 h-2 w-2 rounded-full bg-primary ring-1 ring-background" />
            )}
            <span className="sr-only">Notifications</span>
          </Button>

          {/* User Profile */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="relative h-8 w-8 rounded-full">
                <Avatar className="h-8 w-8">
                  <AvatarImage
                    src={`https://avatar.vercel.sh/${user?.email}`}
                    alt={user?.name || ""}
                  />
                  <AvatarFallback>
                    {user?.name?.charAt(0).toUpperCase() || "U"}
                  </AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuLabel>
                <div className="flex flex-col space-y-1">
                  <p className="text-sm font-medium leading-none">
                    {user?.name}
                  </p>
                  <p className="text-xs leading-none text-muted-foreground">
                    {user?.email}
                  </p>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem onClick={() => navigate("/dashboard/settings")}>
                <User className="mr-2 h-4 w-4" />
                Profile
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => navigate("/dashboard/settings")}>
                <Settings className="mr-2 h-4 w-4" />
                Settings
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                className="text-red-500 focus:text-red-500 focus:bg-red-50"
                onClick={handleLogout}
              >
                <LogOut className="mr-2 h-4 w-4" />
                Log out
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  );
};
