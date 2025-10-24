// User entity
export interface User {
  user_id?: number;
  uuid?: string;
  phone_number?: string;
  username: string;
  email: string;
  password: string;
  verified_at?: Date | null;
  is_active?: boolean;
  created_at?: Date;
  created_by?: string;
}

// BO (Backoffice) response
export interface GetUsersBOResponse {
  user_id: number;
  uuid: string;
  phone_number?: string;
  email: string;
  username: string;
  is_active: boolean;
  account_type: string;
  membership_expired_date?: Date | null;
  created_at: Date;
  created_by: string;
}
