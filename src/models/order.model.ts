export interface Order {
  id?: number; // int64 → number, optional
  user_id: number; // required
  package_id: number; // required
  status: string; // required
  payment_code: string; // required
  payment_url: string; // required
  created_by?: string;
  updated_by?: string;
  created_at?: string;
  updated_at?: string;
}
