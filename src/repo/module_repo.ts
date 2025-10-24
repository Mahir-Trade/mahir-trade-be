import { Pool } from "pg";

import { Module } from "../models/module.model";

import db from "../db/db.config";
import { ModuleQueries } from "./queries/module_queries";
import { PaginationRequest } from "@utils/response";

export class ModuleRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- CREATE MODULE ---
  async createModule(req: Module): Promise<number> {
    let query: string;
    let params: any[];

    if (req.group_id && req.tag) {
      query = ModuleQueries.CreateModuleWithGroupIDAndTag;
      params = [
        req.group_id,
        req.module_name,
        req.thumbnail_url,
        req.tag,
        req.created_by,
      ];
    } else if (req.group_id) {
      query = ModuleQueries.CreateModuleWithGroupID;
      params = [
        req.group_id,
        req.module_name,
        req.thumbnail_url,
        req.created_by,
      ];
    } else if (req.tag) {
      query = ModuleQueries.CreateModuleWithoutGroupID;
      params = [req.module_name, req.thumbnail_url, req.tag, req.created_by];
    } else {
      query = ModuleQueries.CreateModuleWithoutGroupIDAndTag;
      params = [req.module_name, req.thumbnail_url, req.created_by];
    }

    const result = await this.pg.query<{ id: number }>(query, params);

    if (result.rowCount === 0) {
      throw new Error("failed to create module");
    }

    return result.rows[0].id;
  }

  // --- GET MODULE BY ID ---
  async getModuleByID(id: number): Promise<Module> {
    const result = await this.pg.query<{
      id: number;
      uuid: string;
      group_id: number | null;
      module_name: string;
      thumbnail_url: string | null;
      tag: string | null;
      created_by: string | null;
      created_at: string;
      updated_at: string | null;
      updated_by: string | null;
    }>(ModuleQueries.GetModuleByID, [id]);

    if (result.rowCount === 0) {
      throw new Error(`module with id ${id} not found`);
    }

    const row = result.rows[0];
    return {
      id: row.id,
      uuid: row.uuid,
      group_id: row.group_id,
      module_name: row.module_name,
      thumbnail_url: row.thumbnail_url,
      tag: row.tag,
      created_by: row.created_by,
      created_at: row.created_at,
      updated_at: row.updated_at,
      updated_by: row.updated_by,
    };
  }

  // --- UPDATE MODULE ---
  async updateModule(req: Module): Promise<boolean> {
    let query: string;
    let params: any[];

    if (req.thumbnail_url && req.tag) {
      query = ModuleQueries.UpdateModuleWithThumbnailAndTag;
      params = [
        req.module_name,
        req.thumbnail_url,
        req.tag,
        req.updated_by,
        req.id,
      ];
    } else if (req.thumbnail_url) {
      query = ModuleQueries.UpdateModuleWithThumbnail;
      params = [req.module_name, req.thumbnail_url, req.updated_by, req.id];
    } else if (req.tag) {
      query = ModuleQueries.UpdateModuleWithTag;
      params = [req.module_name, req.tag, req.updated_by, req.id];
    } else {
      query = ModuleQueries.UpdateModule;
      params = [req.module_name, req.updated_by, req.id];
    }

    const result = await this.pg.query(query, params);

    return (result.rowCount ?? 0) > 0;
  }

  // --- GET MODULES (pagination) ---
  async getModules(
    req: PaginationRequest
  ): Promise<{ modules: Module[]; totalCount: number }> {
    let query = ModuleQueries.GetModules;
    const params: any[] = [];

    if (!req.showAll) {
      if (req.search) {
        query += ` AND m.module_name ILIKE '%' || $1 || '%'`;
        params.push(req.search);
      }

      const offset = (req.page - 1) * req.limit;
      params.push(req.limit, offset);
      query += ` ORDER BY m.created_at DESC LIMIT $${
        params.length - 1
      } OFFSET $${params.length}`;
    }

    const result = await this.pg.query<{
      total_count: string;
      id: number;
      uuid: string;
      group_id: number | null;
      module_name: string;
      thumbnail_url: string | null;
      tag: string | null;
      created_by: string | null;
      created_at: string;
      updated_at: string | null;
      updated_by: string | null;
    }>(query, params);

    const modules: Module[] = result.rows.map((row) => ({
      id: row.id,
      uuid: row.uuid,
      group_id: row.group_id,
      module_name: row.module_name,
      thumbnail_url: row.thumbnail_url,
      tag: row.tag,
      created_by: row.created_by,
      created_at: row.created_at,
      updated_at: row.updated_at,
      updated_by: row.updated_by,
    }));

    const totalCount =
      result.rows.length > 0 && result.rows[0].total_count
        ? parseInt(result.rows[0].total_count, 10)
        : 0;

    return { modules, totalCount };
  }

  // --- GET MODULES BY GROUP ID ---
  async getModulesByGroupID(groupId: number): Promise<Module[]> {
    const result = await this.pg.query<{
      id: number;
      uuid: string;
      group_id: number | null;
      module_name: string;
      thumbnail_url: string | null;
      created_by: string | null;
      created_at: string;
      updated_at: string | null;
      updated_by: string | null;
    }>(ModuleQueries.GetModulesByGroupID, [groupId]);

    return result.rows.map((row) => ({
      id: row.id,
      uuid: row.uuid,
      group_id: row.group_id,
      module_name: row.module_name,
      thumbnail_url: row.thumbnail_url,
      created_by: row.created_by,
      created_at: row.created_at,
      updated_at: row.updated_at,
      updated_by: row.updated_by,
    }));
  }

  // --- SOFT DELETE MODULE ---
  async softDeleteModule(moduleId: number, operator: string): Promise<boolean> {
    const result = await this.pg.query(ModuleQueries.SoftDeleteModule, [
      operator,
      operator,
      moduleId,
    ]);

    return (result.rowCount ?? 0) > 0;
  }

  // --- REMOVE GROUP ID FROM MODULES ---
  async removeGroupIDFromModules(
    groupId: number,
    operator: string
  ): Promise<boolean> {
    const result = await this.pg.query(ModuleQueries.RemoveGroupIDFromModules, [
      operator,
      groupId,
    ]);

    return (result.rowCount ?? 0) > 0;
  }

  // --- GET PERCENTAGE MARK WATCHED ---
  async getPercentageMarkWatchedModulesUser(
    userId: number,
    moduleId: number
  ): Promise<number> {
    const result = await this.pg.query<{ percentage: number }>(
      ModuleQueries.GetPercentageModulesUser,
      [userId, moduleId, userId, moduleId]
    );

    if (result.rowCount === 0) {
      return 0;
    }

    return Math.ceil(result.rows[0].percentage);
  }
}
