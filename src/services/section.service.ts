import { SectionRepository } from "../repo/section_repo";
import { CreateSectionRequest, UpdateSectionRequest } from "../DTO/section.dto";
import { Section } from "../models/section.model";
import { SectionWithItems } from "../DTO/sectionItem.dto";

export class SectionService {
  private sectionRepo: SectionRepository;

  constructor() {
    this.sectionRepo = new SectionRepository();
  }

  // --- CREATE SECTION ---
  async saveSection(req: CreateSectionRequest): Promise<Section> {
    const section: Section = {
      slug: req.slug,
      type: req.type,
      title: req.title,
      subtitle: req.subtitle ?? "",
      order: req.order ?? 0,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };

    await this.sectionRepo.save(section);
    return section;
  }

  // --- UPDATE BY SLUG ---
  async updateSectionBySlug(
    slug: string,
    req: UpdateSectionRequest
  ): Promise<Section | null> {
    const existing = await this.sectionRepo.getBySlug(slug);
    if (!existing) {
      throw new Error("Section not found");
    }

    if (req.type) existing.type = req.type;
    if (req.title) existing.title = req.title;
    if (req.subtitle !== undefined) existing.subtitle = req.subtitle;
    if (req.order !== undefined) existing.order = req.order;

    existing.updated_at = new Date().toISOString();

    await this.sectionRepo.update(existing);
    return existing;
  }

  // --- GET BY TYPE ---
  async getSection(sectionType: string): Promise<Section | null> {
    return this.sectionRepo.getByType(sectionType);
  }

  // --- GET ALL WITH ITEMS ---
  async getAllSectionsWithItems(): Promise<SectionWithItems[]> {
    return this.sectionRepo.getAllSectionsWithItems();
  }

  // --- UPDATE BY ID ---
  async updateSection(id: string, req: UpdateSectionRequest): Promise<Section> {
    const existing = await this.sectionRepo.getByID(id);
    if (!existing) {
      throw new Error("Section not found");
    }

    if (req.slug) existing.slug = req.slug;
    if (req.type) existing.type = req.type;
    if (req.title) existing.title = req.title;
    if (req.subtitle !== undefined) existing.subtitle = req.subtitle;
    if (req.order !== undefined) existing.order = req.order;

    existing.updated_at = new Date().toISOString();

    await this.sectionRepo.update(existing);
    return existing;
  }

  // --- DELETE ---
  async deleteSection(id: string): Promise<any> {
    return this.sectionRepo.delete(id);
  }
}
