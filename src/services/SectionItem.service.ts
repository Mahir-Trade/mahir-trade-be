import {
  CreateSectionItemRequest,
  UpdateSectionItemRequest,
} from "../DTO/sectionItem.dto";
import { SectionItem } from "../models/sectionItems.model";
import { SectionRepository } from "../repo/section_repo";
import { SectionItemRepository } from "../repo/sectionItem_repo";

export class SectionItemService {
  private sectionItemRepo: SectionItemRepository;
  private sectionRepo: SectionRepository;

  constructor() {
    this.sectionItemRepo = new SectionItemRepository();
    this.sectionRepo = new SectionRepository();
  }

  // --- GET ALL ITEMS ---
  async getItems(): Promise<SectionItem[]> {
    return this.sectionItemRepo.getAll();
  }

  // --- GET ITEMS BY SECTION ID ---
  async getItemsBySectionID(sectionId: string): Promise<SectionItem[]> {
    return this.sectionItemRepo.getBySectionID(sectionId);
  }

  // --- CREATE ITEM ---
  async createItem(req: CreateSectionItemRequest): Promise<SectionItem> {
    // 1️⃣ Validasi apakah section ada
    const section = await this.sectionRepo.findOneByID(req.section_id);
    if (!section) {
      throw new Error("Invalid section_id: not found");
    }

    // 2️⃣ Insert item baru
    const item: SectionItem = {
      section_id: req.section_id, // convert camelCase → snake_case
      title: req.title,
      subtitle: req.subtitle ?? "",
      subjek: req.subjek ?? "",
      image_url: req.image_url ?? "",
      icon_url: req.icon_url ?? "",
      order: req.order ?? 0,
      extra_data: JSON.stringify(req.extra_data ?? {}),
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };

    return this.sectionItemRepo.create(item);
  }

  // --- UPDATE ITEM ---
  async updateItem(
    id: string,
    req: UpdateSectionItemRequest
  ): Promise<SectionItem> {
    const item = await this.sectionItemRepo.getByID(id);
    if (!item) throw new Error("Item not found");

    if (req.title) item.title = req.title;
    if (req.subtitle) item.subtitle = req.subtitle;
    if (req.subjek) item.subjek = req.subjek;
    if (req.image_url) item.image_url = req.image_url;
    if (req.icon_url) item.icon_url = req.icon_url;
    if (req.order !== undefined) item.order = req.order;
    if (req.extra_data) item.extra_data = req.extra_data;

    item.updated_at = new Date().toISOString();

    await this.sectionItemRepo.update(item);
    return item;
  }

  // --- DELETE ITEM ---
  async deleteItem(id: string): Promise<any> {
    return this.sectionItemRepo.delete(id);
  }
}
