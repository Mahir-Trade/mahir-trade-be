export interface PaginationRequest {
  limit: number; // >= 1, <= 100
  page: number; // >= 1
  search?: string; // optional (omitempty)
  filters?: string[]; // optional
  showAll: boolean;
}
