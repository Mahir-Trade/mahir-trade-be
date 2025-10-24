import path from "path";
import mime from "mime";
import { ModuleRepository } from "../repo/module_repo";
import { SubModuleRepository } from "../repo/subModule_repo";
import { BucketRepository } from "../repo/google/bucket";
import { UserSubModuleRepository } from "../repo/userSubModule_repo";
import { GoogleCfg } from "../repo/google/config";
import {
  DefaultResponse,
  DefaultPaginationResponseData,
  PaginationRequest,
} from "../utils/response";
import {
  SubModuleRequest,
  MarkSubModuleAsWatchedRequest,
  SubModuleResponse,
} from "../DTO/subModule.dto";
import { FileUpload } from "../DTO/google.dto";
import { Logger } from "../utils/logger";
import { SubModule } from "../models/subModule.model"; // pastikan ada

export class SubModuleService {
  private moduleRepo = new ModuleRepository();
  private subModuleRepo = new SubModuleRepository();
  private bucketRepo = new BucketRepository();
  private userSubModuleRepo = new UserSubModuleRepository();
  private googleCfg = new GoogleCfg();

  // ======================================================
  // CREATE SUB MODULE
  // ======================================================
  async createSubModule(req: SubModuleRequest): Promise<DefaultResponse> {
    const entry: SubModule = {
      sub_module_name: req.sub_module_name,
      title: req.title,
      video_url: req.video_url,
      module_id: req.module_id ?? null,
      created_by: "system",
    };

    if (req.module_id) {
      const module = await this.moduleRepo.getModuleByID(req.module_id);
      if (!module) return { code: 404, message: "Module not found" };
    }

    const id = await this.subModuleRepo.createSubModule(entry);

    if (req.video_url) {
      const parsed = this.bucketRepo.urlParser(req.video_url);
      const fileName = parsed.path.split("/")[2];
      this.bucketRepo.startTranscodingJob(
        this.googleCfg.videoBucketName,
        fileName
      );
    }

    return { code: 200, message: "Success", data: { id } };
  }

  // ======================================================
  // GET ALL SUB MODULES
  // ======================================================
  // async getSubModules(req: {
  //   limit: number;
  //   page: number;
  // }): Promise<DefaultResponse> {
  //   const { subModules, totalCount } = await this.subModuleRepo.getSubModules(
  //     req
  //   );
  //   const data = [];

  //   for (const sm of subModules) {
  //     const item = { ...sm };
  //     if (sm.video_url) {
  //       item.video_url = await this.bucketRepo.presignedURL(
  //         this.googleCfg.videoBucketName,
  //         sm.video_url
  //       );
  //     }
  //     data.push(item);
  //   }

  //   return {
  //     code: 200,
  //     message: "Success",
  //     data: {
  //       page: req.page,
  //       limit: req.limit,
  //       totalPages: Math.ceil(totalCount / req.limit),
  //       totalItems: totalCount,
  //       hasNext: req.page * req.limit < totalCount,
  //       hasPrevious: req.page > 1,
  //       results: data,
  //     },
  //   };
  // }

  async getSubModules(req: {
    limit: number;
    page: number;
  }): Promise<DefaultPaginationResponseData<any>> {
    const resp: DefaultPaginationResponseData<any> = {
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
      const { subModules, totalCount } = await this.subModuleRepo.getSubModules(
        req
      );

      const subModuleList = [] as any[];

      for (const sm of subModules) {
        const item = { ...sm };
        if (sm.video_url) {
          item.video_url = await this.bucketRepo.presignedURL(
            this.googleCfg.videoBucketName,
            sm.video_url
          );
        }
        subModuleList.push(item);
      }

      // Meta info
      resp.results = subModuleList;
      resp.meta.totalItems = totalCount;
      resp.meta.totalPages = Math.ceil(totalCount / req.limit);
      resp.meta.hasNext = req.page < resp.meta.totalPages;
      resp.meta.hasPrevious = req.page > 1;
    } catch (err: any) {
      Logger.error(`[service][get_submodules] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // GET SUB MODULE BY ID
  // ======================================================
  async getSubModuleByID(id: number): Promise<DefaultResponse> {
    const subModule = await this.subModuleRepo.getSubModuleByID(id);
    if (!subModule) return { code: 404, message: "Sub module not found" };

    const resolutions = ["240p", "360p", "480p", "720p", "1080p"];
    const sources = [];

    for (const res of resolutions) {
      let videoUrl = subModule.video_url;
      let bucketName = this.googleCfg.videoBucketName;

      if (res !== "1080p") {
        const parsed = this.bucketRepo.urlParser(subModule.video_url);
        const idx = parsed.path.lastIndexOf(".mp4");
        if (idx === -1) continue;

        let newPath =
          parsed.path.slice(0, idx) + `-${res}` + parsed.path.slice(idx);
        const pathSplit = newPath.split("/");
        pathSplit[1] += "_transcoded";
        newPath = pathSplit.join("/");

        videoUrl = `https://${parsed.host}${newPath}`;
        bucketName += "_transcoded";
      }

      const signedUrl = await this.bucketRepo.presignedURL(
        bucketName,
        videoUrl
      );
      sources.push({ resolution: res, url: signedUrl });
    }

    return {
      code: 200,
      message: "Success",
      data: {
        id: subModule.id!,
        sub_module_name: subModule.sub_module_name!,
        sources,
      },
    };
  }

  // ======================================================
  // GET SUB MODULES BY MODULE ID
  // ======================================================
  async getSubModulesByModuleID(
    module_id: number,
    req: PaginationRequest,
    user_id?: number
  ): Promise<DefaultPaginationResponseData<SubModuleResponse>> {
    const resp: DefaultPaginationResponseData<SubModuleResponse> = {
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
      const { subModules, totalCount } =
        await this.subModuleRepo.getSubModulesByModuleID(
          module_id,
          user_id ?? 0,
          req
        );

      const subModulesData: SubModuleResponse[] = [];

      for (const subModule of subModules) {
        const data: SubModuleResponse = {
          id: subModule.id!,
          uuid: subModule.uuid!,
          sub_module_name: subModule.sub_module_name!,
          title: subModule.title!,
          status: subModule.status!,
          created_by: subModule.created_by!,
          updated_by: subModule.updated_by!,
          created_at: subModule.created_at!,
          updated_at: subModule.updated_at!,
        };

        // === Relasi ke Module ===
        if (subModule.module_id) {
          try {
            const module = await this.moduleRepo.getModuleByID(
              subModule.module_id
            );
            if (module) {
              data.module_id = module.id!;
              data.module_name = module.module_name;
            }
          } catch (err: any) {
            Logger.error(
              `[service][getSubModulesByModuleID] error get module by id: ${err.message}`
            );
          }
        }

        // === Dapatkan presigned URL ===
        if (subModule.video_url) {
          try {
            const url = await this.bucketRepo.presignedURL(
              this.googleCfg.videoBucketName,
              subModule.video_url
            );
            data.video_url = url;
          } catch (err: any) {
            Logger.error(
              `[service][getSubModulesByModuleID] error get presigned url: ${err.message}`
            );
          }
        }

        subModulesData.push(data);
      }

      const total_pages = Math.ceil(totalCount / (req.limit || 10));

      resp.results = subModulesData;
      resp.meta.totalItems = totalCount;
      resp.meta.totalPages = total_pages;
      resp.meta.hasNext = (req.page || 1) < total_pages;
      resp.meta.hasPrevious = (req.page || 1) > 1;
    } catch (err: any) {
      Logger.error(
        `[service][getSubModulesByModuleID] error while fetching data: ${err.message}`
      );
    }

    return resp;
  }

  // ======================================================
  // UPDATE & DELETE
  // ======================================================
  async updateSubModule(
    id: number,
    req: SubModuleRequest
  ): Promise<DefaultResponse> {
    const entry: SubModule = {
      id,
      sub_module_name: req.sub_module_name,
      title: req.title,
      video_url: req.video_url,
      module_id: req.module_id ?? null,
    };
    await this.subModuleRepo.updateSubModule(entry);
    return { code: 200, message: "Success" };
  }

  async softDeleteSubModule(subModuleId: number): Promise<DefaultResponse> {
    await this.subModuleRepo.softDeleteSubModule(subModuleId, "system");
    return { code: 200, message: "Success" };
  }

  // ======================================================
  // MARK AS WATCHED
  // ======================================================
  async markSubModuleAsWatched(
    req: MarkSubModuleAsWatchedRequest
  ): Promise<DefaultResponse> {
    const subModule = await this.subModuleRepo.getSubModuleByID(
      req.subModuleId
    );
    if (!subModule) return { code: 404, message: "Sub Module Not Found" };

    const existing =
      await this.userSubModuleRepo.getUserSubModuleBySubModuleIDAndUserID(
        0,
        req.subModuleId
      );
    if (existing) return { code: 200, message: "Already watched" };

    await this.userSubModuleRepo.createUserSubModule({
      user_id: 0,
      sub_module_id: req.subModuleId,
    });

    return { code: 200, message: "Success" };
  }

  // ======================================================
  // FILE UPLOAD
  // ======================================================
  async uploadFile(
    req: FileUpload,
    file: Express.Multer.File
  ): Promise<DefaultResponse> {
    const ext = path.extname(file.originalname);
    const contentType = mime.lookup(ext) || "application/octet-stream";

    const fileUpload: FileUpload = {
      file_name: file.originalname,
      size: file.size,
      localFilePath: file.path,
      file_content_type: contentType,
      bucket_name: "",
    };

    if (contentType.includes("video"))
      fileUpload.bucket_name = this.googleCfg.videoBucketName;
    else if (["image/jpeg", "image/png", "image/jpg"].includes(contentType))
      fileUpload.bucket_name = this.googleCfg.imageBucketName;
    else fileUpload.bucket_name = this.googleCfg.fileBucketName;

    const url = await this.bucketRepo.uploadFile(fileUpload);
    return { code: 200, message: "Success", data: { url } };
  }
}
