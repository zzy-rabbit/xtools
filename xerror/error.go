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
)
