package queries

const (
	QueryCreateOrder = `
		INSERT INTO orders (user_id, package_id, status, payment_code, payment_url, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $6)
		RETURNING id
	`

	QueryGetOrderByID = `
		SELECT id, user_id, package_id, status, payment_code, payment_url, created_by, updated_by, created_at, updated_at
		FROM orders
		WHERE id = $1
			AND deleted_at IS NULL
	`

	QueryGetOrderByPaymentCode = `
		SELECT id, user_id, package_id, status, payment_code, payment_url, created_by, updated_by, created_at, updated_at
		FROM orders
		WHERE payment_code = $1
			AND deleted_at IS NULL
	`

	QueryGetOrders = `
		SELECT id, user_id, package_id, status, payment_code, payment_url, created_by, updated_by, created_at, updated_at
		FROM orders
		WHERE deleted_at IS NULL
		LIMIT $1 OFFSET $2
	`

	QueryUpdateOrderStatus = `
		UPDATE orders
		SET status = $1, updated_by = $2, updated_at = NOW()
		WHERE id = $3
			AND deleted_at IS NULL
	`
)
