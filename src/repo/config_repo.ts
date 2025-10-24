import { Pool } from "pg";
import db from "../db/db.config";

export class ConfigRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  /**
   * 🧩 GetConfigByKey
   * Sama seperti func (c *ConfigRepoImpl) GetConfigByKey(key string)
   */
  async getConfigByKey(key: string): Promise<string> {
    try {
      const result = await this.pg.query(
        "SELECT value FROM config WHERE key = $1",
        [key]
      );
      if (result.rowCount === 0) {
        throw new Error(`Config not found for key: ${key}`);
      }
      return result.rows[0].value;
    } catch (err: any) {
      console.error(`[repo][config][getConfigByKey] error: ${err.message}`);
      throw err;
    }
  }

  /**
   * 🧩 UpdateConfigByKey
   * Sama seperti func (c *ConfigRepoImpl) UpdateConfigByKey(key, value, operatedBy)
   */
  async updateConfigByKey(
    key: string,
    value: string,
    operatedBy: string
  ): Promise<void> {
    try {
      await this.pg.query(
        "UPDATE config SET value = $1, updated_at = NOW(), updated_by = $3 WHERE key = $2",
        [value, key, operatedBy]
      );
    } catch (err: any) {
      console.error(`[repo][config][updateConfigByKey] error: ${err.message}`);
      throw err;
    }
  }
}
