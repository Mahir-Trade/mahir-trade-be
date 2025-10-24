// ===============================
// SECTION MODEL
// ===============================

import { SectionItem } from "./sectionItems.model";

export interface Section {
  id?: string;
  slug?: string;
  type?: string;
  title: string;
  subtitle: string;
  subjek?: string;
  order?: number;
  created_at: string;
  updated_at: string;
  items?: SectionItem[];
}
