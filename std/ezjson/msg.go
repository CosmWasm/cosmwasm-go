package ezjson

// Set this as EmptyStruct{Seen: true} so it will serialize, otherwise it is missing
type EmptyStruct struct {
	Seen bool `json:"do_not_set_this_field"`
}

func (e EmptyStruct) WasSet() bool {
	return e.Seen
}

//// always serialize as an empty struct
//func (e EmptyStruct) MarshalJSON() ([]byte, error) {
//	return []byte(`{}`), nil
//}