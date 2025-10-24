import { Pool, QueryResult } from "pg";

import db from "../db/db.config";
import { SectionQueries } from "./queries/section_queries";
import { SectionWithItems } from "../DTO/sectionItem.dto";
import { Section } from "../models/section.model";
import { SectionItem } from "../models/sectionItems.model";

interface SectionWithItemRow {
  section_id: string;
  section_slug: string;
  section_type: string;
  section_title: string;
  section_subtitle: string;
  section_order: number;
  section_created_at: string;
  section_updated_at: string;
  item_id?: string | null;
  item_title?: string | null;
  item_subtitle?: string | null;
  item_subjek?: string | null;
  image_url?: string | null;
  icon_url?: string | null;
  item_order?: number | null;
  extra_data?: string | null;
  item_created_at?: string | null;
  item_updated_at?: string | null;
}

export class SectionRepository {
  private pg: Pool;

  constructor() {
    this.pg = db;
  }

  // --- SAVE SECTION ---
  async save(section: Section): Promise<boolean> {
    try {
      const result = await this.pg.query(SectionQueries.InsertSection, [
        section.slug,
        section.type,
        section.title,
        section.subtitle,
        section.order,
      ]);
      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[SectionRepository][save] error: ${err.message}`);
      throw err;
    }
  }

  // --- UPDATE SECTION ---
  async update(section: Section): Promise<boolean> {
    try {
      const result = await this.pg.query(SectionQueries.UpdateSection, [
        section.slug,
        section.type,
        section.title,
        section.subtitle,
        section.order,
        section.id,
      ]);
      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[SectionRepository][update] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET BY SLUG ---
  async getBySlug(slug: string): Promise<Section | null> {
    try {
      const result = await this.pg.query(
        `
        SELECT id, slug, type, title, subtitle, "order", created_at, updated_at
        FROM sections
        WHERE slug = $1
        `,
        [slug]
      );

      if (result.rowCount === 0) return null;
      const row = result.rows[0];

      return {
        id: row.id,
        slug: row.slug,
        type: row.type,
        title: row.title,
        subtitle: row.subtitle,
        order: row.order,
        created_at: row.created_at,
        updated_at: row.updated_at,
      };
    } catch (err: any) {
      console.error(`[SectionRepository][getBySlug] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET BY TYPE ---
  async getByType(sectionType: string): Promise<Section | null> {
    try {
      const result = await this.pg.query(SectionQueries.GetSectionByType, [
        sectionType,
      ]);
      if (result.rowCount === 0) return null;
      const row = result.rows[0];

      return {
        id: row.id,
        slug: row.slug,
        type: row.type,
        title: row.title,
        subtitle: row.subtitle,
        order: row.order,
        created_at: row.created_at,
        updated_at: row.updated_at,
      };
    } catch (err: any) {
      console.error(`[SectionRepository][getByType] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET ALL SECTIONS WITH ITEMS ---
  async getAllSectionsWithItems(): Promise<SectionWithItems[]> {
    try {
      const result: QueryResult<SectionWithItemRow> = await this.pg.query(
        SectionQueries.GetAllSectionsWithItems
      );

      const sectionMap: Record<string, SectionWithItems> = {};

      result.rows.forEach((row) => {
        // 🔹 Inisialisasi Section jika belum ada di map
        if (!sectionMap[row.section_id]) {
          sectionMap[row.section_id] = {
            id: row.section_id,
            slug: row.section_slug,
            type: row.section_type,
            title: row.section_title,
            subtitle: row.section_subtitle,
            order: row.section_order,
            created_at: row.section_created_at,
            updated_at: row.section_updated_at,
            items: [],
          };
        }

        // 🔹 Jika ada item di section tersebut
        if (row.item_id) {
          let extra_data: Record<string, any> = {};
          if (row.extra_data) {
            try {
              extra_data = JSON.parse(row.extra_data);
            } catch (err) {
              console.error(
                `[SectionRepository][getAllSectionsWithItems] JSON parse error: ${err}`
              );
            }
          }

          const item: SectionItem = {
            id: row.item_id,
            section_id: row.section_id,
            title: row.item_title ?? "",
            subtitle: row.item_subtitle ?? "",
            subjek: row.item_subjek ?? "",
            image_url: row.image_url ?? "",
            icon_url: row.icon_url ?? "",
            order: row.item_order ?? 0,
            extra_data,
            created_at: row.item_created_at ?? new Date().toISOString(),
            updated_at: row.item_updated_at ?? new Date().toISOString(),
          };

          sectionMap[row.section_id].items.push(item);
        }
      });

      return Object.values(sectionMap);
    } catch (err: any) {
      console.error(
        `[SectionRepository][getAllSectionsWithItems] error: ${err.message}`
      );
      throw err;
    }
  }

  // --- FIND ONE BY ID ---
  async findOneByID(id: string): Promise<Section | null> {
    try {
      const result = await this.pg.query(SectionQueries.FindSectionByID, [id]);
      if (result.rowCount === 0) return null;
      const row = result.rows[0];

      return {
        id: row.id,
        title: row.title,
        subtitle: row.subtitle,
        created_at: row.created_at,
        updated_at: row.updated_at,
      };
    } catch (err: any) {
      console.error(`[SectionRepository][findOneByID] error: ${err.message}`);
      throw err;
    }
  }

  // --- GET BY ID ---
  async getByID(id: string): Promise<Section | null> {
    try {
      const result = await this.pg.query(SectionQueries.GetSectionByID, [id]);
      if (result.rowCount === 0) return null;
      const row = result.rows[0];

      return {
        id: row.id,
        slug: row.slug,
        type: row.type,
        title: row.title,
        subtitle: row.subtitle,
        order: row.order,
        created_at: row.created_at,
        updated_at: row.updated_at,
      };
    } catch (err: any) {
      console.error(`[SectionRepository][getByID] error: ${err.message}`);
      throw err;
    }
  }

  // --- DELETE SECTION ---
  async delete(id: string): Promise<boolean> {
    try {
      const result = await this.pg.query(SectionQueries.DeleteSection, [id]);
      return (result.rowCount ?? 0) > 0;
    } catch (err: any) {
      console.error(`[SectionRepository][delete] error: ${err.message}`);
      throw err;
    }
  }
}
