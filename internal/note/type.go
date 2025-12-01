package note

//go:generate go-enum --marshal --file $GOFILE --names --values --mustparse

// Type of a note.
//
// ENUM(bibliographic, fleeting, index, literature, permanent)
type Type string
