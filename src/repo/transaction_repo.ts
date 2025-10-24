import { Pool } from "pg";
import { Transaction } from "../models/transaction.model";

import db from "../db/db.config";
import { TransactionQueries } from "./queries/transaction_queries";

export class TransactionRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- CREATE TRANSACTION ---
  async createTransaction(req: Transaction): Promise<number> {
    try {
      const result = await this.pg.query(TransactionQueries.CreateTransaction, [
        req.order_id,
        req.amount,
        req.settlement_date,
        req.webhook_id,
        req.created_by,
      ]);

      if (result.rowCount === 0) {
        throw new Error("failed to create transaction");
      }

      return result.rows[0].id as number;
    } catch (err: any) {
      console.error(
        `[TransactionRepository][createTransaction] error: ${err.message}`
      );
      throw err;
    }
  }
}
