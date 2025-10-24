import { GeneralLog } from "@dto/generaLog.dto";
import db from "../db/db.config";
import { Pool } from "pg";

export class GeneralLogRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- CREATE GENERAL LOG ---
  async createGeneralLog(req: GeneralLog): Promise<boolean> {
    try {
      const query = `
        INSERT INTO general_logs (user_id, raw_body, created_by)
        VALUES ($1, $2, $3)
      `;

      const result = await this.pg.query(query, [
        req.user_id,
        req.raw_body,
        req.created_by,
      ]);

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(
        `[repo][general_log][createGeneralLog] error: ${err.message}`
      );
      throw err;
    }
  }
}
