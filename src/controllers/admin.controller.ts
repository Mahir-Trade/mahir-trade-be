import { Request, Response } from "express";

import { UserService } from "../services/user.service";

import {
  AdminLoginRequest,
  UpdateTypeUserRequest,
} from "../models/admin.model";
import { DefaultResponse, PaginationRequest } from "../utils/response";
import { AdminService } from "../services/admin.service";
import { UserContext } from "../config/userContext";
import { validateRequest } from "@middlewares/validateRequest";
import { StartMembershipProgramRequest } from "@dto/startMembership";
import { ToggleUserMembershipRequest } from "@dto/toggleUserMembershipRequest";
import { sendResponse } from "../utils/response";

export class AdminController {
  private adminService: AdminService;
  private userService: UserService;

  constructor() {
    this.adminService = new AdminService();
    this.userService = new UserService();

    this.adminLogin = this.adminLogin.bind(this);
    this.adminRegistration = this.adminRegistration.bind(this);
    this.updateTypeUser = this.updateTypeUser.bind(this);
    this.getDetailUserForBO = this.getDetailUserForBO.bind(this);
    this.getDetailAdminInfo = this.getDetailAdminInfo.bind(this);
    this.getAllUsers = this.getAllUsers.bind(this);

    this.toggleInactiveUserMembership =
      this.toggleInactiveUserMembership.bind(this);
    this.startMembershipProgram = this.startMembershipProgram.bind(this);
    this.getMembershipProgramDate = this.getMembershipProgramDate.bind(this);
    this.updateMembershipProgramDate =
      this.updateMembershipProgramDate.bind(this);
  }

  // --- ADMIN LOGIN ---
  async adminLogin(req: Request, res: Response): Promise<Response> {
    try {
      const body: AdminLoginRequest = req.body;

      // Validation (Go pakai validator)
      if (!body.identity || !body.password) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request body",
          error: "username and password are required",
        } as DefaultResponse);
      }

      const result = await this.adminService.adminLogin(body);
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[controller][AdminLogin] error:", err.message);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- ADMIN REGISTRATION ---
  async adminRegistration(req: Request, res: Response): Promise<Response> {
    try {
      const body = req.body;

      if (!body.email || !body.username || !body.password) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request body",
          error: "email, username and password are required",
        } as DefaultResponse);
      }

      const result = await this.adminService.adminRegistration(body);
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[controller][AdminRegistration] error:", err.message);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- UPDATE TYPE USER ---
  async updateTypeUser(req: Request, res: Response): Promise<Response> {
    try {
      const userId = req.params.user_id;
      if (!userId) {
        return res.status(400).json({
          code: 400,
          message: "user id is required",
        } as DefaultResponse);
      }

      const body: UpdateTypeUserRequest = req.body;
      if (body.isActive === undefined) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request body",
          error: "isActive is required",
        } as DefaultResponse);
      }

      const result = await this.adminService.updateTypeUser(
        parseInt(userId, 10),
        body.isActive,
        (req as any).userData
      );

      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[controller][UpdateTypeUser] error:", err.message);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- GET DETAIL USER FOR BO ---
  async getDetailUserForBO(req: Request, res: Response): Promise<Response> {
    try {
      const userId = parseInt(req.params.user_id, 10);
      if (isNaN(userId)) {
        return res.status(400).json({
          code: 400,
          message: "bad request",
          error: "user_id must be a number",
        } as DefaultResponse);
      }

      const result = await this.userService.getDetailUserForBO(userId);
      return res.status(200).json(result);
    } catch (err: any) {
      console.error("[controller][GetDetailUserForBO] error:", err.message);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- GET DETAIL ADMIN INFO ---
  async getDetailAdminInfo(req: Request, res: Response): Promise<Response> {
    try {
      // Ambil data user dari middleware auth
      const userData = (req as any).userData;
      console.log("🔑 [DEBUG] Token decoded userData:", userData);

      if (!userData || !userData.username) {
        console.error(
          "❌ [ERROR] userData tidak ditemukan di request (kemungkinan token belum diparse)"
        );
        return res.status(401).json({
          code: 401,
          message: "Unauthorized: invalid or missing token",
        } as DefaultResponse);
      }

      // Panggil service untuk ambil data admin
      const result = await this.adminService.getDetailAdminInfo(
        userData.username
      );
      console.log("📦 [DEBUG] Hasil getDetailAdminInfo:", result);

      return res.status(result.code ?? 200).json(result);
    } catch (err: any) {
      console.error("🔥 [controller][getDetailAdminInfo] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err?.message || "Unknown error",
      } as DefaultResponse);
    }
  }

  // --- GET ALL USERS ---
  async getAllUsers(req: Request, res: Response): Promise<Response> {
    try {
      const limit = parseInt(req.query.limit as string, 10) || 10;
      const page = parseInt(req.query.page as string, 10) || 1;
      const search = (req.query.search as string) || "";
      const membershipStatus = (req.query.membershipStatus as string) || "";
      const sortBy = (req.query.sortBy as string) || "DESC";

      // Ambil user context
      const userCtx = UserContext.get();
      console.log("📦 [Controller] Dijalankan oleh:", userCtx);

      const pagination: PaginationRequest = {
        limit,
        page,
        search,
        membership_status: membershipStatus,
        sort_by: sortBy,
      };

      const result = await this.adminService.getAllUsers(pagination);
      return sendResponse(res, 200, result);
    } catch (err: any) {
      console.error("[controller][GetAllUsers] error:", err.message);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  async toggleInactiveUserMembership(
    req: Request,
    res: Response
  ): Promise<Response> {
    try {
      console.log("📦 [Controller] ToggleInactiveUserMembership executed");

      // const errors = validationResult(req);
      // if (!errors.isEmpty()) {
      //   return res.status(400).json({
      //     code: 400,
      //     message: "Invalid request body",
      //     error: errors.array(),
      //   } as DefaultResponse);
      // }

      const body = req.body as ToggleUserMembershipRequest;
      const result = await this.adminService.toggleInactiveUserMembership(body);

      return res.status(200).json(result);
    } catch (err: any) {
      console.error(
        "[controller][ToggleInactiveUserMembership] error:",
        err.message
      );
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  /**
   * 🚀 Start Membership Program
   */
  async startMembershipProgram(req: Request, res: Response): Promise<Response> {
    try {
      console.log("📦 [Controller] StartMembershipProgram executed");

      // const errors = validationResult(req);
      // if (!errors.isEmpty()) {
      //   return res.status(400).json({
      //     code: 400,
      //     message: "Invalid request body",
      //     error: errors.array(),
      //   } as DefaultResponse);
      // }

      const body = req.body as StartMembershipProgramRequest;
      const result = await this.adminService.startMembershipProgram(body);

      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[controller][StartMembershipProgram] error:", err.message);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  /**
   * 📅 Get Membership Program Date
   */
  async getMembershipProgramDate(
    req: Request,
    res: Response
  ): Promise<Response> {
    try {
      console.log("📦 [Controller] GetMembershipProgramDate executed");

      const result = await this.adminService.getMembershipProgramDate();
      return res.status(200).json(result);
    } catch (err: any) {
      console.error(
        "[controller][GetMembershipProgramDate] error:",
        err.message
      );
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  /**
   * ✏️ Update Membership Program Date
   */
  async updateMembershipProgramDate(
    req: Request,
    res: Response
  ): Promise<Response> {
    try {
      console.log("📦 [Controller] UpdateMembershipProgramDate executed");

      // const errors = validateRequest(req);
      // if (!errors.isEmpty()) {
      //   return res.status(400).json({
      //     code: 400,
      //     message: "Invalid request body",
      //     error: errors.array(),
      //   } as DefaultResponse);
      // }

      const body = req.body as StartMembershipProgramRequest;
      const result = await this.adminService.updateMembershipProgramDate(body);

      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error(
        "[controller][UpdateMembershipProgramDate] error:",
        err.message
      );
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      } as DefaultResponse);
    }
  }
}
