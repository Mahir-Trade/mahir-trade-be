export interface ModuleRequest {
  module_name: string; // required
  group_id?: number; // optional
  thumbnail_url?: string; // optional
  tag?: string; // optional
  created_by?: string; // optional
  updated_by?: string; // optional
}

export interface ModuleResponse {
  id: number;
  uuid: string;
  group_id: number;
  group_name: string;
  module_name: string;
  thumbnail_url: string;
  tag: string;
  created_by: string;
  updated_by: string;
  created_at?: string;
  updated_at?: string;
}
