package eoserr

import "fmt"

type Error struct {
	Name string
	Code int
}

func (e Error) Error() string {
	return fmt.Sprintf("eos error: %q, code: %d", e.Name, e.Code)
}

var ErrUnspecifiedException = Error{"unspecified_exception_code", 3990000}
var ErrUnhandledException = Error{"unhandled_exception_code", 3990001}
var ErrTimeoutException = Error{"timeout_exception_code", 3990002}
var ErrFileNotFoundException = Error{"file_not_found_exception_code", 3990003}
var ErrParseErrorException = Error{"parse_error_exception_code", 3990004}
var ErrInvalidArgException = Error{"invalid_arg_exception_code", 3990005}
var ErrKeyNotFoundException = Error{"key_not_found_exception_code", 3990006}
var ErrBadCastException = Error{"bad_cast_exception_code", 3990007}
var ErrOutOfRangeException = Error{"out_of_range_exception_code", 3990008}
var ErrCanceledException = Error{"canceled_exception_code", 3990009}
var ErrAssertException = Error{"assert_exception_code", 3990010}
var ErrEOFException = Error{"eof_exception_code", 3990011}
var ErrStdException = Error{"std_exception_code", 3990013}
var ErrInvalidOperationException = Error{"invalid_operation_exception_code", 3990014}
var ErrUnknownHostException = Error{"unknown_host_exception_code", 3990015}
var ErrNullOptional = Error{"null_optional_code", 3990016}
var ErrUDTError = Error{"udt_error_code", 3990017}
var ErrAESError = Error{"aes_error_code", 3990018}
var ErrOverflow = Error{"overflow_code", 3990019}
var ErrUnderflow = Error{"underflow_code", 3990020}
var ErrDivideByZero = Error{"divide_by_zero_code", 3990021}
