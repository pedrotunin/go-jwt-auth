package validators

import "github.com/pedrotunin/go-jwt-auth/internal/utils"

func IsValidPassword(password string) error {
	if len(password) < 8 {
		return utils.ErrPasswordTooShort
	}
	return nil
}
