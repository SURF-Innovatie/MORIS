import { useState, useEffect } from 'react';
import { NavLink, useNavigate } from 'react-router-dom';
import { Grid, Star, Zap, LogOut, Menu, X, ChevronLeft, ChevronRight } from 'lucide-react';

import { Button } from '../ui/button';
import { useAuth } from '../../hooks/useAuth';
import { useLogout } from '../../hooks/useLogout';

const NAV_ITEMS = [
  { to: '/dashboard', label: 'Dashboard', icon: Grid },
];

interface ExpandableNavbarProps {
  onExpandChange?: (isExpanded: boolean) => void;
}

export const ExpandableNavbar = ({ onExpandChange }: ExpandableNavbarProps) => {
  const [isExpanded, setIsExpanded] = useState(true);
  const [isMobileOpen, setIsMobileOpen] = useState(false);
  const { user } = useAuth();
  const { mutate: logout } = useLogout();
  const navigate = useNavigate();

  useEffect(() => {
    onExpandChange?.(isExpanded);
  }, [isExpanded, onExpandChange]);

  const handleLogout = () => {
    logout(undefined, {
      onSuccess: () => {
        navigate('/', { replace: true });
      },
    });
  };

  const toggleSidebar = () => {
    setIsExpanded(!isExpanded);
  };

  const toggleMobileMenu = () => {
    setIsMobileOpen(!isMobileOpen);
  };

  return (
    <>
      {/* Mobile Header */}
      <div className="lg:hidden fixed top-0 left-0 right-0 z-50 flex items-center justify-between bg-background border-b border-white/10 px-4 py-3">
        <NavLink to="/dashboard" className="flex items-center gap-2 text-lg font-semibold">
          <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
            <Star className="h-6 w-6" aria-hidden />
          </span>
          <span className="font-display text-lg tracking-tight text-foreground">MORIS UI</span>
        </NavLink>
        <Button variant="ghost" size="sm" onClick={toggleMobileMenu}>
          {isMobileOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
        </Button>
      </div>

      {/* Mobile Menu Overlay */}
      {isMobileOpen && (
        <div className="lg:hidden fixed inset-0 z-40 bg-background/95 backdrop-blur-sm pt-16">
          <div className="flex flex-col h-full p-4">
            <nav className="flex flex-col gap-2 mb-6">
              {NAV_ITEMS.map(({ to, label, icon: Icon }) => (
                <NavLink
                  key={to}
                  to={to}
                  end={to === '/dashboard'}
                  onClick={() => setIsMobileOpen(false)}
                  className={({ isActive }) =>
                    [
                      'flex items-center gap-3 rounded-xl px-4 py-3 transition-colors duration-200',
                      isActive
                        ? 'bg-primary/15 text-primary'
                        : 'text-muted-foreground hover:bg-white/10 hover:text-foreground',
                    ].join(' ')
                  }
                >
                  <Icon className="h-5 w-5" aria-hidden />
                  <span className="font-medium">{label}</span>
                </NavLink>
              ))}
            </nav>

            <div className="mt-auto space-y-4">
              {user && (
                <div className="px-4 py-2 text-sm text-muted-foreground">
                  Logged in as <span className="font-medium text-foreground">{user.email}</span>
                </div>
              )}
              <Button variant="ghost" size="sm" className="w-full justify-start" asChild>
                <a href="https://ui.shadcn.com/components" target="_blank" rel="noreferrer">
                  <Zap className="mr-2 h-4 w-4" aria-hidden /> Browse components
                </a>
              </Button>
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
      <aside
        className={`hidden lg:flex fixed left-0 top-0 bottom-0 z-40 flex-col bg-background border-r border-white/10 transition-all duration-300 ${
          isExpanded ? 'w-64' : 'w-20'
        }`}
      >
        {/* Logo Section */}
        <div className="flex items-center justify-between p-4 border-b border-white/10">
          {isExpanded ? (
            <NavLink to="/dashboard" className="flex items-center gap-3">
              <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
                <Star className="h-6 w-6" aria-hidden />
              </span>
              <div className="flex flex-col">
                <span className="font-display text-lg tracking-tight text-foreground">MORIS UI</span>
                <span className="text-[10px] uppercase tracking-[0.3em] text-muted-foreground">
                  Powered by shadcn
                </span>
              </div>
            </NavLink>
          ) : (
            <NavLink to="/dashboard" className="mx-auto">
              <span className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary">
                <Star className="h-6 w-6" aria-hidden />
              </span>
            </NavLink>
          )}
        </div>

        {/* Toggle Button */}
        <button
          onClick={toggleSidebar}
          className="absolute -right-3 top-20 flex h-6 w-6 items-center justify-center rounded-full bg-primary text-primary-foreground shadow-lg hover:bg-primary/90 transition-colors"
          aria-label={isExpanded ? 'Collapse sidebar' : 'Expand sidebar'}
        >
          {isExpanded ? <ChevronLeft className="h-4 w-4" /> : <ChevronRight className="h-4 w-4" />}
        </button>

        {/* Navigation */}
        <nav className="flex flex-col gap-2 p-3 flex-1">
          {NAV_ITEMS.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/dashboard'}
              className={({ isActive }) =>
                [
                  'flex items-center gap-3 rounded-xl px-3 py-2 transition-colors duration-200',
                  isActive
                    ? 'bg-primary/15 text-primary'
                    : 'text-muted-foreground hover:bg-white/10 hover:text-foreground',
                  !isExpanded && 'justify-center',
                ].join(' ')
              }
              title={!isExpanded ? label : undefined}
            >
              <Icon className="h-5 w-5 flex-shrink-0" aria-hidden />
              {isExpanded && <span className="font-medium">{label}</span>}
            </NavLink>
          ))}
        </nav>

        {/* Bottom Section */}
        <div className="p-3 border-t border-white/10 space-y-2">
          {user && isExpanded && (
            <div className="px-3 py-2 text-sm text-muted-foreground truncate">
              {user.email}
            </div>
          )}
          <Button
            variant="ghost"
            size="sm"
            className={`w-full ${isExpanded ? 'justify-start' : 'justify-center'}`}
            asChild
            title={!isExpanded ? 'Browse components' : undefined}
          >
            <a href="https://ui.shadcn.com/components" target="_blank" rel="noreferrer">
              <Zap className={`h-4 w-4 ${isExpanded ? 'mr-2' : ''}`} aria-hidden />
              {isExpanded && 'Browse components'}
            </a>
          </Button>
          <Button
            size="sm"
            variant="destructive"
            className={`w-full ${isExpanded ? 'justify-start' : 'justify-center'}`}
            onClick={handleLogout}
            title={!isExpanded ? 'Logout' : undefined}
          >
            <LogOut className={`h-4 w-4 ${isExpanded ? 'mr-2' : ''}`} aria-hidden />
            {isExpanded && 'Logout'}
          </Button>
        </div>
      </aside>
    </>
  );
};
