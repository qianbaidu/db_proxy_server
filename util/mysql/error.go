package mysql

import "errors"

var (
	ProxyConnectClosed = errors.New("Mysql proxy connection closed.")
	ErrAccessDenied    = errors.New("access denied.")
	ErrUserNotExists   = errors.New("user not exists.")
	ErrPasswordWrong   = errors.New("Wrong user name or password.")
)
