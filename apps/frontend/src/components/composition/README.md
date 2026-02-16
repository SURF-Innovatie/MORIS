# Composition Components

Higher-level components that provide consistent patterns across the MORIS application. These components are built on top of Shadcn UI primitives and encapsulate common UI patterns.

## Philosophy

- **Consistency**: Unified patterns reduce cognitive load and improve UX
- **Reusability**: One component, many use cases
- **Composability**: Works well with other components
- **Accessibility**: Built on accessible Shadcn primitives

## Components

### StatusBadge

Standardized status indicator for pending events, approval states, and other status indicators.

```tsx
import { StatusBadge } from "@/components/composition";

// Pending event
<StatusBadge status="pending" size="xs" />

// Approved state
<StatusBadge status="approved" showIcon={false} />

// Custom label
<StatusBadge status="rejected" label="Denied" />
```

**Props:**
- `status`: "pending" | "approved" | "rejected" | "info"
- `size`: "xs" | "sm" | "md" (default: "xs")
- `showIcon`: boolean (default: true)
- `label`: string (optional, overrides default)

### EmptyState

Consistent empty state for lists, tables, and content areas.

```tsx
import { EmptyState } from "@/components/composition";
import { Building2 } from "lucide-react";

<EmptyState
  icon={Building2}
  title="No projects found"
  description="Create your first project to get started"
  action={{ label: "Create Project", onClick: handleCreate }}
/>
```

**Props:**
- `icon`: LucideIcon (optional)
- `title`: string (required)
- `description`: string (optional)
- `action`: { label, onClick } (optional)
- `size`: "sm" | "default" | "lg"

### ActionMenu

Unified dropdown menu with standardized MoreHorizontal trigger.

```tsx
import { ActionMenu } from "@/components/composition";
import { Edit, Trash } from "lucide-react";

// Simple menu
<ActionMenu
  sections={[
    { label: "Edit", onClick: handleEdit },
    { label: "Delete", onClick: handleDelete, destructive: true }
  ]}
/>

// Sectioned menu
<ActionMenu
  sections={[
    {
      label: "Actions",
      items: [
        { label: "Edit", icon: Edit, onClick: handleEdit },
        { label: "Share", icon: Share, onClick: handleShare }
      ]
    },
    {
      items: [
        { label: "Delete", icon: Trash, onClick: handleDelete, destructive: true }
      ]
    }
  ]}
/>
```

**Props:**
- `sections`: ActionMenuSection[] | ActionMenuItem[]
- `align`: "start" | "center" | "end" (default: "end")
- `size`: "default" | "sm" | "icon"
- `trigger`: React.ReactNode (optional custom trigger)

### ListSkeleton

Loading states that match ListItem variants.

```tsx
import { ListSkeleton } from "@/components/composition";

// Compact list loading
<ListSkeleton variant="compact" count={5} />

// Default list loading
<ListSkeleton variant="default" count={3} />

// Detailed card grid loading
<div className="grid gap-4 md:grid-cols-2">
  <ListSkeleton variant="detailed" count={4} />
</div>
```

**Props:**
- `variant`: "compact" | "default" | "detailed"
- `count`: number (default: 3)

### ListItem

Unified list item component with three variants.

```tsx
import { ListItem } from "@/components/composition";
import { Book } from "lucide-react";

// Compact (sidebar navigation)
<ListItem
  variant="compact"
  title="Project Alpha"
  icon={Book}
  onClick={() => navigate('/project/1')}
/>

// Default (member list)
<ListItem
  title="John Doe"
  subtitle="john@example.com"
  avatarUrl="..."
  badges={[{ label: "Lead", variant: "secondary" }]}
  pending={true}
  action={<ActionMenu sections={...} />}
/>

// Detailed (project card)
<ListItem
  variant="detailed"
  title="Research Project"
  subtitle="Active project with 5 members"
  badges={[{ label: "Active" }]}
  onClick={() => navigate('/project/1')}
>
  <div className="mt-4">Additional details...</div>
</ListItem>
```

**Variants:**
- **compact**: Minimal, single-line (sidebars, navigation)
- **default**: Standard with avatar and description (lists)
- **detailed**: Rich card-style with additional content (grids)

**Props:**
- `variant`: "compact" | "default" | "detailed"
- `title`: string (required)
- `subtitle`: string (optional)
- `avatarUrl`: string (optional)
- `avatarFallback`: string (optional)
- `icon`: LucideIcon (optional, replaces avatar)
- `badges`: Array of badge objects
- `pending`: boolean
- `onClick`: () => void (optional)
- `action`: React.ReactNode (optional)
- `children`: React.ReactNode (detailed variant only)

## Design Tokens

See `src/lib/design-tokens.ts` for standardized sizing:

```typescript
import { designTokens } from "@/lib/design-tokens";

// Button sizes
designTokens.button.xs    // h-7
designTokens.button.sm    // h-8
designTokens.button.iconSm // h-8 w-8

// Avatar sizes
designTokens.avatar.xs    // h-5 w-5
designTokens.avatar.md    // h-9 w-9

// Badge sizes
designTokens.badge.xs     // { height: 'h-5', text: 'text-[10px]', ... }
```

## Migration Guide

When migrating existing components:

1. **Identify the pattern** - Is it a list item? Status indicator? Action menu?
2. **Choose the right composition component** - Match variant to use case
3. **Replace custom implementation** - Use composition component with appropriate props
4. **Test thoroughly** - Verify visual consistency and behavior

### Example Migration

**Before:**
```tsx
<div className="flex items-center gap-4 p-4 border rounded-lg">
  <Avatar className="h-10 w-10">
    <AvatarImage src={member.avatarUrl} />
    <AvatarFallback>{member.initials}</AvatarFallback>
  </Avatar>
  <div>
    <p className="font-semibold">{member.name}</p>
    <p className="text-sm text-muted-foreground">{member.email}</p>
  </div>
  {member.pending && (
    <Badge variant="outline" className="border-yellow-500">Pending</Badge>
  )}
</div>
```

**After:**
```tsx
<ListItem
  title={member.name}
  subtitle={member.email}
  avatarUrl={member.avatarUrl}
  avatarFallback={member.initials}
  pending={member.pending}
/>
```

## Best Practices

1. **Always use composition components** when the pattern exists
2. **Consistent sizing** - Use design tokens for custom components
3. **Don't override styles** unless absolutely necessary
4. **Extend via composition** - Wrap composition components rather than forking them
5. **Document new patterns** - If you create a new pattern, consider adding it here

## Contributing

When adding new composition components:

1. Follow existing patterns (StatusBadge, ListItem as examples)
2. Use design tokens for sizing
3. Support multiple variants where appropriate
4. Include comprehensive JSDoc with examples
5. Update this README with usage documentation
6. Export from `index.ts`
