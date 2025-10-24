import { Request, Response } from "express";
import {
  LoginReq,
  GoogleLoginReq,
  ForgotPasswordReq,
  ResetPasswordRequest,
  UserRegistrationRequest,
} from "../DTO/user.dto";
import { UserService } from "../services/user.service";
import { Console } from "console";

export class AuthController {
  private userService: UserService;

  constructor() {
    this.userService = new UserService();

    this.userRegistration = this.userRegistration.bind(this);
    this.userLogin = this.userLogin.bind(this);
    this.loginWithGoogle = this.loginWithGoogle.bind(this);
    this.callbackGoogle = this.callbackGoogle.bind(this);
    this.assignRoleDiscordToUser = this.assignRoleDiscordToUser.bind(this);
    this.removeRoleDiscordUser = this.removeRoleDiscordUser.bind(this);
    this.inviteDiscordUserToGuild = this.inviteDiscordUserToGuild.bind(this);
    this.connectDiscordAccountAndAssignRole =
      this.connectDiscordAccountAndAssignRole.bind(this);
    this.connectDiscordAccountAndRemoveRole =
      this.connectDiscordAccountAndRemoveRole.bind(this);
    this.getDetailUser = this.getDetailUser.bind(this);
    this.forgotPassword = this.forgotPassword.bind(this);
    this.requestResetPassword = this.requestResetPassword.bind(this);
  }

  // POST /auth/register
  async userRegistration(req: Request, res: Response): Promise<Response> {
    try {
      const data = req.body as UserRegistrationRequest;
      const result = await this.userService.userRegistration({ ...data });

      console.log("Result: ", result);

      return res.status(result.code).json(result);
    } catch (err: any) {
      return res.status(400).json({
        code: 400,
        message: "Invalid request body",
        error: err.message,
      });
    }
  }

  // POST /auth/login
  async userLogin(req: Request, res: Response): Promise<Response> {
    try {
      const data = req.body as LoginReq;

      console.log("Data user login: ", data);

      const result = await this.userService.userLogin({ ...data });

      console.log("Result Login: ", result);

      return res.status(result.code).json(result);
    } catch (err: any) {
      return res.status(400).json({
        code: 400,
        message: "Invalid request body",
        error: err.message,
      });
    }
  }

  // GET /auth/google
  async loginWithGoogle(req: Request, res: Response): Promise<void> {
    try {
      const url = await this.userService.loginWithGoogle();
      res.redirect(307, url);
    } catch (err: any) {
      res
        .status(400)
        .json({ code: 400, message: "bad request", error: err.message });
    }
  }

  // GET /auth/google/callback
  async callbackGoogle(req: Request, res: Response): Promise<any> {
    try {
      const { state, code } = req.query as { state: string; code: string };
      if (!state || !code) {
        return res.status(400).json({
          code: 400,
          message: "bad request",
          error: "state or code is empty",
        });
      }

      const data = { state, code } as GoogleLoginReq;
      const result = await this.userService.callbackGoogle({ ...data });

      const baseUrl = process.env.GOOGLE_FRONTEND_REDIRECT_URL || "";
      const redirectUrl = `${baseUrl}?token=${result.token}&expiredAt=${result.expire}`;

      res.redirect(307, redirectUrl);
    } catch (err: any) {
      res
        .status(502)
        .json({ code: 502, message: "bad gateway", error: err.message });
    }
  }

  // GET /auth/discord/assign-role?code=xxx
  async assignRoleDiscordToUser(
    req: Request,
    res: Response
  ): Promise<Response> {
    try {
      const code = req.query.code as string;
      const { userId } = req.body; // ambil userId dari body
      if (!userId || !code) throw new Error("userId or code is missing");

      const result = await this.userService.assignRoleDiscordToUser({
        userId,
        code,
      });

      return res.status(result.code).json(result);
    } catch (err: any) {
      return res.status(400).json({
        code: 400,
        message: "bad request",
        error: err.message,
      });
    }
  }

  // DELETE /auth/discord/remove-role
  async removeRoleDiscordUser(req: Request, res: Response): Promise<Response> {
    try {
      const { userId } = req.body;
      if (!userId) throw new Error("userId is required");

      const result = await this.userService.removeRoleDiscordToUser({ userId });
      return res.status(result.code).json(result);
    } catch (err: any) {
      return res.status(400).json({
        code: 400,
        message: "bad request",
        error: err.message,
      });
    }
  }

  // GET /auth/discord/invite?code=xxx
  async inviteDiscordUserToGuild(
    req: Request,
    res: Response
  ): Promise<Response> {
    try {
      const code = req.query.code as string;
      const redirectURI = process.env.DISCORD_REDIRECT_URI!;
      const result = await this.userService.inviteDiscordUserToGuild({
        code,
        redirectURI,
      });
      return res.status(result.code).json(result);
    } catch (err: any) {
      return res.status(400).json({
        code: 400,
        message: "bad request",
        error: err.message,
      });
    }
  }

  // GET /auth/discord/connect-assign?code=xxx
  async connectDiscordAccountAndAssignRole(
    req: Request,
    res: Response
  ): Promise<Response> {
    try {
      const code = req.query.code as string;
      const result = await this.userService.connectDiscordAccountAndAssignRole({
        code,
      });
      return res.status(result.code).json(result);
    } catch (err: any) {
      return res.status(400).json({
        code: 400,
        message: "bad request",
        error: err.message,
      });
    }
  }

  // GET /auth/discord/connect-remove?code=xxx
  async connectDiscordAccountAndRemoveRole(
    req: Request,
    res: Response
  ): Promise<Response> {
    try {
      const code = req.query.code as string;
      const result = await this.userService.connectDiscordAccountAndRemoveRole({
        code,
      });
      return res.status(result.code).json(result);
    } catch (err: any) {
      return res.status(400).json({
        code: 400,
        message: "bad request",
        error: err.message,
      });
    }
  }

  // GET /auth/me
  async getDetailUser(req: Request, res: Response): Promise<Response> {
    try {
      const { userId } = req.body; // ambil dari body atau token decode
      if (!userId) throw new Error("userId is required");
      const result = await this.userService.getDetailUser({ userId });
      return res.status(result.code).json(result);
    } catch (err: any) {
      return res.status(400).json({
        code: 400,
        message: "bad request",
        error: err.message,
      });
    }
  }

  // POST /auth/forgot-password
  async forgotPassword(req: Request, res: Response): Promise<Response> {
    try {
      const data = req.body as ForgotPasswordReq;
      const result = await this.userService.forgotPasswordUser({
        email: data.email,
      });
      return res.status(result.code).json(result);
    } catch (err: any) {
      return res.status(400).json({
        code: 400,
        message: "Invalid request body",
        error: err.message,
      });
    }
  }

  // POST /auth/reset-password
  async requestResetPassword(req: Request, res: Response): Promise<Response> {
    try {
      const data = req.body as ResetPasswordRequest;
      const result = await this.userService.resetPasswordUser({ ...data });
      return res.status(result.code).json(result);
    } catch (err: any) {
      return res.status(400).json({
        code: 400,
        message: "Invalid request body",
        error: err.message,
      });
    }
  }
}
