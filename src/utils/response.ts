import { Response } from "express";

// ======================================================
// 📦 Common Response Interfaces (TypeScript version of Go structs)
// ======================================================

/**
 * DefaultResponse<T>
 * Generic supaya data bisa bertipe apa pun (kayak struct response di Go)
 */
export interface DefaultResponse<T = any> {
  code: number;
  message: string;
  data?: T; // optional — seperti `omitempty`
  meta?: any; // optional — seperti `omitempty`
  error?: any; // optional — seperti `omitempty`
}

/**
 * DefaultMetaData
 * Untuk pagination meta info (page, total, limit, dsb)
 */
export interface DefaultMetaData {
  page: number;
  totalPages: number;
  totalItems: number;
  limit: number;
  hasNext: boolean;
  hasPrevious: boolean;
}

/**
 * DefaultPaginationResponseData<T>
 * Response standar untuk list data dengan pagination
 */
export interface DefaultPaginationResponseData<T = any> {
  results: T[];
  meta: DefaultMetaData;
}

/**
 * PaginationRequest
 * Struktur request standar untuk pagination
 */
export interface PaginationRequest {
  limit: number;
  page: number;
  search?: string;
  sort_by?: string;
  membership_status?: string;
  filters?: [string, string][];
  showAll?: boolean; // ✅ pakai camelCase biar match dengan req.query.showAll
}

/**
 * SendResponse<T>
 * Struktur standar untuk API response final
 */
export interface SendResponse<T = any> {
  code: number;
  message: string;
  data?: T;
  meta?: {
    page?: number;
    limit?: number;
    totalPages?: number;
    totalItems?: number;
    hasNext?: boolean;
    hasPrevious?: boolean;
  };
  error?: any;
}

// ======================================================
// 🧠 Helper utama untuk kirim response terstandarisasi
// ======================================================

export function sendResponse<T>(
  res: Response,
  status: number,
  payload: any
): Response {
  // 🔄 Jika payload punya 'results', ubah jadi 'data'
  if (payload?.results) {
    payload.data = payload.results;
    delete payload.results;
  }

  // ✅ Bangun struktur final
  const finalResponse: SendResponse<T> = {
    code: payload.code ?? status,
    message: payload.message ?? "Success",
    data: payload.data,
    meta: payload.meta,
    error: payload.error,
  };

  return res.status(status).json(finalResponse);
}

/**
 * 💥 Helper opsional untuk error handling
 */
export function sendError(
  res: Response,
  status: number,
  message: string,
  error?: any
): Response {
  const finalResponse: SendResponse = {
    code: status,
    message,
    error,
  };

  return res.status(status).json(finalResponse);
}
