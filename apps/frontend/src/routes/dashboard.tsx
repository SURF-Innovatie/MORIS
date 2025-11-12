import { ArrowUpRight, ChartColumn, Users as UsersIcon } from '@mynaui/icons-react';

import { Badge } from '../components/ui/badge';
import { Button } from '../components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card';

const highlightMetrics = [
  {
    title: 'Active members',
    value: '248',
    change: '+18% vs last week',
    icon: UsersIcon,
  },
  {
    title: 'Engagement rate',
    value: '64%',
    change: '+8 pts',
    icon: ChartColumn,
  },
  {
    title: 'Weekly growth',
    value: '12.3%',
    change: 'Steady climb',
    icon: ArrowUpRight,
  },
];

const DashboardRoute = () => {
  return (
    <div className="flex flex-col gap-12">
      <section className="glass-panel overflow-hidden rounded-3xl border border-white/10 bg-gradient-to-br from-primary/15 via-background to-accent/10 px-10 py-12 shadow-mynaui-md">
        <Badge className="mb-6 w-fit" variant="success">
          Beta workspace
        </Badge>
        <h1 className="max-w-2xl font-display text-4xl tracking-tight text-foreground sm:text-5xl">
          MynaUI foundations wired up for your MORIS frontend.
        </h1>
        <p className="mt-4 max-w-2xl text-lg text-muted-foreground">
          Vite, React Router v7, TanStack Query, and Orval are ready to go. Build fast, stay consistent,
          and ship beautiful interfaces without hunting for boilerplate.
        </p>
        <div className="mt-8 flex flex-wrap items-center gap-4">
          <Button size="lg">Create your first flow</Button>
          <Button size="lg" variant="secondary" asChild>
            <a href="https://orval.dev" target="_blank" rel="noreferrer">
              Explore Orval docs
            </a>
          </Button>
        </div>
      </section>

      <section className="grid gap-6 md:grid-cols-2 xl:grid-cols-3">
        {highlightMetrics.map(({ title, value, change, icon: Icon }) => (
          <Card key={title}>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>{title}</CardTitle>
                <span className="inline-flex h-10 w-10 items-center justify-center rounded-2xl bg-primary/15 text-primary">
                  <Icon className="h-5 w-5" aria-hidden />
                </span>
              </div>
              <CardDescription>{change}</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="text-3xl font-semibold text-foreground">{value}</p>
            </CardContent>
          </Card>
        ))}
      </section>
    </div>
  );
};

export default DashboardRoute;
