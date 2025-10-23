package helpers

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
