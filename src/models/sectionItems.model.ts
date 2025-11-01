export interface SectionItem {
  id?: string; // optional karena biasanya auto dari DB
  section_id?: string; // optional karena biasanya di-set dari luar
  title: string;
  subtitle?: string;
  subjek?: string;
  image_url?: string;
  icon_url?: string;
  order?: number;
  extra_data?: any;
  created_at: string;
  updated_at: string;
}
