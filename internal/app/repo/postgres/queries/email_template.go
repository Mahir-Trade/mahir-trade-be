package queries

const (
	QueryGetByKey = `
		SELECT body FROM email_templates
		WHERE key = $1
	`
)
