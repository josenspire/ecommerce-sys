package utils

import "errors"

// Error
var ErrEof = errors.New("EOF")
var ErrClosedPipe = errors.New("io: read/write on closed pipe")
var ErrNoProgress = errors.New("multiple Read calls return no data or error")

var ErrTelOrPswInvalid = errors.New("telephone or password is invalid")
var ErrCurrentUserIsExist = errors.New("current user is already exist, please login")
var ErrParamsMissing = errors.New("sorry, your request params missing")
var ErrDecrypt = errors.New("sorry, your password verification failed")

var ErrInvitationCodeInvalid = errors.New("this invitation code is invalid, please check it")

var ErrAddressNotFound = errors.New("this address was not found, please check it")

// Warning
var WarnParamsMissing = errors.New("params missing in user checking")

// common
var ErrRecordNotFound = errors.New("this record was not found, please check it")
var ErrCreateRecordsIsEmpty = errors.New("insert records is empty, please check your option")
var ErrParamsInValid = errors.New("sorry, your params are invalid, please check it")
