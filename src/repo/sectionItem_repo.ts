import { Pool } from "pg";

import db from "../db/db.config";
import { SectionItemQueries } from "./queries/sectionItem_queries";
import { SectionItem } from "../models/sectionItems.model";

export class SectionItemRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- GET ITEMS BY SECTION ID ---
  async getBySectionID(section_id: string): Promise<SectionItem[]> {
    try {
      const result = await this.pg.query(SectionItemQueries.GetItemsBySection, [
        section_id,
      ]);

      return result.rows.map((row) => {
        let extra_data: Record<string, any> = {};
        if (row.extra_data) {
          try {
            extra_data = JSON.parse(row.extra_data);
          } catch (err) {
            console.error(
              `[SectionItemRepository][getBySectionID] JSON parse error: ${err}`
            );
          }
        }

        return {
          id: row.id,
          section_id: row.section_id,
          title: row.title,
          subtitle: row.subtitle ?? "",
          subjek: row.subjek ?? "",
          image_url: row.image_url ?? "",
          icon_url: row.icon_url ?? "",
          order: row.order ?? 0,
          extra_data,
          created_at: row.created_at,
          updated_at: row.updated_at,
        };
      });
    } catch (err: any) {
      console.error(
        `[SectionItemRepository][getBySectionID] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- GET ALL SECTION ITEMS ---
  async getAll(): Promise<SectionItem[]> {
    try {
      const result = await this.pg.query(SectionItemQueries.GetAllSectionItems);

      return result.rows.map((row) => {
        let extra_data: Record<string, any> = {};
        if (row.extra_data) {
          try {
            extra_data = JSON.parse(row.extra_data);
          } catch (err) {
            console.error(
              `[SectionItemRepository][getAll] JSON parse error: ${err}`
            );
          }
        }

        return {
          id: row.id,
          section_id: row.section_id,
          title: row.title,
          subtitle: row.subtitle ?? "",
          subjek: row.subjek ?? "",
          image_url: row.image_url ?? "",
          icon_url: row.icon_url ?? "",
          order: row.order ?? 0,
          extra_data,
          created_at: row.created_at,
          updated_at: row.updated_at,
        };
      });
    } catch (err: any) {
      console.error(`[SectionItemRepository][getAll] error: ${err.message}`);
      throw err;
    }
  }

  // --- CREATE SECTION ITEM ---
  async create(item: SectionItem): Promise<SectionItem> {
    try {
      const extra_data = JSON.stringify(item.extra_data || {});

      const result = await this.pg.query(SectionItemQueries.InsertSectionItem, [
        item.section_id,
        item.title,
        item.subtitle,
        item.subjek,
        item.image_url,
        item.icon_url,
        item.order,
        extra_data,
      ]);

      const row = result.rows[0];

      return {
        id: row.id,
        section_id: row.section_id,
        title: row.title,
        subtitle: row.subtitle ?? "",
        subjek: row.subjek ?? "",
        image_url: row.image_url ?? "",
        icon_url: row.icon_url ?? "",
        order: row.order ?? 0,
        extra_data: JSON.parse(row.extra_data ?? "{}"),
        created_at: row.created_at,
        updated_at: row.updated_at,
      };
    } catch (err: any) {
      console.error(`[SectionItemRepository][create] error: ${err.message}`);
      throw err;
    }
  }

  // --- UPDATE SECTION ITEM ---
  async update(item: SectionItem): Promise<boolean> {
    try {
      const extra_data = JSON.stringify(item.extra_data || {});

      const result = await this.pg.query(SectionItemQueries.UpdateSectionItem, [
        item.title,
        item.subtitle,
        item.subjek,
        item.image_url,
        item.icon_url,
        item.order,
        extra_data,
        item.id,
      ]);

      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[SectionItemRepository][update] error: ${err.message}`);
      throw err;
    }
  }

  // --- DELETE SECTION ITEM ---
  async delete(id: string): Promise<boolean> {
    try {
      const result = await this.pg.query(SectionItemQueries.DeleteSectionItem, [
        id,
      ]);
      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[SectionItemRepository][delete] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET SECTION ITEM BY ID ---
  async getByID(id: string): Promise<SectionItem | null> {
    try {
      const result = await this.pg.query(
        SectionItemQueries.GetSectionItemByID,
        [id]
      );
      if (result.rowCount === 0) return null;

      const row = result.rows[0];
      let extra_data: Record<string, any> = {};
      if (row.extra_data) {
        try {
          extra_data = JSON.parse(row.extra_data);
        } catch (err) {
          console.error(
            `[SectionItemRepository][getByID] JSON parse error: ${err}`
          );
        }
      }

      return {
        id: row.id,
        section_id: row.section_id,
        title: row.title,
        subtitle: row.subtitle ?? "",
        subjek: row.subjek ?? "",
        image_url: row.image_url ?? "",
        icon_url: row.icon_url ?? "",
        order: row.order ?? 0,
        extra_data,
        created_at: row.created_at,
        updated_at: row.updated_at,
      };
    } catch (err: any) {
      console.error(`[SectionItemRepository][getByID] error: ${err.message}`);
      throw err;
    }
  }

  // --- CHECK DUPLICATE TITLE ---
  async existsBySectionAndTitle(
    section_id: string,
    title: string
  ): Promise<boolean> {
    try {
      const result = await this.pg.query(
        SectionItemQueries.CheckDuplicateSectionItemTitle,
        [section_id, title]
      );

      return result.rows[0]?.exists ?? false;
    } catch (err: any) {
      console.error(
        `[SectionItemRepository][existsBySectionAndTitle] error: ${err.message}`
      );
      throw err;
    }
  }
}
