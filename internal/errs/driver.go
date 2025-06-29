package errs

import "errors"

var (
	EmptyToken           = errors.New("empty token")
	MountPathCannotEmpty = errors.New("mounting path cannot be empty")
)
