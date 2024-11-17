import { getFullUrl, replaceParams } from "@/lib/utils";
import { PaginationQuery, PaginationResponse } from "@/types/api";

export class ApiService<T extends { id: number }> {
  constructor(private baseUrl: string) {}

  protected async request<R>(
    path: string,
    options: RequestInit = {}
  ): Promise<R> {
    const url = getFullUrl(path);
    const response = await fetch(url, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
    });

    if (!response.ok) {
      throw new Error(`API Error: ${response.statusText}`);
    }

    return response.json();
  }

  async list(query?: PaginationQuery): Promise<PaginationResponse<T>> {
    const queryString = query ? `?${new URLSearchParams(query as any)}` : "";
    return this.request<PaginationResponse<T>>(`${this.baseUrl}${queryString}`);
  }

  async view(id: number): Promise<T> {
    return this.request<T>(replaceParams(this.baseUrl + "/:id", { id }));
  }

  async create(data: Omit<T, "id">): Promise<T> {
    return this.request<T>(this.baseUrl, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async update(id: number, data: Partial<T>): Promise<T> {
    return this.request<T>(replaceParams(this.baseUrl + "/:id", { id }), {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async delete(id: number): Promise<void> {
    return this.request(replaceParams(this.baseUrl + "/:id", { id }), {
      method: "DELETE",
    });
  }

  protected async get<R>(path: string): Promise<R> {
    return this.request<R>(path);
  }

  protected async post<R>(path: string, data: unknown): Promise<R> {
    return this.request<R>(path, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }
}
