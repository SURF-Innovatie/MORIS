import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Plus, Trash2, ArrowLeft, Save } from "lucide-react";

import {
  useGetProjectsProjectIdBudget,
  usePostBudgetsBudgetIdLineItems,
  useDeleteBudgetsBudgetIdLineItemsLineItemId,
} from "@api/moris";
import { BudgetCategory, FundingSource } from "@api/model";
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
    "other",
  ] as const),
  budgetedAmount: z.coerce.number().min(0, "Amount must be positive"),
  year: z.coerce.number().int().min(2020, "Invalid year"),
  fundingSource: z.enum([
    "subsidy",
    "cofinancing_cash",
    "cofinancing_inkind",
  ] as const),
});

export function BudgetEditor({ projectId, onDone }: BudgetEditorProps) {
  const { toast } = useToast();
  const queryClient = useQueryClient();
  const [isAdding, setIsAdding] = useState(false);

  const { data: budget } = useGetProjectsProjectIdBudget(projectId);

  const form = useForm<z.infer<typeof lineItemSchema>>({
    resolver: zodResolver(lineItemSchema),
    defaultValues: {
      description: "",
      category: "personnel",
      budgetedAmount: 0 as any, // Cast for input handling
      year: new Date().getFullYear() as any,
      fundingSource: "subsidy",
    },
  });

  const { mutate: addLineItem } = usePostBudgetsBudgetIdLineItems({
    mutation: {
      onSuccess: () => {
        // Invalidate all possible query key variations to ensure refresh
        queryClient.invalidateQueries({ queryKey: ["budget", projectId] });
        queryClient.invalidateQueries({
          queryKey: ["/projects", projectId, "budget"],
        });
        queryClient.invalidateQueries({
          queryKey: [`/projects/${projectId}/budget`],
        });

        toast({ title: "Line Item Added" });
        setIsAdding(false);
        form.reset();
      },
      onError: () => {
        toast({
          variant: "destructive",
          title: "Error",
          description: "Failed to add line item",
        });
      },
    },
  });

  const { mutate: removeLineItem } =
    useDeleteBudgetsBudgetIdLineItemsLineItemId({
      mutation: {
        onSuccess: () => {
          queryClient.invalidateQueries({ queryKey: ["budget", projectId] });
          queryClient.invalidateQueries({
            queryKey: ["/projects", projectId, "budget"],
          });
          queryClient.invalidateQueries({
            queryKey: [`/projects/${projectId}/budget`],
          });
          toast({ title: "Line Item Removed" });
        },
      },
    });

  const onSubmit = (values: z.infer<typeof lineItemSchema>) => {
    if (!budget?.id) return;
    addLineItem({ budgetId: budget.id, data: values });
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
          {!isAdding && (
            <Button size="sm" onClick={() => setIsAdding(true)}>
              <Plus className="mr-2 h-4 w-4" /> Add Item
            </Button>
          )}
        </CardHeader>
        <CardContent>
          {isAdding && (
            <div className="mb-6 p-4 border rounded-md bg-muted/50">
              <Form {...form}>
                <form
                  onSubmit={form.handleSubmit(onSubmit)}
                  className="grid gap-4 md:grid-cols-12 items-end"
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
                          defaultValue={field.value}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Category" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {Object.entries(categoryLabels).map(
                              ([key, label]) => (
                                <SelectItem key={key} value={key}>
                                  {label}
                                </SelectItem>
                              )
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
                          defaultValue={field.value}
                        >
                          <FormControl>
                            <SelectTrigger>
                              <SelectValue placeholder="Funding" />
                            </SelectTrigger>
                          </FormControl>
                          <SelectContent>
                            {Object.entries(fundingSourceLabels).map(
                              ([key, label]) => (
                                <SelectItem key={key} value={key}>
                                  {label}
                                </SelectItem>
                              )
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
                      <FormItem className="col-span-6 md:col-span-1">
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
                      <FormItem className="col-span-6 md:col-span-2">
                        <FormControl>
                          <Input
                            type="number"
                            placeholder="Amount"
                            className="w-full"
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                  <div className="flex gap-2 col-span-12 md:col-span-1">
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
                <TableHead>Funding</TableHead>
                <TableHead>Year</TableHead>
                <TableHead className="text-right">Amount</TableHead>
                <TableHead className="w-[50px]"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {(budget.lineItems || []).map((item) => (
                <TableRow key={item.id}>
                  <TableCell>{item.description || ""}</TableCell>
                  <TableCell>
                    {item.category
                      ? categoryLabels[item.category as BudgetCategory] ||
                        item.category
                      : ""}
                  </TableCell>
                  <TableCell>
                    {item.fundingSource
                      ? fundingSourceLabels[
                          item.fundingSource as FundingSource
                        ] || item.fundingSource
                      : ""}
                  </TableCell>
                  <TableCell>{item.year || ""}</TableCell>
                  <TableCell className="text-right">
                    €{(item.budgetedAmount || 0).toLocaleString()}
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() =>
                        item.id &&
                        budget?.id &&
                        removeLineItem({
                          budgetId: budget.id,
                          lineItemId: item.id,
                        })
                      }
                    >
                      <Trash2 className="h-4 w-4 text-destructive" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
              {(budget.lineItems || []).length === 0 && (
                <TableRow>
                  <TableCell
                    colSpan={6}
                    className="text-center text-muted-foreground h-24"
                  >
                    No items found. Add one to start.
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
