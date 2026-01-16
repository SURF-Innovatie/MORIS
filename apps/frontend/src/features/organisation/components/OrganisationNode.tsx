import { useState } from "react";
import { useGetOrganisationNodesIdChildren } from "@api/moris";
import { OrganisationResponse } from "@api/model";
import { ChevronRight, ChevronDown } from "lucide-react";
import { RorDisplay } from "./RorDisplay";

interface OrganisationNodeProps {
  node: OrganisationResponse;
  renderActions?: (node: OrganisationResponse) => React.ReactNode;
  defaultExpanded?: boolean;
}

export const OrganisationNode = ({
  node,
  renderActions,
  defaultExpanded = false,
}: OrganisationNodeProps) => {
  const [isExpanded, setIsExpanded] = useState(defaultExpanded);
  const { data: children, isFetched } = useGetOrganisationNodesIdChildren(
    node.id!
  );

  const isEmpty = isFetched && children?.length === 0;

  return (
    <div className="ml-4 border-l pl-4">
      <div className="flex items-center gap-2 py-2">
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className={`p-1 hover:bg-gray-100 rounded ${isEmpty ? "invisible" : ""}`}
          disabled={isEmpty}
        >
          {isExpanded ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
        </button>
        <span className="font-medium">{node.name}</span>
        {node.rorId && (
            <div className="ml-2">
                <RorDisplay rorId={node.rorId} />
            </div>
        )}

        {renderActions && (
          <div className="ml-auto flex gap-2">{renderActions(node)}</div>
        )}
      </div>
      {isExpanded && !isEmpty && (
        <div className="ml-4">
          {children?.map((child: OrganisationResponse) => (
            <OrganisationNode
              key={child.id}
              node={child}
              renderActions={renderActions}
            />
          ))}
        </div>
      )}
    </div>
  );
};
