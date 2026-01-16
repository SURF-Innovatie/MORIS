import { z } from "zod";

export const projectFormSchema = z.object({
  title: z.string().min(1, "Title is required"),
  description: z.string().min(1, "Description is required"),
  startDate: z.date({ message: "Start date is required" }),
  endDate: z.date({ message: "End date is required" }),
  organisationID: z.string().uuid("Invalid organisation ID"),
  customFields: z.record(z.string(), z.any()).optional(),
});

export type ProjectFormValues = z.infer<typeof projectFormSchema>;
