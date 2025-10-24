export interface Admin {
  adminId?: number; // int64 -> number
  uuid?: string; // string
  email: string; // required, must be email
  username: string; // required
  password?: string; // required
  createdAt?: string; // string timestamp
  updatedAt?: string; // string timestamp
}

export interface AdminLoginRequest {
  identity: string;
  password: string;
}

export interface UpdateTypeUserRequest {
  isActive: boolean;
}
