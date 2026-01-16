import { useState } from "react";
import RorIcon from "@/components/icons/rorIcon";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { Check, Copy } from "lucide-react";

interface RorDisplayProps {
  rorId: string;
}

export const RorDisplay = ({ rorId }: RorDisplayProps) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = (e: React.MouseEvent) => {
    e.stopPropagation(); // Prevent triggering row/node actions
    navigator.clipboard.writeText(rorId);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger asChild>
          <div
            className="flex cursor-pointer items-center justify-center rounded-sm p-1 hover:bg-muted"
            onClick={handleCopy}
          >
            <RorIcon width={16} height={16} />
            {copied && (
                <span className="sr-only">Copied</span>
            )}
          </div>
        </TooltipTrigger>
        <TooltipContent side="right">
          <div className="flex items-center gap-2">
            <span>{rorId}</span>
            {copied ? <Check size={12} className="text-green-500" /> : <Copy size={12} className="opacity-50" />}
          </div>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};
