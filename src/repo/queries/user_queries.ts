export const UserQueries = {
  CreateUser: `
    INSERT INTO users (email, phone_number, username, password, is_active) 
    VALUES ($1, $2, $3, $4, TRUE) 
    RETURNING user_id
  `,

  FindUserByEmailOrUsername: `
    SELECT user_id, uuid, phone_number, username, email, password
    FROM users 
    WHERE email = $1 OR username = $1
  `,

  FindUserByEmailAndUsername: `
    SELECT user_id, uuid, phone_number, username, email, password 
    FROM users 
    WHERE email = $1 AND username = $2
  `,

  GetUserByID: `
    SELECT user_id, uuid, phone_number, username, email, password, is_active 
    FROM users 
    WHERE user_id = $1 AND deleted_at IS NULL
  `,

  GetUserByUUID: `
    SELECT user_id, uuid, phone_number, username, email, password 
    FROM users 
    WHERE uuid = $1 AND deleted_at IS NULL
  `,

  UpdateTypeUser: `
    UPDATE users
    SET is_active = $1, updated_by = $2, updated_at = NOW()
    WHERE user_id = $3
  `,

  GetUsers: `
  SELECT
    COUNT(*) OVER() as total_count,
    u.user_id,
    u.uuid,
    u.phone_number,
    u.email,
    u.username,
    u.is_active,
    CASE
      WHEN um.is_membership_active = true THEN 'Premium'
      WHEN um.is_membership_active = false THEN 'Expired'
      ELSE 'Standard'
    END as account_type,
    um.expired_at as membership_expired_date,
    u.created_at,
    u.created_by
  FROM
    users u
  LEFT JOIN
    user_memberships um ON um.user_id = u.user_id
  WHERE
    u.deleted_at IS NULL
`,

  UpdatePassword: `
    UPDATE users
    SET password = $1, updated_by = $2, updated_at = NOW()
    WHERE user_id = $3
  `,

  SetUserVerified: `
    UPDATE users
    SET verified_at = NOW()
    WHERE user_id = $1
  `,
};
