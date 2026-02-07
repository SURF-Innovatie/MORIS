import { z } from "zod";

export const projectFormSchema = z.object({
  title: z.string().min(1, "Title is required"),
  slug: z
    .string()
    .min(1, "Slug is required")
    .regex(
      /^[a-z0-9]+(?:-[a-z0-9]+)*$/,
      "Invalid slug format (lowercase, hyphens)",
    ),
  description: z.string().min(1, "Description is required"),
  startDate: z.date({ message: "Start date is required" }),
  endDate: z.date({ message: "End date is required" }),
  organisationID: z.string().uuid("Invalid organisation ID"),
  customFields: z.record(z.string(), z.any()).optional(),
});

export type ProjectFormValues = z.infer<typeof projectFormSchema>;
