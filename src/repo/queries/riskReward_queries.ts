export const RiskRewardQueries = {
  InsertRiskReward: `
    INSERT INTO risk_rewards (
      symbol, sl_pips, tp_pips, lot_size,
      total_risk, total_reward, risk_reward_ratio,
      rounded_lot, rounded_risk, rounded_tp,
      nett_sl, nett_tp
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    RETURNING symbol
  `,

  GetRiskRewardBySymbol: `
    SELECT symbol, sl_pips, tp_pips, lot_size,
           total_risk, total_reward, risk_reward_ratio,
           rounded_lot, rounded_risk, rounded_tp,
           nett_sl, nett_tp
    FROM risk_rewards
    WHERE symbol = $1
  `,
};
