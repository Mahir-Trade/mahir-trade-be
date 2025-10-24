export interface Transaction {
  id?: number; // int64 → number
  order_id: number; // required
  amount: number; // required
  settlement_date?: string;
  webhook_id?: string;
  created_by?: string;
  updated_by?: string;
  created_at?: string;
  updated_at?: string;
}
