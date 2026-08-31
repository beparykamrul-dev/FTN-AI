package security

import (
    "errors"
    "strings"
)

var ErrApprovalRequired = errors.New("approval required")
var ErrInvalidAction = errors.New("invalid privileged action")

// Gate is the policy boundary for privileged ISP mutations.
// It deliberately does not execute the mutation; callers must obtain an
// approved request before invoking an infrastructure adapter.
type Gate struct{}

func (Gate) ValidateAction(action string) error {
    action = strings.TrimSpace(action)
    if action == "" {
        return ErrInvalidAction
    }
    return nil
}

func (Gate) RequireApproval(action string, approved bool) error {
    if err := (Gate{}).ValidateAction(action); err != nil {
        return err
    }
    if !approved {
        return ErrApprovalRequired
    }
    return nil
}
