package feat_test

func consumeChannel[T any](c <-chan T) []T {
	buf := []T{}
	for v := range c {
		buf = append(buf, v)
	}
	return buf
}
