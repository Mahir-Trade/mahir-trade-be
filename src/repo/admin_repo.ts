import { Pool } from "pg";
import { Admin } from "../models/admin.model";

import db from "../db/db.config";
import { AdminQueries } from "./queries/admin_queries";

export class AdminRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- FIND BY USERNAME ---
  async findByUsername(username: string): Promise<Admin | null> {
    try {
      const result = await this.pg.query(AdminQueries.FindByEmail, [username]);

      if (result.rowCount === 0) {
        return null; // sama seperti return kosong di Go
      }

      return result.rows[0] as Admin;
    } catch (err: any) {
      console.error(`[repo][admin][findByUsername] error: ${err.message}`);
      throw err;
    }
  }

  // --- CREATE ADMIN ---
  async createAdmin(req: Admin): Promise<number> {
    try {
      const result = await this.pg.query(AdminQueries.CreateAdmin, [
        req.email,
        req.username,
        req.password,
      ]);

      if (result.rowCount === 0) {
        throw new Error("failed to create admin");
      }

      return result.rows[0].id as number; // sama dengan rows.Scan(&id) di Go
    } catch (err: any) {
      if (
        err.message.includes("duplicate key value") &&
        err.message.includes("unique constraint")
      ) {
        throw new Error("email or username already exist");
      }
      console.error(`[repo][admin][createAdmin] error: ${err.message}`);
      throw err;
    }
  }

  // --- SOFT DELETE ADMIN ---
  async softDeleteAdmin(id: number, operator: string): Promise<boolean> {
    try {
      const result = await this.pg.query(AdminQueries.SoftDeleteAdmin, [
        operator,
        operator,
        id,
      ]);

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[repo][admin][softDeleteAdmin] error: ${err.message}`);
      throw err;
    }
  }
}
