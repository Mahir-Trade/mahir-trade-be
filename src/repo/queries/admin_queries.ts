export const AdminQueries = {
  CreateAdmin: `
    INSERT INTO admins (email, username, password)
    VALUES ($1, $2, $3)
    RETURNING admin_id
  `,

  FindByEmail: `
    SELECT admin_id, uuid, email, username, password
    FROM admins
    WHERE username = $1
  `,

  SoftDeleteAdmin: `
    UPDATE admins
    SET deleted_at = NOW(),
        deleted_by = $1,
        updated_at = NOW(),
        updated_by = $2
    WHERE admin_id = $3
  `,
};
