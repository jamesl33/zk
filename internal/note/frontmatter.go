package note

// Frontmatter - TODO
type Frontmatter struct {
	// Type - TODO
	Type Type `yaml:"type"`

	// Title - TODO
	Title string `yaml:"title"`

	// Date - TODO
	Date string `yaml:"date"`

	// Title - TODO
	Tags []string `yaml:"tags"`
}
