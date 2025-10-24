export interface SubModuleRequest {
  module_id?: number;
  sub_module_name: string;
  title: string;
  video_url: string;
}

export interface SubModuleResponse {
  id: number;
  uuid: string;
  module_id?: number;
  module_name?: string;
  sub_module_name: string;
  title: string;
  video_url?: string;
  status?: string;
  created_by: string;
  updated_by: string;
  created_at?: string;
  updated_at?: string;
}

export interface MarkSubModuleAsWatchedRequest {
  subModuleId: number;
}
