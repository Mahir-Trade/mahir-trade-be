package postgres

import (
	"database/sql"

	"go.uber.org/dig"
)

type (
	ConfigRepo interface {
		GetConfigByKey(key string) (string, error)
		UpdateConfigByKey(key string, value string, operatedBy string) error
	}

	ConfigRepoImpl struct {
		dig.In

		*sql.DB
	}
)

func NewConfigRepo(impl ConfigRepoImpl) ConfigRepo {
	return &impl
}

func (c *ConfigRepoImpl) GetConfigByKey(key string) (string, error) {
	var value string
	err := c.QueryRow("SELECT value FROM config WHERE key = $1", key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (c *ConfigRepoImpl) UpdateConfigByKey(key string, value string, operatedBy string) error {
	_, err := c.Exec("UPDATE config SET value = $1, updated_at = NOW(), updated_by = $3 WHERE key = $2", value, key, operatedBy)
	if err != nil {
		return err
	}
	return nil
}
