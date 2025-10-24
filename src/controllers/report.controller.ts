import { Request, Response } from "express";
import { ReportService } from "../services/report.service";
import { validate } from "class-validator";
import { DefaultResponse, PaginationRequest } from "../utils/response";
import { Report } from "../models/report.model"; // ✅ FIX: Import yang benar

export class ReportController {
  private reportService: ReportService;

  constructor() {
    this.reportService = new ReportService();

    this.createReport = this.createReport.bind(this);
    this.getReports = this.getReports.bind(this);
    this.getReportByID = this.getReportByID.bind(this);
    this.updateReport = this.updateReport.bind(this);
    this.deleteReport = this.deleteReport.bind(this);
  }

  // --- CREATE REPORT ---
  async createReport(req: Request, res: Response): Promise<Response> {
    try {
      const report = req.body as Report;

      const errors = await validate(report);
      if (errors.length > 0) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request body",
          error: errors,
        } as DefaultResponse);
      }

      const resp = await this.reportService.createReport(report);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[ReportController][createReport] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal server error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- GET REPORTS ---
  async getReports(req: Request, res: Response): Promise<Response> {
    try {
      const pagination: PaginationRequest = {
        limit: parseInt(req.query.limit as string) || 10,
        page: parseInt(req.query.page as string) || 1,
        search: (req.query.search as string) || "",
      };

      const resp = await this.reportService.getReports(pagination);

      // ✅ FIX: `resp.results` adalah array, bukan object dengan `code`
      return res.status(200).json(resp);
    } catch (err: any) {
      console.error("[ReportController][getReports] error:", err);
      return res.status(400).json({
        code: 400,
        message: "Bad request",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- GET REPORT BY ID ---
  async getReportByID(req: Request, res: Response): Promise<Response> {
    try {
      const id = parseInt(req.params.id);
      if (!id || isNaN(id)) {
        return res.status(400).json({
          code: 400,
          message: "Invalid report id",
          error: "Invalid report id",
        } as DefaultResponse);
      }

      const resp = await this.reportService.getReportByID(id);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[ReportController][getReportByID] error:", err);
      return res.status(400).json({
        code: 400,
        message: "Bad request",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- UPDATE REPORT ---
  async updateReport(req: Request, res: Response): Promise<Response> {
    try {
      const id = parseInt(req.params.id);
      if (!id || isNaN(id)) {
        return res.status(400).json({
          code: 400,
          message: "Invalid report id",
          error: "Invalid report id",
        } as DefaultResponse);
      }

      const report = req.body as Report;
      report.id = id; // ✅ Pastikan di model `id?: number` biar valid

      const errors = await validate(report);
      if (errors.length > 0) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request body",
          error: errors,
        } as DefaultResponse);
      }

      const resp = await this.reportService.updateReport(report);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[ReportController][updateReport] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal server error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- DELETE REPORT ---
  async deleteReport(req: Request, res: Response): Promise<Response> {
    try {
      const id = parseInt(req.params.id);
      if (!id || isNaN(id)) {
        return res.status(400).json({
          code: 400,
          message: "Invalid report id",
          error: "Invalid report id",
        } as DefaultResponse);
      }

      const { deleted_by } = req.body;
      const resp = await this.reportService.deleteReport(id, deleted_by);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[ReportController][deleteReport] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal server error",
        error: err.message,
      } as DefaultResponse);
    }
  }
}
