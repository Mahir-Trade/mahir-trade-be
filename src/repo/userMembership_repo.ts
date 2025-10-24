import { Pool } from "pg";
import { UserMembership } from "../models/userMembership.model";

import db from "../db/db.config";
import { UserMembershipQueries } from "./queries/userMembership_queries";

export class UserMembershipRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- CREATE USER MEMBERSHIP ---
  async createUserMembership(req: UserMembership): Promise<number> {
    try {
      const result = await this.pg.query(
        UserMembershipQueries.CreateUserMembership,
        [
          req.user_id,
          req.package_id,
          req.expired_at,
          req.is_membership_active,
          req.created_by,
        ]
      );

      if (result.rowCount === 0) {
        throw new Error("failed to create user membership");
      }

      return result.rows[0].id as number;
    } catch (err: any) {
      console.error(
        `[UserMembershipRepository][createUserMembership] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- BULK UPDATE USER MEMBERSHIPS ---
  async updateBulkUserMembership(): Promise<boolean> {
    const client = await this.pg.connect();
    try {
      await client.query("BEGIN");

      const result = await client.query(
        UserMembershipQueries.BulkUpdateUserMembershipExpired
      );

      await client.query("COMMIT");

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      await client.query("ROLLBACK");
      console.error(
        `[UserMembershipRepository][updateBulkUserMembership] error: ${err.message}`
      );
      throw err;
    } finally {
      client.release();
    }
  }

  // --- UPDATE USER MEMBERSHIP BY USER ID ---
  async updateUserMembershipByUserID(req: UserMembership): Promise<boolean> {
    try {
      const result = await this.pg.query(
        UserMembershipQueries.UpdateUserMembershipExpired,
        [req.expired_at, req.is_membership_active, req.updated_by, req.user_id]
      );

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(
        `[UserMembershipRepository][updateUserMembershipByUserID] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- GET USER MEMBERSHIP BY USER ID ---
  async getUserMembershipByUserID(
    userId: number
  ): Promise<UserMembership | null> {
    try {
      const result = await this.pg.query(
        UserMembershipQueries.GetUserMembershipByUserID,
        [userId]
      );

      if (result.rowCount === 0) {
        return null;
      }

      return result.rows[0] as UserMembership;
    } catch (err: any) {
      console.error(
        `[UserMembershipRepository][getUserMembershipByUserID] error: ${err.message}`
      );
      throw err;
    }
  }

  async getUserMembershipExpired(): Promise<UserMembership[]> {
    try {
      const result = await db.query(
        UserMembershipQueries.GetUserMembershipExpired
      );
      return result.rows.map((row) => ({
        id: row.id,
        user_id: row.user_id,
        expired_at: row.expired_at,
        is_membership_active: row.is_membership_active,
        exclusive_expired_at: row.exclusive_expired_at,
        package_id: row.package_id,
      }));
    } catch (err: any) {
      console.error(
        "[userMembershipRepoImpl][getUserMembershipExpired] error while QueryContext:",
        err.message
      );
      throw err;
    }
  }

  /**
   * 🔁 Sama persis dengan UpdateUserMembershipExpired di Go
   */
  async updateUserMembershipExpired(
    userIds: number[],
    updatedBy: string
  ): Promise<any> {
    try {
      let queryBuilder = UserMembershipQueries.UpdateUserMembershipsByUserIDs;
      queryBuilder += " WHERE deleted_at IS NULL AND user_id IN (";

      const placeholders = userIds.map((_, i) => `$${i + 2}`).join(", ");
      queryBuilder += placeholders + ")";

      const params = [updatedBy, ...userIds];

      await db.query(queryBuilder, params);
    } catch (err: any) {
      console.error(
        "[userMembershipRepoImpl][updateUserMembershipExpired] error while ExecContext:",
        err.message
      );
      throw err;
    }
  }

  /**
   * 🚀 Sama persis dengan BulkUpdateMembershipPreOrderActivation di Go
   */
  async bulkUpdateMembershipPreOrderActivation(): Promise<any> {
    const client = await db.connect();
    try {
      await client.query("BEGIN");
      await client.query(
        UserMembershipQueries.BulkUpdateMembershipPreOrderActivation
      );
      await client.query("COMMIT");
    } catch (err: any) {
      await client.query("ROLLBACK");
      console.error(
        "[userMembershipRepoImpl][bulkUpdateMembershipPreOrderActivation] error while ExecContext:",
        err.message
      );
      throw err;
    } finally {
      client.release();
    }
  }
}
