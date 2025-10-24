import { Request, Response } from "express";

import { Group, GetGroupsRequest } from "../models/group.model";
import { DefaultResponse } from "../utils/response";
import { validate } from "class-validator";
import { GroupService } from "../services/group.service";

export class GroupController {
  private groupService: GroupService;

  constructor() {
    this.groupService = new GroupService();

    this.createGroup = this.createGroup.bind(this);
    this.getGroupByID = this.getGroupByID.bind(this);
    this.getGroups = this.getGroups.bind(this);
    this.updateGroup = this.updateGroup.bind(this);
    this.deleteGroup = this.deleteGroup.bind(this);
  }

  // --- POST /groups ---
  async createGroup(req: Request, res: Response): Promise<Response> {
    try {
      const data = req.body as Group;
      // const errors = await validate(data);

      // if (errors.length > 0) {
      //   return res.status(400).json({
      //     code: 400,
      //     message: "Invalid request body",
      //     error: errors,
      //   } as DefaultResponse);
      // }

      const result = await this.groupService.createGroup(data);
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[GroupController][createGroup] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // --- GET /groups/:id ---
  async getGroupByID(req: Request, res: Response): Promise<Response> {
    try {
      const id = parseInt(req.params.id, 10);
      if (isNaN(id) || id <= 0) {
        return res.status(400).json({
          code: 400,
          message: "group id is required or invalid",
        });
      }

      const result = await this.groupService.getGroupById(id);
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[GroupController][getGroupByID] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Failed to get group",
        error: err.message,
      });
    }
  }

  // --- GET /groups ---
  async getGroups(req: Request, res: Response): Promise<Response> {
    try {
      const params: GetGroupsRequest = {
        limit: req.query.limit ? parseInt(req.query.limit as string, 10) : 10,
        page: req.query.page ? parseInt(req.query.page as string, 10) : 1,
        showAll: req.query.showAll === "true",
      };

      const result = await this.groupService.getGroups(params);
      return res.status(200).json(result);
    } catch (err: any) {
      console.error("[GroupController][getGroups] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Failed to fetch groups",
        error: err.message,
      });
    }
  }

  // --- PUT /groups/:id ---
  async updateGroup(req: Request, res: Response): Promise<Response> {
    try {
      const id = parseInt(req.params.id, 10);
      if (isNaN(id) || id <= 0) {
        return res.status(400).json({
          code: 400,
          message: "group id is required or invalid",
        });
      }

      const data = req.body as Group;
      data.id = id;

      // const errors = await validate(data);
      // if (errors.length > 0) {
      //   return res.status(400).json({
      //     code: 400,
      //     message: "Invalid request body",
      //     error: errors,
      //   });
      // }

      const result = await this.groupService.updateGroup(data);
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[GroupController][updateGroup] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Failed to update group",
        error: err.message,
      });
    }
  }

  // --- DELETE /groups/:id ---
  async deleteGroup(req: Request, res: Response): Promise<Response> {
    try {
      const id = parseInt(req.params.id, 10);
      const username = req.body.username || "system";

      if (isNaN(id) || id <= 0) {
        return res.status(400).json({
          code: 400,
          message: "group id is required or invalid",
        });
      }

      const result = await this.groupService.deleteGroup(id, username);
      return res.status(result.code).json(result);
    } catch (err: any) {
      console.error("[GroupController][deleteGroup] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Failed to delete group",
        error: err.message,
      });
    }
  }
}
