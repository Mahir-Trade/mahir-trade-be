import axios from "axios";
import { EconomicCalendarItem } from "../models/economicCalendar.model";

export class EconomicCalendarService {
  constructor() {
    // bisa tambahkan dependency di sini kalau perlu (misalnya logger)
  }

  async getCalendar(from: string, to: string): Promise<EconomicCalendarItem[]> {
    const apiKey = process.env.ECONOMIC_CALENDAR;
    if (!apiKey) {
      throw new Error("missing ECONOMIC_CALENDAR in env");
    }

    // sama persis dengan Go → country = "united states"
    const url = `https://api.tradingeconomics.com/calendar/country/united states/${from}/${to}?c=${apiKey}`;

    try {
      const response = await axios.get(url);
      return response.data as EconomicCalendarItem[];
    } catch (err: any) {
      throw new Error(`failed to fetch calendar data: ${err.message}`);
    }
  }
}
