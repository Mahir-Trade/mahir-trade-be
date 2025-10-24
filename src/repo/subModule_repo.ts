import db from "../db/db.config";
import { Pool } from "pg";
import { SubModuleQueries } from "./queries/subModule_queries";
import { SubModule } from "@models/subModule.model";
import { PaginationRequest } from "@utils/response";

export class SubModuleRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- CREATE SUB MODULE ---
  async createSubModule(req: SubModule): Promise<number> {
    try {
      const query = req.module_id
        ? SubModuleQueries.CreateSubModule
        : SubModuleQueries.CreateSubModuleWithoutModuleID;

      const params = req.module_id
        ? [
            req.module_id,
            req.sub_module_name,
            req.title,
            req.video_url,
            req.created_by,
          ]
        : [req.sub_module_name, req.title, req.video_url, req.created_by];

      const result = await this.pg.query(query, params);

      if (result.rowCount === 0) {
        throw new Error("failed to create submodule");
      }

      return result.rows[0].id as number;
    } catch (err: any) {
      console.error(
        `[SubModuleRepository][createSubModule] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- GET SUB MODULE BY ID ---
  async getSubModuleByID(id: number): Promise<SubModule | null> {
    try {
      const result = await this.pg.query(SubModuleQueries.GetSubModuleByID, [
        id,
      ]);

      if (result.rowCount === 0) return null;
      return result.rows[0] as SubModule;
    } catch (err: any) {
      console.error(
        `[SubModuleRepository][getSubModuleByID] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- GET SUB MODULES (paginated) ---
  async getSubModules(
    req: PaginationRequest
  ): Promise<{ subModules: SubModule[]; totalCount: number }> {
    try {
      const limit = req.limit || 10;
      const offset = req.page > 1 ? (req.page - 1) * limit : 0;

      const result = await this.pg.query(SubModuleQueries.GetSubModules, [
        limit,
        offset,
      ]);

      const subModules: SubModule[] = result.rows.map((row) => ({
        ...row,
      }));

      const totalCount =
        result.rows.length > 0 ? parseInt(result.rows[0].total_count, 10) : 0;

      return { subModules, totalCount };
    } catch (err: any) {
      console.error(
        `[SubModuleRepository][getSubModules] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- GET SUB MODULES BY MODULE ID ---
  async getSubModulesByModuleID(
    moduleID: number,
    userID: number,
    req: PaginationRequest
  ): Promise<{ subModules: SubModule[]; totalCount: number }> {
    try {
      const limit = req.limit || 10;
      const offset = req.page > 1 ? (req.page - 1) * limit : 0;

      const result = await this.pg.query(
        SubModuleQueries.GetSubModuleByModuleID,
        [moduleID, userID, limit, offset]
      );

      if (result.rowCount === 0) {
        throw new Error(`sub module with module id ${moduleID} not found`);
      }

      const subModules: SubModule[] = result.rows.map((row) => ({
        ...row,
      }));

      const totalCount =
        result.rows.length > 0 ? parseInt(result.rows[0].total_count, 10) : 0;

      return { subModules, totalCount };
    } catch (err: any) {
      console.error(
        `[SubModuleRepository][getSubModulesByModuleID] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- UPDATE SUB MODULE ---
  async updateSubModule(req: SubModule): Promise<boolean> {
    try {
      const result = await this.pg.query(SubModuleQueries.UpdateSubModule, [
        req.sub_module_name,
        req.title,
        req.video_url,
        req.updated_by,
        req.module_id,
        req.id,
      ]);

      if ((result.rowCount ?? 0) === 0) {
        throw new Error(`sub module with id ${req.id} not found`);
      }

      return true;
    } catch (err: any) {
      console.error(
        `[SubModuleRepository][updateSubModule] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- SOFT DELETE SUB MODULE ---
  async softDeleteSubModule(
    subModuleId: number,
    operator: string
  ): Promise<boolean> {
    try {
      const result = await this.pg.query(SubModuleQueries.SoftDeleteSubModule, [
        operator,
        operator,
        subModuleId,
      ]);

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(
        `[SubModuleRepository][softDeleteSubModule] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- REMOVE MODULE ID FROM SUB MODULES ---
  async removeModuleIDFromSubModules(
    moduleID: number,
    operator: string
  ): Promise<boolean> {
    try {
      const result = await this.pg.query(
        SubModuleQueries.RemoveModuleIDFromSubModules,
        [operator, moduleID]
      );

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(
        `[SubModuleRepository][removeModuleIDFromSubModules] error: ${err.message}`
      );
      throw err;
    }
  }
}
