import { ApiService } from "../apiService";
import { User } from "@/types/models";
import { API_ENDPOINTS } from "../constants";

class UserService extends ApiService<User> {
  constructor() {
    super(API_ENDPOINTS.users.list.path);
  }

  // Add user-specific methods here
  async getCurrentUser() {
    return this.get("/api/users/me");
  }
}

export const userService = new UserService();
