import { Request, Response } from "express";
import { RiskRewardService } from "../services/riskReward.service";
import { DefaultResponse } from "../utils/response";
import { RiskRewardRequest } from "../models/riskReward.model";
import { validate } from "class-validator";

export class RiskRewardController {
  private risk_reward_service: RiskRewardService;

  constructor() {
    this.risk_reward_service = new RiskRewardService();

    this.risk_reward = this.risk_reward.bind(this);
  }

  async risk_reward(req: Request, res: Response): Promise<Response> {
    try {
      const body = req.body as RiskRewardRequest[];

      if (!Array.isArray(body) || body.length === 0) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request body",
          error: "Request body must be a non-empty array of RiskRewardRequest",
        } as DefaultResponse);
      }

      // ✅ Call service
      const result = await this.risk_reward_service.risk_reward(body);

      return res.status(200).json({
        code: 200,
        message: "Risk Reward calculation successful",
        data: result,
      } as DefaultResponse);
    } catch (err: any) {
      console.error("[RiskRewardController][risk_reward] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal server error",
        error: err.message,
      } as DefaultResponse);
    }
  }
}
