package validators

import (
	"regexp"

	"github.com/pedrotunin/go-jwt-auth/internal/utils"
)

func IsValidEmail(email string) error {
	const emailRegex = `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]*[a-zA9])?)*$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return utils.ErrInvalidEmail
	}
	return nil
}
