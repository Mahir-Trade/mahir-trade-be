import { Pool } from "pg";

import db from "../db/db.config";
import { EmailQueries } from "./queries/email_queries";

export class EmailRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  async getByKey(key: string): Promise<string> {
    try {
      const result = await this.pg.query(EmailQueries.GetByKey, [key]);

      if (result.rowCount === 0) {
        throw new Error(`email template not found for key: ${key}`);
      }

      return result.rows[0].body as string;
    } catch (err: any) {
      console.error(`[repo][email][getByKey] error: ${err.message}`);
      throw err;
    }
  }
}
