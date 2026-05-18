package textfield

import (
	"github.com/LjungErik/zetra-lucene/lucene/document"
	"github.com/LjungErik/zetra-lucene/lucene/document/field"
)

type TextField struct {
	*field.Field
}

var _ document.DocumentField = (*TextField)(nil)

func New(name, value string, stored bool) *TextField {
	f := field.New(name,
		field.WithValue(value),
		field.WithStored(stored),
		field.WithIndexable(true),
		field.WithType(document.Text),
	)

	return &TextField{f}
}
