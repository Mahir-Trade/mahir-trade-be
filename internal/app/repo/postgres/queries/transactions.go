package queries

const (
	QueryCreateTransaction = `
		INSERT INTO transactions (order_id, amount, settlement_date, webhook_id, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $5)
		RETURNING id
	`
)
