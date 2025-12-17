import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Loader2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/use-toast";
import { usePostProjects } from "@api/moris";

import { projectFormSchema } from "@/lib/schemas/project";
import { ProjectFields } from "@/components/project-form/ProjectFields";
import { EMPTY_UUID } from "@/lib/constants";

export default function CreateProjectRoute() {
  const navigate = useNavigate();
  const { toast } = useToast();

  const { mutateAsync: createProject, isPending: isCreating } =
    usePostProjects();

  const form = useForm<z.infer<typeof projectFormSchema>>({
    resolver: zodResolver(projectFormSchema),
    defaultValues: {
      title: "",
      description: "",
      // TODO: This should be dynamic or selected from a list
      organisationID: EMPTY_UUID,
    },
  });

  async function onSubmit(values: z.infer<typeof projectFormSchema>) {
    try {
      await createProject({
        data: {
          title: values.title,
          description: values.description,
          start_date: values.startDate.toISOString(),
          end_date: values.endDate.toISOString(),
          owning_org_node_id: values.organisationID,
        },
      });
      toast({
        title: "Project created",
        description: "The new project has been successfully created.",
      });
      navigate("/dashboard");
    } catch (error) {
      toast({
        variant: "destructive",
        title: "Error",
        description: "Something went wrong. Please try again.",
      });
    }
  }

  return (
    <div className="mx-auto max-w-2xl py-8">
      <div className="mb-8">
        <h1 className="text-3xl font-bold tracking-tight">Create Project</h1>
        <p className="text-muted-foreground">
          Fill in the details to start a new project.
        </p>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <ProjectFields form={form} />

          {/* Organisation ID is still specific here as it might be editable during creation but not update */}
          <FormField
            control={form.control}
            name="organisationID"
            render={({ field }) => (
              <FormItem className="max-w-2xl">
                <FormLabel>Organisation ID</FormLabel>
                <FormControl>
                  <Input placeholder="UUID" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex justify-end gap-4">
            <Button
              type="button"
              variant="ghost"
              onClick={() => navigate("/dashboard")}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isCreating}>
              {isCreating ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Saving...
                </>
              ) : (
                "Create Project"
              )}
            </Button>
          </div>
        </form>
      </Form>
    </div>
  );
}
