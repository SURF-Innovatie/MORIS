import { useState } from "react";
import { UseFormReturn } from "react-hook-form";
import { format } from "date-fns";
import { CalendarIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { cn } from "@/lib/utils";
import { ProjectFormValues } from "@/lib/schemas/project";

interface ProjectFieldsProps {
  form: UseFormReturn<ProjectFormValues>;
  disabledFields?: {
    title?: boolean;
    description?: boolean;
    startDate?: boolean;
    endDate?: boolean;
  };
  pendingFields?: {
    title?: boolean;
    description?: boolean;
    startDate?: boolean;
    endDate?: boolean;
  };
}

export function ProjectFields({
  form,
  disabledFields,
  pendingFields,
}: ProjectFieldsProps) {
  const [isStartDateOpen, setIsStartDateOpen] = useState(false);
  const [isEndDateOpen, setIsEndDateOpen] = useState(false);
  const startDate = form.watch("startDate");

  const isPending = (key: keyof NonNullable<typeof pendingFields>) =>
    pendingFields?.[key];

  return (
    <div className="space-y-6">
      <FormField
        control={form.control}
        name="title"
        render={({ field }) => (
          <FormItem className="max-w-2xl">
            <FormLabel className="flex items-center gap-2">
              Title
              {isPending("title") && (
                <Badge variant="secondary" className="h-5 text-[10px]">
                  Pending Approval
                </Badge>
              )}
            </FormLabel>
            <FormControl>
              <Input
                placeholder="Project title"
                {...field}
                disabled={disabledFields?.title || isPending("title")}
              />
            </FormControl>
            <FormMessage />
          </FormItem>
        )}
      />

      <FormField
        control={form.control}
        name="description"
        render={({ field }) => (
          <FormItem className="max-w-2xl">
            <FormLabel className="flex items-center gap-2">
              Description
              {isPending("description") && (
                <Badge variant="secondary" className="h-5 text-[10px]">
                  Pending Approval
                </Badge>
              )}
            </FormLabel>
            <FormControl>
              <Textarea
                placeholder="Describe the project..."
                className="min-h-[120px] resize-none"
                {...field}
                disabled={
                  disabledFields?.description || isPending("description")
                }
              />
            </FormControl>
            <FormDescription>
              A brief summary of what this project is about.
            </FormDescription>
            <FormMessage />
          </FormItem>
        )}
      />

      <div className="grid gap-4 sm:grid-cols-2 max-w-2xl">
        <FormField
          control={form.control}
          name="startDate"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel className="flex items-center gap-2">
                Start Date
                {isPending("startDate") && (
                  <Badge variant="secondary" className="h-5 text-[10px]">
                    Pending
                  </Badge>
                )}
              </FormLabel>
              <Popover open={isStartDateOpen} onOpenChange={setIsStartDateOpen}>
                <PopoverTrigger asChild>
                  <FormControl>
                    <Button
                      variant={"outline"}
                      className={cn(
                        "w-full pl-3 text-left font-normal",
                        !field.value && "text-muted-foreground"
                      )}
                      disabled={
                        disabledFields?.startDate || isPending("startDate")
                      }
                    >
                      {field.value ? (
                        format(field.value, "PPP")
                      ) : (
                        <span>Pick a date</span>
                      )}
                      <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={field.value}
                    onSelect={(date) => {
                      field.onChange(date);
                      setIsStartDateOpen(false);
                    }}
                    disabled={(date) => date < new Date("1900-01-01")}
                    initialFocus
                    captionLayout="dropdown"
                    fromYear={1900}
                    toYear={2100}
                  />
                </PopoverContent>
              </Popover>
              <FormMessage />
            </FormItem>
          )}
        />

        <FormField
          control={form.control}
          name="endDate"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel className="flex items-center gap-2">
                End Date
                {isPending("endDate") && (
                  <Badge variant="secondary" className="h-5 text-[10px]">
                    Pending
                  </Badge>
                )}
              </FormLabel>
              <Popover open={isEndDateOpen} onOpenChange={setIsEndDateOpen}>
                <PopoverTrigger asChild>
                  <FormControl>
                    <Button
                      variant={"outline"}
                      className={cn(
                        "w-full pl-3 text-left font-normal",
                        !field.value && "text-muted-foreground"
                      )}
                      disabled={disabledFields?.endDate || isPending("endDate")}
                    >
                      {field.value ? (
                        format(field.value, "PPP")
                      ) : (
                        <span>Pick a date</span>
                      )}
                      <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={field.value}
                    onSelect={(date) => {
                      field.onChange(date);
                      setIsEndDateOpen(false);
                    }}
                    disabled={(date) =>
                      date < new Date("1900-01-01") ||
                      (startDate ? date < startDate : false)
                    }
                    initialFocus
                    captionLayout="dropdown"
                    fromYear={1900}
                    toYear={2100}
                  />
                </PopoverContent>
              </Popover>
              <FormMessage />
            </FormItem>
          )}
        />
      </div>
    </div>
  );
}
