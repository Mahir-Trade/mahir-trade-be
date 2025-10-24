export interface CreateSectionRequest {
  slug: string;
  type: string;
  title: string;
  subtitle?: string;
  order?: number;
}

export interface UpdateSectionRequest {
  slug?: string;
  type?: string;
  title?: string;
  subtitle?: string;
  order?: number;
}
