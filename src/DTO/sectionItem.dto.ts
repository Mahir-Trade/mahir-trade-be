import { SectionItem } from "../models/sectionItems.model";

export interface SectionWithItems {
  id: string;
  slug: string;
  type: string;
  title: string;
  subtitle?: string;
  order?: number;
  created_at: string;
  updated_at: string;
  items: SectionItem[];
}

export interface CreateSectionItemRequest {
  section_id: string; // camelCase utk consistency di TS
  title: string;
  subtitle?: string;
  subjek?: string;
  image_url?: string;
  icon_url?: string;
  order?: number;
  extra_data?: Record<string, any>;
  created_at?: string;
  updated_at?: string;
}

export interface UpdateSectionItemRequest {
  title?: string;
  subtitle?: string;
  subjek?: string;
  image_url?: string;
  icon_url?: string;
  order?: number;
  extra_data?: Record<string, any>;
  updated_at?: string;
}
