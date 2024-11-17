import { BASE_URL } from "@/api/constants";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const replaceParams = (
  url: string,
  params: Record<string, string | number>
) => {
  let result = url;
  Object.entries(params).forEach(([key, value]) => {
    result = result.replace(`:${key}`, String(value));
  });
  return result;
};

export const getFullUrl = (path: string) => `${BASE_URL}${path}`;
