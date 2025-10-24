export interface UserMembership {
  id?: number; // int64 → number
  user_id: number; // required
  package_id: number; // required
  expired_at: string; // required
  is_membership_active: boolean; // required
  created_by?: string;
  updated_by?: string;
  created_at?: string;
  updated_at?: string;
}
