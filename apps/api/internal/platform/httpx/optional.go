package httpx

// Optional returns nil for an empty string, so optional response fields are
// left out of the JSON instead of sent empty.
func Optional(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
