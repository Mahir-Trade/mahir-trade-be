import { Request, Response } from "express";
import { SubModuleService } from "../services/subModule.service";
import {
  SubModuleRequest,
  MarkSubModuleAsWatchedRequest,
} from "@dto/subModule.dto";
import { FileUpload } from "../DTO/google.dto";
import { DefaultResponse } from "../utils/response";

export class SubModuleController {
  private subModuleService: SubModuleService;

  constructor() {
    this.subModuleService = new SubModuleService();

    this.createSubModule = this.createSubModule.bind(this);
    this.getSubModules = this.getSubModules.bind(this);
    this.getSubModuleByID = this.getSubModuleByID.bind(this);
    this.getSubModulesByModuleID = this.getSubModulesByModuleID.bind(this);
    this.updateSubModule = this.updateSubModule.bind(this);
    this.softDeleteSubModule = this.softDeleteSubModule.bind(this);
    this.markSubModuleAsWatched = this.markSubModuleAsWatched.bind(this);
    this.uploadFile = this.uploadFile.bind(this);
  }

  // --- POST /sub-modules ---
  async createSubModule(req: Request, res: Response): Promise<Response> {
    try {
      const body: SubModuleRequest = req.body;
      const result = await this.subModuleService.createSubModule(body);
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[SubModuleController][createSubModule]", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // --- GET /sub-modules ---
  async getSubModules(req: Request, res: Response): Promise<Response> {
    try {
      const { limit = 10, page = 1 } = req.query;
      const pagination = { limit: Number(limit), page: Number(page) };

      const result = await this.subModuleService.getSubModules(pagination);
      return res.status(200).json(result);
    } catch (err: any) {
      console.error("[SubModuleController][getSubModules]", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // --- GET /sub-modules/:sub_module_id ---
  async getSubModuleByID(req: Request, res: Response): Promise<Response> {
    try {
      const id = Number(req.params.sub_module_id);
      if (isNaN(id))
        return res.status(400).json({ code: 400, message: "Invalid ID" });

      const result = await this.subModuleService.getSubModuleByID(id);
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[SubModuleController][getSubModuleByID]", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // GET SUBMODULES BY MODULE ID
  // ======================================================
  async getSubModulesByModuleID(
    req: Request,
    res: Response
  ): Promise<Response> {
    try {
      const idParam = req.params.module_id;

      if (!idParam) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request body",
          error: "module id is required",
        } as DefaultResponse);
      }

      const module_id = parseInt(idParam, 10);
      if (isNaN(module_id) || module_id <= 0) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request body",
          error: "module id is required",
        } as DefaultResponse);
      }

      const limit = parseInt(req.query.limit as string, 10) || 10;
      const page = parseInt(req.query.page as string, 10) || 1;

      const paginationReq = {
        limit,
        page,
      };

      const result = await this.subModuleService.getSubModulesByModuleID(
        module_id,
        paginationReq
      );

      // Jika service return error object, kita tangani
      if ((result as any).error) {
        return res.status(500).json(result);
      }

      return res.status(200).json(result);
    } catch (err: any) {
      return res.status(500).json({
        code: 500,
        message: "Internal server error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- PUT /sub-modules/:sub_module_id ---
  async updateSubModule(req: Request, res: Response): Promise<Response> {
    try {
      const id = Number(req.params.sub_module_id);
      if (isNaN(id))
        return res.status(400).json({ code: 400, message: "Invalid ID" });

      const body: SubModuleRequest = req.body;
      const result = await this.subModuleService.updateSubModule(id, body);
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[SubModuleController][updateSubModule]", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // --- DELETE /sub-modules/:sub_module_id ---
  async softDeleteSubModule(req: Request, res: Response): Promise<Response> {
    try {
      const id = Number(req.params.sub_module_id);
      if (isNaN(id))
        return res.status(400).json({ code: 400, message: "Invalid ID" });

      const result = await this.subModuleService.softDeleteSubModule(id);
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[SubModuleController][softDeleteSubModule]", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // --- POST /sub-modules/mark-watched ---
  async markSubModuleAsWatched(req: Request, res: Response): Promise<Response> {
    try {
      const body: MarkSubModuleAsWatchedRequest = req.body;
      const result = await this.subModuleService.markSubModuleAsWatched(body);
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[SubModuleController][markSubModuleAsWatched]", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // --- POST /sub-modules/upload ---
  async uploadFile(req: Request, res: Response): Promise<Response> {
    try {
      if (!req.file)
        return res.status(400).json({ code: 400, message: "file is required" });

      const fileUpload: FileUpload = {
        file_name: req.file.originalname,
        size: req.file.size,
        localFilePath: "",
        file_content_type: "",
        bucket_name: "",
      };

      const result = await this.subModuleService.uploadFile(
        fileUpload,
        req.file
      );
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[SubModuleController][uploadFile]", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }
}
