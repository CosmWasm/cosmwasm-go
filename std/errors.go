package std

import (
	"reflect"
)

type ErrorType uintptr

const (
	Success ErrorType = iota
	GenericError
	InvalidBase64Error
	InvalidUtf8Error
	NotFoundError
	NullPointerError
	ParseError
	SerializeError
	UnauthorizedError
	UnderflowError
)

func package_StdError(err string) CosmosResponseError {
	return CosmosResponseError(`{"Err":` + err + `}`)
}

func GenerateError(errType ErrorType, msg string, msg_plus string) CosmosResponseError {
	var errorMsg string
	switch errType {
	case GenericError:
		errorMsg = GenericErr{Msg: msg}.Error()
	case InvalidBase64Error:
		errorMsg = InvalidBase64{Msg: msg}.Error()
	case InvalidUtf8Error:
		errorMsg = InvalidUtf8{Msg: msg}.Error()
	case NotFoundError:
		errorMsg = NotFound{Kind: msg}.Error()
	case NullPointerError:
		errorMsg = NullPointer{Msg: msg}.Error()
	case ParseError:
		errorMsg = ParseErr{Msg: msg}.Error()
	case SerializeError:
		errorMsg = SerializeErr{Msg: msg}.Error()
	case UnauthorizedError:
		errorMsg = Unauthorized{Msg: msg}.Error()
	case UnderflowError:
		errorMsg = Underflow{Minuend: msg, Subtrahend: msg_plus}.Error()
	default:
		errorMsg = GenericErr{Msg: msg + "Success ???"}.Error()
	}
	return package_StdError(errorMsg)
}

// StdError captures all errors returned from the Rust code as StdError.
// Exactly one of the fields should be set.
type StdError struct {
	SuccessRet    SuccessResult `json:"success_result"`
	GenericErr    GenericErr    `json:"generic_err,omitempty"`
	InvalidBase64 InvalidBase64 `json:"invalid_base64,omitempty"`
	InvalidUtf8   InvalidUtf8   `json:"invalid_utf8,omitempty"`
	NotFound      NotFound      `json:"not_found,omitempty"`
	NullPointer   NullPointer   `json:"null_pointer,omitempty"`
	ParseErr      ParseErr      `json:"parse_err,omitempty"`
	SerializeErr  SerializeErr  `json:"serialize_err,omitempty"`
	Unauthorized  Unauthorized  `json:"unauthorized,omitempty"`
	Underflow     Underflow     `json:"underflow,omitempty"`
}

var (
	_ error = GenericErr{}
	_ error = InvalidBase64{}
	_ error = InvalidUtf8{}
	_ error = NotFound{}
	_ error = NullPointer{}
	_ error = ParseErr{}
	_ error = SerializeErr{}
	_ error = Unauthorized{}
	_ error = Underflow{}
)

type GenericErr struct {
	Msg string `json:"msg,omitempty"`
}

func (e GenericErr) Error() string {
	return `{"generic_err":{"msg":"` + e.Msg + `"}}`
}

type InvalidBase64 struct {
	Msg string `json:"msg,omitempty"`
}

func (e InvalidBase64) Error() string {
	return `{"invalid_base64":{"msg":"` + e.Msg + `"}}`
}

type InvalidUtf8 struct {
	Msg string `json:"msg,omitempty"`
}

func (e InvalidUtf8) Error() string {
	return `{"invalid_utf8":{"msg":"` + e.Msg + `"}}`
}

type NotFound struct {
	Kind string `json:"kind,omitempty"`
}

func (e NotFound) Error() string {
	return `{"not_found":{"kind":"` + e.Kind + `"}}`
}

type NullPointer struct {
	Msg string `json:"msg"`
}

func (e NullPointer) Error() string {
	return `{"null_pointer": nil}`
}

type ParseErr struct {
	Target string `json:"target,omitempty"`
	Msg    string `json:"msg,omitempty"`
}

func (e ParseErr) Error() string {
	return `{"parse_err":{"target":"` + e.Target + `","msg":"` + e.Msg + `"}}`
}

type SerializeErr struct {
	Source string `json:"source,omitempty"`
	Msg    string `json:"msg,omitempty"`
}

func (e SerializeErr) Error() string {
	return `{"serializing":{"source":"` + e.Source + `","msg":"` + e.Msg + `"}}`
}

type Unauthorized struct {
	Msg string `json:"msg"`
}

func (e Unauthorized) Error() string {
	return `{"unauthorized": nil}`
}

type Underflow struct {
	Minuend    string `json:"minuend,omitempty"`
	Subtrahend string `json:"subtrahend,omitempty"`
}

func (e Underflow) Error() string {
	return `{"underflow":{"minuend":"` + e.Minuend + `","subtrahend":"` + e.Subtrahend + `"}}`
}

// check if an interface is nil (even if it has type info)
func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		// IsNil panics if you try it on a struct (not a pointer)
		return reflect.ValueOf(i).IsNil()
	}
	// if we aren't a pointer, can't be nil, can we?
	return false
}

// SystemError captures all errors returned from the Rust code as SystemError.
// Exactly one of the fields should be set.
type SystemError struct {
	SuccessRet         SuccessResult      `json:"success_result"`
	InvalidRequest     InvalidRequest     `json:"invalid_request,omitempty"`
	InvalidResponse    InvalidResponse    `json:"invalid_response,omitempty"`
	NoSuchContract     NoSuchContract     `json:"no_such_contract,omitempty"`
	Unknown            Unknown            `json:"unknown,omitempty"`
	UnsupportedRequest UnsupportedRequest `json:"unsupported_request,omitempty"`
}

var (
	_ error = InvalidRequest{}
	_ error = InvalidResponse{}
	_ error = NoSuchContract{}
	_ error = Unknown{}
	_ error = UnsupportedRequest{}
)

type SuccessResult struct {
	Msg string
}

func (e SuccessResult) Error() string {
	return "success"
}

type InvalidRequest struct {
	Err     string `json:"error"`
	Request []byte `json:"request"`
}

func (e InvalidRequest) Error() string {
	return `{"invalid_request":{"error":"` + e.Err + `","request":"` + string(e.Request) + `"}}`
}

type InvalidResponse struct {
	Err      string `json:"error"`
	Response []byte `json:"response"`
}

func (e InvalidResponse) Error() string {
	return `{"invalid_response":{"error":"` + e.Err + `","response":"` + string(e.Response) + `"}}`
}

type NoSuchContract struct {
	Addr string `json:"addr,omitempty"`
}

func (e NoSuchContract) Error() string {
	return `{"no_such_contract":{"addr":"` + e.Addr + `"}}`
}

type Unknown struct {
	Msg string `json:"msg"`
}

func (e Unknown) Error() string {
	return `{"unknow":nil}`
}

type UnsupportedRequest struct {
	Kind string `json:"kind,omitempty"`
}

func (e UnsupportedRequest) Error() string {
	return `{"unsupported_request":{"kind":"` + e.Kind + `"}}`
}

// ToSystemError will try to convert the given error to an SystemError.
// This is important to returning any Go error back to Rust.
//
// If it is already StdError, return self.
// If it is an error, which could be a sub-field of StdError, embed it.
// If it is anything else, **return nil**
//
// This may return nil on an unknown error, whereas ToStdError will always create
// a valid error type.
func ToSystemError(err error) SystemError {
	if isNil(err) {
		return SystemError{SuccessRet: SuccessResult{Msg: "success"}}
	}
	switch t := err.(type) {
	case InvalidRequest:
		return SystemError{InvalidRequest: t}
	case *InvalidRequest:
		return SystemError{InvalidRequest: *t}
	case InvalidResponse:
		return SystemError{InvalidResponse: t}
	case *InvalidResponse:
		return SystemError{InvalidResponse: *t}
	case NoSuchContract:
		return SystemError{NoSuchContract: t}
	case *NoSuchContract:
		return SystemError{NoSuchContract: *t}
	case Unknown:
		return SystemError{Unknown: t}
	case *Unknown:
		return SystemError{Unknown: *t}
	case UnsupportedRequest:
		return SystemError{UnsupportedRequest: t}
	case *UnsupportedRequest:
		return SystemError{UnsupportedRequest: *t}
	default:
		return SystemError{Unknown: Unknown{Msg: "Unknow System Error"}}
	}
}
