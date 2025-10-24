import { RiskResult, RiskRewardRequest } from "../models/riskReward.model";

export class RiskRewardService {
  private pair_configs: Record<
    string,
    { pip_digits: number; pip_value_per_lot: number }
  >;

  constructor() {
    this.pair_configs = {
      EURUSD: { pip_digits: 4, pip_value_per_lot: 10 },
      GBPUSD: { pip_digits: 4, pip_value_per_lot: 10 },
      USDJPY: { pip_digits: 2, pip_value_per_lot: 9 },
      XAUUSD: { pip_digits: 1, pip_value_per_lot: 1 },
      USDCAD: { pip_digits: 4, pip_value_per_lot: 10 },
      AUDUSD: { pip_digits: 4, pip_value_per_lot: 10 },
      NZDUSD: { pip_digits: 4, pip_value_per_lot: 10 },
      USDCHF: { pip_digits: 4, pip_value_per_lot: 10 },
      EURJPY: { pip_digits: 2, pip_value_per_lot: 9 },
      GBPJPY: { pip_digits: 2, pip_value_per_lot: 9 },
      EURGBP: { pip_digits: 4, pip_value_per_lot: 10 },
      GBPCHF: { pip_digits: 4, pip_value_per_lot: 10 },
      BTCUSD: { pip_digits: 2, pip_value_per_lot: 1 },
      ETHUSD: { pip_digits: 2, pip_value_per_lot: 1 },
      XAGUSD: { pip_digits: 3, pip_value_per_lot: 0.5 },
    };
  }

  async risk_reward(reqs: RiskRewardRequest[]): Promise<RiskResult[]> {
    const results: RiskResult[] = [];

    for (const req of reqs) {
      const symbol = req.symbol.toUpperCase();
      const conf = this.pair_configs[symbol] || {
        pip_digits: 4,
        pip_value_per_lot: 10,
      };

      let sl_pips: number;
      let tp_pips: number;

      if (req.type === "buy") {
        sl_pips =
          (req.sl_price - req.open_price) * Math.pow(10, conf.pip_digits);
        tp_pips =
          (req.tp_price - req.open_price) * Math.pow(10, conf.pip_digits);
      } else {
        sl_pips =
          (req.open_price - req.sl_price) * Math.pow(10, conf.pip_digits);
        tp_pips =
          (req.open_price - req.tp_price) * Math.pow(10, conf.pip_digits);
      }

      sl_pips = this.round_to(sl_pips, 1);
      tp_pips = this.round_to(tp_pips, 1);

      const risk_amount = (req.balance * req.risk_pct) / 100;
      const lot_size =
        risk_amount / (Math.abs(sl_pips) * conf.pip_value_per_lot);

      const total_risk = lot_size * Math.abs(sl_pips) * conf.pip_value_per_lot;
      const total_reward = lot_size * tp_pips * conf.pip_value_per_lot;
      const risk_reward_ratio = tp_pips / Math.abs(sl_pips);

      const rounded_lot = this.round_to(lot_size, 2);
      const rounded_risk = this.round_to(total_risk, 2);
      const rounded_tp = this.round_to(total_reward, 2);

      const nett_sl = `-$${rounded_risk.toFixed(2)}`;
      const nett_tp = `$${rounded_tp.toFixed(2)}`;

      const result: RiskResult = {
        symbol: req.symbol,
        sl_pips,
        tp_pips,
        lot_size,
        total_risk,
        total_reward,
        risk_reward_ratio,
        rounded_lot,
        rounded_risk,
        rounded_tp,
        nett_sl,
        nett_tp,
      };

      results.push(result);
    }

    return results;
  }

  private round_to(value: number, decimals: number): number {
    const factor = Math.pow(10, decimals);
    return Math.round(value * factor) / factor;
  }
}
