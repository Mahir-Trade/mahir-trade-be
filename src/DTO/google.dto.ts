export interface FileUpload {
  file_name: string;
  size: number;
  localFilePath: string;
  file_content_type: string;
  bucket_name: string;
}

export interface URLParserResponse {
  URL: string;
  host: string;
  path: string;
  bucket_name: string;
}

export interface GoogleCallbackRequest {
  state: string;
  code: string;
}
