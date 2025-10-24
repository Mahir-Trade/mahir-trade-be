export interface UserSubModule {
  id?: number; // int64 → number
  uuid?: string; // optional
  user_id: number; // required
  sub_module_id: number; // required
  created_by?: string;
  updated_by?: string;
  created_at?: string;
  updated_at?: string;
}
