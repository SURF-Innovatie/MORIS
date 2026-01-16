import { useState } from "react";
import { Download, ChevronDown, Loader2, Check, AlertCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useGetAdapters } from "@api/moris";
import { AXIOS_INSTANCE } from "@/api/custom-axios";

interface ExportButtonProps {
    projectId: string;
}

export function ExportButton({ projectId }: ExportButtonProps) {
    const [exportStatus, setExportStatus] = useState<"idle" | "loading" | "success" | "error">("idle");

    const { data: adapters, isLoading: adaptersLoading } = useGetAdapters();

    // Filter sinks that support "project" data type
    const projectSinks = adapters?.sinks?.filter(
        (sink) => sink.supported_types?.includes("project")
    ) ?? [];

    const handleExport = async (sinkName: string, outputType: string) => {
        setExportStatus("loading");
        try {
            if (outputType === "file") {
                // File-based export - download the file
                const response = await AXIOS_INSTANCE.post(
                    `/projects/${projectId}/export/${sinkName}`,
                    {},
                    { responseType: "blob" }
                );
                
                // Extract filename from Content-Disposition header
                const contentDisposition = response.headers["content-disposition"];
                let filename = `export_${projectId}.json`;
                if (contentDisposition) {
                    const match = contentDisposition.match(/filename="(.+)"/);
                    if (match) filename = match[1];
                }

                // Create download link
                const blob = new Blob([response.data], { type: response.headers["content-type"] });
                const url = window.URL.createObjectURL(blob);
                const link = document.createElement("a");
                link.href = url;
                link.download = filename;
                document.body.appendChild(link);
                link.click();
                document.body.removeChild(link);
                window.URL.revokeObjectURL(url);
            } else {
                // API-based export
                await AXIOS_INSTANCE.post(`/projects/${projectId}/export/${sinkName}`);
            }
            setExportStatus("success");
            setTimeout(() => setExportStatus("idle"), 3000);
        } catch (error) {
            setExportStatus("error");
            setTimeout(() => setExportStatus("idle"), 3000);
        }
    };

    if (adaptersLoading || projectSinks.length === 0) {
        return null;
    }

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button
                    variant="outline"
                    disabled={exportStatus === "loading"}
                    className="gap-2"
                >
                    {exportStatus === "loading" ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                    ) : exportStatus === "success" ? (
                        <Check className="h-4 w-4 text-green-500" />
                    ) : exportStatus === "error" ? (
                        <AlertCircle className="h-4 w-4 text-red-500" />
                    ) : (
                        <Download className="h-4 w-4" />
                    )}
                    Export
                    <ChevronDown className="h-3 w-3" />
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
                {projectSinks.map((sink) => (
                    <DropdownMenuItem
                        key={sink.name}
                        onClick={() => handleExport(sink.name!, sink.output?.type || "api")}
                        disabled={exportStatus === "loading"}
                    >
                        <div className="flex flex-col gap-0.5">
                            <span className="font-medium">{sink.output?.label || sink.display_name}</span>
                            {sink.output?.description && (
                                <span className="text-xs text-muted-foreground">
                                    {sink.output.description}
                                </span>
                            )}
                        </div>
                    </DropdownMenuItem>
                ))}
            </DropdownMenuContent>
        </DropdownMenu>
    );
}
