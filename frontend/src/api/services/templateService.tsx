import { ApiService } from "../apiService";
import { Template, ConvertUrlResponse } from "@/types/models";
import { API_ENDPOINTS } from "../constants";

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
  async fetchContent(htmlPath: string): Promise<string> {
    // This method should fetch the HTML content using the html_path received from convertUrl
    const response = await this.get<string>(htmlPath);
    return response;
  }
}

export const templateService = new TemplateService();
