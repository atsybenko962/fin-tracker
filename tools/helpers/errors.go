package helpers

import "fmt"

const (
	StatusSuccess                 = ""
	StatusDatabaseBuildQueryError = "failed to build select query"
	StatusDatabaseMakeQueryError  = "make db query error in search"
	StatusDatabaseScanRecordError = "db read item record error"
	StatusItemNotFound            = "item not found"
	StatusDeadlineExceeded        = "context deadline exceeded"
	StatusInvalidArgumentError    = "invalid request argument"
	StatusForbiddenError          = "forbidden error"
	StatusUnauthorizedError       = "unauthorized user request"
	StatusInternalError           = "internal error"
)

// Wrap оборачивает ошибки для прокидывания наверх по стеку вызова
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
