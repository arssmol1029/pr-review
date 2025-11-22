package errors

import (
	"errors"
	"fmt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrTeamNotFound = errors.New("team not found")
	ErrPRNotFound   = errors.New("pull request not found")
	ErrTeamExists   = errors.New("team already exists")
	ErrPRExists     = errors.New("pull request already exists")
	ErrPRMerged     = errors.New("pull request already merged")
	ErrNotAssigned  = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate  = errors.New("no active replacement candidate in team")
)

func WrapError(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
