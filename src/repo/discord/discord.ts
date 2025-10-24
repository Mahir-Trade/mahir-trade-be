// src/repos/DiscordRepo.ts
import { Client, GatewayIntentBits } from "discord.js";
import axios from "axios";
import {
  DiscordRoleRequest,
  DiscordUser,
  GetDiscordUserRequest,
} from "@dto/discord.dto";

export interface TokenResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  scope: string;
  refresh_token: string;
  error?: string;
  error_description?: string;
}

export interface IDiscordRepo {
  addRoleToMember(req: DiscordRoleRequest): Promise<void>;
  removeRoleFromMember(req: DiscordRoleRequest): Promise<void>;
  getDiscordUser(req: GetDiscordUserRequest): Promise<DiscordUser>;
  exchangeCodeForToken(
    code: string,
    redirectURI: string
  ): Promise<TokenResponse>;
  getUserDataByAccessToken(accessToken: string): Promise<DiscordUser>;
  inviteUserToGuild(
    userId: string,
    guildId: string,
    accessToken: string
  ): Promise<void>;
}

export class DiscordRepository {
  private token: string;
  private guildId: string;
  private roleId: string;
  private baseUrl: string;
  private clientId: string;
  private clientSecret: string;

  constructor() {
    this.token = process.env.DISCORD_TOKEN || "";
    this.guildId = process.env.DISCORD_GUILD_ID || "";
    this.roleId = process.env.DISCORD_ROLE_ID || "";
    this.baseUrl = process.env.DISCORD_BASE_URL || "https://discord.com/api";
    this.clientId = process.env.DISCORD_CLIENT_ID || "";
    this.clientSecret = process.env.DISCORD_CLIENT_SECRET || "";
  }

  private createClient() {
    return new Client({
      intents: [GatewayIntentBits.Guilds, GatewayIntentBits.GuildMembers],
    });
  }

  async addRoleToMember(req: DiscordRoleRequest): Promise<void> {
    const client = this.createClient();
    await client.login(this.token);

    const guild = await client.guilds.fetch(this.guildId);
    const member = await guild.members.fetch(req.user_id);

    await member.roles.add(this.roleId);

    client.destroy();
  }

  async removeRoleFromMember(req: DiscordRoleRequest): Promise<void> {
    const client = this.createClient();
    await client.login(this.token);

    const guild = await client.guilds.fetch(this.guildId);
    const member = await guild.members.fetch(req.user_id);

    await member.roles.remove(this.roleId);

    client.destroy();
  }

  async getDiscordUser(req: GetDiscordUserRequest): Promise<DiscordUser> {
    const client = this.createClient();
    await client.login(this.token);

    const guild = await client.guilds.fetch(this.guildId);
    const members = await guild.members.search({
      query: req.username,
      limit: 1,
    });

    client.destroy();

    if (members.size === 0) {
      throw new Error("user not found");
    }

    const user = members.first()!.user;
    return {
      id: user.id,
      email: (user as any).email || "",
      username: user.username,
    };
  }

  async exchangeCodeForToken(
    code: string,
    redirectURI: string
  ): Promise<TokenResponse> {
    const tokenURL = `${this.baseUrl}/v10/oauth2/token`;

    const params = new URLSearchParams();
    params.append("grant_type", "authorization_code");
    params.append("code", code);
    params.append("redirect_uri", redirectURI);

    const res = await axios.post<TokenResponse>(tokenURL, params, {
      auth: { username: this.clientId, password: this.clientSecret },
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
    });

    return res.data;
  }

  async getUserDataByAccessToken(accessToken: string): Promise<DiscordUser> {
    const userURL = `${this.baseUrl}/users/@me`;

    const res = await axios.get(userURL, {
      headers: { Authorization: `Bearer ${accessToken}` },
    });

    return res.data;
  }

  async inviteUserToGuild(
    userId: string,
    guildId: string,
    accessToken: string
  ): Promise<void> {
    const inviteGuildURL = `${this.baseUrl}/guilds/${guildId}/members/${userId}`;

    await axios.put(
      inviteGuildURL,
      { access_token: accessToken },
      {
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bot ${this.token}`,
        },
      }
    );
  }
}
