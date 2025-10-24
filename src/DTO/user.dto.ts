// --- DTO ---
export interface LoginReq {
  identity: string;
  password: string;
}

export interface GoogleLoginReq {
  state: string;
  code: string;
}

export interface ForgotPasswordReq {
  email: string;
}

export interface JWTData {
  email: string;
  user_id: number;
  username: string;
}

export interface LoginResponse {
  token: string;
  expire: Date;
}

export interface ResetPasswordRequest {
  queryOne: string;
  queryTwo: string;
  password: string;
  password_confirmation: string;
}

// ---- Auth related ----
export interface LoginReq {
  identity: string;
  password: string;
}

export interface GoogleLoginReq {
  state: string;
  code: string;
}

export interface ForgotPasswordReq {
  email: string;
}

export interface ResetPasswordRequest {
  q1: string;
  q2: string;
  password: string;
  password_confirmation: string;
}

export interface UserVerificationReq {
  email: string;
  expired_time: string;
}

export interface LoginResponse {
  token: string;
  expire: Date;
}

export interface UserRegistrationRequest {
  email: string;
  phone_number?: string;
  username: string;
  password: string;
  password_confirmation: string;
  referral_code?: string;
}

// ---- User Context Request ----
export interface UserCtxReq {
  user_id: number;
  email: string;
  username: string;
}

// ---- Update Profile ----
export interface UpdateProfileRequest {
  fullname?: string;
  phone_number?: string;
  address?: string;
  city?: string;
  country?: string;
  postal_code?: string;
  about_me?: string;
  job_title?: string;
  company?: string;
  facebook_url?: string;
  twitter_url?: string;
  linkedin_url?: string;
  instagram_url?: string;
  profile_picture_url?: string;
}

// ---- Change Password ----
export interface ChangePasswordRequest {
  current_password: string;
  new_password: string;
  new_password_confirmation: string;
}
