package matcher

// Tags - TODO
func Tags(i, e []string) (Matcher, error) {
	matchers := make([]Matcher, 0)

	for _, tag := range i {
		matchers = append(matchers, Tagged(tag))
	}

	for _, tag := range e {
		matchers = append(matchers, Not(Tagged(tag)))
	}

	// TODO
	all := And(matchers...)

	return all, nil
}
