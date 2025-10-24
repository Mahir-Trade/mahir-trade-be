export interface SubModule {
  id?: number; // int64 → number, optional karena omitempty
  uuid?: string; // string, optional karena omitempty
  module_id: number | null; // sql.NullInt64 → bisa number atau null
  sub_module_name: string; // required
  title: string; // required
  video_url: string; // required
  status?: string; // optional
  created_by?: string; // optional
  updated_by?: string; // optional
  created_at?: string; // optional (datetime → ISO string)
  updated_at?: string; // optional (datetime → ISO string)
}
