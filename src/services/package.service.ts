import { PackageRepository } from "../repo/package_repo";
import { Logger } from "../utils/logger";
import { Package } from "../models/package.model";
import {
  DefaultPaginationResponseData,
  DefaultResponse,
  PaginationRequest,
} from "../utils/response";

export class PackageService {
  private packageRepo: PackageRepository;

  constructor() {
    this.packageRepo = new PackageRepository();
  }

  // ======================================================
  // PRIVATE HELPER
  // ======================================================
  private getDiscountExpired(now: Date): Date {
    const firstDayNextMonth = new Date(
      now.getFullYear(),
      now.getMonth() + 1,
      1,
      0,
      0,
      0
    );
    return new Date(firstDayNextMonth.getTime() - 1000); // minus 1 second
  }

  // ======================================================
  // CREATE PACKAGE
  // ======================================================
  async createPackage(req: Package): Promise<DefaultResponse<{ id: number }>> {
    const resp: DefaultResponse<{ id: number }> = {
      code: 201,
      message: "Success",
    };

    try {
      const discountExpired = this.getDiscountExpired(new Date());

      // ✅ fix: handle undefined safely
      if (typeof req.discountedPrice === "number" && req.discountedPrice >= 0) {
        await this.packageRepo.updatePackageDiscountExpired(discountExpired);
      }

      req.discountExpired = discountExpired.toISOString();
      const packageId = await this.packageRepo.createPackage(req);

      resp.data = { id: packageId };
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][CreatePackage] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // GET PACKAGES
  // ======================================================
  async getPackages(
    req: PaginationRequest
  ): Promise<DefaultPaginationResponseData<Package>> {
    // ✅ fix: use correct structure (results + meta)
    const resp: DefaultPaginationResponseData<Package> = {
      results: [],
      meta: {
        page: req.page,
        limit: req.limit,
        totalPages: 0,
        totalItems: 0,
        hasNext: false,
        hasPrevious: false,
      },
    };

    try {
      const { packages, totalCount } = await this.packageRepo.getPackages(req);

      const totalPages = Math.ceil(totalCount / req.limit);

      resp.results = packages;
      resp.meta = {
        page: req.page,
        limit: req.limit,
        totalPages,
        totalItems: totalCount,
        hasNext: req.page < totalPages,
        hasPrevious: req.page > 1,
      };
    } catch (err: any) {
      Logger.error(`[service][GetPackages] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // GET PACKAGE BY ID
  // ======================================================
  async getPackageByID(packageId: number): Promise<DefaultResponse<Package>> {
    const resp: DefaultResponse<Package> = { code: 200, message: "Success" };

    try {
      // ✅ fix typo: method name from getPackageById → getPackageByID
      const packageData = await this.packageRepo.getPackageByID(packageId);
      resp.data = packageData;
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][GetPackageByID] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // UPDATE PACKAGE
  // ======================================================
  async updatePackage(req: Package): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "Success" };

    try {
      if (!req.id) {
        throw new Error("package id is required");
      }

      const packageData = await this.packageRepo.getPackageByID(req.id);
      if (!packageData || packageData.id === 0) {
        throw new Error(`package with id ${req.id} not found`);
      }

      if (
        typeof req.discountedPrice === "number" &&
        req.discountedPrice !== 0 &&
        req.discountedPrice !== packageData.discountedPrice
      ) {
        const newDiscountExpired = this.getDiscountExpired(new Date());
        await this.packageRepo.updatePackageDiscountExpired(newDiscountExpired);
      }

      await this.packageRepo.updatePackage(req);
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][UpdatePackage] ${err.message}`);
    }

    return resp;
  }

  // ======================================================
  // DELETE PACKAGE
  // ======================================================
  async deletePackage(
    packageId: number,
    deletedBy: string
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "Success" };

    try {
      const packageData = await this.packageRepo.getPackageByID(packageId);
      if (!packageData || packageData.id === 0) {
        throw new Error(`package with id ${packageId} not found`);
      }

      await this.packageRepo.softDeletePackage(packageId, deletedBy);
    } catch (err: any) {
      resp.code = 400;
      resp.message = "Bad request";
      resp.error = err.message;
      Logger.error(`[service][DeletePackage] ${err.message}`);
    }

    return resp;
  }
}
