package matcher

// Tags matcher for the given include/exclude list.
func Tags(i, e []string) (Matcher, error) {
	matchers := make([]Matcher, 0)

	for _, tag := range i {
		matchers = append(matchers, Tagged(tag))
	}

	for _, tag := range e {
		matchers = append(matchers, Not(Tagged(tag)))
	}

	// Combine the matchers
	all := And(matchers...)

	return all, nil
}
