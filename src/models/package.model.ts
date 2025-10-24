export interface Package {
  id?: number; // int64 → number, optional
  price: number; // required
  durationInMonth: number; // required
  description: string; // required
  discountedPrice?: number; // optional
  discountExpired?: string; // optional
  createdBy?: string; // optional
  updatedBy?: string; // optional
  createdAt?: string; // optional
  updatedAt?: string; // optional
}
