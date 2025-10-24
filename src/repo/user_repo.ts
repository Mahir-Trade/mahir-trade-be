import { UserQueries } from "./queries/user_queries";
import { User, GetUsersBOResponse } from "../models/user.model";
import db from "../db/db.config";
import { PaginationRequest } from "../utils/response";

export class UserRepository {
  async createUser(user: User): Promise<number> {
    try {
      const result = await db.query(UserQueries.CreateUser, [
        user.email,
        user.phone_number,
        user.username,
        user.password,
      ]);
      return result.rows[0].user_id;
    } catch (err: any) {
      if (err.message.includes("duplicate key value")) {
        throw new Error("email or username already exist");
      }
      throw err;
    }
  }

  async findUserByEmailOrUsername(identity: string): Promise<User | null> {
    console.log("Identity data: ", identity);

    const result = await db.query(UserQueries.FindUserByEmailOrUsername, [
      identity,
    ]);

    console.log("Result Hasil: ", result);

    return result.rows[0] || null;
  }

  async findUserByEmailAndUsername(
    email: string,
    username: string
  ): Promise<User | null> {
    const result = await db.query(UserQueries.FindUserByEmailAndUsername, [
      email,
      username,
    ]);
    return result.rows[0] || null;
  }

  async getUserById(id: number): Promise<User> {
    const result = await db.query(UserQueries.GetUserByID, [id]);
    if (result.rowCount === 0) {
      throw new Error("user not found");
    }
    return result.rows[0];
  }

  async getUserByUUID(uuid: string): Promise<User | null> {
    const result = await db.query(UserQueries.GetUserByUUID, [uuid]);
    return result.rows[0] || null;
  }

  async updateTypeUser(
    id: number,
    isActive: boolean,
    operator: string
  ): Promise<boolean> {
    const result = await db.query(UserQueries.UpdateTypeUser, [
      isActive,
      operator,
      id,
    ]);

    // rowCount pasti ada, defaultnya 0 kalau tidak ada row yang kena update
    return result.rowCount! > 0;
  }

  async getAllUser(
    req: PaginationRequest
  ): Promise<{ users: any[]; totalCount: number }> {
    let query = UserQueries.GetUsers;
    const params: any[] = [];

    // 🔍 Filter search by email
    if (req.search) {
      query += ` AND LOWER(u.email) ILIKE '%' || $${params.length + 1} || '%'`;
      params.push(req.search.toLowerCase());
    }

    // 🧩 Filter membership status
    if (req.membership_status) {
      const status = req.membership_status.toUpperCase();

      if (status === "ACTIVE") {
        query += ` AND um.is_membership_active = $${params.length + 1}`;
        params.push(true);
      } else if (status === "EXPIRED") {
        query += ` AND um.is_membership_active = $${params.length + 1}`;
        params.push(false);
      } else if (status === "PRE_ORDER") {
        query += ` AND um.is_membership_active IS NULL`;
      }
    }

    // 🔢 Sort + Pagination (⚠️ tambahin spasi di depan ORDER)
    const sortDirection = req.sort_by?.toUpperCase() === "ASC" ? "ASC" : "DESC";
    query += ` ORDER BY u.created_at ${sortDirection} LIMIT $${
      params.length + 1
    } OFFSET $${params.length + 2}`;
    params.push(req.limit, (req.page - 1) * req.limit);

    // Debug log
    console.log("[DEBUG] Final Query:", query);
    console.log("[DEBUG] Params:", params);

    // 🧠 Execute query
    const result = await db.query(query, params);

    // 🧾 Map data result
    const users = result.rows.map((row: any) => ({ ...row }));
    const totalCount =
      result.rows.length > 0 ? parseInt(result.rows[0].total_count, 10) : 0;

    return { users, totalCount };
  }

  async updatePassword(
    userId: number,
    password: string,
    operator: string
  ): Promise<boolean> {
    const result = await db.query(UserQueries.UpdatePassword, [
      password,
      operator,
      userId,
    ]);
    return result.rowCount! > 0;
  }

  async setUserVerified(userId: number): Promise<boolean> {
    const result = await db.query(UserQueries.SetUserVerified, [userId]);
    return result.rowCount! > 0;
  }
}
