import { ApiService } from "../apiService";
import { Template, ConvertUrlResponse } from "@/types/models";
import { API_ENDPOINTS } from "../constants";
import { replaceParams } from "@/lib/utils";

class TemplateService extends ApiService<Template> {
  constructor() {
    super(API_ENDPOINTS.templates.list.path);
  }

  async convertUrl(url: string): Promise<ConvertUrlResponse> {
    const response = await this.post<ConvertUrlResponse>(
      API_ENDPOINTS.templates.convert.path,
      { url }
    );
    return response;
  }
  async fetchContent(id: number): Promise<string> {
    const response = await this.get<string>(
      replaceParams(API_ENDPOINTS.templates.fetchContent.path, { id })
    );
    return response;
  }
}

export const templateService = new TemplateService();
