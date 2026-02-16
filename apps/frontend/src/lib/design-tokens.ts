/**
 * Design Tokens - Centralized size and spacing conventions
 *
 * This module provides consistent sizing tokens across the application
 * to ensure visual harmony and maintainability. Use these tokens when
 * creating or updating components.
 */

export const designTokens = {
  /**
   * Button size variants
   * - xs: Extra small (h-7) - Compact actions, tight spaces
   * - sm: Small (h-8) - Secondary actions
   * - default: Default (h-9) - Primary actions
   * - lg: Large (h-10) - Hero CTAs
   * - iconSm: Icon button (h-8 w-8) - Icon-only buttons
   */
  button: {
    xs: "h-7",
    sm: "h-8",
    default: "h-9",
    lg: "h-10",
    iconSm: "h-8 w-8",
  },

  /**
   * Avatar size variants
   * - xs: Extra small (h-5 w-5) - Inline text, chips
   * - sm: Small (h-6 w-6) - Compact lists
   * - md: Medium (h-9 w-9) - Default lists
   * - lg: Large (h-10 w-10) - Featured items, headers
   * - xl: Extra large (h-12 w-12) - Profile headers
   */
  avatar: {
    xs: "h-5 w-5",
    sm: "h-6 w-6",
    md: "h-9 w-9",
    lg: "h-10 w-10",
    xl: "h-12 w-12",
  },

  /**
   * Badge size variants
   * - xs: Extra small (h-5 text-[10px]) - Inline status, counts
   * - sm: Small (h-6 text-xs) - List items
   * - md: Medium (h-7 text-sm) - Cards, default
   */
  badge: {
    xs: { height: "h-5", text: "text-[10px]", padding: "px-1.5" },
    sm: { height: "h-6", text: "text-xs", padding: "px-2" },
    md: { height: "h-7", text: "text-sm", padding: "px-2.5" },
  },

  /**
   * Spacing scale
   * Consistent spacing for margins, padding, gaps
   */
  spacing: {
    xs: "0.5rem", // 8px
    sm: "0.75rem", // 12px
    md: "1rem", // 16px
    lg: "1.5rem", // 24px
    xl: "2rem", // 32px
  },

  /**
   * Border radius scale
   */
  radius: {
    sm: "rounded-sm", // 2px
    default: "rounded-md", // 4px
    lg: "rounded-lg", // 8px
    full: "rounded-full", // 9999px
  },
} as const;

/**
 * Helper to combine badge size classes
 */
export function getBadgeSizeClasses(size: keyof typeof designTokens.badge) {
  const { height, text, padding } = designTokens.badge[size];
  return `${height} ${text} ${padding}`;
}
