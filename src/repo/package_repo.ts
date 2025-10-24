import { Pool } from "pg";

import { Package } from "../models/package.model";

import db from "../db/db.config";
import { PackageQueries } from "./queries/package_queries";
import { PaginationRequest } from "@utils/response";

export class PackageRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- CREATE PACKAGE ---
  async createPackage(req: Package): Promise<number> {
    try {
      const result = await this.pg.query<{ id: number }>(
        PackageQueries.CreatePackage,
        [
          req.price,
          req.durationInMonth,
          req.description,
          req.discountedPrice,
          req.discountExpired,
          req.createdBy,
        ]
      );

      if (result.rowCount === 0) {
        throw new Error("failed to create package");
      }

      return result.rows[0].id;
    } catch (err: any) {
      console.error(`[packageRepo][createPackage] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET PACKAGES ---
  async getPackages(
    req: PaginationRequest
  ): Promise<{ packages: Package[]; totalCount: number }> {
    try {
      if (req.limit === 0) {
        req.limit = 10;
      }

      let offset = 0;
      if (req.page > 1) {
        offset = (req.page - 1) * req.limit;
      }

      const result = await this.pg.query(PackageQueries.GetPackages, [
        req.limit,
        offset,
      ]);

      const packages: Package[] = (result.rows as any[]).map((row) => ({
        id: row.id,
        price: row.price,
        durationInMonth: row.duration_in_month, // snake_case dari DB
        description: row.description,
        discountedPrice: row.discounted_price,
        discountExpired: row.discount_expired,
        createdBy: row.created_by,
        updatedBy: row.updated_by,
        createdAt: row.created_at,
        updatedAt: row.updated_at,
      }));

      const totalCount =
        result.rows.length > 0 ? parseInt(result.rows[0].total_count, 10) : 0;

      return { packages, totalCount };
    } catch (err: any) {
      console.error(`[packageRepo][getPackages] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET PACKAGE BY ID ---
  async getPackageByID(id: number): Promise<Package> {
    try {
      const result = await this.pg.query<Package>(
        PackageQueries.GetPackageByID,
        [id]
      );

      if (result.rowCount === 0) {
        throw new Error(`package with id ${id} not found`);
      }

      return result.rows[0];
    } catch (err: any) {
      console.error(`[packageRepo][getPackageByID] error: ${err.message}`);
      throw err;
    }
  }

  // --- UPDATE PACKAGE ---
  async updatePackage(req: Package): Promise<boolean> {
    try {
      const result = await this.pg.query(PackageQueries.UpdatePackage, [
        req.price,
        req.durationInMonth,
        req.description,
        req.discountedPrice,
        req.id,
        req.updatedBy,
      ]);

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[packageRepo][updatePackage] error: ${err.message}`);
      throw err;
    }
  }

  // --- SOFT DELETE PACKAGE ---
  async softDeletePackage(id: number, deletedBy: string): Promise<boolean> {
    try {
      const result = await this.pg.query(PackageQueries.SoftDeletePackage, [
        deletedBy,
        deletedBy,
        id,
      ]);

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[packageRepo][softDeletePackage] error: ${err.message}`);
      throw err;
    }
  }

  // --- UPDATE PACKAGE DISCOUNT EXPIRED ---
  async updatePackageDiscountExpired(time: Date): Promise<boolean> {
    try {
      const result = await this.pg.query(
        PackageQueries.UpdatePackageDiscountExpired,
        [time]
      );

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(
        `[packageRepo][updatePackageDiscountExpired] error: ${err.message}`
      );
      throw err;
    }
  }
}
