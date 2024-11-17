import { ApiService } from "../apiService";
import {
  Template,
  ConvertUrlRequest,
  ConvertUrlResponse,
} from "@/types/models";
import { API_ENDPOINTS } from "../constants";

class TemplateService extends ApiService<Template> {
  constructor() {
    super(API_ENDPOINTS.templates.list.path);
  }

  async convertUrl(url: string): Promise<ConvertUrlResponse> {
    return this.post<ConvertUrlResponse>(API_ENDPOINTS.templates.convert.path, {
      url,
    });
  }

  async fetchContent(path: string): Promise<string> {
    return this.get<string>(path);
  }
}

export const templateService = new TemplateService();
