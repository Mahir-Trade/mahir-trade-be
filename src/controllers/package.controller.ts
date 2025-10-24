import { Request, Response } from "express";

import { Package } from "../models/package.model";
import { validate } from "class-validator";
import { DefaultResponse, PaginationRequest } from "../utils/response";
import { PackageService } from "../services/package.service";

export class PackageController {
  private packageService: PackageService;

  constructor() {
    this.packageService = new PackageService();

    this.createPackage = this.createPackage.bind(this);
    this.getPackages = this.getPackages.bind(this);
    this.getPackageByID = this.getPackageByID.bind(this);
    this.updatePackage = this.updatePackage.bind(this);
    this.deletePackage = this.deletePackage.bind(this);
  }

  // --- CREATE PACKAGE ---
  async createPackage(req: Request, res: Response): Promise<Response> {
    try {
      const pkg = req.body as Package;
      // const errors = await validate(pkg);
      // if (errors) {
      //   return res.status(400).json({
      //     code: 400,
      //     message: "Invalid request body",
      //     error: errors,
      //   } as DefaultResponse);
      // }

      const resp = await this.packageService.createPackage(pkg);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[PackageController][createPackage] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal server error",
        error: err.message,
      });
    }
  }

  // --- GET PACKAGES ---
  async getPackages(req: Request, res: Response): Promise<Response> {
    try {
      const limit = parseInt(req.query.limit as string, 10) || 10;
      const page = parseInt(req.query.page as string, 10) || 1;

      const paginationReq: PaginationRequest = { limit, page };
      const resp = await this.packageService.getPackages(paginationReq);
      return res.status(200).json(resp);
    } catch (err: any) {
      console.error("[PackageController][getPackages] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal server error",
        error: err.message,
      });
    }
  }

  // --- GET PACKAGE BY ID ---
  async getPackageByID(req: Request, res: Response): Promise<Response> {
    try {
      const idParam = req.params.id;
      if (!idParam) {
        return res.status(400).json({
          code: 400,
          message: "package id is required",
        });
      }

      const packageID = parseInt(idParam, 10);
      if (isNaN(packageID) || packageID <= 0) {
        return res.status(400).json({
          code: 400,
          message: "invalid package id",
        });
      }

      const resp = await this.packageService.getPackageByID(packageID);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[PackageController][getPackageByID] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal server error",
        error: err.message,
      });
    }
  }

  // --- UPDATE PACKAGE ---
  async updatePackage(req: Request, res: Response): Promise<Response> {
    try {
      const idParam = req.params.id;
      if (!idParam) {
        return res.status(400).json({
          code: 400,
          message: "package id is required",
        });
      }

      const packageID = parseInt(idParam, 10);
      if (isNaN(packageID) || packageID <= 0) {
        return res.status(400).json({
          code: 400,
          message: "invalid package id",
        });
      }

      const pkg = req.body as Package;
      pkg.id = packageID;

      const errors = await validate(pkg);
      if (errors) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request body",
          error: errors,
        });
      }

      const resp = await this.packageService.updatePackage(pkg);
      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[PackageController][updatePackage] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal server error",
        error: err.message,
      });
    }
  }

  // --- DELETE PACKAGE ---
  async deletePackage(req: Request, res: Response): Promise<Response> {
    console.log("[DEBUG] Incoming DELETE /packages/:id");
    console.log("Params:", req.params);

    try {
      const idParam = req.params.id;
      if (!idParam) {
        return res.status(400).json({
          code: 400,
          message: "package id is required",
        });
      }

      const packageID = parseInt(idParam, 10);
      if (isNaN(packageID) || packageID <= 0) {
        return res.status(400).json({
          code: 400,
          message: "invalid package id",
        });
      }

      // Ambil user dari context atau dari req.userData
      const userData = (req as any).userData;
      const deletedBy = userData?.email || userData?.username || "unknown";

      console.log("[DEBUG] User context:", userData);
      console.log("[DEBUG] Deleting package:", { packageID, deletedBy });

      const resp = await this.packageService.deletePackage(
        packageID,
        deletedBy
      );
      console.log("[DEBUG] Service Response:", resp);

      return res.status(resp.code).json(resp);
    } catch (err: any) {
      console.error("[ERROR] [PackageController][deletePackage] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal server error",
        error: err.message,
        stack: err.stack,
      });
    }
  }
}
