package matcher

import "fmt"

// Entire returns a matcher for the entire note, using the given fixed/glob/regex patterns.
func Entire(f, g, r string) (Matcher, error) {
	fm, err := Frontmatter(f, g, r)
	if err != nil {
		return nil, fmt.Errorf("failed to create frontmatter matcher: %w", err)
	}

	// The use of 'Or' means we defer loading the note body if the frontmatter already matches
	body, err := Body(f, g, r)
	if err != nil {
		return nil, fmt.Errorf("failed to create body matcher: %w", err)
	}

	return Or(fm, body), nil
}
