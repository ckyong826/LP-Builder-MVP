// Common API-related types
export interface PaginationQuery {
  page: number;
  pageSize: number;
  orderBy?: string;
  sort?: "asc" | "desc";
}

export interface PaginationResponse<T> {
  data: T[];
  total: number;
  page: number;
  size: number;
}

// You can also add other common API types here
export interface ApiError {
  error: string;
  details?: string;
}