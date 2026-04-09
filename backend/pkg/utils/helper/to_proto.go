package helper

func ToProto[T any, P any](src *T, f func(*T) *P) *P {
	if src == nil {
		return nil
	}
	return f(src)
}
