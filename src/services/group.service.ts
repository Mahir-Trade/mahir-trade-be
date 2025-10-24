import {
  DefaultResponse,
  DefaultPaginationResponseData,
} from "../utils/response";
import { Group, GetGroupsRequest } from "../models/group.model";

import { Logger } from "../utils/logger";
import { GroupRepository } from "../repo/group_repo";
import { ModuleRepository } from "../repo/module_repo";

export class GroupService {
  private groupRepo: GroupRepository;
  private moduleRepo: ModuleRepository;

  constructor() {
    this.groupRepo = new GroupRepository();
    this.moduleRepo = new ModuleRepository();
  }

  // ======================================================
  // CREATE GROUP
  // ======================================================
  async createGroup(data: Group): Promise<DefaultResponse<{ id: number }>> {
    const resp: DefaultResponse<{ id: number }> = {
      code: 201,
      message: "Success",
    };

    try {
      const id = await this.groupRepo.createGroup(data);
      resp.data = { id };
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][createGroup] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // GET GROUP BY ID
  // ======================================================
  async getGroupById(groupId: number): Promise<DefaultResponse<Group>> {
    const resp: DefaultResponse<Group> = {
      code: 200,
      message: "Success",
    };

    try {
      const group = await this.groupRepo.getGroupByID(groupId);
      if (!group) {
        const errMsg = "Group not found";
        resp.code = 404;
        resp.message = errMsg;
        resp.error = errMsg;
        Logger.error(`[service][getGroupById] ${errMsg}`);
        return resp;
      }

      resp.data = group;
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][getGroupById] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // GET GROUPS (with pagination)
  // ======================================================
  async getGroups(
    req: GetGroupsRequest
  ): Promise<DefaultPaginationResponseData<Group>> {
    const resp: DefaultPaginationResponseData<Group> = {
      results: [],
      meta: {
        page: req.page || 1,
        limit: req.limit || 10,
        totalItems: 0,
        totalPages: 0,
        hasNext: false,
        hasPrevious: false,
      },
    };

    try {
      const { groups, totalCount } = await this.groupRepo.getGroups(req);

      const totalPages = Math.ceil(totalCount / (req.limit || 10));
      const hasNext = (req.page || 1) < totalPages;
      const hasPrevious = (req.page || 1) > 1;

      resp.results = groups;
      resp.meta = {
        page: req.page || 1,
        limit: req.limit || 10,
        totalItems: totalCount,
        totalPages,
        hasNext,
        hasPrevious,
      };
    } catch (err: any) {
      Logger.error(`[service][getGroups] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // UPDATE GROUP
  // ======================================================
  async updateGroup(data: Group): Promise<DefaultResponse> {
    const resp: DefaultResponse = {
      code: 200,
      message: "Success",
    };

    try {
      const existing = await this.groupRepo.getGroupByID(data.id!);
      if (!existing) {
        const errMsg = "Group not found";
        resp.code = 404;
        resp.message = errMsg;
        resp.error = errMsg;
        Logger.error(`[service][updateGroup] ${errMsg}`);
        return resp;
      }

      if (data.group_name === existing.group_name) {
        const errMsg = "Group name is the same as before";
        resp.code = 400;
        resp.message = errMsg;
        resp.error = errMsg;
        Logger.error(`[service][updateGroup] ${errMsg}`);
        return resp;
      }

      await this.groupRepo.updateGroup(data);
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][updateGroup] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // DELETE GROUP
  // ======================================================
  async deleteGroup(
    groupId: number,
    username: string = "system"
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = {
      code: 200,
      message: "Success",
    };

    try {
      const existing = await this.groupRepo.getGroupByID(groupId);
      if (!existing) {
        const errMsg = "Group not found";
        resp.code = 404;
        resp.message = errMsg;
        resp.error = errMsg;
        Logger.error(`[service][deleteGroup] ${errMsg}`);
        return resp;
      }

      await this.groupRepo.softDeleteGroup(groupId, username);
      await this.moduleRepo.removeGroupIDFromModules(groupId, username);
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][deleteGroup] ${err.message}`);
    }

    return resp;
  }
}
