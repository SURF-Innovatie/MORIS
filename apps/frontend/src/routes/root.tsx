import { NavLink, Outlet } from 'react-router-dom';
import { Grid, Users, Star, Zap } from 'lucide-react';

import { Button } from '../components/ui/button';
import QueryIndicator from '../components/status/query-indicator';

const NAV_ITEMS = [
  { to: '/', label: 'Overview', icon: Grid },
  { to: '/users', label: 'Users', icon: Users },
];

const RootLayout = () => {
  return (
    <div className="min-h-screen bg-background">
      <div className="mx-auto flex min-h-screen w-full max-w-6xl flex-col px-6 py-10">
        <header className="flex items-center justify-between gap-6">
          <NavLink to="/" className="flex items-center gap-3 text-lg font-semibold">
              <span className="flex h-12 w-12 items-center justify-center rounded-2xl bg-primary/10 text-primary">
              <Star className="h-7 w-7" aria-hidden />
            </span>
            <div className="flex flex-col">
              <span className="font-display text-xl tracking-tight text-foreground">MORIS UI</span>
              <span className="text-xs uppercase tracking-[0.3em] text-muted-foreground">
                Powered by shadcn UI
              </span>
            </div>
          </NavLink>
          <div className="flex items-center gap-3">
            <Button variant="ghost" size="sm" asChild>
        <a href="https://ui.shadcn.com/components" target="_blank" rel="noreferrer">
          <Zap className="mr-1 h-4 w-4" aria-hidden /> Browse components
              </a>
            </Button>
            <Button size="sm">Deploy</Button>
          </div>
        </header>

  <nav className="mt-10 flex items-center gap-3 rounded-2xl border border-white/10 bg-white/5 p-2 text-sm shadow-sm backdrop-blur">
          {NAV_ITEMS.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/'}
              className={({ isActive }) =>
                [
                  'flex items-center gap-2 rounded-xl px-4 py-2 transition-colors duration-200',
                  isActive
                    ? 'bg-primary/15 text-primary'
                    : 'text-muted-foreground hover:bg-white/10 hover:text-foreground',
                ].join(' ')
              }
            >
              <Icon className="h-4 w-4" aria-hidden />
              <span className="font-medium">{label}</span>
            </NavLink>
          ))}
        </nav>

        <main className="flex flex-1 flex-col py-12">
          <Outlet />
        </main>

        <footer className="mt-auto flex items-center justify-between border-t border-white/10 pt-6 text-xs text-muted-foreground">
          <span>Built with love and the shadcn UI design kit.</span>
          <span>&copy; {new Date().getFullYear()} MORIS</span>
        </footer>
      </div>
      <QueryIndicator />
    </div>
  );
};

export default RootLayout;
