import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { standardSchemaResolver } from "@hookform/resolvers/standard-schema";
import {
  Loader2,
  Plus,
  Search,
  Building2,
  Trash2,
  ExternalLink,
  MapPin,
  CheckCircle,
  XCircle,
} from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useToast } from "@/hooks/use-toast";
import {
  usePostAffiliatedOrganisations,
  useGetOrganisationNodesRorSearch,
  useGetAffiliatedOrganisationsVatLookup,
} from "@api/moris";
import { RORItem } from "@api/model";
import {
  createAffiliatedOrganisationAddedEvent,
  createAffiliatedOrganisationRemovedEvent,
  ProjectEventType,
} from "@/api/events";
import { AffiliatedOrganisationResponse } from "@api/model";
import { ConfirmationModal } from "@/components/ui/confirmation-modal";
import { Allowed } from "@/components/auth/Allowed";
import { useAccess } from "@/context/AccessContext";

// Schema for manual organisation entry
const manualOrgSchema = z.object({
  name: z.string().min(1, "Name is required"),
  kvkNumber: z.string().optional(),
  rorId: z.string().optional(),
  vatNumber: z.string().optional(),
  city: z.string().optional(),
  country: z.string().optional(),
});

// Schema for KVK search
const kvkSearchSchema = z.object({
  query: z.string().min(2, "Enter at least 2 characters"),
});

// Schema for ROR search
const rorSearchSchema = z.object({
  query: z.string().min(2, "Enter at least 2 characters"),
});

// Schema for VAT search
const vatSearchSchema = z.object({
  vatNumber: z
    .string()
    .min(5, "Enter a valid VAT number (e.g., NL822655287B01)"),
});

interface AffiliatedOrganisationsTabProps {
  projectId: string;
  affiliatedOrganisations: AffiliatedOrganisationResponse[];
  onRefresh: () => void;
}

export function AffiliatedOrganisationsTab({
  projectId,
  affiliatedOrganisations,
  onRefresh,
}: AffiliatedOrganisationsTabProps) {
  const { toast } = useToast();
  const { hasAccess } = useAccess();

  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [orgToDelete, setOrgToDelete] = useState<string | null>(null);
  const [activeSearchTab, setActiveSearchTab] = useState("manual");

  // KVK Search state - disabled for now as KVK API needs configuration
  const [kvkSearchQuery, setKvkSearchQuery] = useState<string | null>(null);
  const kvkResults: any = null;
  const isSearchingKvk = false;

  // ROR Search state - uses existing working endpoint
  const [rorSearchQuery, setRorSearchQuery] = useState<string | null>(null);
  const { data: rorResults, isLoading: isSearchingRor } =
    useGetOrganisationNodesRorSearch(
      { q: rorSearchQuery! },
      { query: { enabled: !!rorSearchQuery && rorSearchQuery.length > 2 } },
    );

  // VAT Search state - uses generated hook
  const [vatSearchQuery, setVatSearchQuery] = useState<string | null>(null);
  const {
    data: vatResult,
    isLoading: isSearchingVat,
    error: vatQueryError,
  } = useGetAffiliatedOrganisationsVatLookup(
    { vat_number: vatSearchQuery! },
    { query: { enabled: !!vatSearchQuery && vatSearchQuery.length > 4 } },
  );
  const vatError = vatQueryError
    ? "Failed to lookup VAT number"
    : vatResult && !vatResult.valid
      ? "VAT number is not valid"
      : null;

  const { mutateAsync: createOrganisation, isPending: isCreating } =
    usePostAffiliatedOrganisations();

  // Forms
  const manualForm = useForm<z.infer<typeof manualOrgSchema>>({
    resolver: standardSchemaResolver(manualOrgSchema),
    defaultValues: {
      name: "",
      kvkNumber: "",
      rorId: "",
      vatNumber: "",
      city: "",
      country: "",
    },
  });

  const kvkForm = useForm<z.infer<typeof kvkSearchSchema>>({
    resolver: standardSchemaResolver(kvkSearchSchema),
    defaultValues: { query: "" },
  });

  const rorForm = useForm<z.infer<typeof rorSearchSchema>>({
    resolver: standardSchemaResolver(rorSearchSchema),
    defaultValues: { query: "" },
  });

  const vatForm = useForm<z.infer<typeof vatSearchSchema>>({
    resolver: standardSchemaResolver(vatSearchSchema),
    defaultValues: { vatNumber: "" },
  });

  async function onManualSubmit(values: z.infer<typeof manualOrgSchema>) {
    try {
      const newOrg = await createOrganisation({
        data: {
          name: values.name,
          kvk_number: values.kvkNumber,
          ror_id: values.rorId,
          vat_number: values.vatNumber,
          city: values.city,
          country: values.country,
        },
      });

      if (newOrg?.id) {
        await createAffiliatedOrganisationAddedEvent(projectId, {
          affiliated_organisation_id: newOrg.id,
        });

        toast({
          title: "Organisation added",
          description: "The organisation has been affiliated with the project.",
        });
        setIsDialogOpen(false);
        manualForm.reset();
        onRefresh();
      }
    } catch (error) {
      console.error(error);
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to add organisation. Please try again.",
      });
    }
  }

  function onKvkSearch(values: z.infer<typeof kvkSearchSchema>) {
    setKvkSearchQuery(values.query);
  }

  function onRorSearch(values: z.infer<typeof rorSearchSchema>) {
    setRorSearchQuery(values.query);
  }

  async function addFromKvk(result: any) {
    try {
      const newOrg = await createOrganisation({
        data: {
          name: result.handelsnaam || result.naam || "Unknown",
          kvk_number: result.kvkNummer,
          city: result.adres?.binnenlandsAdres?.plaats || "",
          country: "Netherlands",
        },
      });

      if (newOrg?.id) {
        await createAffiliatedOrganisationAddedEvent(projectId, {
          affiliated_organisation_id: newOrg.id,
        });

        toast({
          title: "Organisation added",
          description:
            "The KVK organisation has been affiliated with the project.",
        });
        setIsDialogOpen(false);
        setKvkSearchQuery(null);
        kvkForm.reset();
        onRefresh();
      }
    } catch (error) {
      console.error(error);
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to add organisation. Please try again.",
      });
    }
  }

  async function addFromRor(item: RORItem) {
    try {
      const newOrg = await createOrganisation({
        data: {
          name: item.name || "Unknown",
          ror_id: item.id,
          city: item.addresses?.[0]?.city || "",
          country: item.country?.country_name || "",
        },
      });

      if (newOrg?.id) {
        await createAffiliatedOrganisationAddedEvent(projectId, {
          affiliated_organisation_id: newOrg.id,
        });

        toast({
          title: "Organisation added",
          description:
            "The ROR organisation has been affiliated with the project.",
        });
        setIsDialogOpen(false);
        setRorSearchQuery(null);
        rorForm.reset();
        onRefresh();
      }
    } catch (error) {
      console.error(error);
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to add organisation. Please try again.",
      });
    }
  }

  function onVatLookup(values: z.infer<typeof vatSearchSchema>) {
    setVatSearchQuery(values.vatNumber);
  }

  async function addFromVat() {
    if (!vatResult || !vatResult.valid) return;

    try {
      const newOrg = await createOrganisation({
        data: {
          name: vatResult.name || "Unknown",
          vat_number: vatResult.vat_number,
          city: vatResult.city || "",
          country: vatResult.country_code || "",
        },
      });

      if (newOrg?.id) {
        await createAffiliatedOrganisationAddedEvent(projectId, {
          affiliated_organisation_id: newOrg.id,
        });

        toast({
          title: "Organisation added",
          description: "The organisation has been affiliated with the project.",
        });
        setIsDialogOpen(false);
        setVatSearchQuery(null);
        vatForm.reset();
        onRefresh();
      }
    } catch (error) {
      console.error(error);
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to add organisation. Please try again.",
      });
    }
  }

  async function onRemoveOrganisation() {
    if (!orgToDelete) return;

    try {
      await createAffiliatedOrganisationRemovedEvent(projectId, {
        affiliated_organisation_id: orgToDelete,
      });

      toast({
        title: "Organisation removed",
        description: "The organisation has been removed from the project.",
      });
      setOrgToDelete(null);
      onRefresh();
    } catch (error) {
      console.error(error);
      toast({
        variant: "destructive",
        title: "Error",
        description: "Failed to remove organisation.",
      });
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-medium">Affiliated Organisations</h3>
          <p className="text-sm text-muted-foreground">
            Manage organisations affiliated with this project.
          </p>
        </div>
        <Allowed event={ProjectEventType.AffiliatedOrganisationAdded}>
          <Button onClick={() => setIsDialogOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Add Organisation
          </Button>
        </Allowed>
      </div>

      {/* Organisation Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {affiliatedOrganisations.map((org) => (
          <Card key={org.id} className="relative">
            <CardHeader className="pb-3">
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-2">
                  <Building2 className="h-5 w-5 text-muted-foreground" />
                  <CardTitle className="text-base">{org.name}</CardTitle>
                </div>
                {hasAccess(ProjectEventType.AffiliatedOrganisationRemoved) && (
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-muted-foreground hover:text-destructive"
                    onClick={() => setOrgToDelete(org.id!)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                )}
              </div>
              {(org.city || org.country) && (
                <CardDescription className="flex items-center gap-1">
                  <MapPin className="h-3 w-3" />
                  {[org.city, org.country].filter(Boolean).join(", ")}
                </CardDescription>
              )}
            </CardHeader>
            <CardContent className="pt-0">
              <div className="flex flex-wrap gap-2">
                {org.kvk_number && (
                  <Badge variant="secondary">KVK: {org.kvk_number}</Badge>
                )}
                {org.ror_id && (
                  <a
                    href={
                      org.ror_id.startsWith("http")
                        ? org.ror_id
                        : `https://ror.org/${org.ror_id}`
                    }
                    target="_blank"
                    rel="noopener noreferrer"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <Badge
                      variant="outline"
                      className="gap-1 cursor-pointer hover:bg-accent"
                    >
                      ROR
                      <ExternalLink className="h-3 w-3" />
                    </Badge>
                  </a>
                )}
                {org.vat_number && (
                  <Badge variant="secondary">VAT: {org.vat_number}</Badge>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
        {affiliatedOrganisations.length === 0 && (
          <div className="col-span-full flex flex-col items-center justify-center p-8 text-center border rounded-lg border-dashed text-muted-foreground">
            <Building2 className="h-10 w-10 mb-2 opacity-50" />
            <p>No affiliated organisations yet.</p>
          </div>
        )}
      </div>

      {/* Add Organisation Dialog */}
      <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle>Add Affiliated Organisation</DialogTitle>
            <DialogDescription>
              Add an organisation manually or search via VAT, ROR, or KVK.
            </DialogDescription>
          </DialogHeader>

          <Tabs value={activeSearchTab} onValueChange={setActiveSearchTab}>
            <TabsList className="grid w-full grid-cols-4">
              <TabsTrigger value="manual">Manual</TabsTrigger>
              <TabsTrigger value="vat">VAT</TabsTrigger>
              <TabsTrigger value="ror">ROR</TabsTrigger>
              <TabsTrigger value="kvk">KVK</TabsTrigger>
            </TabsList>

            {/* Manual Entry Tab */}
            <TabsContent value="manual" className="space-y-4 mt-4">
              <Form {...manualForm}>
                <form
                  onSubmit={manualForm.handleSubmit(onManualSubmit)}
                  className="space-y-4"
                >
                  <FormField
                    control={manualForm.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Organisation Name *</FormLabel>
                        <FormControl>
                          <Input
                            placeholder="Enter organisation name"
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <div className="grid grid-cols-2 gap-4">
                    <FormField
                      control={manualForm.control}
                      name="kvkNumber"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>KVK Number</FormLabel>
                          <FormControl>
                            <Input placeholder="12345678" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={manualForm.control}
                      name="rorId"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>ROR ID</FormLabel>
                          <FormControl>
                            <Input
                              placeholder="https://ror.org/..."
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <FormField
                    control={manualForm.control}
                    name="vatNumber"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>VAT Number</FormLabel>
                        <FormControl>
                          <Input placeholder="NL123456789B01" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <div className="grid grid-cols-2 gap-4">
                    <FormField
                      control={manualForm.control}
                      name="city"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>City</FormLabel>
                          <FormControl>
                            <Input placeholder="Amsterdam" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={manualForm.control}
                      name="country"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Country</FormLabel>
                          <FormControl>
                            <Input placeholder="Netherlands" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </div>
                  <DialogFooter>
                    <Button
                      type="button"
                      variant="secondary"
                      onClick={() => setIsDialogOpen(false)}
                    >
                      Cancel
                    </Button>
                    <Button type="submit" disabled={isCreating}>
                      {isCreating && (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      )}
                      Add Organisation
                    </Button>
                  </DialogFooter>
                </form>
              </Form>
            </TabsContent>

            {/* VAT Search Tab */}
            <TabsContent value="vat" className="space-y-4 mt-4">
              <Form {...vatForm}>
                <form
                  onSubmit={vatForm.handleSubmit(onVatLookup)}
                  className="flex gap-2"
                >
                  <FormField
                    control={vatForm.control}
                    name="vatNumber"
                    render={({ field }) => (
                      <FormItem className="flex-1">
                        <FormControl>
                          <Input
                            placeholder="Enter VAT number (e.g., NL822655287B01)"
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <Button type="submit" disabled={isSearchingVat}>
                    {isSearchingVat ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Search className="h-4 w-4" />
                    )}
                  </Button>
                </form>
              </Form>

              {vatError && (
                <div className="flex items-center gap-2 p-3 border rounded-lg bg-destructive/10 text-destructive">
                  <XCircle className="h-4 w-4" />
                  <p className="text-sm">{vatError}</p>
                </div>
              )}

              {vatResult && vatResult.valid && (
                <div className="space-y-3">
                  <div className="flex items-center gap-2 p-3 border rounded-lg bg-green-50 dark:bg-green-950 text-green-700 dark:text-green-300">
                    <CheckCircle className="h-4 w-4" />
                    <p className="text-sm font-medium">VAT number is valid</p>
                  </div>
                  <div className="p-4 border rounded-lg space-y-2">
                    <h4 className="font-medium">{vatResult.name}</h4>
                    <p className="text-sm text-muted-foreground whitespace-pre-line">
                      {vatResult.address?.trim()}
                    </p>
                    <p className="text-sm text-muted-foreground">
                      VAT: {vatResult.vat_number}
                    </p>
                  </div>
                  <Button
                    onClick={addFromVat}
                    disabled={isCreating}
                    className="w-full"
                  >
                    {isCreating && (
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    )}
                    <Plus className="mr-2 h-4 w-4" />
                    Add This Organisation
                  </Button>
                </div>
              )}

              {!vatResult && !vatError && !isSearchingVat && (
                <p className="text-sm text-muted-foreground text-center py-4">
                  Enter a VAT number with country code (e.g., NL822655287B01)
                  and click search to validate.
                </p>
              )}
            </TabsContent>

            {/* KVK Search Tab */}
            <TabsContent value="kvk" className="space-y-4 mt-4">
              <Form {...kvkForm}>
                <form
                  onSubmit={kvkForm.handleSubmit(onKvkSearch)}
                  className="flex gap-2"
                >
                  <FormField
                    control={kvkForm.control}
                    name="query"
                    render={({ field }) => (
                      <FormItem className="flex-1">
                        <FormControl>
                          <Input
                            placeholder="Search by name or KVK number..."
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <Button type="submit" disabled={isSearchingKvk}>
                    {isSearchingKvk ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Search className="h-4 w-4" />
                    )}
                  </Button>
                </form>
              </Form>

              {kvkResults?.resultaten && kvkResults.resultaten.length > 0 && (
                <div className="max-h-[300px] overflow-y-auto space-y-2">
                  {kvkResults.resultaten.map((result: any, idx: number) => (
                    <div
                      key={`${result.kvkNummer}-${idx}`}
                      className="flex items-center justify-between p-3 border rounded-lg hover:bg-accent cursor-pointer"
                      onClick={() => addFromKvk(result)}
                    >
                      <div>
                        <p className="font-medium">{result.handelsnaam}</p>
                        <p className="text-sm text-muted-foreground">
                          KVK: {result.kvkNummer}
                          {result.adres?.binnenlandsAdres?.plaats &&
                            ` • ${result.adres.binnenlandsAdres.plaats}`}
                        </p>
                      </div>
                      <Plus className="h-4 w-4 text-muted-foreground" />
                    </div>
                  ))}
                </div>
              )}

              {kvkSearchQuery &&
                !isSearchingKvk &&
                !kvkResults?.resultaten?.length && (
                  <p className="text-sm text-muted-foreground text-center py-4">
                    No results found for "{kvkSearchQuery}"
                  </p>
                )}
            </TabsContent>

            {/* ROR Search Tab */}
            <TabsContent value="ror" className="space-y-4 mt-4">
              <Form {...rorForm}>
                <form
                  onSubmit={rorForm.handleSubmit(onRorSearch)}
                  className="flex gap-2"
                >
                  <FormField
                    control={rorForm.control}
                    name="query"
                    render={({ field }) => (
                      <FormItem className="flex-1">
                        <FormControl>
                          <Input
                            placeholder="Search by organisation name..."
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <Button type="submit" disabled={isSearchingRor}>
                    {isSearchingRor ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <Search className="h-4 w-4" />
                    )}
                  </Button>
                </form>
              </Form>

              {rorResults && rorResults.length > 0 && (
                <div className="max-h-[300px] overflow-y-auto space-y-2">
                  {rorResults.map((item: RORItem) => (
                    <div
                      key={item.id}
                      className="flex items-center justify-between p-3 border rounded-lg hover:bg-accent cursor-pointer"
                      onClick={() => addFromRor(item)}
                    >
                      <div>
                        <p className="font-medium">{item.name}</p>
                        <p className="text-sm text-muted-foreground">
                          {item.country?.country_name || ""}
                          {item.addresses?.[0]?.city &&
                            ` • ${item.addresses[0].city}`}
                        </p>
                      </div>
                      <Plus className="h-4 w-4 text-muted-foreground" />
                    </div>
                  ))}
                </div>
              )}

              {rorSearchQuery &&
                !isSearchingRor &&
                (!rorResults || rorResults.length === 0) && (
                  <p className="text-sm text-muted-foreground text-center py-4">
                    No results found for "{rorSearchQuery}"
                  </p>
                )}
            </TabsContent>
          </Tabs>
        </DialogContent>
      </Dialog>

      {/* Confirmation Modal for Remove */}
      <ConfirmationModal
        isOpen={!!orgToDelete}
        onClose={() => setOrgToDelete(null)}
        onConfirm={onRemoveOrganisation}
        title="Remove Organisation"
        description="Are you sure you want to remove this organisation from the project? This action cannot be undone."
        confirmLabel="Remove"
        variant="destructive"
      />
    </div>
  );
}
