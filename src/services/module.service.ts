import { ModuleRequest, ModuleResponse } from "../DTO/module.dto";
import { Module } from "../models/module.model";
import { BucketRepository } from "../repo/google/bucket";
import { GoogleCfg } from "../repo/google/config";
import { GroupRepository } from "../repo/group_repo";
import { ModuleRepository } from "../repo/module_repo";
import { SubModuleRepository } from "../repo/subModule_repo";
import { Logger } from "../utils/logger";
import {
  DefaultPaginationResponseData,
  DefaultResponse,
  PaginationRequest,
} from "../utils/response";

export class ModuleService {
  private module_repo: ModuleRepository;
  private group_repo: GroupRepository;
  private sub_module_repo: SubModuleRepository;
  private bucket_repo: BucketRepository;
  private google_cfg: GoogleCfg;

  constructor() {
    this.module_repo = new ModuleRepository();
    this.group_repo = new GroupRepository();
    this.sub_module_repo = new SubModuleRepository();
    this.bucket_repo = new BucketRepository();
    this.google_cfg = new GoogleCfg();
  }

  // ======================================================
  // CREATE MODULE
  // ======================================================
  async createModule(req: ModuleRequest): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 201, message: "Success" };

    try {
      const model: Module = {
        module_name: req.module_name,
        created_by: req.created_by ?? "system",
        group_id: req.group_id ?? null,
        thumbnail_url: req.thumbnail_url ?? null,
        tag: req.tag ?? null,
      };

      if (req.group_id) {
        const group = await this.group_repo.getGroupByID(req.group_id);
        if (!group) throw new Error("group not found");
      }

      const module_id = await this.module_repo.createModule(model);
      resp.data = { id: module_id };
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][create_module] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // GET MODULE BY ID
  // ======================================================
  async getModuleByID(module_id: number): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "Success" };

    try {
      const module = await this.module_repo.getModuleByID(module_id);
      if (!module) throw new Error("module not found");

      const data: ModuleResponse = {
        id: module.id ?? 0,
        uuid: module.uuid ?? "",
        module_name: module.module_name,
        created_by: module.created_by ?? "",
        updated_by: module.updated_by ?? "",
        created_at: module.created_at ?? undefined,
        updated_at: module.updated_at ?? undefined,
        tag: module.tag ?? "",
        group_id: module.group_id ?? 0,
        group_name: "",
        thumbnail_url: "",
      };

      if (module.group_id) {
        const group = await this.group_repo.getGroupByID(module.group_id);
        if (group) {
          data.group_name = group.group_name;
        }
      }

      if (module.thumbnail_url) {
        data.thumbnail_url = await this.bucket_repo.presignedURL(
          this.google_cfg.imageBucketName,
          module.thumbnail_url
        );
      }

      resp.data = data;
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][get_module_by_id] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // GET MODULES (pagination)
  // ======================================================
  async getModules(
    req: PaginationRequest
  ): Promise<DefaultPaginationResponseData<ModuleResponse>> {
    const resp: DefaultPaginationResponseData<ModuleResponse> = {
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
      Logger.debug("[service][get_modules] start", { req });

      const { modules, totalCount } = await this.module_repo.getModules(req);

      Logger.debug("[service][get_modules] DB result:", {
        totalCount,
        modulesPreview: modules?.slice(0, 3) ?? [],
      });

      const module_list: ModuleResponse[] = [];

      for (const module of modules) {
        Logger.debug("[service][get_modules] processing module:", {
          id: module.id,
          name: module.module_name,
          group_id: module.group_id,
          thumb: module.thumbnail_url,
        });

        const data: ModuleResponse = {
          id: module.id ?? 0,
          uuid: module.uuid ?? "",
          module_name: module.module_name,
          created_by: module.created_by ?? "",
          updated_by: module.updated_by ?? "",
          created_at: module.created_at ?? undefined,
          updated_at: module.updated_at ?? undefined,
          tag: module.tag ?? "",
          group_id: module.group_id ?? 0,
          group_name: "",
          thumbnail_url: "",
        };

        // Group lookup
        if (module.group_id) {
          try {
            const group = await this.group_repo.getGroupByID(module.group_id);
            if (group) {
              data.group_name = group.group_name;
            } else {
              Logger.warn("[service][get_modules] group not found:", {
                group_id: module.group_id,
              });
            }
          } catch (err: any) {
            Logger.error(
              `[service][get_modules][group_repo] ${err.message}`,
              err
            );
          }
        }

        // Thumbnail URL
        if (module.thumbnail_url && module.thumbnail_url.trim() !== "") {
          try {
            data.thumbnail_url = await this.bucket_repo.presignedURL(
              this.google_cfg.imageBucketName,
              module.thumbnail_url
            );
          } catch (err: any) {
            Logger.error(
              `[service][get_modules][presignedURL] ${err.message}`,
              err
            );
            data.thumbnail_url = module.thumbnail_url; // fallback biar gak crash
          }
        } else {
          Logger.debug("[service][get_modules] module tanpa thumbnail:", {
            id: module.id,
            module_name: module.module_name,
          });
        }

        module_list.push(data);
      }

      resp.results = module_list;
      resp.meta.totalItems = totalCount;
      resp.meta.totalPages = Math.ceil(totalCount / req.limit);
      resp.meta.hasNext = req.page < resp.meta.totalPages;
      resp.meta.hasPrevious = req.page > 1;

      Logger.debug("[service][get_modules] finish:", {
        totalItems: totalCount,
        returned: module_list.length,
      });
    } catch (err: any) {
      Logger.error(`[service][get_modules] ${err.message}`, err);
    }

    return resp;
  }

  // ======================================================
  // UPDATE MODULE
  // ======================================================
  async updateModule(
    module_id: number,
    req: ModuleRequest
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "Success" };

    try {
      const existing = await this.module_repo.getModuleByID(module_id);
      if (!existing) throw new Error("module not found");
      if (req.module_name === existing.module_name)
        throw new Error("module name is same as before");

      const model: Module = {
        id: module_id,
        module_name: req.module_name,
        updated_by: req.updated_by ?? "system",
        thumbnail_url: req.thumbnail_url ?? null,
        tag: req.tag ?? null,
      };

      await this.module_repo.updateModule(model);
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][update_module] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // GET MODULES BY GROUP ID
  // ======================================================
  async getModuleByGroupID(group_id: number): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "Success" };

    try {
      const modules = await this.module_repo.getModulesByGroupID(group_id);
      const result: ModuleResponse[] = [];

      for (const module of modules) {
        const data: ModuleResponse = {
          id: module.id ?? 0,
          uuid: module.uuid ?? "",
          module_name: module.module_name,
          created_by: module.created_by ?? "",
          updated_by: module.updated_by ?? "",
          created_at: module.created_at ?? undefined,
          updated_at: module.updated_at ?? undefined,
          tag: module.tag ?? "",
          group_id: module.group_id ?? 0,
          group_name: "",
          thumbnail_url: "",
        };

        if (module.thumbnail_url) {
          data.thumbnail_url = await this.bucket_repo.presignedURL(
            this.google_cfg.imageBucketName,
            module.thumbnail_url
          );
        }

        result.push(data);
      }

      resp.data = result;
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][get_modules_by_group_id] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // DELETE MODULE
  // ======================================================
  async deleteModule(
    module_id: number,
    username: string
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "Success" };

    try {
      const module = await this.module_repo.getModuleByID(module_id);
      if (!module) throw new Error("module not found");

      await this.module_repo.softDeleteModule(module_id, username);
      await this.sub_module_repo.removeModuleIDFromSubModules(
        module_id,
        username
      );
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][delete_module] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // GET WATCHED PERCENTAGE
  // ======================================================
  async getPercentageMarkWatchedModulesUser(
    user_id: number,
    module_id: number
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "Success" };

    try {
      const percentage =
        await this.module_repo.getPercentageMarkWatchedModulesUser(
          user_id,
          module_id
        );
      resp.data = { percentage };
    } catch (err: any) {
      resp.code = 500;
      resp.message = "Internal server error";
      resp.error = err.message;
      Logger.error(
        `[service][get_percentage_mark_watched_modules_user] ${err.message}`
      );
    }

    return resp;
  }
}
