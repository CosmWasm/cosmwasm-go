package safe_math

type SafeMath interface {
	SafeAdd(a, b uint64) (uint64, error)
	SafeSub(a, b uint64) (uint64, error)
	SafeMul(a, b uint64) (uint64, error)
	SafeDiv(a, b uint64) (uint64, error)
}
