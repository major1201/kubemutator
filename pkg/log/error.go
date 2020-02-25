package log

import (
	"fmt"
	"go.uber.org/zap"
)

// NamedError extends the zap filed to log the error message more beautified
func NamedError(key string, err error) zap.Field {
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	return zap.NamedError(key, err)
}

// Error extends the zap field to log the error message more beautified
func Error(err error) zap.Field {
	return NamedError("error", err)
}

// Errors extends the zap field to log multiple errors more beautified
func Errors(key string, errs []error) zap.Field {
	fmt.Printf("%s\n", key)
	if errs != nil {
		for _, err := range errs {
			fmt.Printf("%+v\n", err)
		}
	}
	return zap.Errors(key, errs)
}
