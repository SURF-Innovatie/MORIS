import { useGetOrganisationNodesTree } from "@api/moris";
import { OrganisationTreeNode } from "@api/model";
import Tree from "react-d3-tree";
import { useMemo, useState, useEffect, useRef } from "react";
import { Loader2, ZoomIn, ZoomOut, RotateCcw } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { useNavigate } from "react-router-dom";
import { Label } from "@/components/ui/label";

interface OrganisationTreeViewProps {
  height?: number;
}

// Extend TreeDatum to include id explicitly for type safety/clarity in our transform
interface CustomTreeDatum {
  name: string;
  attributes?: Record<string, string | number | boolean>;
  children?: CustomTreeDatum[];
  // Extra fields for custom rendering
  id: string;
  rorId?: string;
}

const transformNode = (node: OrganisationTreeNode): CustomTreeDatum => {
  const attributes: Record<string, string | number | boolean> = {};
  if (node.rorId) {
    attributes["ROR ID"] = node.rorId;
  }
  return {
    name: node.name || "Unnamed",
    attributes,
    children: node.children?.map(transformNode) || [],
    id: node.id || "",
    rorId: node.rorId || undefined,
  };
};

export const OrganisationTreeView = ({
  height = 600,
}: OrganisationTreeViewProps) => {
  const { data: treeData, isLoading, error } = useGetOrganisationNodesTree();
  const [selectedRootId, setSelectedRootId] = useState<string | null>(null);
  const [translate, setTranslate] = useState({ x: 400, y: 50 });
  const [zoom, setZoom] = useState(1);
  const containerRef = useRef<HTMLDivElement>(null);
  const navigate = useNavigate();

  useEffect(() => {
    if (treeData && treeData.length > 0) {
      if (!selectedRootId) {
        // Default to first root if none selected
        setSelectedRootId(treeData[0].id || null);
      } else {
        // Verify selected root still exists (e.g. after refetch)
        const exists = treeData.find((n) => n.id === selectedRootId);
        if (!exists) {
          setSelectedRootId(treeData[0].id || null);
        }
      }
    }
  }, [treeData, selectedRootId]);

  const selectedNode = useMemo(() => {
    if (!treeData) return null;
    return treeData.find((n) => n.id === selectedRootId);
  }, [treeData, selectedRootId]);

  const formattedData = useMemo(() => {
    if (!selectedNode) return [];
    return [transformNode(selectedNode)];
  }, [selectedNode]);

  const centerView = () => {
    if (containerRef.current) {
      const { width } = containerRef.current.getBoundingClientRect();
      setTranslate({ x: width / 2, y: 50 });
      setZoom(1);
    }
  };

  // Center view only when the selected root changes (and data is ready)
  useEffect(() => {
    if (formattedData.length > 0) {
      centerView();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedRootId]);

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const renderCustomNode = ({ nodeDatum, toggleNode }: any) => {
    // nodeDatum is the data object we passed in, plus internal d3 stuff
    const hasChildren = nodeDatum.children && nodeDatum.children.length > 0;
    const isLeaf = !hasChildren;

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const handleNavigate = (e: any) => {
      e.stopPropagation();
      if (nodeDatum.id) {
        navigate(`/dashboard/admin/organisations/${nodeDatum.id}/roles`);
      }
    };

    return (
      <g>
        {/* Node circle - Navigates on click */}
        <circle
          r={12}
          fill={isLeaf ? "#fff" : "#777"}
          stroke={isLeaf ? "#777" : "#555"}
          strokeWidth={1}
          onClick={handleNavigate}
          className="cursor-pointer hover:stroke-primary hover:stroke-2 transition-all"
        />

        {/* ForeignObject for HTML content (text wrapping) */}
        <foreignObject x="20" y="-20" width="200" height="100">
          <div
            className="flex flex-col justify-center text-xs border border-transparent hover:border-border hover:bg-accent/50 rounded p-1 cursor-pointer transition-colors"
            onClick={handleNavigate}
            title="Click to edit organisation"
          >
            <span className="font-semibold break-words whitespace-normal leading-tight">
              {nodeDatum.name}
            </span>
            {nodeDatum.attributes?.["ROR ID"] && (
              <span className="text-[10px] text-muted-foreground mt-0.5">
                {nodeDatum.attributes["ROR ID"]}
              </span>
            )}
          </div>
        </foreignObject>
      </g>
    );
  };

  if (isLoading) {
    return (
      <div className="flex h-full w-full items-center justify-center p-8">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex h-full w-full items-center justify-center p-8 text-destructive">
        Failed to load organisation tree
      </div>
    );
  }

  if (!treeData || treeData.length === 0) {
    return (
      <div className="flex h-full w-full items-center justify-center p-8 text-muted-foreground">
        No organisation data available
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-2 h-full">
      <div className="flex items-center justify-between gap-2 px-1">
        <div className="flex items-center gap-2">
          {treeData.length > 1 && (
            <div className="flex flex-col gap-2">
              <Label
                htmlFor="organisation-select"
                className="text-xs font-medium"
              >
                Select Organisation
              </Label>
              <Select
                value={selectedRootId || undefined}
                onValueChange={setSelectedRootId}
              >
                <SelectTrigger id="organisation-select" className="w-[280px]">
                  <SelectValue placeholder="Select Organisation" />
                </SelectTrigger>
                <SelectContent>
                  {treeData.map((node) => (
                    <SelectItem key={node.id} value={node.id || ""}>
                      {node.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          )}
        </div>
        <div className="flex items-center gap-1 bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60 p-1 rounded-md border mt-6">
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={() => setZoom((z) => Math.min(z * 1.2, 3))}
            title="Zoom In"
          >
            <ZoomIn className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={() => setZoom((z) => Math.max(z / 1.2, 0.3))}
            title="Zoom Out"
          >
            <ZoomOut className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8"
            onClick={centerView}
            title="Reset View"
          >
            <RotateCcw className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <div
        id="treeWrapper"
        ref={containerRef}
        className="w-full flex-1 rounded-md border bg-background overflow-hidden relative"
        style={{ height: height ? `${height - 50}px` : "500px" }}
      >
        {formattedData.length > 0 && (
          <Tree
            data={formattedData}
            orientation="vertical"
            pathFunc="step"
            translate={translate}
            zoom={zoom}
            onUpdate={(state) => {
              // Sync internal d3 state with our react state to allow manual pan/zoom to coexist with buttons
              setTranslate(state.translate);
              setZoom(state.zoom);
            }}
            nodeSize={{ x: 220, y: 120 }}
            separation={{ siblings: 1.2, nonSiblings: 1.5 }}
            renderCustomNodeElement={renderCustomNode}
            enableLegacyTransitions={true}
            transitionDuration={300}
            draggable={true}
            collapsible={false}
          />
        )}
      </div>
    </div>
  );
};
