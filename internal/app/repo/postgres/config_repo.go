package postgres

import (
	"database/sql"

	"go.uber.org/dig"
)

type (
	ConfigRepo interface {
		GetConfigByKey(key string) (string, error)
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
