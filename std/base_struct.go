package std

//easyjson:skip
type MemRegion struct {
	Offset   uint32
	Capacity uint32
	Length   uint32
}

const REGION_HEAD_SIZE uint32 = 12

//easyjson:skip
type SliceHeader_tinyGo struct {
	Data uintptr
	Len  uintptr
	Cap  uintptr
}

//easyjson:skip
type SliceHeader_Go struct {
	Data uintptr
	Len  int
	Cap  int
}
