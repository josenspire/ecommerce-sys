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
var ErrSecurityCodeInvalid = "security code is invalid or expired, please check it or obtain another one"
var ErrPrivKeyParseFailedFromPEM = errors.New("failed to parse private key from PEM file")
var ErrPubKeyParseFailedFromPEM = errors.New("failed to parse public key from PEM file")

var ErrInvitationCodeInvalid = errors.New("this invitation code is invalid, please check it")

var ErrAddressNotFound = errors.New("this address was not found, please check it")
var ErrOrderNotFound = "operation failure, this order was not exist or already completed"
var ErrProductNotFound = "sorry, this product was not found or has been removed"

// Warning
var WarnParamsMissing = errors.New("params missing in user checking")
var WarnClassifiesMissing = "sorry, there are not have any classifies"
var WarnUserTeamMissing = "sorry, your agent teams is not initial, please call the system's administrator"
var WarnAccountNeedVerify = errors.New("your account needs to be verified by mobile phone and then set a password")
var WarnTelephoneAlreadyRegistered = errors.New("this phone number has been registered, please login with this telephone")
var WarnTelephoneNotRegistered = errors.New("sorry, the phone number has not been registered yet, please register first")

// common
var ErrRecordNotFound = errors.New("this record was not found, please check it")
var ErrCreateRecordsIsEmpty = errors.New("insert records is empty, please check your option")
var ErrParamsInValid = errors.New("sorry, your params are invalid, please check it")
var ErrPEMIsNotExist = errors.New("sorry, parse key fail, the PEM file is empty")
var ErrMysqlInitFailure = errors.New("error, server initial [mysql] pool failure, please check it")
var ErrRedisInitFailure = errors.New("error, server initial [redis] pool failure, please check it")
