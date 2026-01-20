import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import {
  Shield,
  Lock,
  Smartphone,
  Loader2,
  AlertTriangle,
  Copy,
  Check,
} from "lucide-react";
import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  useGetProfileApiKeys,
  usePostProfileApiKeys,
  useDeleteProfileApiKeysKeyId,
} from "@/api/generated-orval/moris";
import { useQueryClient } from "@tanstack/react-query";

export function SecuritySettings() {
  const queryClient = useQueryClient();
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [newKeyName, setNewKeyName] = useState("");
  const [newlyCreatedKey, setNewlyCreatedKey] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    if (newlyCreatedKey) {
      await navigator.clipboard.writeText(newlyCreatedKey);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  // Queries
  const { data: keysData, isLoading: isLoadingKeys } = useGetProfileApiKeys();
  const keys = keysData?.apiKeys || [];

  // Mutations
  const createKeyMutation = usePostProfileApiKeys({
    mutation: {
      onSuccess: (data) => {
        setNewlyCreatedKey(data.plainKey || null);
        setNewKeyName("");
        setIsCreateOpen(false);
        queryClient.invalidateQueries({
          queryKey: ["/profile/api-keys"], // Can also use getGetProfileApiKeysQueryKey() from moris.ts if exported
        });
      },
      onError: (error) => {
        console.error("Failed to create API key", error);
      },
    },
  });

  const revokeKeyMutation = useDeleteProfileApiKeysKeyId({
    mutation: {
      onSuccess: () => {
        queryClient.invalidateQueries({
          queryKey: ["/profile/api-keys"],
        });
      },
      onError: (error) => {
        console.error("Failed to revoke API key", error);
      },
    },
  });

  const createKey = () => {
    createKeyMutation.mutate({ data: { name: newKeyName } });
  };

  const revokeKey = (id: string) => {
    if (!confirm("Are you sure you want to revoke this API key?")) return;
    revokeKeyMutation.mutate({ keyId: id });
  };

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <Lock className="h-5 w-5 text-muted-foreground" />
            <div className="space-y-1">
              <CardTitle>Password</CardTitle>
              <CardDescription>
                Change your password to keep your account secure
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-2">
            <Label htmlFor="current-password">Current Password</Label>
            <Input id="current-password" type="password" />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="new-password">New Password</Label>
            <Input id="new-password" type="password" />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="confirm-password">Confirm Password</Label>
            <Input id="confirm-password" type="password" />
          </div>
          <div className="flex justify-end">
            <Button disabled>Update Password</Button>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <Smartphone className="h-5 w-5 text-muted-foreground" />
            <div className="space-y-1">
              <CardTitle>Two-Factor Authentication</CardTitle>
              <CardDescription>
                Add an extra layer of security to your account
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between space-x-2">
            <div className="flex flex-col space-y-1">
              <span className="font-medium">Enable 2FA</span>
              <span className="text-sm text-muted-foreground">
                Secure your account with TOTP (Authenticator App)
              </span>
            </div>
            <Switch disabled />
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <div className="flex items-center gap-2">
            <Shield className="h-5 w-5 text-muted-foreground" />
            <div className="space-y-1">
              <CardTitle>API Keys</CardTitle>
              <CardDescription>
                Manage personal access tokens for external tools and
                integrations
              </CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          {isLoadingKeys ? (
            <div className="flex justify-center py-4">
              <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
            </div>
          ) : keys.length === 0 ? (
            <div className="text-sm text-muted-foreground py-4 text-center">
              No API keys found. Create one to get started.
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Prefix</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead>Last Used</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {keys.map((key) => (
                  <TableRow key={key.id}>
                    <TableCell className="font-medium">{key.name}</TableCell>
                    <TableCell className="font-mono text-xs">
                      {key.keyPrefix}
                    </TableCell>
                    <TableCell>
                      {key.createdAt
                        ? new Date(key.createdAt).toLocaleDateString()
                        : "-"}
                    </TableCell>
                    <TableCell>
                      {key.lastUsedAt
                        ? new Date(key.lastUsedAt).toLocaleDateString()
                        : "Never"}
                    </TableCell>
                    <TableCell>
                      <span
                        className={`inline-flex items-center rounded-full px-2 py-1 text-xs font-medium ${
                          key.isActive
                            ? "bg-green-50 text-green-700 ring-1 ring-inset ring-green-600/20"
                            : "bg-red-50 text-red-700 ring-1 ring-inset ring-red-600/10"
                        }`}
                      >
                        {key.isActive ? "Active" : "Revoked"}
                      </span>
                    </TableCell>
                    <TableCell className="text-right">
                      {key.isActive && (
                        <Button
                          variant="ghost"
                          size="sm"
                          className="text-red-600 hover:text-red-700 hover:bg-red-50"
                          onClick={() => key.id && revokeKey(key.id)}
                          disabled={revokeKeyMutation.isPending}
                        >
                          Revoke
                        </Button>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}

          <div className="mt-4 flex justify-end">
            <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
              <DialogTrigger asChild>
                <Button>Create New API Key</Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Create New API Key</DialogTitle>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                  <div className="grid gap-2">
                    <Label htmlFor="key-name">Name</Label>
                    <Input
                      id="key-name"
                      placeholder="e.g. Power BI Integration"
                      value={newKeyName}
                      onChange={(e) => setNewKeyName(e.target.value)}
                    />
                  </div>
                </div>
                <div className="flex justify-end">
                  <Button
                    onClick={createKey}
                    disabled={!newKeyName || createKeyMutation.isPending}
                  >
                    {createKeyMutation.isPending ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Creating...
                      </>
                    ) : (
                      "Create Key"
                    )}
                  </Button>
                </div>
              </DialogContent>
            </Dialog>
          </div>
        </CardContent>
      </Card>

      <Dialog
        open={!!newlyCreatedKey}
        onOpenChange={() => {
          setNewlyCreatedKey(null);
          setCopied(false);
        }}
      >
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2">
              <Shield className="h-5 w-5 text-green-600" />
              API Key Created
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-6">
            <div className="rounded-lg bg-amber-50 border border-amber-200 p-4">
              <div className="flex gap-3">
                <AlertTriangle className="h-5 w-5 text-amber-600 shrink-0" />
                <div className="space-y-1">
                  <h3 className="text-sm font-semibold text-amber-900">
                    Save your secret key
                  </h3>
                  <p className="text-sm text-amber-800 leading-relaxed">
                    Please copy your API key now. For security reasons, you
                    won&apos;t be able to see it again.
                  </p>
                </div>
              </div>
            </div>

            <div className="space-y-2">
              <Label className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">
                Your API Key
              </Label>
              <div className="relative group">
                <code className="block w-full p-4 pr-12 rounded-lg bg-muted font-mono text-sm break-all border border-input whitespace-pre-wrap">
                  {newlyCreatedKey}
                </code>
                <Button
                  variant="ghost"
                  size="icon"
                  className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 hover:bg-background/80"
                  onClick={handleCopy}
                >
                  {copied ? (
                    <Check className="h-4 w-4 text-green-600" />
                  ) : (
                    <Copy className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>

            <div className="flex justify-end pt-2">
              <Button
                className="w-full sm:w-auto"
                onClick={() => {
                  setNewlyCreatedKey(null);
                  setCopied(false);
                }}
              >
                Done
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
