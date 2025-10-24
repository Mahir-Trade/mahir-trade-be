export interface DefaultResponse<T = any> {
  code: number;
  message: string;
  data?: T; // optional (omitempty)
  error?: any; // optional (omitempty)
}

export interface DefaultMetaData {
  page: number;
  totalPages: number;
  totalItems: number;
  limit: number;
  hasNext: boolean;
  hasPrevious: boolean;
}

export interface DefaultPaginationResponseData<T = any> {
  results: T[];
  meta: DefaultMetaData;
}
