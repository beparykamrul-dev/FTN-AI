package account

import (
	"errors"
	"strings"
)

var (
	ErrAccountNumberRequired = errors.New("account number is required")
	ErrAccountNameRequired   = errors.New("account name is required")
	ErrAccountStatusInvalid  = errors.New("invalid account status")
)

func Validate(a Account) error {
	if strings.TrimSpace(a.AccountNo) == "" {
		return ErrAccountNumberRequired
	}
	if strings.TrimSpace(a.Name) == "" {
		return ErrAccountNameRequired
	}
	if a.Status != StatusActive && a.Status != StatusSuspended && a.Status != StatusDisabled {
		return ErrAccountStatusInvalid
	}
	return nil
}
