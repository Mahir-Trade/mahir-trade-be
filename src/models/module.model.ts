export interface Module {
  id?: number;
  uuid?: string;
  group_id?: number | null;
  module_name: string;
  thumbnail_url?: string | null;
  tag?: string | null;
  created_by?: string | null; // ✅ tambahkan | null
  updated_by?: string | null; // ✅ tambahkan | null
  created_at?: string | null; // ✅ tambahkan | null
  updated_at?: string | null; // ✅ tambahkan | null
}

export interface GetModulesRequest {
  limit: number; // int64 → number
  page: number; // int64 → number
}
