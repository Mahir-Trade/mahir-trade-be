export interface Report {
  id: number;
  report_name: string; // required
  report_thumbnail_url: string; // required
  report_file_url: string; // required
  created_by: string;
  updated_by: string;
  created_at: string;
  updated_at: string;
}
