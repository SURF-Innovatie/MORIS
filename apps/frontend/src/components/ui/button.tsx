import * as React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cva, type VariantProps } from 'class-variance-authority';

import { cn } from '../../lib/utils';

const buttonVariants = cva(
  'relative inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-lg font-medium transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-60',
  {
    variants: {
      variant: {
        default:
          'btn-gradient text-primary-foreground shadow-mynaui-sm hover:shadow-mynaui-md hover:translate-y-[-1px]',
        secondary:
          'bg-secondary/60 text-secondary-foreground border border-border hover:bg-secondary/80 shadow-mynaui-sm',
        outline:
          'bg-transparent text-foreground border border-border hover:bg-white/5',
        ghost:
          'bg-transparent text-foreground hover:bg-white/10 hover:text-primary transition-colors',
        destructive: 'bg-destructive text-destructive-foreground shadow-mynaui-sm hover:bg-destructive/90',
      },
      size: {
        default: 'h-11 px-5 text-sm',
        sm: 'h-9 rounded-md px-3 text-xs',
        lg: 'h-12 rounded-xl px-6 text-base',
        icon: 'h-10 w-10 rounded-xl',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : 'button';

    return (
      <Comp className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props} />
    );
  },
);
Button.displayName = 'Button';

export { Button, buttonVariants };
