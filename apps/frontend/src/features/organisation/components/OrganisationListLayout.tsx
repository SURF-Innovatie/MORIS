import { ReactNode } from "react";

interface OrganisationListLayoutProps {
  title: string;
  headerActions?: ReactNode;
  isLoading?: boolean;
  isEmpty?: boolean;
  emptyMessage?: ReactNode;
  children: ReactNode;
}

export const OrganisationListLayout = ({
  title,
  headerActions,
  isLoading,
  isEmpty,
  emptyMessage,
  children,
}: OrganisationListLayoutProps) => {
  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">{title}</h1>
        {headerActions}
      </div>
      <div className="space-y-2">
        {children}
        {isEmpty && (
          <div className="text-gray-500">
            {emptyMessage || "No organizations found."}
          </div>
        )}
      </div>
    </div>
  );
};
