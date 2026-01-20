import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Plus, Trash2, ArrowLeft, Save } from "lucide-react";
import { AddNWOSubsidyDialog } from "./AddNWOSubsidyDialog";
import NwoIcon from "@/components/icons/nwoIcon";

import {
  useGetProjectsProjectIdBudget,
  usePostBudgetsBudgetIdLineItems,
  useDeleteBudgetsBudgetIdLineItemsLineItemId,
} from "@api/moris";
import { BudgetCategory, FundingSource, Project } from "@api/model";
import { categoryLabels, fundingSourceLabels } from "@/lib/constants";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useToast } from "@/hooks/use-toast";
import { Badge } from "@/components/ui/badge";

interface BudgetEditorProps {
  projectId: string;
  onDone: () => void;
}

const lineItemSchema = z.object({
  description: z.string().min(1, "Description is required"),
  category: z.enum([
    "personnel",
    "material",
    "investment",
    "travel",
    "management",
    "grant",
    "other",
  ] as const),
  budgetedAmount: z.coerce.number().min(0, "Amount must be positive"),
  year: z.coerce.number().int().min(2020, "Invalid year"),
  fundingSource: z.enum([
    "subsidy",
    "cofinancing_cash",
    "cofinancing_inkind",
  ] as const),
  nwoGrantId: z.string().optional().nullable(),
});

export function BudgetEditor({ projectId, onDone }: BudgetEditorProps) {
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [isAdding, setIsAdding] = useState(false);
  const [isAddingNWO, setIsAddingNWO] = useState(false);

  const { data: budget } = useGetProjectsProjectIdBudget(projectId);

  const form = useForm<z.infer<typeof lineItemSchema>>({
    resolver: zodResolver(lineItemSchema),
    defaultValues: {
      description: "",
      category: "personnel",
      budgetedAmount: 0 as any, // Cast for input handling
      year: new Date().getFullYear() as any,
      fundingSource: "subsidy",
      nwoGrantId: null,
    },
  });

  const { mutate: addLineItem, isPending: isAddingLineItem } =
    usePostBudgetsBudgetIdLineItems({
      mutation: {
        onSuccess: () => {
          // Invalidate all possible query key variations to ensure refresh
          queryClient.invalidateQueries({ queryKey: ["budget", projectId] });
          queryClient.invalidateQueries({
            queryKey: [`/projects/${projectId}/budget`],
          });

          toast({ title: "Line Item Added" });
          setIsAdding(false);
          setIsAddingNWO(false);
          form.reset();
        },
        onError: () => {
          toast({
            variant: "destructive",
            title: "Error",
            description: "Failed to add line item",
          });
          setIsAddingNWO(false);
        },
      },
    });

  const { mutate: removeLineItem } =
    useDeleteBudgetsBudgetIdLineItemsLineItemId({
      mutation: {
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["budget", projectId] });
          queryClient.invalidateQueries({
            queryKey: [`/projects/${projectId}/budget`],
          });
          toast({ title: "Line Item Removed" });
        },
      },
    });

  const onSubmit = (values: z.infer<typeof lineItemSchema>) => {
    if (!budget?.id) return;
    addLineItem({
      budgetId: budget.id,
      data: {
        ...values,
        nwoGrantId: values.nwoGrantId || undefined,
      },
    });
  };

  const handleNWOSelect = (project: Project) => {
    if (!budget?.id) return;
    setIsAddingNWO(true);

    let year = new Date().getFullYear();
    if (project.start_date) {
      // start_date is likely a string based on generated types for dates
      const date = new Date(project.start_date as unknown as string);
      if (!isNaN(date.getFullYear())) {
        year = date.getFullYear();
      }
    }

    console.log(project);
    addLineItem({
      budgetId: budget.id,
      data: {
        description: project.title || "NWO Grant",
        category: "grant",
        budgetedAmount: project.award_amount || 0,
        year: year,
        fundingSource: "subsidy",
        nwoGrantId: project.project_id,
      },
    });
  };

  if (!budget) return <div>Loading...</div>;

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold tracking-tight">Edit Budget</h2>
        <Button variant="outline" onClick={onDone}>
          <ArrowLeft className="mr-2 h-4 w-4" /> Done
        </Button>
      </div>

      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <CardTitle>Line Items</CardTitle>
          <div className="flex gap-2">
            {!isAdding && (
              <>
                <AddNWOSubsidyDialog
                  onSelect={handleNWOSelect}
                  isSubmitting={isAddingNWO}
                  disabled={isAddingLineItem}
                />
                <Button size="sm" onClick={() => setIsAdding(true)}>
                  <Plus className="mr-2 h-4 w-4" /> Add Item
                </Button>
              </>
            )}
          </div>
        </CardHeader>
        <CardContent>
          {isAdding && (
            <div className="mb-6 p-4 border rounded-md bg-muted/50">
              <Form {...form}>
                <form
                  onSubmit={form.handleSubmit(onSubmit)}
                  className="grid gap-4 md:grid-cols-12"
                >
                  <FormField
                    control={form.control}
                    name="description"
                    render={({ field }) => (
                      <FormItem className="col-span-12 md:col-span-4">
                        <FormControl>
                          <Input placeholder="Description" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="category"
                    render={({ field }) => (
                      <FormItem className="col-span-6 md:col-span-2">
                        <Select
                          onValueChange={field.onChange}
                          value={field.value}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Category" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {Object.entries(categoryLabels).map(
                              ([value, label]) => (
                                <SelectItem key={value} value={value}>
                                  {label}
                                </SelectItem>
                              ),
                            )}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="fundingSource"
                    render={({ field }) => (
                      <FormItem className="col-span-6 md:col-span-2">
                        <Select
                          onValueChange={field.onChange}
                          value={field.value}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Source" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {Object.entries(fundingSourceLabels).map(
                              ([value, label]) => (
                                <SelectItem key={value} value={value}>
                                  {label}
                                </SelectItem>
                              ),
                            )}
                          </SelectContent>
                        </Select>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="year"
                    render={({ field }) => (
                      <FormItem className="col-span-4 md:col-span-1">
                        <FormControl>
                          <Input type="number" placeholder="Year" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="budgetedAmount"
                    render={({ field }) => (
                      <FormItem className="col-span-4 md:col-span-2">
                        <FormControl>
                          <Input
                            type="number"
                            placeholder="Amount"
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="col-span-4 md:col-span-1 flex gap-2">
                    <Button type="submit" size="icon">
                      <Save className="h-4 w-4" />
                    </Button>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      onClick={() => setIsAdding(false)}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </form>
              </Form>
            </div>
          )}

          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Description</TableHead>
                <TableHead>Category</TableHead>
                <TableHead>Source</TableHead>
                <TableHead>Year</TableHead>
                <TableHead className="text-right">Amount</TableHead>
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {budget.lineItems?.map((item) => (
                <TableRow key={item.id}>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      {item.description}
                      {item.nwoGrantId && (
                        <a // 483.20.036 -> 48320036
                          href={`https://www.nwo.nl/en/projects/${item.nwoGrantId.replaceAll(".", "")}`}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="inline-flex items-center justify-center"
                          title="View NWO Project"
                        >
                          <Button
                            size="icon"
                            variant="ghost"
                            className="h-6 w-6 ml-2"
                          >
                            <NwoIcon width={16} height={16} />
                          </Button>
                        </a>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant="outline">
                      {categoryLabels[item.category as BudgetCategory] ||
                        item.category}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Badge variant="secondary">
                      {fundingSourceLabels[
                        item.fundingSource as FundingSource
                      ] || item.fundingSource}
                    </Badge>
                  </TableCell>
                  <TableCell>{item.year}</TableCell>
                  <TableCell className="text-right">
                    {new Intl.NumberFormat("nl-NL", {
                      style: "currency",
                      currency: "EUR",
                    }).format(item.budgetedAmount || 0)}
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() =>
                        removeLineItem({
                          budgetId: budget.id!,
                          lineItemId: item.id!,
                        })
                      }
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
              {(!budget.lineItems || budget.lineItems.length === 0) && (
                <TableRow>
                  <TableCell
                    colSpan={6}
                    className="text-center text-muted-foreground py-8"
                  >
                    No line items found. Add one manually or link an NWO
                    subsidy.
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
