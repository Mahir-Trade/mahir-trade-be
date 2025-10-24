export const ReportQueries = {
  CreateReport: `
    INSERT INTO reports (report_name, report_thumbnail_url, report_file_url, created_by, updated_by)
    VALUES ($1, $2, $3, $4, $4)
    RETURNING id
  `,

  GetReports: `
    SELECT COUNT(*) OVER() as total_count, id, report_name, report_thumbnail_url, report_file_url, created_by, updated_by, created_at, updated_at
    FROM reports
    WHERE deleted_at IS NULL
  `,

  GetReportByID: `
    SELECT id, report_name, report_thumbnail_url, report_file_url, created_by, updated_by, created_at, updated_at
    FROM reports
    WHERE id = $1
    AND deleted_at IS NULL
  `,

  UpdateReport: `
    UPDATE reports
    SET report_name = $1, report_thumbnail_url = $2, report_file_url = $3, updated_at = NOW(), updated_by = $4
    WHERE id = $5
  `,

  SoftDeleteReport: `
    UPDATE reports
    SET deleted_at = NOW(),
        deleted_by = $1,
        updated_at = NOW(),
        updated_by = $2
    WHERE id = $3
  `,
};
