import * as dotenv from "dotenv";
import { DefaultResponse } from "../utils/response";
import { UserRepository } from "../repo/user_repo";
import {
  UserRegistrationRequest,
  LoginReq,
  GoogleLoginReq,
  ResetPasswordRequest,
} from "../DTO/user.dto";

import { DiscordRoleRequest, DiscordUser } from "../DTO/discord.dto";
import { GoogleCallbackRequest } from "../DTO/google.dto";
import { User } from "../models/user.model";
import {
  generateRandomPassword,
  hashPassword,
  verifyPassword,
} from "../utils/bcrypt";
import { signJWT } from "../utils/jwt";
import { decryptAES256CBC, encryptAES256CBC } from "../utils/security";
import { mappingValuesToTemplate } from "../utils/email";
import { DiscordAccountRepository } from "../repo/discord_repo";
import { DiscordRepository } from "../repo/discord/discord";
import { GoogleRepository } from "../repo/google/google";
import { UserMembershipRepository } from "../repo/userMembership_repo";
import { EmailRepository } from "../repo/email_repo";
import { SendgridRepository } from "../repo/sendgrid/sendgrid";

dotenv.config();

export class UserService {
  private userRepo: UserRepository;
  private discordRepo: DiscordRepository;
  private discordAccountRepo: DiscordAccountRepository;
  private googleRepo: GoogleRepository;
  private userMembershipRepo: UserMembershipRepository;
  private emailTemplateRepo: EmailRepository;
  private sendgridRepo: SendgridRepository;

  constructor() {
    this.userRepo = new UserRepository();
    this.discordRepo = new DiscordRepository();
    this.discordAccountRepo = new DiscordAccountRepository();
    this.googleRepo = new GoogleRepository();
    this.userMembershipRepo = new UserMembershipRepository();
    this.emailTemplateRepo = new EmailRepository();
    this.sendgridRepo = new SendgridRepository();
  }

  async userRegistration(
    data: UserRegistrationRequest
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 201, message: "success" };
    if (data.password !== data.password_confirmation)
      throw new Error("password and password confirmation must be same");

    const email = data.email.trim().toLowerCase();
    const phone_number = data.phone_number?.trim();
    const username = data.username.trim().toLowerCase();
    const password = await hashPassword(data.password);

    console.log("email:", email);
    console.log("phone_number:", phone_number);
    console.log("username: ", username);
    console.log("password: ", password);

    const userId = await this.userRepo.createUser({
      email,
      phone_number,
      username,
      password,
    });

    console.log("Data: ", userId);

    const { token, exp } = await signJWT({
      email,
      user_id: userId,
      username,
    });

    resp.data = { token, expire: exp };

    console.log("Resp Data: ", resp.data);
    console.log("Resp: ", resp);

    return resp;
  }

  async userLogin(data: LoginReq): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "success", data: {} };

    console.log("data identity:", data.identity);

    const identity = data.identity.trim().toLowerCase();

    console.log("Identity login: ", identity);

    const user = await this.userRepo.findUserByEmailOrUsername(identity);

    console.log("user: ", user);

    if (!user || !(await verifyPassword(data.password, user.password)))
      throw new Error("invalid email or password");

    const { token, exp } = await signJWT({
      email: user.email,
      user_id: user.user_id,
      username: user.username,
    });
    resp.data = { token, expire: exp };
    return resp;
  }

  async assignRoleDiscordToUser(data: {
    userId: number;
    code: string;
  }): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "success", data: {} };
    const redirectURI = process.env.DISCORD_REDIRECT_URI_ASSIGN_ROLE!;
    const result = await this.inviteDiscordUserToGuild({
      code: data.code,
      redirectURI,
    });

    const membership = await this.userMembershipRepo.getUserMembershipByUserID(
      data.userId
    );
    if (!membership || !membership.is_membership_active)
      throw new Error("user membership is not active");

    const existing = await this.discordAccountRepo.getDiscordAccountByUserID(
      data.userId
    );
    if (existing?.id) throw new Error("discord account already registered");

    const discordUser = result.data as DiscordUser;

    await this.discordRepo.addRoleToMember({
      user_id: discordUser.id,
    } as DiscordRoleRequest);

    await this.discordAccountRepo.createDiscordAccount({
      userId: data.userId,
      discordAccountId: discordUser.id, // ⬅ pakai ID (huruf besar)
      username: discordUser.username, // ⬅ pakai Username (huruf besar)
      email: discordUser.email, // ⬅ pakai Email (huruf besar)
    });

    return result;
  }

  async removeRoleDiscordToUser(data: {
    userId: number;
  }): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "success", data: {} };
    const discordAccount =
      await this.discordAccountRepo.getDiscordAccountByUserID(data.userId);
    if (!discordAccount?.id) throw new Error("discord account not found");

    await this.discordRepo.removeRoleFromMember({
      user_id: discordAccount.discordAccountId,
    });
    await this.discordAccountRepo.deleteDiscordAccountByUserID(
      discordAccount.id
    );
    resp.data = discordAccount;
    return resp;
  }

  async inviteDiscordUserToGuild(data: {
    code: string;
    redirectURI: string;
  }): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "success", data: {} };
    const tokenResp = await this.discordRepo.exchangeCodeForToken(
      data.code,
      data.redirectURI
    );
    if (tokenResp.error) throw new Error(tokenResp.error_description);

    const discordUser = await this.discordRepo.getUserDataByAccessToken(
      tokenResp.access_token
    );
    await this.discordRepo.inviteUserToGuild(
      discordUser.id,
      process.env.DISCORD_GUILD_ID!,
      tokenResp.access_token
    );
    resp.data = discordUser;
    return resp;
  }

  async connectDiscordAccountAndAssignRole(data: {
    code: string;
  }): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "success", data: {} };
    const tokenResp = await this.discordRepo.exchangeCodeForToken(
      data.code,
      process.env.DISCORD_REDIRECT_URI_ASSIGN_ROLE!
    );
    if (tokenResp.error) throw new Error(tokenResp.error_description);

    const discordUser = await this.discordRepo.getUserDataByAccessToken(
      tokenResp.access_token
    );
    await this.discordRepo.inviteUserToGuild(
      discordUser.id,
      process.env.DISCORD_GUILD_ID!,
      tokenResp.access_token
    );
    await this.discordRepo.addRoleToMember({ user_id: discordUser.id });
    return resp;
  }

  async connectDiscordAccountAndRemoveRole(data: {
    code: string;
  }): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "success", data: {} };
    const tokenResp = await this.discordRepo.exchangeCodeForToken(
      data.code,
      process.env.DISCORD_REDIRECT_URI_REMOVE_ROLE!
    );
    if (tokenResp.error) throw new Error(tokenResp.error_description);

    const discordUser = await this.discordRepo.getUserDataByAccessToken(
      tokenResp.access_token
    );
    await this.discordRepo.inviteUserToGuild(
      discordUser.id,
      process.env.DISCORD_GUILD_ID!,
      tokenResp.access_token
    );
    await this.discordRepo.removeRoleFromMember({ user_id: discordUser.id });
    return resp;
  }

  async loginWithGoogle(): Promise<string> {
    return await this.googleRepo.login();
  }

  async callbackGoogle(
    data: GoogleLoginReq
  ): Promise<{ token: string; expire: Date }> {
    const decodedCode = decodeURIComponent(data.code);

    const userInfo = await this.googleRepo.callback({
      state: data.state,
      code: decodedCode,
    } as GoogleCallbackRequest);

    const email = userInfo.email.trim().toLowerCase();
    const username = userInfo.name.trim().toLowerCase();

    let user = await this.userRepo.findUserByEmailOrUsername(email);

    if (user) {
      const { token, exp } = await signJWT({
        email: user.email,
        user_id: user.user_id,
        username: user.username,
      });

      return { token, expire: exp }; // ✅ rename exp → expire
    }

    const passwordHash = await hashPassword(generateRandomPassword(12));
    const newUser: User = { email, username, password: passwordHash };
    const userId = await this.userRepo.createUser(newUser);

    const { token, exp } = await signJWT({
      email,
      user_id: userId,
      username,
    });

    return { token, expire: exp }; // ✅ sama juga
  }

  async getDetailUser(data: { userId: number }): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "success" };
    const user = await this.userRepo.getUserById(data.userId);
    if (!user) throw new Error("user not found");

    const membership = await this.userMembershipRepo.getUserMembershipByUserID(
      data.userId
    );
    const discordAccount =
      await this.discordAccountRepo.getDiscordAccountByUserID(data.userId);

    resp.data = {
      userId: user.user_id,
      uuid: user.uuid,
      email: user.email,
      phoneNumber: user.phone_number,
      username: user.username,
      isActive: user.is_active,
      isMembershipActive: membership?.is_membership_active || false,
      discordUsername: discordAccount?.username || undefined,
    };
    return resp;
  }

  async getDetailUserForBO(userId: number): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "success" };
    const user = await this.userRepo.getUserById(userId);
    if (!user) throw new Error("user not found");
    user.password = "";
    resp.data = user;
    return resp;
  }

  async updateMembership(): Promise<void> {
    await this.userMembershipRepo.updateBulkUserMembership();
  }

  async forgotPasswordUser(data: { email: string }): Promise<DefaultResponse> {
    const resp: DefaultResponse = {
      code: 200,
      message: "please check your email",
    };
    const email = data.email.trim().toLowerCase();
    const user = await this.userRepo.findUserByEmailOrUsername(email);
    if (!user) return resp;

    const encryptedEmail = encodeURIComponent(
      encryptAES256CBC(
        email,
        process.env.JWT_ENCRYPT_KEY!,
        process.env.JWT_ENCRYPT_IV!
      )
    );
    const encryptedExpiration = encodeURIComponent(
      encryptAES256CBC(
        new Date(Date.now() + 4 * 3600 * 1000).toISOString(),
        process.env.JWT_ENCRYPT_KEY!,
        process.env.JWT_ENCRYPT_IV!
      )
    );

    const verifyLink = `${process.env.FORGOT_PASSWORD_URL}?q1=${encryptedEmail}&q2=${encryptedExpiration}`;
    const emailTemplate = await this.emailTemplateRepo.getByKey(
      "forgot_password"
    );
    const parsedHtml = mappingValuesToTemplate(
      { username: user.username, verifyLink },
      emailTemplate
    );

    await this.sendgridRepo.sendEmail({
      from: process.env.SENDGRID_SENDER_EMAIL!,
      to: email,
      subject: "Mahir Trade Password Recovery",
      body: parsedHtml,
      senderName: process.env.SENDGRID_SENDER_NAME!,
      recepientName: user.username,
    });
    return resp;
  }

  async resetPasswordUser(
    data: ResetPasswordRequest
  ): Promise<DefaultResponse> {
    const resp: DefaultResponse = { code: 200, message: "success" };
    const email = decryptAES256CBC(
      decodeURIComponent(data.q1),
      process.env.JWT_ENCRYPT_KEY!,
      process.env.JWT_ENCRYPT_IV!
    );
    const expiration = decryptAES256CBC(
      decodeURIComponent(data.q2),
      process.env.JWT_ENCRYPT_KEY!,
      process.env.JWT_ENCRYPT_IV!
    );

    const expTime = new Date(expiration);
    if (new Date() > expTime) throw new Error("token expired");

    const user = await this.userRepo.findUserByEmailOrUsername(email);
    if (!user) throw new Error("illegal token");
    if (data.password !== data.password_confirmation)
      throw new Error("password and confirmation mismatch");

    if (!user.user_id) {
      throw new Error("User ID is missing");
    }

    const passwordHash = await hashPassword(data.password);
    const updated = await this.userRepo.updatePassword(
      user.user_id,
      passwordHash,
      user.username
    );
    if (!updated) throw new Error("failed to update password");
    return resp;
  }
}
