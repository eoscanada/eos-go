package errorcodes

type EOSError struct {
	Name string
	Code int
}

var UnspecifiedException = EOSError{"unspecified_exception_code", 3990000}
var UnhandledException = EOSError{"unhandled_exception_code", 3990001}
var TimeoutException = EOSError{"timeout_exception_code", 3990002}
var FileNotFoundException = EOSError{"file_not_found_exception_code", 3990003}
var ParseErrorException = EOSError{"parse_error_exception_code", 3990004}
var InvalidArgException = EOSError{"invalid_arg_exception_code", 3990005}
var KeyNotFoundException = EOSError{"key_not_found_exception_code", 3990006}
var BadCastException = EOSError{"bad_cast_exception_code", 3990007}
var OutOfRangeException = EOSError{"out_of_range_exception_code", 3990008}
var CanceledException = EOSError{"canceled_exception_code", 3990009}
var AssertException = EOSError{"assert_exception_code", 3990010}
var EOFException = EOSError{"eof_exception_code", 3990011}
var StdException = EOSError{"std_exception_code", 3990013}
var InvalidOperationException = EOSError{"invalid_operation_exception_code", 3990014}
var UnknownHostException = EOSError{"unknown_host_exception_code", 3990015}
var NullOptional = EOSError{"null_optional_code", 3990016}
var UDTError = EOSError{"udt_error_code", 3990017}
var AESError = EOSError{"aes_error_code", 3990018}
var Overflow = EOSError{"overflow_code", 3990019}
var Underflow = EOSError{"underflow_code", 3990020}
var DivideByZero = EOSError{"divide_by_zero_code", 3990021}
