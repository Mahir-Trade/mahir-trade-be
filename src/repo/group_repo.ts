import { Pool } from "pg";
import db from "../db/db.config";
import { Group, GetGroupsRequest } from "../models/group.model";
import { GroupQueries } from "./queries/group_queries";

export class GroupRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- CREATE GROUP ---
  async createGroup(req: Group): Promise<number> {
    try {
      const result = await this.pg.query(GroupQueries.CreateGroup, [
        req.group_name,
      ]);

      if (result.rowCount === 0) {
        throw new Error("failed to create group");
      }

      return result.rows[0].id as number;
    } catch (err: any) {
      console.error(`[repo][group][createGroup] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET GROUP BY ID ---
  async getGroupByID(id: number): Promise<Group> {
    try {
      const result = await this.pg.query(GroupQueries.GetGroupByID, [id]);

      if (result.rowCount === 0) {
        throw new Error(`group with id ${id} not found`);
      }

      return result.rows[0] as Group;
    } catch (err: any) {
      console.error(`[repo][group][getGroupByID] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET GROUPS WITH PAGINATION ---
  async getGroups(
    req: GetGroupsRequest
  ): Promise<{ groups: Group[]; totalCount: number }> {
    try {
      let query = GroupQueries.GetGroups;
      const args: any[] = [];

      if (!req.showAll) {
        if (req.limit === 0) req.limit = 10;

        let offset = 0;
        if (req.page > 1) {
          offset = (req.page - 1) * req.limit;
        }

        args.push(req.limit, offset);
        query += ` LIMIT $${args.length - 1} OFFSET $${args.length}`;
      }

      const rows = await this.pg.query(query, args);

      const groups: Group[] = rows.rows.map((row) => ({
        ...row,
      }));

      // total count
      const totalRes = await this.pg.query(GroupQueries.GetTotalGroups);
      const totalCount = totalRes.rows[0]?.count
        ? parseInt(totalRes.rows[0].count, 10)
        : 0;

      return { groups, totalCount };
    } catch (err: any) {
      console.error(`[repo][group][getGroups] error: ${err.message}`);
      throw err;
    }
  }

  // --- UPDATE GROUP ---
  async updateGroup(req: Group): Promise<boolean> {
    try {
      const result = await this.pg.query(GroupQueries.UpdateGroup, [
        req.group_name,
        req.id,
      ]);

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[repo][group][updateGroup] error: ${err.message}`);
      throw err;
    }
  }

  // --- SOFT DELETE GROUP ---
  async softDeleteGroup(groupId: number, operator: string): Promise<boolean> {
    try {
      const result = await this.pg.query(GroupQueries.SoftDeleteGroup, [
        operator,
        operator,
        groupId,
      ]);

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[repo][group][softDeleteGroup] error: ${err.message}`);
      throw err;
    }
  }
}
