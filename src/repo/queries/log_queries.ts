export const GeneralLogQueries = {
  CreateGeneralLog: `
    INSERT INTO general_logs (user_id, raw_body, created_by, updated_by)
    VALUES ($1, $2, $3, $3)
    RETURNING id
  `,
};
