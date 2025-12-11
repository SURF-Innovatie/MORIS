import { useState } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import {
  Grid,
  Star,
  LogOut,
  Menu,
  X,
  User,
  Inbox,
  Calendar,
  Folder,
  Package,
} from "lucide-react";

import { Button } from "../ui/button";
import { Badge } from "../ui/badge";
import { useAuth } from "@/hooks/useAuth";
import { useLogout } from "@/hooks/useLogout";
import { useNotifications } from "@/contexts/NotificationContext";

const NAV_GROUPS = [
  {
    label: "Workspace",
    items: [
      { to: "/dashboard", label: "Dashboard", icon: Grid },
      { to: "/dashboard/inbox", label: "Inbox", icon: Inbox, showBadge: true },
      { to: "/dashboard/calendar", label: "Calendar", icon: Calendar },
    ],
  },
  {
    label: "My Research",
    items: [
      { to: "/dashboard/projects", label: "Projects", icon: Folder },
      { to: "/dashboard/products", label: "Products", icon: Package },
      // { to: "/dashboard/datasets", label: "Datasets", icon: Database },
    ],
  },
  {
    label: "Settings",
    items: [
      { to: "/dashboard/profile", label: "Profile", icon: User },
    ],
  },
];

export const ExpandableNavbar = () => {
  const [isMobileOpen, setIsMobileOpen] = useState(false);
  const { user } = useAuth();
  const { mutate: logout } = useLogout();
  const navigate = useNavigate();
  const { unreadCount } = useNotifications();

  const navGroups = [
    ...NAV_GROUPS,
    ...(user?.is_sys_admin
      ? [
        {
          label: "Admin",
          items: [
            { to: "/dashboard/admin/users", label: "Users", icon: User },
          ],
        },
      ]
      : []),
  ];

  const handleLogout = () => {
    logout(undefined, {
      onSuccess: () => {
        navigate("/", { replace: true });
      },
    });
  };

  const toggleMobileMenu = () => {
    setIsMobileOpen(!isMobileOpen);
  };

  return (
    <>
      {/* Mobile Header */}
      <div className="lg:hidden fixed top-0 left-0 right-0 z-50 flex items-center justify-between bg-background border-r border-black/10 px-4 py-3">
        <NavLink
          to="/dashboard"
          className="flex items-center gap-2 text-lg font-semibold"
        >
          <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <Star className="h-6 w-6" aria-hidden />
          </span>
          <span className="font-display text-lg tracking-tight text-foreground">
            MORIS UI
          </span>
        </NavLink>
        <Button variant="ghost" size="sm" onClick={toggleMobileMenu}>
          {isMobileOpen ? (
            <X className="h-6 w-6" />
          ) : (
            <Menu className="h-6 w-6" />
          )}
        </Button>
      </div>

      {/* Mobile Menu Overlay */}
      {isMobileOpen && (
        <div className="lg:hidden fixed inset-0 z-40 bg-background/95 backdrop-blur-sm pt-16">
          <div className="flex flex-col h-full p-4">
            <nav className="flex flex-col gap-6 mb-6 overflow-y-auto">
              {navGroups.map((group) => (
                <div key={group.label} className="flex flex-col gap-2">
                  <div className="px-4 text-xs font-semibold text-muted-foreground/70 uppercase tracking-wider">
                    {group.label}
                  </div>
                  {group.items.map(({ to, label, icon: Icon, showBadge }) => (
                    <NavLink
                      key={to}
                      to={to}
                      end={to === "/dashboard"}
                      onClick={() => setIsMobileOpen(false)}
                      className={({ isActive }) =>
                        [
                          "flex items-center gap-3 rounded-xl px-4 py-3 transition-colors duration-200",
                          isActive
                            ? "bg-primary/15 text-primary"
                            : "text-muted-foreground hover:bg-white/10 hover:text-foreground",
                        ].join(" ")
                      }
                    >
                      <Icon className="h-5 w-5" aria-hidden />
                      <span className="font-medium">{label}</span>
                      {showBadge && unreadCount > 0 && (
                        <Badge
                          variant="secondary"
                          className="ml-auto h-5 px-1.5 text-[10px]"
                        >
                          {unreadCount}
                        </Badge>
                      )}
                    </NavLink>
                  ))}
                </div>
              ))}
            </nav>

            <div className="mt-auto space-y-4">
              {user && (
                <div className="px-4 py-2 text-sm text-muted-foreground">
                  Logged in as{" "}
                  <span className="font-medium text-foreground">
                    {user?.email}
                  </span>
                </div>
              )}
              <Button
                size="sm"
                variant="destructive"
                className="w-full justify-start"
                onClick={handleLogout}
              >
                <LogOut className="mr-2 h-4 w-4" aria-hidden /> Logout
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Desktop Sidebar */}
      <aside className="hidden lg:flex fixed left-0 top-0 bottom-0 z-40 flex-col bg-background border-r border-black/10 w-64">
        {/* Logo Section */}
        <div className="flex items-center justify-between p-4 border-b border-black/10">
          <NavLink to="/dashboard" className="flex items-center gap-3">
            <div className="flex flex-col">
              <span className="font-display text-lg tracking-tight text-foreground">
                MORIS
              </span>
            </div>
          </NavLink>
        </div>

        {/* Navigation */}
        <nav className="flex flex-col gap-6 p-3 flex-1 overflow-y-auto">
          {navGroups.map((group, groupIndex) => (
            <div key={group.label} className="flex flex-col gap-2">
              <div className="px-3 text-xs font-semibold text-muted-foreground/70 uppercase tracking-wider">
                {group.label}
              </div>
              {group.items.map(({ to, label, icon: Icon, showBadge }) => (
                <NavLink
                  key={to}
                  to={to}
                  end={to === "/dashboard"}
                  className={({ isActive }) =>
                    [
                      "flex items-center gap-3 rounded-xl px-3 py-2 transition-colors duration-200",
                      isActive
                        ? "bg-primary/15 text-primary"
                        : "text-muted-foreground hover:bg-white/10 hover:text-foreground",
                    ].join(" ")
                  }
                >
                  <Icon className="h-5 w-5 flex-shrink-0" aria-hidden />
                  <span className="font-medium flex-1">{label}</span>
                  {showBadge && unreadCount > 0 && (
                    <Badge
                      variant="secondary"
                      className="ml-auto h-5 px-1.5 text-[10px]"
                    >
                      {unreadCount}
                    </Badge>
                  )}
                </NavLink>
              ))}
              {/* Add separator between groups if not the last one */}
              {groupIndex < navGroups.length - 1 && (
                <div className="mx-3 my-1 border-t border-white/5" />
              )}
            </div>
          ))}
        </nav>

        {/* Bottom Section */}
        <div className="p-3 border-t border-white/10 space-y-2">
          {user && (
            <div className="px-3 py-2 text-sm text-muted-foreground truncate">
              {user?.email}
            </div>
          )}
          <Button
            size="sm"
            variant="destructive"
            className="w-full justify-start"
            onClick={handleLogout}
          >
            <LogOut className="h-4 w-4 mr-2" aria-hidden />
            Logout
          </Button>
        </div>
      </aside>
    </>
  );
};
