package note

//go:generate go-enum --marshal --file $GOFILE --names --values --mustparse

// Type - TODO
//
// ENUM(bibliographic, fleeting, index, literature, permanent)
type Type string
