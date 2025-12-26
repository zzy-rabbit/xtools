package xerror

const (
	defaultSection = 0
)

var (
	ErrSuccess            = New(defaultSection+0, "success")
	ErrFail               = New(defaultSection+1, "fail")
	ErrInvalidParam       = New(defaultSection+2, "invalid param")
	ErrAuthenticationFail = New(defaultSection+3, "authentication fail")
	ErrPermissionDenied   = New(defaultSection+4, "permission denied")
	ErrNotFound           = New(defaultSection+5, "not found")
	ErrAlreadyExists      = New(defaultSection+6, "already exists")
	ErrTimeout            = New(defaultSection+7, "timeout")
	ErrInternalError      = New(defaultSection+8, "internal error")
	ErrNotImplemented     = New(defaultSection+9, "not implemented")
	ErrNotSupported       = New(defaultSection+10, "not supported")
	ErrFileOperatFail     = New(defaultSection+11, "file operation fail")
)
