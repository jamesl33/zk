package vector

// FindOptions encapsulates the options for finding notes.
type FindOptions struct {
	// limit is the maximum number of notes to find, defaults to 15.
	limit int

	// maxDistance is the maximum cosine distance a note may be from the query, unbounded if zero.
	maxDistance float32
}

// WithLimit allows specifying the maximum number of notes to find.
func WithLimit(limit int) func(*FindOptions) {
	return func(o *FindOptions) {
		o.limit = limit
	}
}

// WithMaxDistance allows specifying the maximum cosine distance a note may be from the query.
func WithMaxDistance(maxDistance float32) func(*FindOptions) {
	return func(o *FindOptions) {
		o.maxDistance = maxDistance
	}
}
