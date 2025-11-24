package note

// Frontmatter - TODO
type Frontmatter struct {
	// Type - TODO
	Type Type `yaml:"type,omitempty"`

	// Title - TODO
	Title string `yaml:"title,omitempty"`

	// Date - TODO
	Date string `yaml:"date,omitempty"`

	// Title - TODO
	Tags []string `yaml:"tags,omitempty"`
}
