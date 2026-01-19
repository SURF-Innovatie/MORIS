import { useState, useRef } from "react";
import { Loader2, Upload, FileUp, Check } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useToast } from "@/hooks/use-toast";
import {
  usePostZenodoDepositions,
  usePutZenodoDepositionsId,
  usePostZenodoDepositionsIdPublish,
  useGetZenodoStatus,
} from "@api/moris";
import { UploadType, AccessRight } from "@api/model";
import { AXIOS_INSTANCE } from "@/api/custom-axios";
import ZenodoIcon from "@/components/icons/zenodoIcon";

interface ZenodoUploadDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: (
    doi: string,
    zenodoUrl: string,
    depositionId: number,
    title: string,
    uploadType: UploadType,
  ) => void;
}

export function ZenodoUploadDialog({
  open,
  onOpenChange,
  onSuccess,
}: ZenodoUploadDialogProps) {
  const { toast } = useToast();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [uploadType, setUploadType] = useState<UploadType>(
    UploadType.UploadTypeDataset,
  );
  const [creatorName, setCreatorName] = useState("");
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [step, setStep] = useState<"form" | "uploading" | "success">("form");
  const [currentAction, setCurrentAction] = useState("");
  const [resultData, setResultData] = useState<{
    doi: string;
    url: string;
    depositionId: number;
  } | null>(null);

  const { data: zenodoStatus } = useGetZenodoStatus();
  const { mutateAsync: createDeposition } = usePostZenodoDepositions();
  const { mutateAsync: updateDeposition } = usePutZenodoDepositionsId();
  const { mutateAsync: publishDeposition } =
    usePostZenodoDepositionsIdPublish();

  const resetForm = () => {
    setTitle("");
    setDescription("");
    setUploadType(UploadType.UploadTypeDataset);
    setCreatorName("");
    setSelectedFile(null);
    setStep("form");
    setCurrentAction("");
    setResultData(null);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setSelectedFile(e.target.files[0]);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      setSelectedFile(e.dataTransfer.files[0]);
    }
  };

  const handleSubmit = async () => {
    if (!selectedFile || !title || !description || !creatorName) {
      toast({
        title: "Missing Fields",
        description: "Please fill in all required fields and select a file.",
        variant: "destructive",
      });
      return;
    }

    setStep("uploading");

    try {
      // Step 1: Create deposition
      setCurrentAction("Creating deposition...");
      const deposition = await createDeposition();

      if (!deposition.id) {
        throw new Error("Failed to create deposition");
      }

      // Step 2: Upload file using custom axios with FormData
      setCurrentAction("Uploading file...");
      const formData = new FormData();
      formData.append("file", selectedFile);

      await AXIOS_INSTANCE.post(
        `/zenodo/depositions/${deposition.id}/files`,
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        },
      );

      // Step 3: Update metadata
      setCurrentAction("Adding metadata...");
      await updateDeposition({
        id: deposition.id,
        data: {
          title,
          description,
          upload_type: uploadType,
          access_right: AccessRight.AccessOpen,
          creators: [{ name: creatorName }],
        },
      });

      // Step 4: Publish
      setCurrentAction("Publishing...");
      const published = await publishDeposition({ id: deposition.id });

      // Success!
      const doi = published.doi || "";
      // Use record_id for URL - check for sandbox env var or default based on dev mode
      const isSandbox =
        import.meta.env.VITE_ZENODO_SANDBOX === "true" || import.meta.env.DEV;
      const baseUrl = isSandbox
        ? "https://sandbox.zenodo.org"
        : "https://zenodo.org";
      const url = `${baseUrl}/records/${published.record_id}`;

      setResultData({ doi, url, depositionId: deposition.id });
      setStep("success");
    } catch (error: any) {
      console.error("Zenodo upload error:", error);
      toast({
        title: "Upload Failed",
        description:
          error.response?.data?.message || "Failed to upload to Zenodo",
        variant: "destructive",
      });
      setStep("form");
    }
  };

  const handleClose = () => {
    if (step !== "uploading") {
      // If closing after success, call the success callback
      if (step === "success" && resultData) {
        onSuccess(
          resultData.doi,
          resultData.url,
          resultData.depositionId,
          title,
          uploadType,
        );
      }
      resetForm();
      onOpenChange(false);
    }
  };

  if (!zenodoStatus?.linked) {
    return (
      <Dialog open={open} onOpenChange={handleClose}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Zenodo Not Connected</DialogTitle>
            <DialogDescription>
              You need to connect your Zenodo account before you can upload
              files. Go to your profile settings to connect Zenodo.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button onClick={handleClose}>Close</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <ZenodoIcon width={24} height={24} />
            {step === "success" ? "Upload Complete!" : "Upload to Zenodo"}
          </DialogTitle>
          <DialogDescription>
            {step === "form" &&
              "Upload a file to Zenodo and get a DOI for your research product."}
            {step === "uploading" &&
              "Please wait while we upload and publish your file..."}
            {step === "success" && "Your file has been published on Zenodo."}
          </DialogDescription>
        </DialogHeader>

        {step === "form" && (
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="title">Title *</Label>
              <Input
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Enter product title"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description *</Label>
              <Textarea
                id="description"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Describe your research product..."
                rows={3}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="creator">Creator Name *</Label>
              <Input
                id="creator"
                value={creatorName}
                onChange={(e) => setCreatorName(e.target.value)}
                placeholder="Family name, Given names"
              />
            </div>

            <div className="space-y-2">
              <Label>Upload Type</Label>
              <Select
                value={uploadType}
                onValueChange={(v) => setUploadType(v as UploadType)}
              >
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value={UploadType.UploadTypeDataset}>
                    Dataset
                  </SelectItem>
                  <SelectItem value={UploadType.UploadTypePublication}>
                    Publication
                  </SelectItem>
                  <SelectItem value={UploadType.UploadTypeSoftware}>
                    Software
                  </SelectItem>
                  <SelectItem value={UploadType.UploadTypeImage}>
                    Image
                  </SelectItem>
                  <SelectItem value={UploadType.UploadTypeVideo}>
                    Video
                  </SelectItem>
                  <SelectItem value={UploadType.UploadTypeOther}>
                    Other
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label>File *</Label>
              <input
                ref={fileInputRef}
                type="file"
                onChange={handleFileChange}
                className="hidden"
              />
              <div
                onClick={() => fileInputRef.current?.click()}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
                className={`flex flex-col items-center justify-center gap-2 p-6 border-2 border-dashed rounded-lg cursor-pointer transition-colors ${
                  isDragging
                    ? "border-primary bg-primary/5"
                    : "hover:border-primary"
                }`}
              >
                {selectedFile ? (
                  <>
                    <FileUp className="h-5 w-5 text-primary" />
                    <span className="text-sm font-medium">
                      {selectedFile.name}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      ({(selectedFile.size / 1024 / 1024).toFixed(2)} MB)
                    </span>
                  </>
                ) : (
                  <>
                    <Upload className="h-8 w-8 text-muted-foreground" />
                    <span className="text-sm text-muted-foreground">
                      {isDragging
                        ? "Drop file here"
                        : "Drag & drop a file or click to select"}
                    </span>
                  </>
                )}
              </div>
            </div>
          </div>
        )}

        {step === "uploading" && (
          <div className="flex flex-col items-center justify-center py-8 space-y-4">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
            <p className="text-sm text-muted-foreground">{currentAction}</p>
          </div>
        )}

        {step === "success" && resultData && (
          <div className="space-y-4 py-4">
            <div className="flex items-center gap-2 text-green-600">
              <Check className="h-5 w-5" />
              <span className="font-medium">Published successfully!</span>
            </div>
            <div className="rounded-lg border p-4 space-y-2 bg-muted/50">
              <div>
                <span className="text-sm font-medium">DOI:</span>
                <span className="ml-2 font-mono text-sm">{resultData.doi}</span>
              </div>
              <div>
                <span className="text-sm font-medium">Zenodo URL:</span>
                <a
                  href={resultData.url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="ml-2 text-sm text-primary hover:underline"
                >
                  {resultData.url}
                </a>
              </div>
            </div>
          </div>
        )}

        <DialogFooter>
          {step === "form" && (
            <>
              <Button variant="outline" onClick={handleClose}>
                Cancel
              </Button>
              <Button
                onClick={handleSubmit}
                disabled={
                  !selectedFile || !title || !description || !creatorName
                }
              >
                <Upload className="mr-2 h-4 w-4" />
                Upload & Publish
              </Button>
            </>
          )}
          {step === "success" && <Button onClick={handleClose}>Done</Button>}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
