// Export services
export { userService } from "./services/userService";
export { templateService } from "./services/templateService";

// Export base service class (if needed for extending in other places)
export { ApiService } from "./apiService";

// Export types
export type { User } from "@/types/models";
export type { Template } from "@/types/models";
export type { PaginationQuery, PaginationResponse } from "@/types/api";

// Export constants
export { API_ENDPOINTS, BASE_URL } from "./constants";
