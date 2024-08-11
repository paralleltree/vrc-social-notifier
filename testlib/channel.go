package testlib

func ConsumeChannel[T any](c <-chan T) []T {
	buf := []T{}
	for v := range c {
		buf = append(buf, v)
	}
	return buf
}
