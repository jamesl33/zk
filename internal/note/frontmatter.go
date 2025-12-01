package note

// Frontmatter is the YAML front-matter embedded at the front of a note.
type Frontmatter struct {
	// Type of the note.
	Type Type `yaml:"type"`

	// Title for the note.
	Title string `yaml:"title"`

	// Date the note was created.
	Date string `yaml:"date"`

	// Tags for the note.
	Tags []string `yaml:"tags"`
}
