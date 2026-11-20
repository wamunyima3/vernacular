package httpx

// Optional returns nil for an empty string, so optional response fields are
// left out of the JSON instead of sent empty.
func Optional[T ~string](s T) *T {
	if s == "" {
		return nil
	}
	return &s
}

// Deref returns the value of an optional request field, or its zero value.
func Deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
