export const BASE_URL = "https://rnwpb-161-139-102-112.a.free.pinggy.link";

// Define endpoint types for better type safety
export type EndpointConfig = {
  path: string;
  method: "GET" | "POST" | "PUT" | "DELETE";
};

export const API_ENDPOINTS = {
  users: {
    list: { path: "/api/users", method: "GET" },
    view: { path: "/api/users/:id", method: "GET" },
    create: { path: "/api/users", method: "POST" },
    update: { path: "/api/users/:id", method: "PUT" },
    delete: { path: "/api/users/:id", method: "DELETE" },
  },
  templates: {
    list: { path: "/api/templates", method: "GET" },
    view: { path: "/api/templates/:id", method: "GET" },
    create: { path: "/api/templates", method: "POST" },
    update: { path: "/api/templates/:id", method: "PUT" },
    delete: { path: "/api/templates/:id", method: "DELETE" },
    convert: { path: "/api/templates/convert", method: "POST" },
    fetchContent: { path: "/api/templates/:id/content", method: "GET" },
  },
} as const;
