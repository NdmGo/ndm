package utils

import (
	"context"
	"fmt"
)

// NewErr create new error, supports wrapping original error
func NewErr(err error, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", msg, err)
	}
	return fmt.Errorf("%s", msg)
}

// Throw throws out an exception, which can be caught be TryCatch or recover.
func Throw(exception interface{}) {
	panic(exception)
}

// Try implements try... logistics using internal panic...recover.
// It returns error if any exception occurs, or else it returns nil.
func Try(ctx context.Context, try func(ctx context.Context)) (err error) {
	if try == nil {
		return
	}
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok {
				err = v
			} else {
				err = fmt.Errorf("%s", exception)
			}
		}
	}()
	try(ctx)
	return
}

// TryCatch implements `try...catch..`. logistics using internal `panic...recover`.
// It automatically calls function `catch` if any exception occurs and passes the exception as an error.
// If `catch` is given nil, it ignores the panic from `try` and no panic will throw to parent goroutine.
//
// But, note that, if function `catch` also throws panic, the current goroutine will panic.
func TryCatch(ctx context.Context, try func(ctx context.Context), catch func(ctx context.Context, exception error)) {
	if try == nil {
		return
	}
	if exception := Try(ctx, try); exception != nil && catch != nil {
		catch(ctx, exception)
	}
}
