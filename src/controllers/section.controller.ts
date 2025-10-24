import { Request, Response } from "express";
import { SectionService } from "../services/section.service";
import { CreateSectionRequest, UpdateSectionRequest } from "../DTO/section.dto";

export class SectionController {
  private sectionService: SectionService;

  constructor() {
    this.sectionService = new SectionService();

    this.createSection = this.createSection.bind(this);
    this.updateSectionBySlug = this.updateSectionBySlug.bind(this);
    this.getSectionByType = this.getSectionByType.bind(this);
    this.getFullSection = this.getFullSection.bind(this);
    this.updateSection = this.updateSection.bind(this);
    this.deleteSection = this.deleteSection.bind(this);
  }

  // POST /api/landing-page
  async createSection(req: Request, res: Response): Promise<Response> {
    try {
      const body: { sections: CreateSectionRequest } = req.body;

      // Validasi payload
      if (!body || !body.sections) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request",
        });
      }

      await this.sectionService.saveSection(body.sections);

      return res.status(200).json({
        code: 200,
        message: "Section created successfully",
        data: body.sections,
      });
    } catch (err: any) {
      return res.status(500).json({
        code: 500,
        message: "Failed to save",
        error: err.message,
      });
    }
  }

  // PUT /api/landing-page/:slug
  async updateSectionBySlug(req: Request, res: Response): Promise<Response> {
    try {
      const slug = req.params.slug;
      const body: UpdateSectionRequest = req.body;

      if (!body.type && !body.title && !body.subtitle && !body.order) {
        return res.status(400).json({
          code: 400,
          message: "At least one field must be provided for update",
        });
      }

      const updated = await this.sectionService.updateSectionBySlug(slug, body);

      return res.status(200).json({
        code: 200,
        message: "Section updated successfully",
        data: updated,
      });
    } catch (err: any) {
      return res.status(500).json({
        code: 500,
        message: "Update failed",
        error: err.message,
      });
    }
  }

  // GET /api/landing-page/:type
  async getSectionByType(req: Request, res: Response): Promise<Response> {
    try {
      const sectionType = req.params.type;

      const result = await this.sectionService.getSection(sectionType);

      if (!result || !result.type) {
        return res.status(404).json({
          code: 404,
          message: "Section not found",
        });
      }

      return res.status(200).json({
        code: 200,
        message: "Success",
        data: result,
      });
    } catch (err: any) {
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // GET /api/landing-page
  async getFullSection(req: Request, res: Response): Promise<Response> {
    try {
      const data = await this.sectionService.getAllSectionsWithItems();

      return res.status(200).json({
        code: 200,
        message: "Success",
        data: {
          sections: data,
        },
      });
    } catch (err: any) {
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      });
    }
  }

  // PUT /api/landing-page/:id
  async updateSection(req: Request, res: Response): Promise<Response> {
    try {
      const id = req.params.id;
      const body: UpdateSectionRequest = req.body;

      if (
        !body.slug &&
        !body.type &&
        !body.title &&
        !body.subtitle &&
        !body.order
      ) {
        return res.status(400).json({
          code: 400,
          message: "At least one field must be provided for update",
        });
      }

      const updated = await this.sectionService.updateSection(id, body);

      return res.status(200).json({
        code: 200,
        message: "Section updated successfully",
        data: updated,
      });
    } catch (err: any) {
      return res.status(500).json({
        code: 500,
        message: "Update failed",
        error: err.message,
      });
    }
  }

  // DELETE /api/landing-page/:id
  async deleteSection(req: Request, res: Response): Promise<Response> {
    try {
      const id = req.params.id;

      await this.sectionService.deleteSection(id);

      return res.status(200).json({
        code: 200,
        message: "Item deleted",
      });
    } catch (err: any) {
      return res.status(500).json({
        code: 500,
        message: "Delete failed",
        error: err.message,
      });
    }
  }
}
