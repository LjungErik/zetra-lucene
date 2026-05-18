package storedfield

import (
	"github.com/LjungErik/zetra-lucene/lucene/document"
	"github.com/LjungErik/zetra-lucene/lucene/document/field"
)

type StoredField struct {
	*field.Field
}

var _ document.DocumentField = (*StoredField)(nil)

func New(name, value string) *StoredField {
	f := field.New(name,
		field.WithValue(value),
		field.WithStored(true),
		field.WithIndexable(false),
		field.WithType(document.Text),
	)

	return &StoredField{f}
}
