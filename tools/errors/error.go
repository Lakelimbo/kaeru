package errors

// A way to return client-facing errors (such as in API responses)
// that is more-or-less structured and predictable.
//
// The idea is that this sort of error should return a readable
// identifier and probably a message explaining. Depending on
// the type of error, you could use the Errors ("fields" in JSON)
// to list different parts of the error (such as in a validation, where
// multiple fields can be errored at once). Use Raw ("raw_error" in
// JSON) if it's useful to show the user the raw error without any
// transformations.
type Error struct {
	Identifier ErrorCode         `json:"identifier"`
	Message    string            `json:"message,omitempty"`
	Errors     map[string]string `json:"fields,omitempty"`
	Raw        error             `json:"raw_error,omitempty"`
}

type ErrorCode string

// general errors
var (
	ErrInternal    ErrorCode = "err_internal"
	ErrInvalidBody ErrorCode = "erro_invalid_body"
)

// user errors
var (
	ErrUserInvalidCredentials ErrorCode = "err_user_invalid_credentials"
	ErrUserNotFound           ErrorCode = "err_user_not_found"
	ErrUserUnauthenticated    ErrorCode = "err_user_unauthenticated"
)

// app errors
var (
	ErrAppCreation ErrorCode = "err_app_creation"
	ErrAppNotFound ErrorCode = "err_app_not_found"
)
