// All your model interfaces
export interface BaseModel {
  id: number;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface User extends BaseModel {
  name: string;
  email: string;
}

export interface Template extends BaseModel {
  original_url: string;
  html_path: string;
  file_paths: string;
  status: string;
  error_message?: string;
}

export interface ConvertUrlRequest {
  url: string;
}

export interface ConvertUrlResponse {
  message: string;
  conversion: Template;
  html_path: string;
  file_paths: string;
}