package std

// Set this as EmptyStruct{Seen: true} so it will serialize, otherwise it is missing
type EmptyStruct struct {
	Seen bool `json:"seen,omitempty"`
}

// always serialize as an empty struct
func (e EmptyStruct) MarshalJSON() ([]byte, error) {
	return []byte(`{}`), nil
}

// we store seen
func (e *EmptyStruct) UnmarshalJSON(data []byte) error {
	e.Seen = true
	return nil
}
