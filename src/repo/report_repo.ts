import { Pool } from "pg";
import { Report } from "../models/report.model";

import db from "../db/db.config";
import { ReportQueries } from "./queries/report_queries";
import { PaginationRequest } from "@utils/response";

export class ReportRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- CREATE REPORT ---
  async createReport(req: Report): Promise<number> {
    try {
      const result = await this.pg.query(ReportQueries.CreateReport, [
        req.report_name,
        req.report_thumbnail_url,
        req.report_file_url,
        req.created_by,
      ]);

      if (result.rowCount === 0) {
        throw new Error("failed to create report");
      }

      return result.rows[0].id as number;
    } catch (err: any) {
      console.error(`[repo][report][createReport] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET REPORTS (with pagination + search) ---
  async getReports(
    req: PaginationRequest
  ): Promise<{ reports: Report[]; totalCount: number }> {
    try {
      const limit = req.limit || 10;
      const offset = req.page > 1 ? (req.page - 1) * limit : 0;

      let query = ReportQueries.GetReports;
      const params: any[] = [];

      if (req.search) {
        query += ` AND report_name ILIKE '%' || $1 || '%'`;
        params.push(req.search);
      }

      query += ` LIMIT $${params.length + 1} OFFSET $${params.length + 2}`;
      params.push(limit, offset);

      const result = await this.pg.query(query, params);

      const reports: Report[] = (result.rows as any[]).map((row) => ({
        id: row.id,
        report_name: row.report_name,
        report_thumbnail_url: row.report_thumbnail_url,
        report_file_url: row.report_file_url,
        created_by: row.created_by,
        updated_by: row.updated_by,
        created_at: row.created_at,
        updated_at: row.updated_at,
      }));

      const totalCount =
        result.rows.length > 0 ? parseInt(result.rows[0].total_count, 10) : 0;

      return { reports, totalCount };
    } catch (err: any) {
      console.error(`[repo][report][getReports] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET REPORT BY ID ---
  async getReportByID(id: number): Promise<Report> {
    try {
      const result = await this.pg.query(ReportQueries.GetReportByID, [id]);

      if (result.rowCount === 0) {
        throw new Error(`report with id ${id} not found`);
      }

      const row = result.rows[0];
      return {
        id: row.id,
        report_name: row.report_name,
        report_thumbnail_url: row.report_thumbnail_url,
        report_file_url: row.report_file_url,
        created_by: row.created_by,
        updated_at: row.updated_by,
        created_at: row.created_at,
        updated_by: row.updated_at,
      };
    } catch (err: any) {
      console.error(`[repo][report][getReportByID] error: ${err.message}`);
      throw err;
    }
  }

  // --- UPDATE REPORT ---
  async updateReport(req: Report): Promise<boolean> {
    try {
      const result = await this.pg.query(ReportQueries.UpdateReport, [
        req.report_name,
        req.report_thumbnail_url,
        req.report_file_url,
        req.updated_by,
        req.id,
      ]);

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[repo][report][updateReport] error: ${err.message}`);
      throw err;
    }
  }

  // --- SOFT DELETE REPORT ---
  async softDeleteReport(id: number, deletedBy: string): Promise<boolean> {
    try {
      const result = await this.pg.query(ReportQueries.SoftDeleteReport, [
        deletedBy,
        deletedBy,
        id,
      ]);

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[repo][report][softDeleteReport] error: ${err.message}`);
      throw err;
    }
  }
}
