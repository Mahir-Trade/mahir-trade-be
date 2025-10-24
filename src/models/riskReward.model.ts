export interface RiskRewardRequest {
  symbol: string;
  type: "buy" | "sell";
  open_price: number; // float64 → number
  sl_price: number; // float64 → number
  tp_price: number; // float64 → number
  balance: number; // float64 → number
  risk_pct: number; // float64 → number
}

export interface RiskResult {
  symbol: string;
  sl_pips: number; // float64 → number
  tp_pips: number; // float64 → number
  lot_size: number;
  total_risk: number;
  total_reward: number;
  risk_reward_ratio: number;
  rounded_lot: number;
  rounded_risk: number;
  rounded_tp: number;
  nett_sl: string;
  nett_tp: string;
}
