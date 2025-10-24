import { ReportRepository } from "../repo/report_repo";
import { BucketRepository } from "../repo/google/bucket";
import { GoogleCfg } from "../repo/google/config";
import { Logger } from "../utils/logger";
import {
  DefaultResponse,
  DefaultPaginationResponseData,
  PaginationRequest,
} from "@utils/response";
import { Report } from "../models/report.model";

export class ReportService {
  private reportRepo = new ReportRepository();
  private bucketRepo = new BucketRepository();
  private googleCfg = new GoogleCfg();

  // --- CREATE ---
  async createReport(report: Report): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 201, message: "Success" };
    try {
      const id = await this.reportRepo.createReport(report);
      resp.data = { id };
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[ReportService][createReport] ${err.message}`);
    }
    return resp;
  }

  // --- GET ALL ---
  async getReports(
    req: PaginationRequest
  ): Promise<DefaultPaginationResponseData<Report>> {
    const resp: DefaultPaginationResponseData<Report> = {
      results: [],
      meta: {
        page: req.page,
        limit: req.limit,
        totalItems: 0,
        totalPages: 0,
        hasNext: false,
        hasPrevious: false,
      },
    };

    try {
      const { reports, totalCount } = await this.reportRepo.getReports(req);

      // Buat signed URL buat file & thumbnail
      for (const report of reports) {
        for (const key of [
          "report_thumbnail_url",
          "report_file_url",
        ] as const) {
          const value = report[key];
          if (!value) continue;
          try {
            const parsed = this.bucketRepo.urlParser(value);
            report[key] = await this.bucketRepo.presignedURL(
              parsed.bucketName,
              value
            );
          } catch (err: any) {
            Logger.error(`[ReportService][getReports] ${key}: ${err.message}`);
          }
        }
      }

      resp.results = reports;
      resp.meta.totalItems = totalCount;
      resp.meta.totalPages = Math.ceil(totalCount / req.limit);
      resp.meta.hasNext = req.page < resp.meta.totalPages;
      resp.meta.hasPrevious = req.page > 1;
    } catch (err: any) {
      Logger.error(`[ReportService][getReports] ${err.message}`);
      // Optional: bisa bungkus error di results
      resp.results = [];
    }

    return resp;
  }

  // --- GET BY ID ---
  async getReportByID(id: number): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "Success" };
    try {
      const report = await this.reportRepo.getReportByID(id);
      if (!report) {
        resp.code = 404;
        resp.message = "Report not found";
        return resp;
      }

      for (const key of ["report_thumbnail_url", "report_file_url"] as const) {
        const value = report[key];
        if (!value) continue;
        try {
          const parsed = this.bucketRepo.urlParser(value);
          report[key] = await this.bucketRepo.presignedURL(
            parsed.bucketName,
            value
          );
        } catch (err: any) {
          Logger.error(`[ReportService][getReportByID] ${key}: ${err.message}`);
        }
      }

      resp.data = report;
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[ReportService][getReportByID] ${err.message}`);
    }
    return resp;
  }

  // --- UPDATE ---
  async updateReport(report: Report): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "Success" };
    try {
      const found = await this.reportRepo.getReportByID(report.id);
      if (!found) {
        resp.code = 404;
        resp.message = "Report not found";
        return resp;
      }

      await this.reportRepo.updateReport(report);
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[ReportService][updateReport] ${err.message}`);
    }
    return resp;
  }

  // --- DELETE ---
  async deleteReport(id: number, deletedBy: string): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "Success" };
    try {
      const found = await this.reportRepo.getReportByID(id);
      if (!found) {
        resp.code = 404;
        resp.message = "Report not found";
        return resp;
      }

      await this.reportRepo.softDeleteReport(id, deletedBy);
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[ReportService][deleteReport] ${err.message}`);
    }
    return resp;
  }
}
