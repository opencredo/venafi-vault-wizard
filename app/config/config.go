package config

import "errors"

var ErrBlankParam = errors.New("cannot be blank")

type Config interface {
	Validate() error
}

func ValidateConfigs(configs ...Config) error {
	for _, config := range configs {
		err := config.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
