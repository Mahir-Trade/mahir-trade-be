package infra

import "time"

func InitTimezone() error {
	local, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return err
	}
	time.Local = local

	return nil
}