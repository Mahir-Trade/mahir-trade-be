import { Admin, AdminLoginRequest } from "../models/admin.model";
import {
  DefaultPaginationResponseData,
  DefaultResponse,
  PaginationRequest,
} from "../utils/response";

import { UserRepository } from "../repo/user_repo";
import { hashPassword, verifyPassword } from "../utils/bcrypt";
import { signJWT } from "../utils/jwt";
import { AdminRepository } from "../repo/admin_repo";
import { UserContext } from "../config/userContext";
import { StartMembershipProgramRequest } from "../DTO//startMembership";
import { ConfigRepository } from "../repo/config_repo";
import { UserMembershipRepository } from "../repo/userMembership_repo";
import { ToggleUserMembershipRequest } from "@dto/toggleUserMembershipRequest";

export class AdminService {
  private adminRepository: AdminRepository;
  private userRepository: UserRepository;
  private userMembershipRepository: UserMembershipRepository;
  private configRepository: ConfigRepository;

  constructor() {
    this.adminRepository = new AdminRepository();
    this.userRepository = new UserRepository();
    this.configRepository = new ConfigRepository();
    this.userMembershipRepository = new UserMembershipRepository();
  }

  async adminLogin(req: AdminLoginRequest): Promise<DefaultResponse> {
    let resp: DefaultResponse = {
      code: 200,
      message: "Login Success",
      data: {},
    };

    const admin = await this.adminRepository.findByUsername(req.identity);
    if (!admin) {
      resp.code = 401;
      resp.message = "Invalid Username or Password";
      return resp;
    }
    if (!admin.password) {
      resp.code = 500;
      resp.message = "Internal Server Error";
      return resp;
    }

    const valid = await verifyPassword(req.password, admin.password);
    if (!valid) {
      resp.code = 401;
      resp.message = "Invalid Username or Password";
      return resp;
    }

    const payload = {
      userId: admin.adminId,
      email: admin.email,
      username: admin.username,
    };

    const { token, exp } = await signJWT(payload);

    resp.data = { token, exp };
    return resp;
  }

  async adminRegistration(req: Admin): Promise<DefaultResponse> {
    let resp: DefaultResponse = {
      code: 201,
      message: "Registration Success",
    };

    req.email = req.email.trim().toLowerCase();
    req.username = req.username.trim().toLowerCase();

    if (!req.password) {
      resp.code = 400;
      resp.message = "Invalid request body";
      resp.data = "password is required";
      return resp;
    }
    req.password = await hashPassword(req.password);

    const adminId = await this.adminRepository.createAdmin(req);

    const payload = {
      userId: adminId,
      email: req.email,
      username: req.username,
    };

    const { token, exp } = await signJWT(payload);

    resp.data = { token, exp };
    return resp;
  }

  async updateTypeUser(
    id: number,
    isActive: boolean,
    currentAdmin: { username: string }
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = {
      code: 200,
      message: "Success",
      data: {},
    };

    const user = await this.userRepository.getUserById(id);
    if (!user) {
      resp.code = 404;
      resp.message = "User Not Found";
      return resp;
    }

    if (user.is_active === isActive) {
      resp.code = 400;
      resp.message = "User already in this state";
      return resp;
    }

    // ✅ parameter order sama dengan repository
    const state = await this.userRepository.updateTypeUser(
      id,
      isActive,
      currentAdmin.username
    );

    resp.data = state;
    return resp;
  }

  async getDetailAdminInfo(username: string): Promise<DefaultResponse> {
    let resp: DefaultResponse = {
      code: 200,
      message: "Success",
      data: {},
    };

    console.log(
      "🔍 [DEBUG] getDetailAdminInfo() dipanggil dengan username:",
      username
    );

    try {
      const admin = await this.adminRepository.findByUsername(username);
      console.log("📦 [DEBUG] Hasil pencarian admin di DB:", admin);

      if (!admin) {
        console.error(
          "❌ [DEBUG] Admin tidak ditemukan untuk username:",
          username
        );
        resp.code = 404;
        resp.message = "Admin not found";
        return resp;
      }

      // Hapus password sebelum dikembalikan
      admin.password = "";
      resp.data = admin;

      console.log("✅ [DEBUG] Response final:", resp);
      return resp;
    } catch (err) {
      console.error("🔥 [ERROR] Terjadi error di getDetailAdminInfo:", err);
      resp.code = 500;
      resp.message = "Internal Server Error";
      return resp;
    }
  }

  async getAllUsers(
    req: PaginationRequest
  ): Promise<DefaultPaginationResponseData> {
    const ctx = UserContext.get();
    console.log("🧠 [Service] Dipanggil oleh:", ctx?.username, ctx?.email);

    const { users, totalCount } = await this.userRepository.getAllUser(req);

    const page = req.page ?? 1;
    const limit = req.limit ?? 10;
    const totalPages = Math.ceil(totalCount / limit);

    const resp: DefaultPaginationResponseData = {
      results: users,
      meta: {
        page,
        limit,
        totalPages,
        totalItems: totalCount,
        hasNext: page < totalPages,
        hasPrevious: page > 1,
      },
    };

    return resp;
  }

  async toggleInactiveUserMembership(
    req: ToggleUserMembershipRequest
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = {
      code: 200,
      message: "Success",
      data: {},
    };

    try {
      // 1️⃣ Validasi userIds wajib
      if (!req.user_ids || req.user_ids.length === 0) {
        const msg = "User IDs are required";
        resp.code = 400;
        resp.message = msg;
        resp.error = msg;
        return resp;
      }

      // 2️⃣ Ambil admin dari context
      const adminData = UserContext.get();
      if (!adminData || !adminData.username) {
        const msg = "failed to get admin data from context";
        console.error("[service][ToggleInactiveUserMembership]", msg);
        resp.code = 403;
        resp.message = "Internal Server Error";
        resp.error = msg;
        return resp;
      }

      // 3️⃣ Update membership expired
      const result =
        await this.userMembershipRepository.updateUserMembershipExpired(
          req.user_ids,
          adminData.username
        );

      if (!result) {
        resp.code = 400;
        resp.message = "Failed to toggle user membership";
      }

      return resp;
    } catch (err: any) {
      console.error(
        "[service][ToggleInactiveUserMembership] error:",
        err.message
      );
      resp.code = 500;
      resp.message = "Internal Server Error";
      resp.error = err.message;
      return resp;
    }
  }

  /**
   * 🚀 StartMembershipProgram
   * (1:1 dari Go)
   */
  async startMembershipProgram(
    req: StartMembershipProgramRequest
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = {
      code: 200,
      message: "Success",
    };

    try {
      // 1️⃣ Ambil admin dari context
      const adminData = UserContext.get();
      if (!adminData || !adminData.username) {
        const msg = "failed to get admin data from context";
        console.error("[service][StartMembershipProgram]", msg);
        resp.code = 500;
        resp.message = "Internal Server Error";
        resp.error = msg;
        return resp;
      }

      // 2️⃣ Ambil config lama
      const currStartStr = await this.configRepository.getConfigByKey(
        "MEMBERSHIP_PROGRAM_START_DATE"
      );
      const currEndStr = await this.configRepository.getConfigByKey(
        "MEMBERSHIP_PROGRAM_END_DATE"
      );

      const today = new Date();
      today.setUTCHours(0, 0, 0, 0); // truncate time

      // 3️⃣ Cek apakah program sedang berjalan
      if (currStartStr && currEndStr) {
        const currStart = new Date(currStartStr);
        const currEnd = new Date(currEndStr);

        if (!isNaN(currStart.getTime()) && !isNaN(currEnd.getTime())) {
          if (today >= currStart && today <= currEnd) {
            const msg = "Membership program already running";
            resp.code = 400;
            resp.message = msg;
            resp.error = "bad request";
            return resp;
          }
        }
      }

      // 4️⃣ Validasi endDate
      if (!req.end_date) {
        const msg = "end_date is required";
        resp.code = 400;
        resp.message = msg;
        resp.error = msg;
        return resp;
      }

      const endDate = new Date(req.end_date);
      if (isNaN(endDate.getTime())) {
        const msg = "Invalid end_date format (use YYYY-MM-DD)";
        resp.code = 400;
        resp.message = msg;
        resp.error = msg;
        return resp;
      }

      const startDate = today.toISOString().split("T")[0];

      // 5️⃣ Update semua user pre-order jadi aktif
      const bulkResult =
        await this.userMembershipRepository.bulkUpdateMembershipPreOrderActivation();
      if (!bulkResult) {
        const msg = "Failed to activate pre-order memberships";
        resp.code = 500;
        resp.message = "Internal Server Error";
        resp.error = msg;
        return resp;
      }

      // 6️⃣ Update config start_date & end_date
      await this.configRepository.updateConfigByKey(
        "MEMBERSHIP_PROGRAM_START_DATE",
        startDate,
        adminData.username
      );

      await this.configRepository.updateConfigByKey(
        "MEMBERSHIP_PROGRAM_END_DATE",
        req.end_date,
        adminData.username
      );

      return resp;
    } catch (err: any) {
      console.error("[service][StartMembershipProgram] error:", err.message);
      resp.code = 500;
      resp.message = "Internal Server Error";
      resp.error = err.message;
      return resp;
    }
  }

  /**
   * 📅 GetMembershipProgramDate
   * Equivalent to: func (a *AdminSvcImpl) GetMembershipProgramDate(ctx context.Context)
   */
  async getMembershipProgramDate(): Promise<DefaultResponse> {
    const resp: DefaultResponse = {
      code: 200,
      message: "Success",
      data: {},
    };

    try {
      const startDate = await this.configRepository.getConfigByKey(
        "MEMBERSHIP_PROGRAM_START_DATE"
      );
      const endDate = await this.configRepository.getConfigByKey(
        "MEMBERSHIP_PROGRAM_END_DATE"
      );

      resp.data = {
        start_date: startDate,
        end_date: endDate,
      };

      return resp;
    } catch (err: any) {
      console.error("[service][GetMembershipProgramDate] error:", err.message);
      resp.code = 500;
      resp.message = "Internal Server Error";
      resp.error = err.message;
      return resp;
    }
  }

  /**
   * ✏️ Update Membership Program Date
   * (1:1 dari UpdateMembershipProgramDate di Go)
   */
  async updateMembershipProgramDate(
    req: StartMembershipProgramRequest
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = {
      code: 200,
      message: "Success",
    };

    try {
      // 🧩 1️⃣ Validasi end_date wajib
      if (!req.end_date) {
        const msg = "end_date is required";
        resp.code = 400;
        resp.message = msg;
        resp.error = msg;
        return resp;
      }

      // 🧠 2️⃣ Ambil admin data dari context
      const adminData = UserContext.get();
      if (!adminData || !adminData.username) {
        const msg = "failed to get admin data from context";
        console.error("[service][UpdateMembershipProgramDate]", msg);
        resp.code = 500;
        resp.message = "Internal Server Error";
        resp.error = msg;
        return resp;
      }

      // 📅 3️⃣ Validasi format tanggal (YYYY-MM-DD)
      const endDate = new Date(req.end_date);
      if (isNaN(endDate.getTime())) {
        const msg = "Invalid end_date format (use YYYY-MM-DD)";
        resp.code = 400;
        resp.message = msg;
        resp.error = msg;
        return resp;
      }

      // 🔍 4️⃣ Ambil start_date dari DB
      const startDateStr = await this.configRepository.getConfigByKey(
        "MEMBERSHIP_PROGRAM_START_DATE"
      );
      const startDate = new Date(startDateStr);

      if (isNaN(startDate.getTime())) {
        const msg = "Invalid start_date format in database";
        resp.code = 500;
        resp.message = msg;
        resp.error = msg;
        return resp;
      }

      // 🔒 5️⃣ Validasi end_date > start_date
      if (endDate <= startDate) {
        const msg = "end_date must be after start date";
        resp.code = 400;
        resp.message = msg;
        resp.error = msg;
        return resp;
      }

      // 💾 6️⃣ Update ke DB
      await this.configRepository.updateConfigByKey(
        "MEMBERSHIP_PROGRAM_END_DATE",
        req.end_date,
        adminData.username
      );

      return resp;
    } catch (err: any) {
      console.error(
        "[service][UpdateMembershipProgramDate] error:",
        err.message
      );
      resp.code = 500;
      resp.message = "Internal Server Error";
      resp.error = err.message;
      return resp;
    }
  }
}
