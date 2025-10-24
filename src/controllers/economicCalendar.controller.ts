import { Request, Response } from "express";
import { DefaultResponse } from "../utils/response";
import { EconomicCalendarService } from "../services/economicCalendar.service";

export class EconomicCalendarController {
  private economicCalendarService: EconomicCalendarService;

  constructor() {
    this.economicCalendarService = new EconomicCalendarService();

    this.getCalendar = this.getCalendar.bind(this);
  }

  // --- GET CALENDAR ---
  async getCalendar(req: Request, res: Response): Promise<Response> {
    try {
      const from = req.query.from as string;
      const to = req.query.to as string;

      if (!from || !to) {
        return res.status(400).json({
          code: 400,
          message: "Missing query parameters: from or to",
        } as DefaultResponse);
      }

      const result = await this.economicCalendarService.getCalendar(from, to);

      return res.status(200).json({
        code: 200,
        message: "Success",
        data: result,
      } as DefaultResponse);
    } catch (err: any) {
      console.error("[EconomicCalendarController][getCalendar] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Failed to fetch calendar data",
        error: err.message,
      } as DefaultResponse);
    }
  }
}
