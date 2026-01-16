import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import {
  useGetProfileApiKeys,
  usePostProfileApiKeys,
  useDeleteProfileApiKeysKeyId,
} from "@api/moris";
import { APIKeyResponse as APIKey } from "@api/model";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";
import { useToast } from "@/hooks/use-toast";
import {
  Loader2,
  Plus,
  Key,
  Trash2,
  Copy,
  Check,
  AlertTriangle,
} from "lucide-react";

export function APIKeysSettings() {
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [newKeyName, setNewKeyName] = useState("");
  const [createdKey, setCreatedKey] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [keyToRevoke, setKeyToRevoke] = useState<APIKey | null>(null);

  const queryClient = useQueryClient();
  const { toast } = useToast();

  const { data: keysData, isLoading } = useGetProfileApiKeys();
  const apiKeys = keysData?.apiKeys || [];

  const { mutate: createKey, isPending: isCreating } = usePostProfileApiKeys({
    mutation: {
      onSuccess: (data) => {
        queryClient.invalidateQueries({ queryKey: ["/profile/api-keys"] });
        setCreatedKey(data.plainKey || null);
        setNewKeyName("");
      },
      onError: () => {
        toast({
          title: "Error",
          description: "Failed to create API key",
          variant: "destructive",
        });
      },
    },
  });

  const { mutate: revokeKey } = useDeleteProfileApiKeysKeyId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({ queryKey: ["/profile/api-keys"] });
        setKeyToRevoke(null);
        toast({
          title: "Success",
          description: "API key revoked",
        });
      },
      onError: () => {
        toast({
          title: "Error",
          description: "Failed to revoke API key",
          variant: "destructive",
        });
      },
    },
  });

  const handleCreate = () => {
    if (newKeyName.trim()) {
      createKey({ data: { name: newKeyName.trim() } });
    }
  };

  const handleCopy = async () => {
    if (createdKey) {
      await navigator.clipboard.writeText(createdKey);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleCloseCreatedDialog = () => {
    setCreatedKey(null);
    setShowCreateDialog(false);
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-24">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="text-sm text-muted-foreground flex items-center gap-2">
          <AlertTriangle className="h-4 w-4 text-yellow-500" />
          API keys act with your full permissions
        </div>
        <Dialog open={showCreateDialog} onOpenChange={setShowCreateDialog}>
          <DialogTrigger asChild>
            <Button size="sm">
              <Plus className="h-4 w-4 mr-2" />
              Create New Key
            </Button>
          </DialogTrigger>
          <DialogContent>
            {createdKey ? (
              <>
                <DialogHeader>
                  <DialogTitle>API Key Created</DialogTitle>
                  <DialogDescription>
                    This is the only time you will see this key. Copy it now!
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                  <div className="flex items-center gap-2">
                    <Input
                      value={createdKey}
                      readOnly
                      className="font-mono text-sm"
                    />
                    <Button variant="outline" size="icon" onClick={handleCopy}>
                      {copied ? (
                        <Check className="h-4 w-4 text-green-500" />
                      ) : (
                        <Copy className="h-4 w-4" />
                      )}
                    </Button>
                  </div>
                  <div className="bg-yellow-50 border border-yellow-200 rounded p-3 text-sm text-yellow-800">
                    <strong>Important:</strong> Store this key securely. You
                    won't be able to see it again.
                  </div>
                </div>
                <DialogFooter>
                  <Button onClick={handleCloseCreatedDialog}>
                    I've copied my key
                  </Button>
                </DialogFooter>
              </>
            ) : (
              <>
                <DialogHeader>
                  <DialogTitle>Create API Key</DialogTitle>
                  <DialogDescription>
                    Create a new API key for Power BI or other external tools
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4 py-4">
                  <div className="space-y-2">
                    <Label htmlFor="keyName">Key Name</Label>
                    <Input
                      id="keyName"
                      placeholder="e.g., Power BI Dashboard"
                      value={newKeyName}
                      onChange={(e) => setNewKeyName(e.target.value)}
                    />
                  </div>
                </div>
                <DialogFooter>
                  <Button
                    variant="outline"
                    onClick={() => setShowCreateDialog(false)}
                  >
                    Cancel
                  </Button>
                  <Button
                    onClick={handleCreate}
                    disabled={!newKeyName.trim() || isCreating}
                  >
                    {isCreating && (
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    )}
                    Create Key
                  </Button>
                </DialogFooter>
              </>
            )}
          </DialogContent>
        </Dialog>
      </div>

      {/* Key List */}
      {apiKeys && apiKeys.length > 0 ? (
        <div className="space-y-2">
          {apiKeys.map((key) => (
            <APIKeyRow
              key={key.id}
              apiKey={key}
              onRevokeClick={() => setKeyToRevoke(key)}
            />
          ))}
        </div>
      ) : (
        <div className="text-center py-8 text-muted-foreground">
          <Key className="h-8 w-8 mx-auto mb-2 opacity-50" />
          <p>No API keys yet</p>
        </div>
      )}

      {/* Revoke Confirmation Modal */}
      <ConfirmationModal
        isOpen={!!keyToRevoke}
        onClose={() => setKeyToRevoke(null)}
        onConfirm={() =>
          keyToRevoke?.id && revokeKey({ keyId: keyToRevoke.id })
        }
        title="Revoke API Key"
        description={`Are you sure you want to revoke "${keyToRevoke?.name}"? This action cannot be undone and any applications using this key will stop working.`}
        confirmLabel="Revoke Key"
        variant="destructive"
      />
    </div>
  );
}

interface APIKeyRowProps {
  apiKey: APIKey;
  onRevokeClick: () => void;
}

function APIKeyRow({ apiKey, onRevokeClick }: APIKeyRowProps) {
  return (
    <div className="flex items-center justify-between p-3 border rounded-lg">
      <div className="flex items-center gap-3">
        <Key className="h-4 w-4 text-muted-foreground" />
        <div>
          <p className="font-medium">{apiKey.name}</p>
          <p className="text-sm text-muted-foreground font-mono">
            {apiKey.keyPrefix}
          </p>
        </div>
      </div>
      <div className="flex items-center gap-4">
        <div className="text-right text-sm text-muted-foreground">
          <p>Created {formatDate(apiKey.createdAt || "")}</p>
          {apiKey.lastUsedAt && (
            <p>Last used {formatDate(apiKey.lastUsedAt)}</p>
          )}
        </div>
        {!apiKey.isActive ? (
          <span className="text-sm text-red-500">Revoked</span>
        ) : (
          <Button variant="ghost" size="icon" onClick={onRevokeClick}>
            <Trash2 className="h-4 w-4 text-red-500" />
          </Button>
        )}
      </div>
    </div>
  );
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString("nl-NL", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}
