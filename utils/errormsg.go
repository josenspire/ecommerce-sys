package utils

import "errors"

// Error
var ErrEof = errors.New("EOF")
var ErrClosedPipe = errors.New("io: read/write on closed pipe")
var ErrNoProgress = errors.New("multiple Read calls return no data or error")

var ErrTelOrPswInvalid = errors.New("telephone or password is invalid")
var ErrCurrentUserIsExist = errors.New("current user is already exist, please login")

// Warning
var WarnParamsMissing = errors.New("params missing in user checking")
