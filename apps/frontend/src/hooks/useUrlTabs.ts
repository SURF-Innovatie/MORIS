import { useSearchParams } from "react-router-dom";
import { useCallback } from "react";

export function useUrlTabs(defaultTab: string) {
  const [searchParams, setSearchParams] = useSearchParams();
  const currentTab = searchParams.get("tab") || defaultTab;

  const setTab = useCallback(
    (tab: string) => setSearchParams({ tab }),
    [setSearchParams],
  );

  return [currentTab, setTab] as const;
}
