import { BudgetCategory, FundingSource } from "../api/generated-orval/model";

// Shared constants for the application

export const EMPTY_UUID = "00000000-0000-0000-0000-000000000000";
export const STORAGE_KEY_AUTH_TOKEN = "moris_auth_token";
export const STORAGE_KEY_AUTH_USER = "moris_auth_user";
export const EVENT_NOTIFICATIONS_SHOULD_REFRESH =
  "notifications:should-refresh";

export const categoryLabels: Record<BudgetCategory, string> = {
  personnel: "Personnel",
  material: "Material",
  investment: "Investment",
  travel: "Travel",
  management: "Management",
  other: "Other",
};

export const fundingSourceLabels: Record<FundingSource, string> = {
  subsidy: "Subsidy",
  cofinancing_cash: "Co-financing (Cash)",
  cofinancing_inkind: "Co-financing (In-Kind)",
};

export const statusLabels: Record<string, string> = {
  draft: "Draft",
  submitted: "Submitted",
  approved: "Approved",
  locked: "Locked",
};

export const healthStatusColors: Record<string, string> = {
  on_track: "#22c55e", // green
  warning: "#eab308", // yellow
  at_risk: "#ef4444", // red
};

export const healthStatusLabels: Record<string, string> = {
  on_track: "On Track",
  warning: "Warning",
  at_risk: "At Risk",
};

export const healthStatusEmoji: Record<string, string> = {
  on_track: "🟢",
  warning: "🟡",
  at_risk: "🔴",
};
