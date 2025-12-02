import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export const EMPTY_UUID = "00000000-0000-0000-0000-000000000000";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
