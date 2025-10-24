import { Pool } from "pg";
import { UserSubModule } from "../models/userSubModule.model";
import db from "../db/db.config";
import { UserSubModuleQueries } from "./queries/userSubModule_queries";

export class UserSubModuleRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- CREATE USER SUB MODULE ---
  async createUserSubModule(req: UserSubModule): Promise<number> {
    try {
      const result = await this.pg.query(
        UserSubModuleQueries.CreateUserSubModule,
        [req.user_id, req.sub_module_id, req.created_by]
      );

      if (result.rowCount === 0) {
        throw new Error("failed to create user sub module");
      }

      return result.rows[0].id as number;
    } catch (err: any) {
      console.error(
        `[UserSubModuleRepository][createUserSubModule] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- GET USER SUB MODULE BY SUBMODULE ID & USER ID ---
  async getUserSubModuleBySubModuleIDAndUserID(
    userId: number,
    subModuleId: number
  ): Promise<UserSubModule | null> {
    try {
      const result = await this.pg.query(
        UserSubModuleQueries.GetUserSubModuleBySubModuleIDAndUserID,
        [subModuleId, userId]
      );

      if (result.rowCount === 0) {
        return null; // sama kayak Go return struct kosong
      }

      return result.rows[0] as UserSubModule;
    } catch (err: any) {
      console.error(
        `[UserSubModuleRepository][getUserSubModuleBySubModuleIDAndUserID] error: ${err.message}`
      );
      throw err;
    }
  }
}
