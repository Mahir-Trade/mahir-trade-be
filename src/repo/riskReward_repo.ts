import { RiskResult } from "@models/riskReward.model";
import db from "../db/db.config";
import { Pool } from "pg";
import { RiskRewardQueries } from "./queries/riskReward_queries";

export class RiskRewardRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- SAVE RISK REWARD RESULT ---
  async save(result: RiskResult): Promise<boolean> {
    try {
      const res = await this.pg.query(RiskRewardQueries.InsertRiskReward, [
        result.symbol,
        result.sl_pips,
        result.tp_pips,
        result.lot_size,
        result.total_risk,
        result.total_reward,
        result.risk_reward_ratio,
        result.rounded_lot,
        result.rounded_risk,
        result.rounded_tp,
        result.nett_sl,
        result.nett_tp,
      ]);

      return (res.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[repo][riskReward][save] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET RISK REWARD BY SYMBOL ---
  async getBySymbol(symbol: string): Promise<RiskResult[]> {
    try {
      const result = await this.pg.query(
        RiskRewardQueries.GetRiskRewardBySymbol,
        [symbol]
      );

      return (result.rows as any[]).map((row) => ({
        symbol: row.symbol,
        sl_pips: parseFloat(row.sl_pips),
        tp_pips: parseFloat(row.tp_pips),
        lot_size: parseFloat(row.lot_size),
        total_risk: parseFloat(row.total_risk),
        total_reward: parseFloat(row.total_reward),
        risk_reward_ratio: parseFloat(row.risk_reward_ratio),
        rounded_lot: parseFloat(row.rounded_lot),
        rounded_risk: parseFloat(row.rounded_risk),
        rounded_tp: parseFloat(row.rounded_tp),
        nett_sl: row.nett_sl,
        nett_tp: row.nett_tp,
      }));
    } catch (err: any) {
      console.error(`[repo][riskReward][getBySymbol] error: ${err.message}`);
      throw err;
    }
  }
}
