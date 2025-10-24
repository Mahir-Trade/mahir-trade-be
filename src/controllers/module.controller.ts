import { ModuleRequest } from "../DTO/module.dto";
import { ModuleService } from "../services/module.service";
import { DefaultResponse, PaginationRequest } from "../utils/response";
import { validate } from "class-validator";
import { Request, Response } from "express";

export class ModuleController {
  private moduleService: ModuleService;

  constructor() {
    this.moduleService = new ModuleService();

    this.createModule = this.createModule.bind(this);
    this.getModuleByID = this.getModuleByID.bind(this);
    this.getModules = this.getModules.bind(this);
    this.getModulesByGroupID = this.getModulesByGroupID.bind(this);
    this.updateModule = this.updateModule.bind(this);
    this.deleteModule = this.deleteModule.bind(this);
    this.getPercentageMarkWatchedModulesUser =
      this.getPercentageMarkWatchedModulesUser.bind(this);
  }

  // ======================================================
  // CREATE MODULE
  // ======================================================
  async createModule(req: Request, res: Response): Promise<Response> {
    try {
      const moduleReq = req.body as ModuleRequest;
      // const errors = await validate(moduleReq);
      // if (errors.length > 0) {
      //   return res.status(400).json({
      //     code: 400,
      //     message: "Invalid request body",
      //     error: errors,
      //   } as DefaultResponse);
      // }

      const resp = await this.moduleService.createModule(moduleReq);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[ModuleController][createModule] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // ======================================================
  // GET MODULE BY ID
  // ======================================================
  async getModuleByID(req: Request, res: Response): Promise<Response> {
    try {
      const idParam = req.params.module_id;
      if (!idParam) {
        return res.status(400).json({
          code: 400,
          message: "module id is required",
        });
      }

      const moduleID = parseInt(idParam, 10);
      if (isNaN(moduleID) || moduleID <= 0) {
        return res.status(400).json({
          code: 400,
          message: "invalid module id",
        });
      }

      const resp = await this.moduleService.getModuleByID(moduleID);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[ModuleController][getModuleByID] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // ======================================================
  // GET MODULES (pagination)
  // ======================================================
  async getModules(req: Request, res: Response): Promise<Response> {
    try {
      const paginationReq: PaginationRequest = {
        page: req.query.page ? parseInt(req.query.page as string, 10) : 1,
        limit: req.query.limit ? parseInt(req.query.limit as string, 10) : 10,
        search: (req.query.search as string) || undefined,
      };

      const resp = await this.moduleService.getModules(paginationReq);
      return res.status(200).json(resp);
    } catch (err: any) {
      console.error("[ModuleController][getModules] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // ======================================================
  // GET MODULES BY GROUP ID
  // ======================================================
  async getModulesByGroupID(req: Request, res: Response): Promise<Response> {
    try {
      const idParam = req.params.group_id;
      if (!idParam) {
        return res.status(400).json({
          code: 400,
          message: "group id is required",
        });
      }

      const groupID = parseInt(idParam, 10);
      if (isNaN(groupID) || groupID <= 0) {
        return res.status(400).json({
          code: 400,
          message: "invalid group id",
        });
      }

      const resp = await this.moduleService.getModuleByGroupID(groupID);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[ModuleController][getModulesByGroupID] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // ======================================================
  // UPDATE MODULE
  // ======================================================
  async updateModule(req: Request, res: Response): Promise<Response> {
    try {
      const idParam = req.params.module_id;
      if (!idParam) {
        return res.status(400).json({
          code: 400,
          message: "module id is required",
        });
      }

      const moduleID = parseInt(idParam, 10);
      if (isNaN(moduleID) || moduleID <= 0) {
        return res.status(400).json({
          code: 400,
          message: "invalid module id",
        });
      }

      const moduleReq = req.body as ModuleRequest;
      // const errors = await validate(moduleReq);
      // if (errors.length > 0) {
      //   return res.status(400).json({
      //     code: 400,
      //     message: "Invalid request body",
      //     error: errors,
      //   });
      // }

      const resp = await this.moduleService.updateModule(moduleID, moduleReq);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[ModuleController][updateModule] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // ======================================================
  // DELETE MODULE
  // ======================================================
  async deleteModule(req: Request, res: Response): Promise<Response> {
    try {
      const idParam = req.params.module_id;
      if (!idParam) {
        return res.status(400).json({
          code: 400,
          message: "module id is required",
        });
      }

      const moduleID = parseInt(idParam, 10);
      if (isNaN(moduleID) || moduleID <= 0) {
        return res.status(400).json({
          code: 400,
          message: "invalid module id",
        });
      }

      const username = req.body.username || "system";
      const resp = await this.moduleService.deleteModule(moduleID, username);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[ModuleController][deleteModule] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // ======================================================
  // GET WATCHED PERCENTAGE
  // ======================================================
  async getPercentageMarkWatchedModulesUser(
    req: Request,
    res: Response
  ): Promise<Response> {
    try {
      const moduleIdParam = req.params.module_id;
      const userIdParam = req.query.user_id as string;

      if (!moduleIdParam || !userIdParam) {
        return res.status(400).json({
          code: 400,
          message: "module_id and user_id are required",
        });
      }

      const moduleId = parseInt(moduleIdParam, 10);
      const userId = parseInt(userIdParam, 10);

      if (isNaN(moduleId) || isNaN(userId)) {
        return res.status(400).json({
          code: 400,
          message: "invalid parameters",
        });
      }

      const resp = await this.moduleService.getPercentageMarkWatchedModulesUser(
        userId,
        moduleId
      );
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error(
        "[ModuleController][getPercentageMarkWatchedModulesUser] error:",
        err
      );
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }
}
