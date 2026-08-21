package mathx

func Abs[T Integer](x T) T {
	if x >= 0 {
		return x
	}
	return -x
}
