package utils

const (
	// 200*
	REQUEST_SUCCESS  uint = 200000
	PARAMS_MISSING   uint = 200400
	TOKEN_IS_EXPIRED uint = 200401
	TELEPHONE_VERIFY uint = 200203

	// 400*
	REQUEST_FAIL uint = 400000

	// 500*
	SERVER_UNKNOW_ERROR   uint = 500000
	SERVER_REJECT_REQUEST uint = 500003

	// business code
	RECORD_NOT_FOUND           uint = 200404
	USER_TELEPHONE_PSW_INVALID uint = 200404
	ADDRESS_NOT_FOUND          uint = 200404
	ORDER_NOT_FOUND            uint = 200404
	SECURITY_CODE_INVALID      uint = 200400
)
