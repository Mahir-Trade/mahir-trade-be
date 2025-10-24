export interface Group {
  id?: number; // int64 → number, omitempty → optional
  uuid?: string; // string, omitempty → optional
  group_name: string; // required
  created_at?: string; // optional
  updated_at?: string; // optional
}

export interface GetGroupsRequest {
  limit: number; // int64 → number
  page: number; // int64 → number
  showAll: boolean; // bool → boolean
}
