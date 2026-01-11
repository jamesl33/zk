package note

// Frontmatter is the YAML front-matter embedded at the front of a note.
type Frontmatter struct {
	// Type of the note.
	Type Type `yaml:"type" json:"type" jsonschema:"The notes type (i.e. bibliographic, fleeting, index, literature, permanent)"`

	// Title for the note.
	Title string `yaml:"title" json:"title" jsonschema:"The title for the note"`

	// Date the note was created.
	Date string `yaml:"date" json:"date" jsonschema:"The date the note was created (e.g. 2006-01-02)"`

	// Tags for the note.
	Tags []string `yaml:"tags" json:"tags" jsonschema:"The notes tags (e.g. short, simple keywords to improve discoverability)"`
}
