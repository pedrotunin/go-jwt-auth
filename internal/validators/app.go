package validators

import (
	"strings"

	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

func IsValidAppName(name string) error {
	if len(strings.TrimSpace(name)) < 4 {
		return utils.ErrAppNameInvalid
	}

	return nil
}

func IsValidAppDescription(description string) error {
	if len(strings.TrimSpace(description)) < 15 {
		return utils.ErrAppDescInvalid
	}

	return nil
}
