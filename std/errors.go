package std

// This is a (temporary?) helper to use in place of errors.New
func NewError(msg string) error {
	return myError{msg: msg}
}

//easyjson:skip
type myError struct {
	msg string
}

func (m myError) Error() string {
	return m.msg
}

//easyjson:skip
type OutOfGasError struct{}

var _ error = OutOfGasError{}

func (o OutOfGasError) Error() string {
	return "Out of gas"
}

// StdError captures all errors returned from the Rust code as StdError.
// Exactly one of the fields should be set.
type StdError struct {
	GenericErr    *GenericErr    `json:"generic_err,omitempty"`
	InvalidBase64 *InvalidBase64 `json:"invalid_base64,omitempty"`
	InvalidUtf8   *InvalidUtf8   `json:"invalid_utf8,omitempty"`
	NotFound      *NotFound      `json:"not_found,omitempty"`
	NullPointer   *NullPointer   `json:"null_pointer,omitempty"`
	ParseErr      *ParseErr      `json:"parse_err,omitempty"`
	SerializeErr  *SerializeErr  `json:"serialize_err,omitempty"`
	Unauthorized  *Unauthorized  `json:"unauthorized,omitempty"`
	Underflow     *Underflow     `json:"underflow,omitempty"`
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
	Msg string
}

func (e GenericErr) Error() string {
	return "GenericErr: " + e.Msg
}

type InvalidBase64 struct {
	Msg string
}

func (e InvalidBase64) Error() string {
	return "InvalidBase64: " + e.Msg
}

type InvalidUtf8 struct {
	Msg string
}

func (e InvalidUtf8) Error() string {
	return "InvalidUtf8: " + e.Msg
}

type NotFound struct {
	Kind string
}

func (e NotFound) Error() string {
	return "NotFound: " + e.Kind
}

type NullPointer struct{}

func (e NullPointer) Error() string {
	return `NullPointer`
}

type ParseErr struct {
	Target string
	Msg    string
}

func (e ParseErr) Error() string {
	return "ParseErr (" + e.Target + "): " + e.Msg
}

type SerializeErr struct {
	Source string
	Msg    string
}

func (e SerializeErr) Error() string {
	return "SerializeErr (" + e.Source + "): " + e.Msg
}

type Unauthorized struct{}

func (e Unauthorized) Error() string {
	return "Unauthorized"
}

type Underflow struct {
	Minuend    string
	Subtrahend string
}

func (e Underflow) Error() string {
	return "Underflow subtract " + e.Minuend + " from " + e.Subtrahend
}
