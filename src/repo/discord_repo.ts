import { Pool } from "pg";

import { DiscordAccount } from "../models/discord.model";
import { DiscordQueries } from "./queries/discord_queries";
import db from "../db/db.config";

export class DiscordAccountRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- CREATE DISCORD ACCOUNT ---
  async createDiscordAccount(req: DiscordAccount): Promise<number> {
    try {
      const result = await this.pg.query(DiscordQueries.CreateDiscordAccount, [
        req.userId,
        req.discordAccountId,
        req.username,
        req.email,
      ]);

      if (result.rowCount === 0) {
        throw new Error("failed to create discord account");
      }

      return result.rows[0].id as number;
    } catch (err: any) {
      console.error(
        `[repo][discord][createDiscordAccount] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- GET DISCORD ACCOUNT BY USER ID ---
  async getDiscordAccountByUserID(
    userID: number
  ): Promise<DiscordAccount | null> {
    try {
      const result = await this.pg.query(
        DiscordQueries.GetDiscordAccountByUserID,
        [userID]
      );

      if (result.rowCount === 0) {
        return null; // sama dengan return kosong di Go
      }

      return result.rows[0] as DiscordAccount;
    } catch (err: any) {
      console.error(
        `[repo][discord][getDiscordAccountByUserID] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- DELETE DISCORD ACCOUNT BY USER ID (soft delete) ---
  async deleteDiscordAccountByUserID(userID: number): Promise<boolean> {
    try {
      const result = await this.pg.query(
        DiscordQueries.SoftDeleteDiscordAccount,
        [userID]
      );

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(
        `[repo][discord][deleteDiscordAccountByUserID] error: ${err.message}`
      );
      throw err;
    }
  }
}
