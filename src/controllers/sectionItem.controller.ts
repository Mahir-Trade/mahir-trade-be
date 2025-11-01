import { Request, Response } from "express";
import { SectionItemService } from "../services/SectionItem.service";

import { DefaultResponse } from "../utils/response";
import {
  CreateSectionItemRequest,
  UpdateSectionItemRequest,
} from "../DTO/sectionItem.dto";

export class SectionItemController {
  private sectionItemService: SectionItemService;

  constructor() {
    this.sectionItemService = new SectionItemService();

    this.createItem = this.createItem.bind(this);
    this.getItemsBySectionID = this.getItemsBySectionID.bind(this);
    this.getItems = this.getItems.bind(this);
    this.updateItem = this.updateItem.bind(this);
    this.deleteItem = this.deleteItem.bind(this);
  }

  // --- GET /section-items/:section_id ---
  async getItemsBySectionID(req: Request, res: Response): Promise<Response> {
    try {
      const sectionId = req.params.section_id;
      const items = await this.sectionItemService.getItemsBySectionID(
        sectionId
      );

      return res.status(200).json({
        code: 200,
        message: "Success",
        data: items,
      } as DefaultResponse);
    } catch (err: any) {
      console.error("[SectionItemController][getItemsBySectionID] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Failed to fetch items",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- GET /section-items ---
  async getItems(req: Request, res: Response): Promise<Response> {
    try {
      const data = await this.sectionItemService.getItems();

      return res.status(200).json({
        code: 200,
        message: "Success",
        data,
      } as DefaultResponse);
    } catch (err: any) {
      console.error("[SectionItemController][getItems] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Internal Server Error",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- POST /section-items ---
  async createItem(req: Request, res: Response): Promise<Response> {
    try {
      const body = req.body as { sections: CreateSectionItemRequest[] };

      if (!body || !Array.isArray(body.sections)) {
        return res.status(400).json({
          code: 400,
          message: "Invalid request body, expected { sections: [...] }",
        } as DefaultResponse);
      }

      const createdItems = [];
      for (const item of body.sections) {
        const created = await this.sectionItemService.createItem(item);
        createdItems.push(created);
      }

      return res.status(201).json({
        code: 201,
        message: "Bulk section item creation successful",
        data: createdItems,
      } as DefaultResponse);
    } catch (err: any) {
      console.error("[SectionItemController][createItem] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Creation failed",
        error: err.message,
      } as DefaultResponse);
    }
  }
  // --- PUT /section-items/:id ---
  async updateItem(req: Request, res: Response): Promise<Response> {
    try {
      const id = req.params.id;
      const body = req.body as UpdateSectionItemRequest;

      if (
        !body.title &&
        !body.subtitle &&
        !body.subjek &&
        !body.image_url &&
        !body.icon_url &&
        !body.extra_data
      ) {
        return res.status(400).json({
          code: 400,
          message: "At least one field must be provided for update",
        } as DefaultResponse);
      }

      const updated = await this.sectionItemService.updateItem(id, body);

      return res.status(200).json({
        code: 200,
        message: "Item updated successfully",
        data: updated,
      } as DefaultResponse);
    } catch (err: any) {
      console.error("[SectionItemController][updateItem] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Update failed",
        error: err.message,
      } as DefaultResponse);
    }
  }

  // --- DELETE /section-items/:id ---
  async deleteItem(req: Request, res: Response): Promise<Response> {
    try {
      const id = req.params.id;
      await this.sectionItemService.deleteItem(id);

      return res.status(200).json({
        code: 200,
        message: "Item deleted successfully",
      } as DefaultResponse);
    } catch (err: any) {
      console.error("[SectionItemController][deleteItem] error:", err);
      return res.status(500).json({
        code: 500,
        message: "Delete failed",
        error: err.message,
      } as DefaultResponse);
    }
  }
}
